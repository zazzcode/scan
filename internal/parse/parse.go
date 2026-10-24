// Package parse turns a JavaScript or TypeScript tree into a route graph.
//
// It records Express-style registrations and the handler chain written at the
// call. It does not resolve imports or decide whether a value is sanitized.
package parse

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unsafe"

	sitter "github.com/tree-sitter/go-tree-sitter"
	js "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	ts "github.com/tree-sitter/tree-sitter-typescript/bindings/go"

	"github.com/zazzcode/scan/internal/graph"
	"github.com/zazzcode/scan/internal/lang"
	"github.com/zazzcode/scan/internal/walk"
)

// methods are Express router methods this slice understands.
var methods = map[string]string{
	"get":    "GET",
	"post":   "POST",
	"put":    "PUT",
	"patch":  "PATCH",
	"delete": "DELETE",
	"all":    "ALL",
	"use":    "USE",
}

// Routes scans root and returns the execution graph.
func Routes(root string) (*graph.Graph, error) {
	parsed, err := walk.Walk(root, walk.Options{})
	if err != nil {
		return nil, err
	}

	parsers, err := newParsers()
	if err != nil {
		return nil, err
	}
	defer parsers.close()

	out := &graph.Graph{
		SchemaVersion: graph.SchemaVersion,
		Root:          root,
		Routes:        []graph.Route{},
		Gaps:          []graph.Gap{},
	}

	for _, file := range parsed.Files {
		if file.Language != lang.JavaScript && file.Language != lang.TypeScript {
			continue
		}
		src, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(file.Path)))
		if err != nil {
			return nil, err
		}
		parser := parsers.forPath(file.Path)
		routes, gaps, err := extractFile(parser, file.Path, src)
		if err != nil {
			return nil, err
		}
		out.Routes = append(out.Routes, routes...)
		out.Gaps = append(out.Gaps, gaps...)
	}

	sort.Slice(out.Routes, func(i, j int) bool {
		if out.Routes[i].Registration.File != out.Routes[j].Registration.File {
			return out.Routes[i].Registration.File < out.Routes[j].Registration.File
		}
		if out.Routes[i].Registration.StartByte != out.Routes[j].Registration.StartByte {
			return out.Routes[i].Registration.StartByte < out.Routes[j].Registration.StartByte
		}
		return out.Routes[i].PathPattern < out.Routes[j].PathPattern
	})
	sort.Slice(out.Gaps, func(i, j int) bool {
		if out.Gaps[i].At.File != out.Gaps[j].At.File {
			return out.Gaps[i].At.File < out.Gaps[j].At.File
		}
		if out.Gaps[i].At.StartByte != out.Gaps[j].At.StartByte {
			return out.Gaps[i].At.StartByte < out.Gaps[j].At.StartByte
		}
		return out.Gaps[i].Kind < out.Gaps[j].Kind
	})
	return out, nil
}

type parsers struct {
	js, ts, tsx *sitter.Parser
}

func newParsers() (*parsers, error) {
	jsParser, err := openParser(js.Language())
	if err != nil {
		return nil, fmt.Errorf("javascript grammar: %w", err)
	}
	tsParser, err := openParser(ts.LanguageTypescript())
	if err != nil {
		jsParser.Close()
		return nil, fmt.Errorf("typescript grammar: %w", err)
	}
	tsxParser, err := openParser(ts.LanguageTSX())
	if err != nil {
		jsParser.Close()
		tsParser.Close()
		return nil, fmt.Errorf("tsx grammar: %w", err)
	}
	return &parsers{js: jsParser, ts: tsParser, tsx: tsxParser}, nil
}

func openParser(ptr unsafe.Pointer) (*sitter.Parser, error) {
	parser := sitter.NewParser()
	if err := parser.SetLanguage(sitter.NewLanguage(ptr)); err != nil {
		parser.Close()
		return nil, err
	}
	return parser, nil
}

func (p *parsers) close() {
	p.js.Close()
	p.ts.Close()
	p.tsx.Close()
}

func (p *parsers) forPath(path string) *sitter.Parser {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".tsx":
		return p.tsx
	case ".ts", ".mts", ".cts":
		return p.ts
	default:
		return p.js
	}
}

