// Package report exports analysis results without changing finding eligibility.
package report

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/zazzcode/scan/internal/graph"
)

const SchemaURI = "https://docs.oasis-open.org/sarif/sarif/v2.1.0/errata01/os/schemas/sarif-schema-2.1.0.json"

// SARIF returns an interoperable view of the canonical evidence report. Advisory
// findings, protected claims and incomplete analysis remain distinct.
func SARIF(report *graph.StringReport) map[string]any {
	Enrich(report)
	rules := []any{}
	results := []any{}
	indexes := map[string]int{}
	notifications := []any{}
	successful := true
	for i, use := range report.Uses {
		id, cwe := ruleIdentity(use)
		index, ok := indexes[id]
		g := guidance(use)
		if !ok {
			index = len(rules)
			indexes[id] = index
			tags := []string{"security"}
			if cwe != "" {
				tags = append(tags, cwe)
			}
			markdown := g.Title + "\n\n"
			for _, action := range g.Actions {
				markdown += "- " + action + "\n"
			}
			if g.Example != "" {
				markdown += "\nIllustrative example (unverified):\n\n```typescript\n" + g.Example + "\n```\n"
			}
			rules = append(rules, map[string]any{"id": id, "shortDescription": map[string]string{"text": g.Title}, "help": map[string]string{"text": g.Title + ". " + strings.Join(g.Actions, " "), "markdown": markdown}, "properties": map[string]any{"tags": tags}})
		}
		kind, level := resultKind(report.Findings[i])
		route := routeName(report, use)
		finding := report.Findings[i]
		message := fmt.Sprintf("%s: %s [%s]. %s %s", route, finding.Title, finding.Status, finding.Summary, finding.Protection)
		props := map[string]any{"scan.finding": report.Findings[i], "scan.route": route, "scan.useId": use.ID, "scan.kind": use.Kind, "scan.assessment": use.Assessment, "scan.reviewStatus": "not_run", "scan.runtimeConfirmation": "not_run", "scan.recommendations": g, "scan.coverageGaps": use.Gaps}
		if use.Check != nil {
			props["scan.check"] = use.Check
		}
		if c := use.Completeness; c != nil {
			props["scan.requiredEvidence"] = c.Required
			props["scan.blockingGaps"] = c.Blockers
			props["scan.coverageGaps"] = c.Coverage
			props["scan.scope"] = c.Scope
		}
		primary := use.Sink.Span
		if primary.File == "" {
			for _, e := range report.Entries {
				if e.Route.ID == use.RouteID {
					primary = e.Route.Registration
					break
				}
			}
		}
		result := map[string]any{"ruleId": id, "ruleIndex": index, "kind": kind, "level": level, "message": map[string]string{"text": message}, "properties": props}
		if primary.File != "" {
			loc := location(primary, g.Title)
			loc["logicalLocations"] = []any{map[string]string{"name": route, "fullyQualifiedName": route, "kind": "function"}}
			result["locations"] = []any{loc}
		}
		related := relatedLocations(use)
		if len(related) > 0 {
			result["relatedLocations"] = related
		}
		flows := codeFlows(report, use)
		if len(flows) > 0 {
			result["codeFlows"] = flows
		}
		control := []any{}
		for _, e := range use.ControlEdges {
			control = append(control, map[string]any{"kind": e.Kind, "from": location(e.From.Span, e.From.Kind), "to": location(e.To.Span, e.To.Kind)})
		}
		if len(control) > 0 {
			props["scan.controlEdges"] = control
		}
		results = append(results, result)
		if use.Assessment.Status != "completed" {
			if use.Assessment.Status == "failed" {
				successful = false
			}
			notifications = append(notifications, map[string]any{"level": "note", "message": map[string]string{"text": route + ": semantic assessment " + use.Assessment.Status}})
		}
	}
	for _, g := range report.Gaps {
		successful = false
		notifications = append(notifications, map[string]any{"level": "warning", "message": map[string]string{"text": g.Kind + ": " + g.Detail}})
	}
	invocation := map[string]any{"executionSuccessful": successful}
	if len(notifications) > 0 {
		invocation["toolExecutionNotifications"] = notifications
	}
	artifacts := []any{}
	for _, e := range report.Files {
		artifacts = append(artifacts, map[string]any{"location": map[string]string{"uri": relativeURI(e.Span.File)}, "hashes": map[string]string{"sha-256": e.SHA256}})
	}
	run := map[string]any{"tool": map[string]any{"driver": map[string]any{"name": "scan", "informationUri": "https://github.com/zazzcode/scan", "rules": rules}}, "results": results, "artifacts": artifacts, "invocations": []any{invocation}, "properties": map[string]any{"scan.findingTypes": report.FindingTypes, "scan.limitations": report.Limitations, "scan.selectedRoute": report.SelectedRoute, "scan.identityPolicy": "snapshot-local use IDs; edit-stable baseline identity is not implemented"}}
	return map[string]any{"$schema": SchemaURI, "version": "2.1.0", "runs": []any{run}}
}

func ruleIdentity(use graph.StringUse) (string, string) {
	kind := use.Kind
	if use.Check != nil {
		kind = use.Check.Kind
	}
	switch kind {
	case "sql_query", "sql_query_candidate":
		return "scan/sql-input-protection", "CWE-89"
	case "nosql_query", "nosql_query_candidate":
		return "scan/nosql-input-protection", "CWE-943"
	case "input_validation", "input_validation_candidate":
		return "scan/input-validation", "CWE-20"
	case "account_creation", "account_creation_candidate":
		return "scan/required-field-validation", "CWE-20"
	case "log_sink", "log_sink_candidate":
		return "scan/sensitive-data-logging", "CWE-532"
	default:
		return "scan/unresolved-string-use", ""
	}
}
func resultKind(f graph.Finding) (string, string) {
	if f.Status == "not_applicable" {
		return "notApplicable", "none"
	}
	switch f.RecordType {
	case "finding":
		return "fail", "warning"
	case "protection":
		return "pass", "none"
	default:
		return "review", "note"
	}
}

func routeName(report *graph.StringReport, use graph.StringUse) string {
	if use.Check != nil {
		return use.Check.Method + " " + use.Check.Path
	}
	for _, e := range report.Entries {
		if e.Route.ID == use.RouteID {
			return e.Route.Method + " " + e.Route.PathPattern
		}
	}
	return use.RouteID
}
func relativeURI(path string) string { return (&url.URL{Path: path}).String() }
func location(span graph.Span, text string) map[string]any {
	loc := map[string]any{"message": map[string]string{"text": text}}
	if span.File != "" {
		region := map[string]any{}
		if span.StartLine > 0 {
			region["startLine"] = span.StartLine
		}
		if span.EndLine >= span.StartLine && span.EndLine > 0 {
			region["endLine"] = span.EndLine
		}
		loc["physicalLocation"] = map[string]any{"artifactLocation": map[string]string{"uri": relativeURI(span.File)}, "region": region}
	}
	return loc
}
