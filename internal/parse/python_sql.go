package parse

import (
	"regexp"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/zazzcode/zazzles/internal/graph"
)

var pythonSelect = regexp.MustCompile(`(?is)^\s*select\s+.+?\s+from\s+([a-zA-Z_][a-zA-Z0-9_]*)\b`)
var pythonInsert = regexp.MustCompile(`(?is)^\s*insert\s+into\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(([^)]+)\)\s*values\s*\(([^)]+)\)`)
var pythonUpdate = regexp.MustCompile(`(?is)^\s*update\s+([a-zA-Z_][a-zA-Z0-9_]*)\s+set\s+(.+?)\s+where\b`)

func (f *pythonFlow) sql(n, object *sitter.Node, values []pythonValue, receiver pythonValue) pythonValue {
	if len(values) == 0 {
		return receiver
	}
	text := values[0]
	all := pythonSlot(text, "sql_text")
	for _, v := range values[1:] {
		all = pyCombine(all, pythonSlot(v, "sql_parameters"))
	}

	receiver.table = ""
	if text.constant && receiver.driver == "sqlite3" {
		if match := pythonSelect.FindStringSubmatch(text.literal); len(match) > 0 {
			receiver.table = strings.ToLower(match[1])
			receiver.evidence = append(receiver.evidence, f.module.evidence("storage_query", n))
			if object != nil && object.Kind() == "identifier" {
				f.env[pyText(f.module, object)] = receiver
			}
		}
		f.storageWrite(n, text.literal, pyArgs(n), values)
	}
	if object != nil && object.Kind() == "identifier" {
		f.env[pyText(f.module, object)] = receiver
	}
	if len(all.deps) == 0 && !text.built {
		return receiver
	}
	u := f.addUse("sql_query_candidate", n, all)
	binding := append([]graph.Evidence{}, receiver.evidence...)
	binding = append(binding, f.module.evidence("sink_binding", n.ChildByFieldName("function")))
	u.Context = append(u.Context, binding...)
	u.APIModel = pythonDBAPIModel(receiver.driver)
	if text.built {
		u.Gaps = append(u.Gaps, f.gap("query_text_interpolated", "Python expression constructs executable SQL text; parameter tuple protects only separate value slots", n))
	}
	for _, d := range all.deps {
		kind := "input_validation_absent"
		if d.protection == "numeric" {
			kind = "validation_type_only"
		}
		u.Gaps = append(u.Gaps, f.gap(kind, "No complete business-policy enforcement established for "+d.container+"."+d.key, n))
	}
	c := &graph.Completeness{ClaimKind: "injection", Required: []string{"request_source", "sink_binding", "route_handler", "control_flow"}, Evidence: map[string][]graph.Evidence{"sink_binding": binding, "route_handler": {f.route.module.evidence("route_handler", f.route.fn)}, "control_flow": append([]graph.Evidence{f.module.evidence("operation", n)}, f.context...)}, Blockers: []graph.Gap{}, Coverage: append([]graph.Gap{}, u.Gaps...), Scope: "Bounded Python DB-API source-to-operation flow; SQL text and bound parameters have separate roles."}
	for _, d := range all.deps {
		if d.dep.Source.Kind == "request_source" {
			c.Evidence["request_source"] = append(c.Evidence["request_source"], d.dep.Source)
		}
	}
	c.Evidence["query_text"] = []graph.Evidence{f.module.evidence("query_text", pyArgs(n)[0])}
	if len(pyArgs(n)) > 1 {
		c.Evidence["bound_parameters"] = []graph.Evidence{f.module.evidence("bound_parameters", pyArgs(n)[1])}
	}
	for _, g := range append(append([]graph.Gap{}, f.gaps...), append(text.gaps, f.guardGaps(text)...)...) {
		switch g.Kind {
		case "python_parameter_unresolved", "python_value_unresolved", "python_call_unresolved", "python_loop_unresolved", "python_comprehension_unresolved", "python_decorator_unresolved", "work_limit", "dependency_limit", "transfer_limit", "input_helper_unresolved", "python_finally_unresolved":
			c.Blockers = append(c.Blockers, g)
		}
	}
	if text.built && len(text.deps) == 0 {
		c.Blockers = append(c.Blockers, f.gap("query_structure_unresolved", "Dynamic SQL structure/key restrictions are not resolved", n))
	}
	fields := []string{}
	for _, d := range all.deps {
		if d.key != "" {
			fields = appendUniqueString(fields, d.key)
		}
	}
	u.Check = &graph.Check{ID: "automatic-python-query", Method: f.route.route.Method, Path: f.route.route.PathPattern, Kind: "sql_query", Fields: fields, Validation: "Assess only requirements established by supplied server code and source-backed form-field policy. A browser form can declare a requirement but does not enforce it on direct HTTP requests. Missing business requirements remain uncertain.", Protection: "Request values must remain data and must not change SQL syntax. Bound values and dynamic identifiers are separate slots."}
	u.Completeness = c
	return receiver
}
func (f *pythonFlow) storageWrite(n *sitter.Node, sql string, nodes []*sitter.Node, values []pythonValue) {
	if len(nodes) < 2 {
		return
	}
	params := nodes[1]
	if params.Kind() != "tuple" && params.Kind() != "list" {
		return
	}
	table := ""
	var fields []string
	if match := pythonInsert.FindStringSubmatch(sql); len(match) > 0 {
		table = strings.ToLower(match[1])
		fields = strings.Split(match[2], ",")
	}
	if match := pythonUpdate.FindStringSubmatch(sql); len(match) > 0 {
		table = strings.ToLower(match[1])
		for _, part := range strings.Split(match[2], ",") {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) != 2 || strings.TrimSpace(kv[1]) != "?" {
				return
			}
			fields = append(fields, kv[0])
		}
	}
	for i, key := range fields {
		if uint(i) >= params.NamedChildCount() {
			break
		}
		v := f.eval(params.NamedChild(uint(i)))
		if len(v.deps) == 0 {
			continue
		}
		storage := "python-sql:" + table + "." + strings.ToLower(strings.Trim(strings.TrimSpace(key), "`\"[]"))
		v = pyTransfer(v, f.module.evidence("storage_write", n))
		f.set.writes[storage] = append(f.set.writes[storage], pythonWrite{v, f.route.route.ID, f.module.evidence("storage_write", n)})
	}
}
func (s *pythonSet) linkStorage(out *graph.StringReport) {
	for i := range out.Uses {
		u := &out.Uses[i]
		if u.XSS == nil {
			continue
		}
		var deps []graph.Dependency
		var paths []graph.XSSPath
		var related []string
		resolved := true
		stored := false
		for _, d := range u.Dependencies {
			if d.StorageKey == "" {
				deps = append(deps, d)
				continue
			}
			stored = true
			u.XSS.StorageKey = d.StorageKey
			writes := s.writes[d.StorageKey]
			if len(writes) == 0 {
				resolved = false
				deps = append(deps, d)
				continue
			}
			for _, w := range writes {
				related = appendUniqueString(related, w.entry)
				for _, source := range w.value.deps {
					nd := source.dep
					nd.StorageKey = d.StorageKey
					nd.Transfers = append(append([]graph.Evidence{}, nd.Transfers...), d.Transfers...)
					deps = append(deps, nd)
					p := u.XSS.Protection
					if source.protection == "unknown" {
						p = "unknown"
					} else if p == "raw" && (source.protection == "encoded" || source.protection == "numeric" || source.protection == "non_content") {
						p = "encoded"
					}
					paths = append(paths, graph.XSSPath{Dependency: len(deps) - 1, Entry: w.entry, Protection: p})
				}
			}
		}
		if stored {
			if len(deps) > 64 {
				deps = deps[:64]
				u.Gaps = append(u.Gaps, graph.Gap{Kind: "dependency_limit", Detail: "Stored origin expansion exceeds 64 paths", At: u.Sink.Span})
				u.XSS.Protection = "unknown"
				resolved = false
			}
			u.Dependencies = deps
			u.XSS.Type = "stored"
			u.XSS.OriginResolved = resolved
			u.XSS.RelatedEntries = related
			if len(paths) > len(deps) {
				paths = paths[:len(deps)]
			}
			u.XSS.Paths = paths
			raw, unknown := false, false
			for _, p := range paths {
				raw = raw || p.Protection == "raw"
				unknown = unknown || p.Protection == "unknown"
			}
			if !raw && unknown {
				u.XSS.Protection = "unknown"
			}
			if len(paths) > 0 && !raw && !unknown && resolved {
				u.XSS.Protection = "encoded"
			}
			if !resolved {
				u.XSS.ProtectionReason += " Stored field origin is unresolved."
			}
		}
	}
}
func (s *pythonSet) attachChecks(out *graph.StringReport, checks []graph.Check) {
	for _, check := range checks {
		matched := false
		for i := range out.Uses {
			u := &out.Uses[i]
			if check.Kind != "sql_query" || u.Kind != "sql_query_candidate" {
				continue
			}
			fieldMatch := false
			for _, d := range u.Dependencies {
				for _, field := range check.Fields {
					for _, t := range d.Transfers {
						fieldMatch = fieldMatch || (t.Kind == "field_projection" && t.Text == field)
					}
					fieldMatch = fieldMatch || d.Source.Text == field || strings.HasPrefix(d.Source.Text, field+":")
				}
			}
			if !fieldMatch {
				continue
			}
			for _, e := range out.Entries {
				if e.Route.ID == u.RouteID && e.Route.PathPattern == check.Path && (check.Method == "" || strings.EqualFold(check.Method, e.Route.Method)) {
					copy := check
					u.Check = &copy
					matched = true
				}
			}
		}
		if !matched {
			out.Gaps = append(out.Gaps, graph.Gap{Kind: "python_check_unresolved", Detail: "No supported operation matched check " + check.ID})
		}
	}
}

func pythonSlot(v pythonValue, slot string) pythonValue {
	v.deps = append([]pythonDep{}, v.deps...)
	for i := range v.deps {
		v.deps[i].dep.SinkSlot = slot
	}
	return v
}
