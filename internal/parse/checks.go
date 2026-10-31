package parse

import (
	"fmt"
	"strings"

	"github.com/zazzcode/scan/internal/graph"
)

func (s *moduleSet) buildChecks(report *graph.StringReport, checks []graph.Check) {
	selected := []graph.StringUse{}
	for _, check := range checks {
		s.contextLimits = nil
		s.contextFields = check.Fields
		context, chainGaps := s.chainEvidence(report, check)
		if len(context) == 0 {
			report.Gaps = append(report.Gaps, graph.Gap{Kind: "check_route_missing", Detail: check.ID + ": declared method/path was not resolved"})
			continue
		}
		if check.Kind == "sql_query" {
			matched := 0
			for _, old := range report.Uses {
				if old.Kind != "sql_query_candidate" || !checkRoute(report, old.RouteID, check) {
					continue
				}
				matched++
				use := old
				use.ID += "/" + check.ID
				use.Check = &check
				use.Context = append(context, use.Context...)
				use.Completeness = &graph.Completeness{Required: []string{"request_source", "sql_query_candidate", "route_handler", "sink_binding"}, Evidence: evidenceRoles(use.Context), Blockers: append([]graph.Gap{}, chainGaps...), Coverage: []graph.Gap{}, Scope: "declared check over the supplied route and local source-to-SQL evidence"}
				use.Completeness.Evidence["sql_query_candidate"] = []graph.Evidence{use.Sink}
				for _, d := range use.Dependencies {
					if fieldDependency(d, check.Fields) {
						use.Completeness.Evidence["request_source"] = append(use.Completeness.Evidence["request_source"], d.Source)
					}
				}
				controlComplete := false
				for _, e := range old.Context {
					if e.Kind != "handler" {
						continue
					}
					module := s.modules[e.Span.File]
					if module == nil {
						continue
					}
					body := nodeAt(module.tree.RootNode(), e.Span.StartByte, e.Span.EndByte)
					if body == nil {
						continue
					}
					edges, missing := controlContext(module, body)
					use.ControlEdges = append(use.ControlEdges, edges...)
					use.Completeness.Blockers = append(use.Completeness.Blockers, missing...)
					controlComplete = len(missing) == 0
					use.Completeness.Evidence["control_flow"] = []graph.Evidence{e}
				}
				binding := s.queryBinding(use)
				if len(binding) > 0 {
					use.Context = append(use.Context, binding...)
					use.Completeness.Evidence["sink_binding"] = binding
				}
				for _, g := range use.Gaps {
					switch g.Kind {
					case "unsupported_control_flow":
						if controlComplete {
							use.Completeness.Coverage = append(use.Completeness.Coverage, g)
						} else {
							use.Completeness.Blockers = append(use.Completeness.Blockers, g)
						}
					case "middleware_unresolved", "query_text_interpolated", "query_document_unvalidated", "input_validation_absent", "input_validation_unresolved", "input_validation_middleware", "validation_insufficient_for_sink", "validation_type_only", "sink_must_bind_free_text", "validation_free_text_bound":
						// Query shape and missing validation stay visible. They do not
						// convert a clear insufficient protection answer into an abstention.
						use.Completeness.Coverage = append(use.Completeness.Coverage, g)

					case "unknown_value":
						// Unknown result-model/options metadata is not a dependency of SQL text.
						if gapTouchesSQLSource(g, use, check.Fields) {
							use.Completeness.Blockers = append(use.Completeness.Blockers, g)
						} else {
							use.Completeness.Coverage = append(use.Completeness.Coverage, g)
						}
					case "unknown_transfer":
						if gapTouchesSQLSource(g, use, check.Fields) {
							use.Completeness.Blockers = append(use.Completeness.Blockers, g)
						} else {
							use.Completeness.Coverage = append(use.Completeness.Coverage, g)
						}
					case "sink_binding_unverified":
						if len(binding) == 0 {
							use.Completeness.Blockers = append(use.Completeness.Blockers, g)
						} else {
							use.Completeness.Coverage = append(use.Completeness.Coverage, g)
						}
					default:
						use.Completeness.Blockers = append(use.Completeness.Blockers, g)
					}
				}
				use.Completeness.Blockers = append(use.Completeness.Blockers, s.contextLimits...)
				selected = append(selected, use)
			}
			if matched == 0 {
				selected = append(selected, missingSQLSink(report, check, context, chainGaps))
			}
		} else if check.Kind == "account_creation" {
			use := s.registrationCheck(report, check, context, chainGaps)
			use.Completeness.Blockers = append(use.Completeness.Blockers, s.contextLimits...)
			selected = append(selected, use)
		}
	}
	for _, old := range report.Uses {
		if old.Kind == "input_validation_candidate" {
			selected = append(selected, old)
		}
	}
	report.Uses = selected
}

