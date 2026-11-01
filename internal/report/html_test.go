package report

import (
	"bytes"
	"encoding/json"
	"github.com/zazzcode/scan/internal/graph"
	"strings"
	"testing"
)

func TestHumanFindingsPreserveEvidenceAndSeparateUncertainty(t *testing.T) {
	cases := []struct{ name, kind, gap, text, record, title string }{
		{"SQL text", "sql_query_candidate", "query_text_interpolated", "db.query(`SELECT ${value}`)", "finding", "SQL query text"},
		{"ORM value", "sql_query_candidate", "query_text_interpolated", "Model.findAll({[Op.like]: `%${value}%`})", "analysis", "library review"},
		{"NoSQL expression", "nosql_query_candidate", "query_text_interpolated", "collection.find({$where: 'id == '+id})", "finding", "NoSQL $where"},
		{"NoSQL document", "nosql_query_candidate", "query_document_unvalidated", "collection.update({_id: req.body.id})", "finding", "modeled type constraint"},
		{"unknown validator", "input_validation_candidate", "input_validation_unresolved", "req.body.email", "analysis", "validation needs review"},
		{"no validator observed", "input_validation_candidate", "input_validation_absent", "req.body.email", "finding", "analyzed code"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			u := graph.StringUse{ID: "u", RouteID: "r", Kind: tt.kind, Sink: graph.Evidence{Text: tt.text}, Gaps: []graph.Gap{{Kind: tt.gap}}, Assessment: graph.Assessment{Status: "not_run"}}
			before, _ := json.Marshal(u)
			r := &graph.StringReport{Uses: []graph.StringUse{u}}
			Enrich(r)
			f := r.Findings[0]
			if f.RecordType != tt.record || !strings.Contains(f.Title, tt.title) || len(f.Recommendations.Actions) == 0 || !strings.Contains(f.Uncertainty, "not run") {
				t.Fatalf("finding=%+v", f)
			}
			after, _ := json.Marshal(r.Uses[0])
			if !bytes.Equal(before, after) {
				t.Fatal("presentation changed detector evidence")
			}
			if tt.gap == "input_validation_unresolved" && !strings.Contains(f.Validation, "not been established") {
				t.Fatal(f.Validation)
			}
		})
	}
}

func TestHumanFindingAssessmentStatus(t *testing.T) {
	for _, tt := range []struct{ assessment, finding, record string }{{"completed", "protected_for_claim", "protection"}, {"completed", "supported_potential", "finding"}, {"failed", "protected_for_claim", "analysis"}, {"not_run", "supported_potential", "analysis"}} {
		u := graph.StringUse{Kind: "outbound_http_candidate", Assessment: graph.Assessment{Status: tt.assessment, FindingStatus: tt.finding}}
		f := describe(u)
		if f.RecordType != tt.record {
			t.Fatalf("%+v -> %+v", tt, f)
		}
	}
}

func TestStandaloneHTMLGroupsAndEscapesSavedEvidence(t *testing.T) {
	hostile := `<script>alert('source')</script><img src="https://example.invalid/beacon">`
	r := &graph.StringReport{SchemaVersion: 1, Root: hostile, Entries: []graph.Entry{
		{Route: graph.Route{ID: "one", Method: "GET", PathPattern: "/same", Registration: graph.Span{File: hostile, StartLine: 1}}},
		{Route: graph.Route{ID: "two", Method: "POST", PathPattern: "/same"}},
		{Route: graph.Route{ID: "empty", Method: "GET", PathPattern: "/empty"}},
	}, Uses: []graph.StringUse{
		{ID: "u1", RouteID: "one", Kind: "nosql_query_candidate", Sink: graph.Evidence{Text: hostile}, Gaps: []graph.Gap{{Kind: "query_text_interpolated"}}},
		{ID: "u2", RouteID: "two", Kind: "input_validation_candidate"},
		{ID: "u3", RouteID: "missing", Kind: "outbound_http_candidate"},
	}}
	var b bytes.Buffer
	if err := HTML(&b, r); err != nil {
		t.Fatal(err)
	}
	html := b.String()
	for _, want := range []string{"GET /same", "POST /same", "GET /empty", "Unresolved entry", "No findings reported", "Recommended changes", "$where", "default-src 'none'"} {
		if !strings.Contains(html, want) {
			t.Fatalf("missing %q", want)
		}
	}
	if strings.Count(html, `class="entry"`) != 4 || strings.Count(html, `<article class="finding">`) != 3 {
		t.Fatal("lost or duplicated rows/findings")
	}
	for _, bad := range []string{"<script", "<img", "<link", "fetch(", "XMLHttpRequest"} {
		if strings.Contains(html, bad) {
			t.Fatalf("active or external content: %s", bad)
		}
	}
	if !strings.Contains(html, "&lt;script&gt;") {
		t.Fatal("hostile source was not escaped")
	}
	var next bytes.Buffer
	if err := HTML(&next, r); err != nil {
		t.Fatal(err)
	}
	if b.String() != next.String() {
		t.Fatal("unstable output")
	}
}
