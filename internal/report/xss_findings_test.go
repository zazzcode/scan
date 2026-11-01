package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/zazzcode/zazzles/internal/graph"
)

func TestXSSClassificationAndExportAgreement(t *testing.T) {
	for _, tt := range []struct{ protection, record, status string }{
		{"raw", "finding", "poor_practice"},
		{"encoded", "protection", "protected_for_claim"},
		{"sanitized", "protection", "protected_for_claim"},
		{"unknown", "analysis", "needs_review"},
	} {
		t.Run(tt.protection, func(t *testing.T) {
			u := graph.StringUse{ID: "x", RouteID: "r", Kind: "xss_candidate",
				Sink:       graph.Evidence{Text: `<%- value %>`, Span: graph.Span{File: "views/page.ejs", StartLine: 7}},
				XSS:        &graph.XSSFacts{Type: "reflected", EntryKind: "server", OutputContext: "html_attribute", Protection: tt.protection},
				Assessment: graph.Assessment{Status: "not_run"}}
			r := &graph.StringReport{SchemaVersion: 1, Entries: []graph.Entry{{Route: graph.Route{ID: "r", Method: "GET", PathPattern: "/page", Registration: graph.Span{File: "app.js", StartLine: 3}}}}, Uses: []graph.StringUse{u}}
			Enrich(r)
			f := r.Findings[0]
			if f.RecordType != tt.record || f.Status != tt.status || len(f.TypeIDs) != 1 || f.TypeIDs[0] != "reflected-xss" {
				t.Fatalf("classification: %+v", f)
			}
			var html bytes.Buffer
			if err := HTML(&html, r); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(html.String(), "views/page.ejs:7") || !strings.Contains(html.String(), "Reflected XSS") || !strings.Contains(html.String(), displayLabel(tt.status)) {
				t.Fatal("HTML lost XSS type, location or status")
			}
			sarif, err := json.Marshal(SARIF(r))
			if err != nil || !strings.Contains(string(sarif), "CWE-79") || !strings.Contains(string(sarif), tt.status) {
				t.Fatal("SARIF lost XSS classification")
			}
		})
	}
}