func extractFile(parser *sitter.Parser, file string, src []byte) ([]graph.Route, []graph.Gap, error) {
	tree := parser.Parse(src, nil)
	if tree == nil {
		return nil, []graph.Gap{{
			Kind:   "parse_error",
			Detail: "parser returned no tree",
			At:     graph.Span{File: file, StartLine: 1, EndLine: 1},
		}}, nil
	}
	defer tree.Close()

	root := tree.RootNode()
	var gaps []graph.Gap
	if root.HasError() {
		gaps = append(gaps, graph.Gap{
			Kind:   "parse_error",
			Detail: "syntax error",
			At:     spanOf(file, root),
		})
	}

	locals := map[string]*sitter.Node{}
	collectFuncs(root, src, locals)

	var routes []graph.Route
	walkCalls(root, func(call *sitter.Node) {
		found, dyn := routesFromCall(file, src, call, locals)
		routes = append(routes, found...)
		gaps = append(gaps, dyn...)
	})
	if routes == nil {
		routes = []graph.Route{}
	}
	if gaps == nil {
		gaps = []graph.Gap{}
	}
	return routes, gaps, nil
}

func walkCalls(node *sitter.Node, visit func(*sitter.Node)) {
	if node == nil {
		return
	}
	if node.Kind() == "call_expression" {
		visit(node)
	}
	for i := uint(0); i < node.NamedChildCount(); i++ {
		walkCalls(node.NamedChild(i), visit)
	}
}

func collectFuncs(node *sitter.Node, src []byte, locals map[string]*sitter.Node) {
	if node == nil {
		return
	}
	switch node.Kind() {
	case "function_declaration":
		if name := node.ChildByFieldName("name"); name != nil && name.Kind() == "identifier" {
			locals[name.Utf8Text(src)] = node.ChildByFieldName("body")
		}
	case "variable_declarator":
		name := node.ChildByFieldName("name")
		value := node.ChildByFieldName("value")
		if name != nil && name.Kind() == "identifier" && value != nil && isFunction(value) {
			locals[name.Utf8Text(src)] = value.ChildByFieldName("body")
		}
	}
	for i := uint(0); i < node.NamedChildCount(); i++ {
		collectFuncs(node.NamedChild(i), src, locals)
	}
}

func isFunction(node *sitter.Node) bool {
	return node.Kind() == "function_expression" || node.Kind() == "arrow_function" || node.Kind() == "function"
}

func routesFromCall(file string, src []byte, call *sitter.Node, locals map[string]*sitter.Node) ([]graph.Route, []graph.Gap) {
	fn := call.ChildByFieldName("function")
	if fn == nil || fn.Kind() != "member_expression" {
		return nil, nil
	}
	prop := fn.ChildByFieldName("property")
	if prop == nil {
		return nil, nil
	}
	method, ok := methods[prop.Utf8Text(src)]
	if !ok {
		return nil, nil
	}
	args := call.ChildByFieldName("arguments")
	if args == nil || args.NamedChildCount() == 0 {
		return nil, nil
	}
	first := args.NamedChild(0)
	// app.use(compression()) is middleware, not a route and not a missing path.
	// app.use(prefix, router) still reports dynamic_path because prefix is a name.
	if method == "USE" && first.Kind() != "string" && first.Kind() != "array" && first.Kind() != "identifier" {
		return nil, nil
	}
	paths, dynamic, matched := pathArgs(first, src)
	// receiver.route('/path').get(handler) keeps the path on route() and the
	// handlers on the method. Any receiver works; this is not application-specific.
	handlersFrom := uint(1)
	if dynamic && method != "USE" {
		if chained, ok := chainedRoutePaths(fn, src); ok {
			paths = chained
			dynamic = false
			matched = true
			handlersFrom = 0
		}
	}
	if !matched {
		return nil, nil
	}
	at := spanOf(file, call)
	regName := fn.Utf8Text(src)
	regSpan := spanOf(file, fn)
	if handlersFrom == 0 {
		// The call node starts at the receiver, so every verb in the chain
		// would otherwise share one span. Point at the verb token instead.
		regName = "route." + prop.Utf8Text(src)
		regSpan = spanOf(file, prop)
		at.StartLine = regSpan.StartLine
		at.StartByte = regSpan.StartByte
	}
	if dynamic {
		return nil, []graph.Gap{{
			Kind:   "dynamic_path",
			Detail: "path is not a string literal starting with /",
			At:     at,
		}}
	}

	reg := graph.Step{
		Kind: "registration",
		Symbol: graph.Symbol{
			Name: regName,
			Kind: "method",
			Span: regSpan,
		},
	}
	var chain []graph.Step
	var gaps []graph.Gap
	for i := handlersFrom; i < args.NamedChildCount(); i++ {
		steps, more := argumentSteps(file, src, args.NamedChild(i), locals)
		chain = append(chain, steps...)
		gaps = append(gaps, more...)
	}

	handler := graph.Symbol{}
	for i := len(chain) - 1; i >= 0; i-- {
		if chain[i].Kind == "handler" {
			handler = chain[i].Symbol
			break
		}
	}

	routes := make([]graph.Route, 0, len(paths))
	for _, path := range paths {
		steps := make([]graph.Step, 0, 1+len(chain))
		steps = append(steps, reg)
		steps = append(steps, chain...)
		routes = append(routes, graph.Route{
			ID:           fmt.Sprintf("route:%s:%d:%s:%s", file, at.StartByte, method, path),
			Method:       method,
			PathPattern:  path,
			Framework:    "express",
			Registration: at,
			Handler:      handler,
			Steps:        steps,
		})
	}
	return routes, gaps
}

