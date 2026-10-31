package parse

import (
	"fmt"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/zazzcode/scan/internal/graph"
)

type value struct {
	deps []graph.Dependency
	gaps []graph.Gap
}
type flow struct {
	module  *sourceModule
	set     *moduleSet
	route   graph.Route
	request string
	env     map[string]value
	context []graph.Evidence
	gaps    []graph.Gap
	uses    []graph.StringUse
	work    int
}

func (f *flow) gap(kind, detail string, n *sitter.Node) graph.Gap {
	return graph.Gap{Kind: kind, Detail: detail, At: spanOf(f.module.file, n)}
}

func (f *flow) statements(body *sitter.Node) {
	if body == nil {
		return
	}
	if body.Kind() != "statement_block" {
		f.eval(body)
		return
	}
	for i := uint(0); i < body.NamedChildCount(); i++ {
		n := body.NamedChild(i)
		switch n.Kind() {
		case "lexical_declaration", "variable_declaration":
			for j := uint(0); j < n.NamedChildCount(); j++ {
				d := n.NamedChild(j)
				name := d.ChildByFieldName("name")
				if name == nil {
					continue
				}
				v := f.eval(d.ChildByFieldName("value"))
				if name.Kind() == "identifier" {
					f.env[name.Utf8Text(f.module.src)] = f.transfer(v, "assignment", d)
					if name.Utf8Text(f.module.src) == f.request {
						f.request = ""
					}
				} else {
					if name.Kind() == "object_pattern" {
						for k := uint(0); k < name.NamedChildCount(); k++ {
							binding := name.NamedChild(k)
							if binding.Kind() == "shorthand_property_identifier_pattern" {
								f.env[binding.Utf8Text(f.module.src)] = f.transfer(v, "field_projection", binding)
							} else {
								f.gaps = append(f.gaps, f.gap("unsupported_binding", "destructuring alias/default not modeled", binding))
							}
						}
					} else {
						f.gaps = append(f.gaps, f.gap("unsupported_binding", "binding not modeled", d))
					}
				}
			}
		case "expression_statement", "return_statement", "throw_statement":
			if n.NamedChildCount() > 0 {
				f.eval(n.NamedChild(0))
			}
			if n.Kind() != "expression_statement" {
				return
			}
		case "function_declaration", "comment", "empty_statement":
		default:
			// Preserve branch evidence without claiming a resolved control-flow path.
			f.gaps = append(f.gaps, f.gap("unsupported_control_flow", "statement not traversed: "+n.Kind(), n))
			// Explore branch bodies as possible paths, preserving the explicit gap.
			f.context = append(f.context, f.module.evidence("control_statement", n))
			if condition := n.ChildByFieldName("condition"); condition != nil {
				f.eval(condition)
			}
			saved := copyValues(f.env)
			for j := uint(0); j < n.NamedChildCount(); j++ {
				c := n.NamedChild(j)
				if c.Kind() == "statement_block" {
					f.env = copyValues(saved)
					f.statements(c)
				}
			}
			// Keep possible dependencies across unsupported control flow, explicitly
			// uncertain. This is a may-depend relationship, not a proven witness.
			merged := copyValues(saved)
			for key, v := range f.env {
				merged[key] = combine(merged[key], v)
			}
			for key, v := range merged {
				v.gaps = append(v.gaps, f.gap("unknown_control_effect", "value may change across unsupported control flow", n))
				merged[key] = v
			}
			f.env = merged
		}
	}
}

