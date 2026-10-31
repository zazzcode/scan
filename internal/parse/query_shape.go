package parse

import (
	"fmt"
	"sort"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/zazzcode/scan/internal/graph"
)

// inputFacts records which request fields were read, passed to some other call,
// or consumed by a query sink. It is shared across helper and callback flows
// for one route so a helper call is not mistaken for a validator.
type inputFacts struct {
	inputs  map[string]graph.Evidence
	touched map[string]bool
	sunk    map[string]bool
}

func newInputFacts() *inputFacts {
	return &inputFacts{inputs: map[string]graph.Evidence{}, touched: map[string]bool{}, sunk: map[string]bool{}}
}

func (f *flow) rememberInput(source graph.Evidence) {
	if f.facts == nil || source.Text == "" {
		return
	}
	if _, ok := f.facts.inputs[source.Text]; !ok {
		f.facts.inputs[source.Text] = source
	}
}

func (f *flow) noteInputTouch(values []value) {
	if f.facts == nil {
		return
	}
	for _, v := range values {
		for _, d := range v.deps {
			if d.Source.Kind == "request_source" {
				f.facts.touched[d.Source.Text] = true
			}
		}
	}
}

func (f *flow) structureGaps(kind string, args *sitter.Node, values []value, call *sitter.Node) []graph.Gap {
	gaps := []graph.Gap{}
	text := ""
	if args != nil {
		text = args.Utf8Text(f.module.src)
	}
	textDeps := len(values) > 0 && len(values[0].deps) > 0
	switch kind {
	case "sql_query_candidate":
		if textDeps || strings.Contains(text, "Op.like") {
			gaps = append(gaps, f.gap("query_text_interpolated", "request-controlled value is built into executable query text", call))
		}
	case "nosql_query_candidate":
		if strings.Contains(text, "$where") {
			gaps = append(gaps, f.gap("query_text_interpolated", "a value is built into a NoSQL $where string", call))
		} else if requestDeps(values) {
			gaps = append(gaps, f.gap("query_document_unvalidated", "request-controlled value enters a query, update, or insert document without a modeled type constraint", call))
		}
	}
	for _, field := range requestFields(values) {
		if f.facts != nil {
			f.facts.sunk[field] = true
		}
		if f.facts != nil && f.facts.touched[field] {
			gaps = append(gaps, f.gap("input_validation_unresolved", "a call receives "+field+" before the sink, and that call is not a proven validator", call))
			continue
		}
		gaps = append(gaps, f.gap("input_validation_absent", field+" reaches a sensitive use with no validator call on its path", call))
	}
	return gaps
}

// requestFieldNode returns the req.body.name / req.query.name / req.params.name
// member inside a longer chain such as req.body.handle.length.
func requestFieldNode(n *sitter.Node, src []byte, req string) *sitter.Node {
	var found *sitter.Node
	var search func(*sitter.Node)
	search = func(cur *sitter.Node) {
		if cur == nil || found != nil {
			return
		}
		if cur.Kind() != "member_expression" && cur.Kind() != "subscript_expression" {
			return
		}
		obj := cur.ChildByFieldName("object")
		prop := cur.ChildByFieldName("property")
		if obj != nil && prop != nil && (obj.Kind() == "member_expression" || obj.Kind() == "subscript_expression") {
			base := obj.ChildByFieldName("object")
			baseProp := obj.ChildByFieldName("property")
			if base != nil && baseProp != nil && base.Kind() == "identifier" && base.Utf8Text(src) == req {
				switch baseProp.Utf8Text(src) {
				case "body", "query", "params", "headers", "cookies":
					found = cur
					return
				}
			}
		}
		search(obj)
	}
	search(n)
	return found
}

func requestDeps(values []value) bool {
	for _, v := range values {
		if len(requestFields([]value{v})) > 0 {
			return true
		}
	}
	return false
}

func requestFields(values []value) []string {
	seen := map[string]bool{}
	for _, v := range values {
		for _, d := range v.deps {
			if d.Source.Kind == "request_source" && d.Source.Text != "" {
				seen[d.Source.Text] = true
			}
		}
	}
	out := make([]string, 0, len(seen))
	for field := range seen {
		out = append(out, field)
	}
	sort.Strings(out)
	return out
}

func queryKind(name string, args *sitter.Node, src []byte, values []value) string {
	if strings.HasSuffix(name, ".query") {
		return "sql_query_candidate"
	}
	text := ""
	if args != nil {
		text = args.Utf8Text(src)
	}
	if strings.Contains(text, "Op.like") && (strings.Contains(text, "${") || strings.Contains(text, "+")) {
		return "sql_query_candidate"
	}
	if !collectionWrite(name) {
		return ""
	}
	if !strings.Contains(name, "Collection") && !strings.Contains(text, "$where") && !strings.Contains(text, "$set") {
		return ""
	}
	if requestDeps(values) || strings.Contains(text, "$where") {
		return "nosql_query_candidate"
	}
	return ""
}

func collectionWrite(name string) bool {
	for _, suffix := range []string{".find", ".findOne", ".update", ".updateOne", ".updateMany", ".insert", ".insertOne", ".insertMany", ".deleteOne", ".deleteMany", ".remove", ".replaceOne"} {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}

// enterQueryCallback analyzes a nested function only when its text builds a
// query. Other callbacks stay unknown so ordinary array iterators are not sinks.
func (f *flow) enterQueryCallback(n *sitter.Node) {
	if f.depth >= 2 || n == nil {
		return
	}
	text := n.Utf8Text(f.module.src)
	if !strings.Contains(text, "$where") && !strings.Contains(text, "Op.like") && !strings.Contains(text, ".query(") {
		return
	}
	body := n.ChildByFieldName("body")
	if body == nil {
		return
	}
	child := flow{
		module:  f.module,
		set:     f.set,
		route:   f.route,
		request: f.request,
		env:     copyValues(f.env),
		context: append([]graph.Evidence{}, f.context...),
		gaps:    append([]graph.Gap{}, f.gaps...),
		depth:   f.depth + 1,
		helpers: f.helpers,
		facts:   f.facts,
	}
	child.statements(body)
	f.uses = append(f.uses, child.uses...)
}

// finishValidation reports request fields that never reached a query sink.
// Fields consumed by a sink already carry their validation gap on that use.
func (f *flow) finishValidation() {
	if f.facts == nil {
		return
	}
	fields := make([]string, 0, len(f.facts.inputs))
	for field := range f.facts.inputs {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	for _, field := range fields {
		if f.facts.sunk[field] {
			continue
		}
		source := f.facts.inputs[field]
		kind := "input_validation_absent"
		detail := field + " is read and no validator call handles it before the handler finishes"
		if f.facts.touched[field] {
			kind = "input_validation_unresolved"
			detail = "a call receives " + field + ", and that call is not a proven validator"
		}
		gap := graph.Gap{Kind: kind, Detail: detail, At: source.Span}
		f.uses = append(f.uses, graph.StringUse{
			ID:           fmt.Sprintf("input:%s:%s:%d", f.route.ID, field, source.Span.StartByte),
			RouteID:      f.route.ID,
			Kind:         "input_validation_candidate",
			Sink:         source,
			Dependencies: []graph.Dependency{{Source: source, Transfers: []graph.Evidence{}}},
			Context:      append([]graph.Evidence{}, f.context...),
			Gaps:         append(append([]graph.Gap{}, f.gaps...), gap),
			Assessment:   graph.Assessment{Status: "not_run"},
		})
	}
}
