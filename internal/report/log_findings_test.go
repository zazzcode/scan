package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/zazzcode/scan/internal/graph"
)

func TestLogOutcomes(t *testing.T) {
	cases := []struct {
		name           string
		gaps           []string
		record, status string
		titleWord      string
	}{
		{"protocol credential", []string{"log_protocol_credential"}, "finding", "poor_practice", "credential"},
		{"credential beside unresolved field", []string{"log_transform_unresolved", "log_protocol_credential"}, "finding", "poor_practice", "credential"},
		{"field-name hint", []string{"log_sensitive_field_hint"}, "analysis", "needs_review", "Possible"},
		{"whole container", []string{"log_request_container"}, "analysis", "needs_review", "Possible"},
		{"unresolved masking", []string{"log_transform_unresolved"}, "analysis", "needs_review", "Possible"},
		{"hint beside non-content", []string{"log_non_content_only", "log_sensitive_field_hint"}, "analysis", "needs_review", "Possible"},
		{"non-content only", []string{"log_non_content_only"}, "protection", "protected_for_claim", "non-revealing"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			u := graph.StringUse{ID: "u", Kind: "log_sink_candidate", Assessment: graph.Assessment{Status: "not_run"}}
			for _, k := range tt.gaps {
				u.Gaps = append(u.Gaps, graph.Gap{Kind: k, Detail: k + " detail"})
			}
			u.Gaps = append(u.Gaps, graph.Gap{Kind: "log_destination_unresolved"})
			r := &graph.StringReport{Uses: []graph.StringUse{u}}
			Enrich(r)
			f := r.Findings[0]
			if f.RecordType != tt.record || f.Status != tt.status || !strings.Contains(f.Title, tt.titleWord) {
				t.Fatalf("finding=%+v", f)
			}
			if !strings.Contains(f.Protection, tt.gaps[0]+" detail") {
				t.Fatalf("logging evidence missing: %s", f.Protection)
			}
			if strings.Contains(f.Uncertainty, "construction practice") {
				t.Fatal("query wording leaked into logging uncertainty")
			}
		})
	}
}

func TestLogSARIFRule(t *testing.T) {
	r := &graph.StringReport{SchemaVersion: 1, Uses: []graph.StringUse{{ID: "u", Kind: "log_sink_candidate", Sink: graph.Evidence{Text: "console.log(x)", Span: graph.Span{File: "server.ts", StartLine: 2}}, Gaps: []graph.Gap{{Kind: "log_protocol_credential"}}, Assessment: graph.Assessment{Status: "not_run"}}}}
	data, err := json.Marshal(SARIF(r))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte(`scan/sensitive-data-logging`)) || !bytes.Contains(data, []byte(`CWE-532`)) {
		t.Fatalf("SARIF rule missing: %s", data)
	}
}

func TestRequestLineOutcomes(t *testing.T) {
	cases := []struct {
		gaps   []string
		status string
	}{
		{[]string{"log_request_line_credential"}, "poor_practice"},
		{[]string{"log_request_line_credential", "log_request_line_conditional"}, "needs_review"},
		{[]string{"log_request_line_hint"}, "needs_review"},
	}
	for _, tt := range cases {
		u := graph.StringUse{ID: "u", Kind: "log_sink_candidate", Assessment: graph.Assessment{Status: "not_run"}}
		for _, k := range tt.gaps {
			u.Gaps = append(u.Gaps, graph.Gap{Kind: k, Detail: k})
		}
		r := &graph.StringReport{Uses: []graph.StringUse{u}}
		Enrich(r)
		if r.Findings[0].Status != tt.status || !strings.Contains(r.Findings[0].Title, "URL query") || (tt.status == "needs_review" && !strings.Contains(r.Findings[0].Title, "sensitive data")) {
			t.Fatalf("%v -> %+v", tt.gaps, r.Findings[0])
		}
	}
}

func TestRouteRecordsListFindingsFirst(t *testing.T) {
	r := &graph.StringReport{SchemaVersion: 1, Entries: []graph.Entry{{Route: graph.Route{ID: "r", Method: "GET", PathPattern: "/x"}}}, Uses: []graph.StringUse{
		{ID: "a", RouteID: "r", Kind: "outbound_http_candidate", Sink: graph.Evidence{Text: "fetch(first)"}, Assessment: graph.Assessment{Status: "not_run"}},
		{ID: "b", RouteID: "r", Kind: "log_sink_candidate", Sink: graph.Evidence{Text: "use(morgan(second))"}, Gaps: []graph.Gap{{Kind: "log_request_line_credential"}}, Assessment: graph.Assessment{Status: "not_run"}},
	}}
	var b bytes.Buffer
	if err := HTML(&b, r); err != nil {
		t.Fatal(err)
	}
	html := b.String()
	if strings.Index(html, "use(morgan(second))") > strings.Index(html, "fetch(first)") {
		t.Fatal("finding is not listed before review within its route")
	}
	if r.Uses[0].ID != "a" {
		t.Fatal("presentation order mutated canonical uses")
	}
}