func (f *flow) eval(n *sitter.Node) value {
	if n == nil {
		return value{}
	}
	f.work++
	if f.work > 10000 {
		gap := f.gap("work_limit", "handler expression work limit reached", n)
		if f.work == 10001 {
			f.gaps = append(f.gaps, gap)
		}
		return value{gaps: []graph.Gap{gap}}
	}
	src := f.module.src
	switch n.Kind() {
	case "identifier":
		if v, ok := f.env[n.Utf8Text(src)]; ok {
			return v
		}
		return value{gaps: []graph.Gap{f.gap("unknown_value", "unresolved local value: "+n.Utf8Text(src), n)}}
	case "member_expression", "subscript_expression":
		if isRequestField(n, src, f.request) {
			return value{deps: []graph.Dependency{{Source: f.module.evidence("request_source", n), Transfers: []graph.Evidence{}}}}
		}
		v := f.children(n)
		return f.transfer(v, "property_read", n)
	case "assignment_expression":
		left := n.ChildByFieldName("left")
		v := f.transfer(f.eval(n.ChildByFieldName("right")), "assignment", n)
		if left != nil && left.Kind() == "identifier" {
			f.env[left.Utf8Text(src)] = v
			if left.Utf8Text(src) == f.request {
				f.request = ""
			}
		} else {
			v.gaps = append(v.gaps, f.gap("unsupported_write", "heap/property assignment is not modeled", n))
		}
		return v
	case "call_expression":
		return f.call(n)
	case "arrow_function", "function_expression", "function":
		return value{gaps: []graph.Gap{f.gap("unsupported_callback", "callback/closure body not traversed", n)}}
	case "number", "string", "true", "false", "null", "undefined", "property_identifier", "type_identifier", "comment", "string_fragment", "escape_sequence":
		return value{}
	case "ternary_expression":
		a := f.eval(n.ChildByFieldName("consequence"))
		b := f.eval(n.ChildByFieldName("alternative"))
		v := combine(a, b)
		v.gaps = append(v.gaps, f.gap("conditional_merge", "branch feasibility is not modeled", n))
		return f.transfer(v, "conditional_value", n)
	case "await_expression":
		v := f.children(n)
		v.gaps = append(v.gaps, f.gap("async_response_unresolved", "await response/continuation semantics not modeled", n))
		return f.transfer(v, "await", n)
	default:
		return f.transfer(f.children(n), "expression", n)
	}
}

func isRequestField(n *sitter.Node, src []byte, req string) bool {
	if n.Kind() != "member_expression" && n.Kind() != "subscript_expression" {
		return false
	}
	obj := n.ChildByFieldName("object")
	prop := n.ChildByFieldName("property")
	if obj != nil && prop != nil && obj.Kind() == "identifier" && obj.Utf8Text(src) == req {
		switch prop.Utf8Text(src) {
		case "body", "query", "params", "headers", "cookies":
			return true
		}
	}
	base := n
	for base != nil && (base.Kind() == "member_expression" || base.Kind() == "subscript_expression") {
		obj := base.ChildByFieldName("object")
		if obj != nil && obj.Kind() == "member_expression" {
			root := obj.ChildByFieldName("object")
			prop := obj.ChildByFieldName("property")
			if root != nil && prop != nil && root.Kind() == "identifier" && root.Utf8Text(src) == req {
				switch prop.Utf8Text(src) {
				case "body", "query", "params", "headers", "cookies":
					return true
				}
			}
		}
		base = obj
	}
	return false
}

func (f *flow) children(n *sitter.Node) value {
	v := value{}
	for i := uint(0); i < n.NamedChildCount(); i++ {
		v = combine(v, f.eval(n.NamedChild(i)))
	}
	return v
}

func (f *flow) transfer(v value, kind string, n *sitter.Node) value {
	out := value{gaps: append([]graph.Gap{}, v.gaps...)}
	for _, d := range v.deps {
		if len(d.Transfers) >= 64 {
			out.gaps = append(out.gaps, f.gap("transfer_limit", "dependency exceeds 64 transfer steps", n))
			out.deps = append(out.deps, d)
			continue
		}
		d.Transfers = append(append([]graph.Evidence{}, d.Transfers...), f.module.evidence(kind, n))
		out.deps = append(out.deps, d)
	}
	return out
}