// pathArgs classifies the first argument of a possible route call.
// matched is false when the call is not a route at all.
func pathArgs(node *sitter.Node, src []byte) (paths []string, dynamic, matched bool) {
	if node == nil {
		return nil, false, false
	}
	if text, ok := stringLiteral(node, src); ok {
		if strings.HasPrefix(text, "/") {
			return []string{text}, false, true
		}
		return nil, false, false
	}
	if node.Kind() == "array" {
		var out []string
		for i := uint(0); i < node.NamedChildCount(); i++ {
			text, ok := stringLiteral(node.NamedChild(i), src)
			if !ok || !strings.HasPrefix(text, "/") {
				return nil, true, true
			}
			out = append(out, text)
		}
		if len(out) == 0 {
			return nil, false, false
		}
		return out, false, true
	}
	return nil, true, true
}

// chainedRoutePaths finds receiver.route(path).verb(...).verb(...).
// The path belongs to the route() call. Each verb call is its own route.
func chainedRoutePaths(member *sitter.Node, src []byte) ([]string, bool) {
	obj := member.ChildByFieldName("object")
	for obj != nil && obj.Kind() == "call_expression" {
		fn := obj.ChildByFieldName("function")
		if fn == nil || fn.Kind() != "member_expression" {
			return nil, false
		}
		prop := fn.ChildByFieldName("property")
		if prop != nil && prop.Utf8Text(src) == "route" {
			args := obj.ChildByFieldName("arguments")
			if args == nil || args.NamedChildCount() == 0 {
				return nil, false
			}
			paths, dynamic, matched := pathArgs(args.NamedChild(0), src)
			if !matched || dynamic || len(paths) == 0 {
				return nil, false
			}
			return paths, true
		}
		obj = fn.ChildByFieldName("object")
	}
	return nil, false
}

func argumentSteps(file string, src []byte, node *sitter.Node, locals map[string]*sitter.Node) ([]graph.Step, []graph.Gap) {
	switch node.Kind() {
	case "identifier":
		return identSteps(file, src, node, locals)
	case "call_expression":
		return callSteps(file, src, node, locals, true)
	case "arrow_function", "function_expression", "function":
		step := graph.Step{
			Kind: "handler",
			Symbol: graph.Symbol{
				Kind: "anonymous",
				Span: spanOf(file, node),
			},
		}
		return append([]graph.Step{step}, bodyCalls(file, src, node.ChildByFieldName("body"), "")...), nil
	default:
		return nil, nil
	}
}

func identSteps(file string, src []byte, node *sitter.Node, locals map[string]*sitter.Node) ([]graph.Step, []graph.Gap) {
	name := node.Utf8Text(src)
	step := graph.Step{
		Kind: "handler",
		Symbol: graph.Symbol{
			Name: name,
			Kind: "identifier",
			Span: spanOf(file, node),
		},
	}
	if body, ok := locals[name]; ok {
		return append([]graph.Step{step}, bodyCalls(file, src, body, name)...), nil
	}
	return []graph.Step{step}, []graph.Gap{{
		Kind:   "unresolved_handler",
		Detail: name,
		At:     spanOf(file, node),
	}}
}

