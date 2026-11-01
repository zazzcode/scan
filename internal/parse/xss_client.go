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
	seen := map[string]bool{}
	for _, u := range uses {
		if u.XSS == nil {
			continue
		}
		key := u.ID
		if seen[key] {
			continue
		}
		seen[key] = true
		report.Uses = append(report.Uses, u)
	}
	if len(seen) > 0 {
		report.ClientEntries = append(report.ClientEntries, client.entry)
	}
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
		u.XSS.ResponsePath, u.XSS.ResponseMethod = v.apiPath, v.apiMethod
	}
}

func clientID(file, name string, offset uint) string {
	return fmt.Sprintf("client:%s:%s:%d", file, name, offset)
}