func combine(a, b value) value {
	a.deps = append(append([]graph.Dependency{}, a.deps...), b.deps...)
	a.gaps = append(append([]graph.Gap{}, a.gaps...), b.gaps...)
	if len(a.deps) > 64 {
		a.deps = a.deps[:64]
		a.gaps = append(a.gaps, graph.Gap{Kind: "dependency_limit", Detail: "more than 64 dependencies; remaining paths unexamined", At: a.deps[0].Source.Span})
	}
	if len(a.gaps) > 64 {
		a.gaps = append(a.gaps[:63], graph.Gap{Kind: "gap_limit", Detail: "additional uncertainty reasons omitted", At: a.gaps[0].At})
	}
	return a
}

func (f *flow) call(n *sitter.Node) value {
	fn := n.ChildByFieldName("function")
	args := n.ChildByFieldName("arguments")
	if fn == nil {
		return value{}
	}
	name := fn.Utf8Text(f.module.src)
	// An outer chained call must still visit its receiver, including a query call.
	if fn.Kind() == "member_expression" {
		obj := fn.ChildByFieldName("object")
		if obj != nil && obj.Kind() == "call_expression" {
			f.eval(obj)
		}
	}
	var values []value
	if args != nil {
		for i := uint(0); i < args.NamedChildCount(); i++ {
			values = append(values, f.eval(args.NamedChild(i)))
		}
	}
	kind := ""
	if strings.HasSuffix(name, ".query") {
		kind = "sql_query_candidate"
	}
	if name == "fetch" || name == "axios" || strings.HasPrefix(name, "axios.") {
		kind = "outbound_http_candidate"
	}
	if kind != "" {
		v := value{}
		for _, a := range values {
			v = combine(v, a)
		}
		if len(v.deps) > 0 || len(v.gaps) > 0 || kind == "outbound_http_candidate" {
			gaps := append(append([]graph.Gap{}, f.gaps...), v.gaps...)
			gaps = append(gaps, f.gap("sink_binding_unverified", "candidate API shape; library binding/argument roles require review", n))
			if kind == "outbound_http_candidate" {
				gaps = append(gaps, f.gap("external_response_unresolved", "remote response and asynchronous continuation are not modeled", n))
			}
			context := append([]graph.Evidence{}, f.context...)
			f.uses = append(f.uses, graph.StringUse{ID: fmt.Sprintf("use:%s:%s:%d", f.route.ID, f.module.file, n.StartByte()), RouteID: f.route.ID, Kind: kind, Sink: f.module.evidence(kind, n), Dependencies: v.deps, Context: context, Gaps: gaps, Assessment: graph.Assessment{Status: "not_run"}})
		}
		// Response values are new external sources, never aliases of request arguments.
		return value{gaps: []graph.Gap{f.gap("unknown_call_result", "sink/transport return value is not modeled", n)}}
	}
	v := value{}
	for _, a := range values {
		v = combine(v, a)
	}
	// Method results may depend on their receiver (for example input.trim()).
	if fn.Kind() == "member_expression" {
		obj := fn.ChildByFieldName("object")
		if obj != nil && obj.Kind() != "call_expression" {
			v = combine(v, f.eval(obj))
		}
	}
	if len(v.deps) > 0 {
		v.gaps = append(v.gaps, f.gap("unknown_transfer", "call may transform or discard its arguments: "+name, n))
		if hm, helper := f.set.resolve(f.module, name); hm != nil && helper != nil && !shadowed(n, name, f.module.src) {
			f.context = append(f.context, hm.evidence("helper_definition", helper))
		}
	}
	if len(v.deps) == 0 {
		v.gaps = append(v.gaps, f.gap("unknown_call_result", "call result is not modeled: "+name, n))
	}
	return f.transfer(v, "call_argument_dependency", n)
}

func copyValues(in map[string]value) map[string]value {
	out := map[string]value{}
	for k, v := range in {
		out[k] = v
	}
	return out
}
