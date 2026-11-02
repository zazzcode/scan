package report

import (
	"github.com/zazzcode/zazzles/internal/graph"
	"strings"
	"testing"
)

func TestSourceSpecificSQLiteRemediation(t *testing.T) {
	for _, tt := range []struct{ name, source, wantSQL, wantParams string }{
		{"fstring with LIKE", `query = f"SELECT name FROM files WHERE owner = {owner} AND name LIKE '%{search}%'"`, "owner = ? AND name LIKE ?", `owner, "%" + str(search) + "%"`},
		{"percent value", `cur.execute("select * from accounts where name = '%s'" % name)`, "name = ?", "name,"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			u := graph.StringUse{Kind: "sql_query_candidate", APIModel: &graph.APIModel{Library: "sqlite3"}, Sink: graph.Evidence{Text: tt.source, Span: graph.Span{File: "app.py", StartLine: 18}}, Gaps: []graph.Gap{{Kind: "query_text_interpolated"}}}
			f := describe(u)
			g := f.Recommendations
			if !strings.Contains(g.Example, tt.wantSQL) || !strings.Contains(g.Example, tt.wantParams) || strings.Contains(g.Example, "sequelize") {
				t.Fatalf("wrong driver/source repair %+v", g)
			}
		})
	}
	if _, _, ok := pythonSQLValueCandidate(`q=f"SELECT * FROM {table}"`); ok {
		t.Fatal("dynamic table name treated as a bindable value")
	}
}
func TestProtectedAndReviewDoNotPrescribeUnprovenFixes(t *testing.T) {
	protected := graph.StringUse{Kind: "xss_candidate", Sink: graph.Evidence{Text: "{{ user.team }}", Span: graph.Span{File: "profile.html", StartLine: 46}}, XSS: &graph.XSSFacts{Protection: "encoded", OutputContext: "html_text", OriginResolved: true}}
	f := describe(protected)
	if f.Recommendations.Example != "" || !strings.Contains(f.Recommendations.Title, "Preserve") {
		t.Fatal("Protected output receives repair template")
	}
	protected.Assessment = graph.Assessment{Status: "completed", FindingStatus: "needs_evidence", Protection: "uncertain"}
	f = describe(protected)
	if !strings.Contains(f.Recommendations.Title, "Resolve") || f.Recommendations.Example != "" {
		t.Fatal("Review output prescribes a patch")
	}
}
func TestRawTemplateAndDOMRepairsNameExactOperation(t *testing.T) {
	for _, source := range []string{`{{ user.country | safe }}`, `$comments.append('<li>' + comment + '</li>')`} {
		u := graph.StringUse{Kind: "xss_candidate", Sink: graph.Evidence{Text: source, Span: graph.Span{File: "view.html", StartLine: 12}}, XSS: &graph.XSSFacts{Type: "reflected", Protection: "raw", OutputContext: "html_text", OriginResolved: true}}
		g := describe(u).Recommendations
		if g.Example == "" || !strings.Contains(strings.Join(g.Actions, " "), "view.html:12") {
			t.Fatalf("source-specific repair missing %+v", g)
		}
	}
}

func TestProtectedDisplayCannotOverrideCanonicalBlocker(t *testing.T) {
	u := graph.StringUse{Kind: "xss_candidate", XSS: &graph.XSSFacts{Protection: "encoded", OutputContext: "html_text"}, Completeness: &graph.Completeness{Blockers: []graph.Gap{{Kind: "unresolved_helper"}}}, Assessment: graph.Assessment{Status: "completed", FindingStatus: "protected_for_claim"}}
	f := describe(u)
	if f.RecordType == "protection" || f.Status != "needs_evidence" || f.Title == "Context-appropriate rendering protection is present" {
		t.Fatalf("unsafe display: %+v", f)
	}
}
