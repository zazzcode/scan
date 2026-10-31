package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/zazzcode/scan/internal/graph"
)

func TestSARIFStatusAndRouteEvidence(t *testing.T) {
	for _, tt := range []struct{ status, kind, level string }{{"supported_potential", "fail", "warning"}, {"protected_for_claim", "pass", "none"}, {"needs_evidence", "review", "note"}} {
		t.Run(tt.status, func(t *testing.T) {
			e := graph.Evidence{Kind: "request_source", Text: "req.body.email", Span: graph.Span{File: "routes/session handler.ts", StartLine: 4, EndLine: 4, StartByte: 10, EndByte: 24}, SHA256: "snapshot"}
			use := graph.StringUse{ID: "use", RouteID: "r", Check: &graph.Check{ID: "email-sql", Method: "POST", Path: "/session", Kind: "sql_query", Fields: []string{"email"}}, Sink: e, Dependencies: []graph.Dependency{{Source: e, Transfers: []graph.Evidence{}}}, Context: []graph.Evidence{e, e}, Assessment: graph.Assessment{Status: "completed", FindingStatus: tt.status, Protection: "insufficient"}, Completeness: &graph.Completeness{Evidence: map[string][]graph.Evidence{"sink_binding": {e}}}}
			report := &graph.StringReport{Entries: []graph.Entry{{Route: graph.Route{ID: "r", Method: "POST", PathPattern: "/session", Registration: e.Span}}}, Uses: []graph.StringUse{use}, Files: []graph.Evidence{}}
			doc := SARIF(report)
			runs := doc["runs"].([]any)
			run := runs[0].(map[string]any)
			result := run["results"].([]any)[0].(map[string]any)
			if result["kind"] != tt.kind || result["level"] != tt.level {
				t.Fatalf("result=%+v", result)
			}
			if _, ok := result["fixes"]; ok {
				t.Fatal("illustrative example was mislabeled as an edit")
			}
			if _, ok := result["codeFlows"]; !ok {
				t.Fatal("source-to-use evidence missing")
			}
			if len(result["relatedLocations"].([]any)) != 1 {
				t.Fatal("duplicate related locations")
			}
			data, _ := json.Marshal(doc)
			if !bytes.Contains(data, []byte("session%20handler.ts")) || !bytes.Contains(data, []byte("POST /session")) || !bytes.Contains(data, []byte("illustrative_template")) {
				t.Fatalf("missing route/escaped source/example status: %s", data)
			}
		})
	}
}

func TestSARIFPreservesIncompleteAssessments(t *testing.T) {
	for _, status := range []string{"not_run", "failed"} {
		report := &graph.StringReport{Uses: []graph.StringUse{{Kind: "sql_query_candidate", Assessment: graph.Assessment{Status: status}}}, Gaps: []graph.Gap{}, Files: []graph.Evidence{}}
		doc := SARIF(report)
		run := doc["runs"].([]any)[0].(map[string]any)
		result := run["results"].([]any)[0].(map[string]any)
		if result["kind"] != "review" {
			t.Fatal("incomplete result was called safe or a vulnerability")
		}
		inv := run["invocations"].([]any)[0].(map[string]any)
		if len(inv["toolExecutionNotifications"].([]any)) != 1 {
			t.Fatal("assessment omission hidden")
		}
	}
}

func TestSARIFDoesNotInventRegistrationValueFlow(t *testing.T) {
	use := graph.StringUse{Kind: "account_creation_candidate", Check: &graph.Check{Method: "POST", Path: "/accounts", Kind: "account_creation"}, Assessment: graph.Assessment{Status: "completed", FindingStatus: "supported_potential"}}
	run := SARIF(&graph.StringReport{Uses: []graph.StringUse{use}})["runs"].([]any)[0].(map[string]any)
	result := run["results"].([]any)[0].(map[string]any)
	if _, ok := result["codeFlows"]; ok {
		t.Fatal("fabricated a value-flow path for a modeled persistence mapping")
	}
}
