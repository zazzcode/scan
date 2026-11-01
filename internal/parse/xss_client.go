package parse

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/zazzcode/zazzles/internal/graph"
)

var indexProjection = regexp.MustCompile(`\[[^\]]*\]`)

type clientInjection struct {
	module, name string
	node         *sitter.Node
}

type clientContext struct {
	entry    graph.ClientEntry
	injected map[string]clientInjection
	methods  map[string]*sitter.Node
	fields   map[string]*sitter.Node
	template *graph.Evidence
	bindings []graph.Evidence
	angular  bool
	react    bool
}

func (s *moduleSet) browser(report *graph.StringReport) error {
	var files []string
	for path := range s.allowed {
		files = append(files, path)
	}
	sort.Strings(files)
	for _, file := range files {
		m, err := s.load(file)
		if err != nil {
			if strings.Contains(err.Error(), "2 MiB") {
				report.Gaps = append(report.Gaps, graph.Gap{Kind: "client_source_limit", Detail: "Browser source exceeds the 2 MiB parser limit", At: graph.Span{File: file}})
				continue
			}
			return err
		}
		text := string(m.src)
		if !strings.Contains(text, "document") && !strings.Contains(text, "location") && !strings.Contains(text, "@angular") && !strings.Contains(text, "react") && !strings.Contains(text, "innerHTML") && !strings.Contains(text, "postMessage") && !strings.Contains(text, "setTimeout") {
			continue
		}
		if m.tree.RootNode().HasError() {
			report.Gaps = append(report.Gaps, graph.Gap{Kind: "client_parse_error", Detail: "Browser-relevant source contains syntax errors", At: spanOf(file, m.tree.RootNode())})
			continue
		}
		base := &clientContext{entry: graph.ClientEntry{ID: "client:" + file, Kind: "script", Name: file, Location: spanOf(file, m.tree.RootNode())},
			injected: map[string]clientInjection{}, methods: map[string]*sitter.Node{}, fields: map[string]*sitter.Node{}}
		for _, imp := range m.imports {
			base.react = base.react || imp.module == "react"
		}
		if base.react {
			base.entry.Framework = "react"
		}
		if strings.HasSuffix(file, ".vue") {
			base.entry.Framework = "vue"
		}
		f := newClientFlow(s, m, base)
		f.statements(m.tree.RootNode())
		var functions []*sitter.Node
		var classes []*sitter.Node
		var collect func(*sitter.Node)
		collect = func(n *sitter.Node) {
			if n.Kind() == "class_declaration" {
				classes = append(classes, n)
				return
			}
			if n.Kind() == "function_declaration" || isFunction(n) {
				functions = append(functions, n)
				return
			}
			for i := uint(0); i < n.NamedChildCount(); i++ {
				collect(n.NamedChild(i))
			}
		}
		collect(m.tree.RootNode())
		for _, fn := range functions {
			child := newClientFlow(s, m, base)
			child.env = copyValues(f.env)
			if params, ok := functionParameters(fn, m.src); ok {
				for _, p := range params {
					child.env[p.name] = value{}
				}
			}
			child.statements(fn.ChildByFieldName("body"))
			f.uses = append(f.uses, child.uses...)
		}
		s.vueOutputs(&f)
		addClientUses(report, base, f.uses)
		for _, class := range classes {
			context := s.clientClass(m, class)
			flow := newClientFlow(s, m, context)
			body := class.ChildByFieldName("body")
			if body == nil {
				continue
			}
			for i := uint(0); i < body.NamedChildCount(); i++ {
				member := body.NamedChild(i)
				if member.Kind() != "method_definition" {
					continue
				}
				child := newClientFlow(s, m, context)
				for key, v := range flow.env {
					if strings.HasPrefix(key, "this.") {
						child.env[key] = v
					}
				}
				if params, ok := functionParameters(member, m.src); ok {
					for _, p := range params {
						child.env[p.name] = value{}
					}
				}
				child.statements(member.ChildByFieldName("body"))
				flow.uses = append(flow.uses, child.uses...)
				for key, v := range child.env {
					if strings.HasPrefix(key, "this.") {
						flow.env[key] = v
					}
				}
			}
			s.angularOutputs(&flow)
			addClientUses(report, context, flow.uses)
		}
	}
	return nil
}

func newClientFlow(s *moduleSet, m *sourceModule, client *clientContext) flow {
	return flow{module: m, set: s, client: client, env: map[string]value{}, helpers: &helperExpansion{budget: 24, stack: map[string]bool{}},
		facts: newInputFacts(), context: append([]graph.Evidence{}, client.bindings...)}
}

