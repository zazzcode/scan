package report

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/zazzcode/zazzles/internal/graph"
	htmlstd "html"
	"strings"
	"testing"
)

func TestHTMLPinsOnlyBundledSortingScript(t *testing.T) {
	r := &graph.StringReport{Entries: []graph.Entry{{Route: graph.Route{ID: "r", Method: "GET", PathPattern: `/path/"<script>bad()</script>`}}}}
	var out bytes.Buffer
	if err := HTML(&out, r); err != nil {
		t.Fatal(err)
	}
	html := out.String()
	start := strings.Index(html, "<script>")
	end := strings.Index(html, "</script>")
	if start < 0 || end < start || strings.Count(html, "<script>") != 1 {
		t.Fatal("expected one bundled script")
	}
	actual := html[start+len("<script>") : end]
	hash := sha256.Sum256([]byte(actual))
	if actual != sortSource || !strings.Contains(htmlstd.UnescapeString(html), "script-src 'sha256-"+base64.StdEncoding.EncodeToString(hash[:])+"'") {
		t.Fatal("script differs from the bundled source or CSP hash")
	}
	for _, want := range []string{`data-sort="path"`, `class="sort-button status-sort" data-sort="findings"`, `data-review="0" data-protected="0"`, `class="categories-heading"`, `data-path="/path/&#34;&lt;script&gt;bad()&lt;/script&gt;"`} {
		if !strings.Contains(html, want) {
			t.Fatalf("missing sort structure or escaped attribute: %s", want)
		}
	}
	if strings.Count(html, "data-sort=") != 2 || strings.Contains(html, "background:#f0f2f2") {
		t.Fatal("unexpected independent sort controls or shaded empty category")
	}
}

func TestInconclusiveDisplayPreservesRawAssessment(t *testing.T) {
	for _, value := range []string{"unknown", "uncertain", "needs_evidence"} {
		t.Run(value, func(t *testing.T) {
			if got := displayLabel(value); got != "Inconclusive" {
				t.Fatalf("displayLabel(%q) = %q", value, got)
			}
			u := graph.StringUse{ID: "u", Kind: "outbound_http_candidate", Assessment: graph.Assessment{
				Status: "completed", FindingStatus: "needs_evidence", Protection: value,
				Answers: map[string]string{"validation": value},
			}}
			r := &graph.StringReport{Uses: []graph.StringUse{u}}
			before, _ := json.Marshal(r.Uses)
			var b bytes.Buffer
			if err := HTML(&b, r); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(b.String(), "Assessed injection/use protection: Inconclusive.") || !strings.Contains(b.String(), "Model validation answer: Inconclusive.") {
				t.Fatal("human-facing judgment was not translated")
			}
			after, _ := json.Marshal(r.Uses)
			if !bytes.Equal(before, after) {
				t.Fatal("display label changed raw assessment")
			}
		})
	}
	if displayLabel("not_run") != "Not assessed" || displayLabel("failed") != "Assessment failed" {
		t.Fatal("execution status confused with an inconclusive judgment")
	}
}

