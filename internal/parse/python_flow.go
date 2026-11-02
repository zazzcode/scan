package parse

import (
	"fmt"
	"sort"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/zazzcode/zazzles/internal/graph"
)

type pythonDep struct {
	dep                               graph.Dependency
	protection, container, key, entry string
}
type pythonValue struct {
	deps            []pythonDep
	gaps            []graph.Gap
	fields          map[string]pythonValue
	literal         string
	constant        bool
	binding, table  string
	driver          string
	templateEscape  string
	evidence        []graph.Evidence
	built, rendered bool
	unknown         bool
}
type pythonWrite struct {
	value    pythonValue
	entry    string
	evidence graph.Evidence
}
type pythonFlow struct {
	set     *pythonSet
	module  *pythonModule
	route   pythonRoute
	report  *graph.StringReport
	env     map[string]pythonValue
	context []graph.Evidence
	gaps    []graph.Gap
	depth   int
	work    *int
	stack   map[string]bool
	returns pythonValue
	imports map[string]pythonImport
	inputs  map[string]pythonValue
	guards  map[string][]graph.Gap
}

func pyBody(n *sitter.Node) *sitter.Node {
	if b := pyField(n, "body"); b != nil {
		return b
	}
	if n != nil {
		for i := uint(0); i < n.NamedChildCount(); i++ {
			if c := n.NamedChild(i); c.Kind() == "block" {
				return c
			}
		}
	}
	return nil
}
func (f *pythonFlow) gap(kind, detail string, n *sitter.Node) graph.Gap {
	return graph.Gap{Kind: kind, Detail: detail, At: spanOf(f.module.file, n)}
}
func (f *pythonFlow) source(n *sitter.Node, container, key string) pythonValue {
	e := f.module.evidence("request_source", n)
	v := pythonValue{deps: []pythonDep{{dep: graph.Dependency{Source: e, Transfers: []graph.Evidence{}}, protection: "raw", container: container, key: key}}}
	if key != "" && container != "session" && f.inputs != nil {
		f.inputs[container+"."+key] = v
	}
	return v
}
func (s *pythonSet) analyze(out *graph.StringReport, r pythonRoute) {
	m := r.module
	handler := m.evidence("route_handler", r.fn)
	out.Entries = append(out.Entries, graph.Entry{Route: r.route, Handler: &handler, Gaps: append([]graph.Gap{}, r.gaps...)})
	work := 0
	f := &pythonFlow{set: s, module: m, route: r, report: out, env: map[string]pythonValue{}, context: append([]graph.Evidence{handler, m.evidence("route_registration", r.decorator)}, r.receiver.evidence...), gaps: append([]graph.Gap{}, r.gaps...), work: &work, stack: map[string]bool{}}
	f.guards = map[string][]graph.Gap{}
	f.inputs = map[string]pythonValue{}
	f.imports = map[string]pythonImport{}
	f.context = append(f.context, s.formPolicies(r)...)
	f.parameters(r.fn)
	f.block(r.fn.ChildByFieldName("body"))
	f.finishValidation()
}
func (f *pythonFlow) parameters(fn *sitter.Node) {
	params := fn.ChildByFieldName("parameters")
	if params == nil {
		return
	}
	for i := uint(0); i < params.NamedChildCount(); i++ {
		p := params.NamedChild(i)
		name := p
		if p.Kind() != "identifier" {
			name = p.ChildByFieldName("name")
			if name == nil {
				for j := uint(0); j < p.NamedChildCount(); j++ {
					if p.NamedChild(j).Kind() == "identifier" {
						name = p.NamedChild(j)
						break
					}
				}
			}
		}
		if name == nil {
			continue
		}
		key := pyText(f.module, name)
		if key == "self" || key == "cls" {
			continue
		}
		if f.route.route.Framework == "fastapi" || strings.Contains(f.route.route.PathPattern, "<"+key+">") || strings.Contains(f.route.route.PathPattern, ":"+key+">") {
			container := "query"
			if strings.Contains(pyText(f.module, p), "Body(") {
				container = "body"
			}
			v := f.source(p, container, key)
			typ := pyText(f.module, p.ChildByFieldName("type"))
			if f.route.route.Framework == "fastapi" && (typ == "int" || typ == "float" || typ == "bool") {
				v = pyProtect(v, "numeric")
			}
			if strings.Contains(f.route.route.PathPattern, "<int:"+key+">") {
				v = pyProtect(v, "numeric")
			}
			f.env[key] = v
		} else {
			f.env[key] = pythonValue{}
		}
	}
}
func (f *pythonFlow) clone() *pythonFlow {
	c := *f
	c.env = map[string]pythonValue{}
	for k, v := range f.env {
		c.env[k] = v
	}
	c.gaps = append([]graph.Gap{}, f.gaps...)
	c.context = append([]graph.Evidence{}, f.context...)
	c.guards = map[string][]graph.Gap{}
	for k, v := range f.guards {
		c.guards[k] = append([]graph.Gap{}, v...)
	}
	c.returns = pythonValue{}
	c.imports = map[string]pythonImport{}
	for k, v := range f.imports {
		c.imports[k] = v
	}
	return &c
}
func (f *pythonFlow) merge(branches ...*pythonFlow) {
	keys := map[string]bool{}
	for _, b := range branches {
		f.mergeGuards(b)
		f.gaps = pythonMergeGaps(f.gaps, b.gaps)
		f.returns = pyJoin(f.returns, b.returns, 0)
		for k := range b.env {
			keys[k] = true
		}
	}
	for key := range keys {
		var v pythonValue
		for _, b := range branches {
			if value, ok := b.env[key]; ok {
				v = pyJoin(v, value, 0)
			} else {
				v = pyJoin(v, pythonValue{unknown: true}, 0)
			}
		}
		f.env[key] = v
	}
}

