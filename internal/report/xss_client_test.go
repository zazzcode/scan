package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/zazzcode/zazzles/internal/analysis"
)

func TestClientResultsHaveSeparateHeadingsAndIdentity(t *testing.T) {
	r := &analysis.AnalysisReport{SchemaVersion: 1,
		ClientEntries: []analysis.ClientEntry{{ID: "client:view", Kind: "component", Name: "SearchComponent", Framework: "angular", Location: analysis.Span{File: "search.component.ts", StartLine: 10}}},
		Uses:          []analysis.CandidateUse{{ID: "x", ClientID: "client:view", Kind: "xss_candidate", XSS: &analysis.XSSFacts{Type: "dom", EntryKind: "client", Protection: "raw", OutputContext: "html_text", OriginResolved: true}, Sink: analysis.Evidence{Text: `[innerHTML]="value"`, Span: analysis.Span{File: "search.component.html", StartLine: 5}}, Assessment: analysis.Assessment{Status: "not_run"}}}}
	var html bytes.Buffer
	if err := HTML(&html, r); err != nil {
		t.Fatal(err)
	}
	text := html.String()
	for _, heading := range []string{"Server / API routes", "Browser / client-side", "Component / source", "Source location", "SearchComponent", "search.component.html:5"} {
		if !strings.Contains(text, heading) {
			t.Fatalf("missing %s", heading)
		}
	}
	if !strings.Contains(text, `id="hide-protected"`) || !strings.Contains(text, `data-record-type="finding"`) || !strings.Contains(text, "Showing open Findings and Review items") {
		t.Fatal("remediation filter or record identity missing")
	}
	start := strings.Index(text, `<section id="client-source-rows"`)
	end := strings.Index(text[start:], "</section>") + start
	client := text[start:end]
	for _, routeLabel := range []string{"Route signature", "Route location", "HTTP route", "Route association unavailable"} {
		if strings.Contains(client, routeLabel) {
			t.Fatalf("client results contain API label: %s", routeLabel)
		}
	}
	if r.Findings[0].ClientID != "client:view" || r.Findings[0].RouteID != "" {
		t.Fatal("finding manufactured a route identity")
	}
	sarif := SARIF(r)
	run := sarif["runs"].([]any)[0].(map[string]any)
	result := run["results"].([]any)[0].(map[string]any)
	props := result["properties"].(map[string]any)
	if props["scan.clientEntry"] != "client:view" || props["scan.route"] != nil {
		t.Fatal("SARIF manufactured a route")
	}
}