// missingSQLSink keeps a declared SQL check visible when discovery found no query.
// An empty use list would be indistinguishable from a protected route.
func missingSQLSink(report *graph.StringReport, c graph.Check, context []graph.Evidence, chainGaps []graph.Gap) graph.StringUse {
	var route graph.Route
	blockers := append([]graph.Gap{}, chainGaps...)
	for _, e := range report.Entries {
		if e.Route.Method != c.Method || e.Route.PathPattern != c.Path {
			continue
		}
		route = e.Route
		for _, g := range e.Gaps {
			switch g.Kind {
			case "helper_call_limit", "recursive_helper", "work_limit":
				blockers = append(blockers, g)
			}
		}
	}
	gap := graph.Gap{Kind: "sql_sink_missing", Detail: c.ID + ": declared SQL check produced no candidate sink", At: route.Registration}
	blockers = append(blockers, gap)
	return graph.StringUse{
		ID:           fmt.Sprintf("check:%s:%s:missing-sink", route.ID, c.ID),
		RouteID:      route.ID,
		Kind:         "sql_sink_missing",
		Check:        &c,
		Context:      context,
		Gaps:         []graph.Gap{gap},
		Assessment:   graph.Assessment{Status: "not_run"},
		Completeness: &graph.Completeness{Required: []string{"request_source", "sql_query_candidate", "route_handler", "sink_binding"}, Evidence: evidenceRoles(context), Blockers: blockers, Coverage: []graph.Gap{}, Scope: "declared SQL check produced no candidate sink in the supported handler and helper slice"},
	}
}

func checkRoute(report *graph.StringReport, id string, c graph.Check) bool {
	for _, e := range report.Entries {
		if e.Route.ID == id {
			return e.Route.Method == c.Method && e.Route.PathPattern == c.Path
		}
	}
	return false
}

func gapTouchesSQLSource(g graph.Gap, use graph.StringUse, fields []string) bool {
	for _, d := range use.Dependencies {
		if !fieldDependency(d, fields) {
			continue
		}
		for _, e := range append([]graph.Evidence{d.Source}, d.Transfers...) {
			// Enclosing concatenations include sibling operands; only an actual
			// transfer on this source's path makes a call-result gap relevant.
			if e.Kind != "request_source" && e.Kind != "call_argument_dependency" {
				continue
			}
			if g.At.File == e.Span.File && g.At.StartByte >= e.Span.StartByte && g.At.EndByte <= e.Span.EndByte {
				return true
			}
		}
	}
	return false
}

func (s *moduleSet) registrationCheck(report *graph.StringReport, c graph.Check, context []graph.Evidence, gaps []graph.Gap) graph.StringUse {
	var route graph.Route
	for _, e := range report.Entries {
		if e.Route.Method == c.Method && e.Route.PathPattern == c.Path {
			route = e.Route
			break
		}
	}
	use := graph.StringUse{ID: fmt.Sprintf("check:%s:%s", route.ID, c.ID), RouteID: route.ID, Kind: "account_creation_candidate", Check: &c, Context: context, Dependencies: []graph.Dependency{}, Gaps: []graph.Gap{}, Assessment: graph.Assessment{Status: "not_run"}}
	use.Completeness = &graph.Completeness{Required: []string{"route_handler", "persistence_mapping", "model_constraints"}, Evidence: evidenceRoles(context), Blockers: append([]graph.Gap{}, gaps...), Coverage: []graph.Gap{}, Scope: "registration middleware and generated persistence under the declared framework model"}
	m := s.modules[route.Registration.File]
	generated, models := s.persistenceEvidence(m, c.Path)
	use.Context = append(use.Context, generated...)
	use.Context = append(use.Context, models...)
	use.Completeness.Evidence["persistence_mapping"] = generated
	use.Completeness.Evidence["model_constraints"] = models
	if len(generated) > 0 {
		use.Sink = generated[len(generated)-1]
	} else {
		use.Completeness.Blockers = append(use.Completeness.Blockers, graph.Gap{Kind: "persistence_unresolved", Detail: "no supported framework creation mapping for this route", At: route.Registration})
	}
	// Each supported middleware body contributes explicit branch/successor edges.
	for _, e := range context {
		if e.Kind != "route_handler" {
			continue
		}
		module := s.modules[e.Span.File]
		if module == nil {
			continue
		}
		body := nodeAt(module.tree.RootNode(), e.Span.StartByte, e.Span.EndByte)
		if body == nil {
			continue
		}
		edges, missing := controlContext(module, body)
		use.ControlEdges = append(use.ControlEdges, edges...)
		use.Completeness.Blockers = append(use.Completeness.Blockers, missing...)
		use.Completeness.Evidence["control_flow"] = append(use.Completeness.Evidence["control_flow"], e)
	}
	use.Completeness.Required = append(use.Completeness.Required, "control_flow")
	use.Gaps = append(use.Gaps, use.Completeness.Blockers...)
	return use
}

func (s *moduleSet) queryBinding(use graph.StringUse) []graph.Evidence {
	m := s.modules[use.Sink.Span.File]
	if m == nil {
		return nil
	}
	call := nodeAt(m.tree.RootNode(), use.Sink.Span.StartByte, use.Sink.Span.EndByte)
	if call == nil {
		return nil
	}
	fn := call.ChildByFieldName("function")
	if fn == nil {
		return nil
	}
	obj := fn.ChildByFieldName("object")
	if obj == nil {
		return nil
	}
	name := obj.Utf8Text(m.src)
	parts := strings.Split(name, ".")
	target := m
	if len(parts) == 2 {
		imp, ok := m.imports[parts[0]]
		if !ok || imp.name != "*" {
			return nil
		}
		target = s.importModule(m, imp.module)
		if target == nil {
			return nil
		}
		name = parts[1]
		if !exportedName(target, name) {
			return nil
		}
	}
	return sequelizeBinding(target, name, 0)
}