func addClientUses(report *graph.StringReport, client *clientContext, uses []graph.StringUse) {
	if len(uses) == 0 {
		return
	}
	// The same sink can be reached from module, function and class flows. Records
	// are merged so an earlier weak candidate never hides a later stronger one.
	index := map[string]int{}
	var merged []graph.StringUse
	for _, u := range uses {
		if u.XSS == nil {
			continue
		}
		if i, ok := index[u.ID]; ok {
			merged[i] = mergeXSSUse(merged[i], u)
			continue
		}
		index[u.ID] = len(merged)
		merged = append(merged, u)
	}
	report.Uses = append(report.Uses, merged...)
	if len(merged) > 0 {
		report.ClientEntries = append(report.ClientEntries, client.entry)
	}
}

func xssRank(x *graph.XSSFacts) int {
	switch x.Protection {
	case "raw", "ineffective":
		return 2
	case "encoded", "sanitized":
		return 0
	}
	return 1
}

// xssDependencyProtections gives each dependency's protection within its own record.
// Origin resolution stays a separate record fact and gap, not a path protection.
func xssDependencyProtections(u graph.StringUse) []string {
	out := make([]string, len(u.Dependencies))
	for i := range out {
		out[i] = u.XSS.Protection
		if i < len(u.XSS.Paths) {
			out[i] = u.XSS.Paths[i].Protection
		}
	}
	return out
}

// mergeXSSUse keeps the strongest supported claim (unsafe, then unknown, then
// protected; a resolved origin breaks ties) and retains every other record's paths,
// context and gaps as siblings. Protection is never raised by a merge.
func mergeXSSUse(a, b graph.StringUse) graph.StringUse {
	first, second := a, b
	if xssRank(b.XSS) > xssRank(a.XSS) || (xssRank(b.XSS) == xssRank(a.XSS) && b.XSS.OriginResolved && !a.XSS.OriginResolved) {
		first, second = b, a
	}
	out := first
	facts := *first.XSS
	out.XSS = &facts
	protections := xssDependencyProtections(first)
	seen := map[string]bool{}
	out.Dependencies = nil
	for _, d := range first.Dependencies {
		seen[dependencyKey(d)] = true
		out.Dependencies = append(out.Dependencies, d)
	}
	for i, d := range second.Dependencies {
		if k := dependencyKey(d); !seen[k] {
			seen[k] = true
			out.Dependencies = append(out.Dependencies, d)
			protections = append(protections, xssDependencyProtections(second)[i])
		}
	}
	out.XSS.Paths = xssPaths(protections, nil)
	evidence := map[string]bool{}
	out.Context = nil
	for _, e := range append(append([]graph.Evidence{}, first.Context...), second.Context...) {
		k := fmt.Sprintf("%s:%s:%d:%d", e.Kind, e.Span.File, e.Span.StartByte, e.Span.EndByte)
		if !evidence[k] {
			evidence[k] = true
			out.Context = append(out.Context, e)
		}
	}
	gaps := map[string]bool{}
	out.Gaps = nil
	for _, g := range append(append([]graph.Gap{}, first.Gaps...), second.Gaps...) {
		k := fmt.Sprintf("%s:%s:%s:%d:%d", g.Kind, g.Detail, g.At.File, g.At.StartByte, g.At.EndByte)
		if !gaps[k] {
			gaps[k] = true
			out.Gaps = append(out.Gaps, g)
		}
	}
	if second.XSS.Type == "cross_tier" && out.XSS.ResponsePath == "" {
		out.Gaps = append(out.Gaps, graph.Gap{Kind: "cross_tier_sibling_unresolved", Detail: "Another path to this sink starts at an unlinked API response: " + second.XSS.ResponseMethod + " " + second.XSS.ResponsePath, At: second.Sink.Span})
	}
	return out
}

func clientKey(text string) string {
	return indexProjection.ReplaceAllString(strings.ReplaceAll(text, "?.", "."), "[]")
}

func (f *flow) clientClassEvidence(node *sitter.Node) graph.Evidence {
	if node == nil {
		return graph.Evidence{}
	}
	return f.module.evidence("render_binding", node)
}

func (f *flow) clientRecord(node *sitter.Node, v value, context string, escaped bool, binding graph.Evidence) {
	before := len(f.uses)
	f.recordXSS(f.module.evidence("browser_output", node), v, context, "text/html", escaped, []graph.Evidence{binding})
	if len(f.uses) == before {
		return
	}
	u := &f.uses[len(f.uses)-1]
	if u.ClientID != "" && v.apiPath != "" {
		u.XSS.Type, u.XSS.OriginResolved = "cross_tier", false
		u.XSS.ResponsePath, u.XSS.ResponseMethod, u.XSS.ResponseField = v.apiPath, v.apiMethod, v.responseField
	}
}

func clientID(file, name string, offset uint) string {
	return fmt.Sprintf("client:%s:%s:%d", file, name, offset)
}