func TestHumanFindingsPreserveEvidenceAndSeparateUncertainty(t *testing.T) {
	cases := []struct{ name, kind, gap, text, record, title string }{
		{"SQL text", "sql_query_candidate", "query_text_interpolated", "db.query(`SELECT ${value}`)", "finding", "SQL query text"},
		{"ORM value", "sql_query_candidate", "query_text_interpolated", "Model.findAll({[Op.like]: `%${value}%`})", "analysis", "library review"},
		{"NoSQL expression", "nosql_query_candidate", "query_text_interpolated", "collection.find({$where: 'id == '+id})", "finding", "NoSQL $where"},
		{"NoSQL document", "nosql_query_candidate", "query_document_unvalidated", "collection.update({_id: req.body.id})", "analysis", "type constraints need review"},
		{"unknown validator", "input_validation_candidate", "input_validation_unresolved", "req.body.email", "analysis", "validation needs review"},
		{"no validator observed", "input_validation_candidate", "input_validation_absent", "req.body.email", "analysis", "evidence is incomplete"},
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
	for _, tt := range []struct{ assessment, finding, record string }{{"completed", "protected_for_claim", "protection"}, {"completed", "supported_potential", "analysis"}, {"failed", "protected_for_claim", "analysis"}, {"not_run", "supported_potential", "analysis"}} {
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
	for _, want := range []string{"GET /same", "POST /same", "GET /empty", "Unresolved entry", "No established findings", "Remediation", "$where", "default-src 'none'"} {
		if !strings.Contains(html, want) {
			t.Fatalf("missing %q", want)
		}
	}
	if strings.Count(html, `<details class="entry`) != 4 || strings.Count(html, `<article class="finding">`) != 3 {
		t.Fatal("lost or duplicated rows/findings")
	}
	for _, bad := range []string{"<img", "<link", "fetch(", "XMLHttpRequest"} {
		if strings.Contains(html, bad) {
			t.Fatalf("active or external content: %s", bad)
		}
	}
	if strings.Count(html, "<script>") != 1 || !strings.Contains(html, "<script>"+sortSource+"</script>") {
		t.Fatal("report included a script other than the bundled sorter")
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

func TestPoorPracticeRemainsActionableAcrossModelOutcomes(t *testing.T) {
	for _, tt := range []struct{ status, finding, protection string }{
		{"not_run", "", ""}, {"failed", "", ""}, {"completed", "needs_evidence", "uncertain"}, {"completed", "supported_potential", "insufficient"}, {"completed", "protected_for_claim", "sufficient"},
	} {
		u := graph.StringUse{ID: "u", Kind: "sql_query_candidate", Sink: graph.Evidence{Text: "db.query(`SELECT ${value}`)"}, Gaps: []graph.Gap{{Kind: "query_text_interpolated"}}, Assessment: graph.Assessment{Status: tt.status, FindingStatus: tt.finding, Protection: tt.protection}}
		original, _ := json.Marshal(u.Assessment)
		r := &graph.StringReport{Uses: []graph.StringUse{u}}
		Enrich(r)
		f := r.Findings[0]
		if f.RecordType != "finding" || f.Status != "poor_practice" {
			t.Fatalf("%+v -> %+v", tt, f)
		}
		after, _ := json.Marshal(r.Uses[0].Assessment)
		if !bytes.Equal(original, after) {
			t.Fatal("changed raw assessment")
		}
		result := SARIF(r)["runs"].([]any)[0].(map[string]any)["results"].([]any)[0].(map[string]any)
		if result["kind"] != "fail" || result["level"] != "warning" {
			t.Fatal("SARIF hid actionable practice")
		}
	}
}

func TestHTMLDefaultsToGroupedCountsThenRoutePath(t *testing.T) {
	r := &graph.StringReport{Entries: []graph.Entry{
		{Route: graph.Route{ID: "z", Method: "DELETE", PathPattern: "/z"}},
		{Route: graph.Route{ID: "b", Method: "GET", PathPattern: "/b"}},
		{Route: graph.Route{ID: "a-post", Method: "POST", PathPattern: "/a"}},
		{Route: graph.Route{ID: "a-get", Method: "GET", PathPattern: "/a"}},
	}, Uses: []graph.StringUse{
		{ID: "finding", RouteID: "z", Kind: "sql_query_candidate", Gaps: []graph.Gap{{Kind: "query_text_interpolated"}}},
		{ID: "orphan", RouteID: "unresolved", Kind: "outbound_http_candidate"},
		{ID: "review", RouteID: "a-post", Kind: "outbound_http_candidate"},
		{ID: "protected", RouteID: "b", Kind: "outbound_http_candidate", Assessment: graph.Assessment{Status: "completed", FindingStatus: "protected_for_claim"}},
	}}
	var out bytes.Buffer
	if err := HTML(&out, r); err != nil {
		t.Fatal(err)
	}
	previous := -1
	for _, id := range []string{"z", "a-post", "unresolved", "b", "a-get"} {
		position := strings.Index(out.String(), `id="`+id+`"`)
		if position <= previous {
			t.Fatalf("route %s missing or out of order", id)
		}
		previous = position
	}
}

func TestConfirmedInjectionKeepsValidationAsReviewCategory(t *testing.T) {
	u := graph.StringUse{ID: "u", RouteID: "r", Kind: "sql_query_candidate", Sink: graph.Evidence{Text: "query text"}, Gaps: []graph.Gap{{Kind: "query_text_interpolated"}, {Kind: "input_validation_absent"}}, APIModel: &graph.APIModel{ID: "sequelize-raw-query-v1"}, Completeness: &graph.Completeness{ClaimKind: "injection"}, Assessment: graph.Assessment{Status: "completed", FindingStatus: "confirmed_vulnerability", Answers: map[string]string{"injection_protection": "vulnerable"}}}
	f := describe(u)
	if f.Status != "confirmed_vulnerability" || f.RecordType != "finding" {
		t.Fatalf("finding=%+v", f)
	}
	pills := rowTypePills([]htmlItem{{Finding: f, Use: u}})
	for _, p := range pills {
		if p.ID == "sql-injection" && p.Tone != "confirmed-vulnerability" {
			t.Fatal(p)
		}
		if p.ID == "input-validation" && p.Status != "review" {
			t.Fatal("injection confirmation leaked into validation", p)
		}
	}
}
