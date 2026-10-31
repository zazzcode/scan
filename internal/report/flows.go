package report

import (
	"fmt"
	"strings"

	"github.com/zazzcode/scan/internal/graph"
)

func relatedLocations(use graph.StringUse) []any {
	out := []any{}
	seen := map[string]bool{}
	for _, e := range use.Context {
		key := fmt.Sprintf("%s:%d:%d", e.Span.File, e.Span.StartByte, e.Span.EndByte)
		if seen[key] || e.Span.File == "" {
			continue
		}
		seen[key] = true
		loc := location(e.Span, e.Kind)
		loc["id"] = len(out)
		out = append(out, loc)
	}
	return out
}

func checkedDependency(d graph.Dependency, use graph.StringUse) bool {
	if use.Check == nil {
		return true
	}
	for _, field := range use.Check.Fields {
		if strings.HasSuffix(d.Source.Text, ".body."+field) || strings.HasSuffix(d.Source.Text, ".query."+field) || strings.HasSuffix(d.Source.Text, ".params."+field) {
			return true
		}
	}
	return false
}

func codeFlows(report *graph.StringReport, use graph.StringUse) []any {
	out := []any{}
	for _, dep := range use.Dependencies {
		if !checkedDependency(dep, use) {
			continue
		}
		records := []graph.Evidence{}
		for _, e := range report.Entries {
			if e.Route.ID == use.RouteID {
				records = append(records, graph.Evidence{Kind: "inbound_route", Span: e.Route.Registration})
				break
			}
		}
		records = append(records, dep.Source)
		records = append(records, dep.Transfers...)
		records = append(records, use.Sink)
		locations := []any{}
		seen := map[string]bool{}
		for _, e := range records {
			if e.Span.File == "" {
				continue
			}
			key := fmt.Sprintf("%s:%d:%d", e.Span.File, e.Span.StartByte, e.Span.EndByte)
			if seen[key] {
				continue
			}
			seen[key] = true
			locations = append(locations, map[string]any{"location": location(e.Span, e.Kind), "executionOrder": len(locations), "properties": map[string]any{"scan.sha256": e.SHA256}})
		}
		if len(locations) > 0 {
			out = append(out, map[string]any{"message": map[string]string{"text": "Possible local source-to-use dependency path; not a runtime trace or a complete execution path."}, "threadFlows": []any{map[string]any{"locations": locations}}})
		}
	}
	return out
}