func callSteps(file string, src []byte, node *sitter.Node, locals map[string]*sitter.Node, unwrap bool) ([]graph.Step, []graph.Gap) {
	callee := node.ChildByFieldName("function")
	name := ""
	kind := "method"
	if callee != nil {
		name = callee.Utf8Text(src)
		kind = symbolKind(callee)
	}
	step := graph.Step{
		Kind: "handler",
		Symbol: graph.Symbol{
			Name: name,
			Kind: kind,
			Span: spanOf(file, node),
		},
	}
	steps := []graph.Step{step}
	var gaps []graph.Gap
	if kind == "identifier" {
		if body, ok := locals[name]; ok {
			steps = append(steps, bodyCalls(file, src, body, name)...)
		} else {
			gaps = append(gaps, graph.Gap{
				Kind:   "unresolved_handler",
				Detail: name,
				At:     spanOf(file, node),
			})
		}
	} else if name != "" {
		gaps = append(gaps, graph.Gap{
			Kind:   "unresolved_handler",
			Detail: name,
			At:     spanOf(file, callee),
		})
	}
	// A member call with one function argument is a wrapper plus the real handler,
	// as in wrapper(handler()) or wrapper(handler). A bare call's arguments are
	// parameters, not another step in the chain.
	if unwrap && callee != nil && callee.Kind() == "member_expression" {
		if args := node.ChildByFieldName("arguments"); args != nil && args.NamedChildCount() == 1 {
			inner := args.NamedChild(0)
			var more []graph.Step
			var moreGaps []graph.Gap
			switch inner.Kind() {
			case "identifier":
				more, moreGaps = identSteps(file, src, inner, locals)
			case "call_expression":
				more, moreGaps = callSteps(file, src, inner, locals, false)
			}
			steps = append(steps, more...)
			gaps = append(gaps, moreGaps...)
		}
	}
	return steps, gaps
}

func bodyCalls(file string, src []byte, body *sitter.Node, caller string) []graph.Step {
	if body == nil {
		return nil
	}
	if body.Kind() == "call_expression" {
		callee := body.ChildByFieldName("function")
		name := ""
		if callee != nil {
			name = callee.Utf8Text(src)
		}
		return []graph.Step{{
			Kind:   "call",
			Caller: caller,
			Symbol: graph.Symbol{Name: name, Kind: symbolKind(callee), Span: spanOf(file, body)},
		}}
	}
	var steps []graph.Step
	for i := uint(0); i < body.NamedChildCount(); i++ {
		stmt := body.NamedChild(i)
		if stmt.Kind() != "expression_statement" || stmt.NamedChildCount() == 0 {
			continue
		}
		expr := stmt.NamedChild(0)
		if expr.Kind() != "call_expression" {
			continue
		}
		callee := expr.ChildByFieldName("function")
		name := ""
		if callee != nil {
			name = callee.Utf8Text(src)
		}
		steps = append(steps, graph.Step{
			Kind:   "call",
			Caller: caller,
			Symbol: graph.Symbol{Name: name, Kind: symbolKind(callee), Span: spanOf(file, expr)},
		})
	}
	return steps
}

func symbolKind(callee *sitter.Node) string {
	if callee != nil && callee.Kind() == "identifier" {
		return "identifier"
	}
	return "method"
}

func stringLiteral(node *sitter.Node, src []byte) (string, bool) {
	if node == nil || node.Kind() != "string" {
		return "", false
	}
	text := node.Utf8Text(src)
	if len(text) < 2 {
		return "", false
	}
	return text[1 : len(text)-1], true
}

func spanOf(file string, node *sitter.Node) graph.Span {
	start := node.StartPosition()
	end := node.EndPosition()
	return graph.Span{
		File:      file,
		StartLine: int(start.Row) + 1,
		EndLine:   int(end.Row) + 1,
		StartByte: int(node.StartByte()),
		EndByte:   int(node.EndByte()),
	}
}