// block returns true when the modeled path cannot fall through.
func (f *pythonFlow) addUse(kind string, n *sitter.Node, v pythonValue) *graph.StringUse {
	v.gaps = append(v.gaps, f.guardGaps(v)...)
	deps := []graph.Dependency{}
	for _, d := range v.deps {
		deps = append(deps, d.dep)
	}
	u := graph.StringUse{ID: fmt.Sprintf("use:%s:%s:%d:%s", f.route.route.ID, f.module.file, n.StartByte(), kind), RouteID: f.route.route.ID, Kind: kind, Sink: f.module.evidence(kind, n), Dependencies: deps, Context: append([]graph.Evidence{}, f.context...), Gaps: append(append([]graph.Gap{}, f.gaps...), v.gaps...), Assessment: graph.Assessment{Status: "not_run"}}
	f.report.Uses = append(f.report.Uses, u)
	return &f.report.Uses[len(f.report.Uses)-1]
}

func (f *pythonFlow) finishValidation() {
	sunk := map[string]bool{}
	for _, u := range f.report.Uses {
		if u.RouteID != f.route.route.ID || u.Kind != "sql_query_candidate" {
			continue
		}
		for _, d := range u.Dependencies {
			for key, v := range f.inputs {
				for _, input := range v.deps {
					if input.dep.Source.Span == d.Source.Span {
						sunk[key] = true
					}
				}
			}
		}
	}
	keys := []string{}
	for key := range f.inputs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		policy := f.fieldPolicy(key)
		if sunk[key] && len(policy) == 0 {
			continue
		}
		v := f.inputs[key]
		if len(v.deps) == 0 {
			continue
		}
		sink := v.deps[0].dep.Source
		u := graph.StringUse{ID: "input:" + f.route.route.ID + ":" + key, RouteID: f.route.route.ID, Kind: "input_validation_candidate", Sink: sink, Context: append([]graph.Evidence{}, f.context...), Gaps: append([]graph.Gap{}, f.gaps...), Assessment: graph.Assessment{Status: "not_run"}, Dependencies: []graph.Dependency{v.deps[0].dep}}
		if len(policy) > 0 {
			f.validationClaim(&u, key, v, policy)
		} else {
			u.Gaps = append(u.Gaps, graph.Gap{Kind: "input_validation_unresolved", Detail: key + ": assess actual enforcing conditions in handler; no external business policy was supplied", At: sink.Span})
		}
		f.report.Uses = append(f.report.Uses, u)
	}
}
