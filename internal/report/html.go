package report

import (
	_ "embed"
	"html/template"
	"io"
	"sort"

	"github.com/zazzcode/scan/internal/graph"
)

//go:embed report.html
var htmlSource string
var htmlReport = template.Must(template.New("report").Funcs(template.FuncMap{
	"label": displayLabel,
}).Parse(htmlSource))

type htmlItem struct {
	Finding graph.Finding
	Use     graph.StringUse
}
type htmlRow struct {
	ID, SourceType, Signature   string
	Location                    graph.Span
	Items                       []htmlItem
	Types                       []htmlTypePill
	Gaps                        []graph.Gap
	Findings, Review, Protected int
}
type htmlView struct {
	Glossary                    []graph.FindingType
	Root                        string
	Rows                        []*htmlRow
	Gaps                        []graph.Gap
	Limitations                 []string
	Findings, Review, Protected int
}

// HTML renders a standalone snapshot with native disclosure rows and no network resources.
func HTML(w io.Writer, r *graph.StringReport) error {
	Enrich(r)
	v := htmlView{Glossary: r.FindingTypes, Root: r.Root, Gaps: r.Gaps, Limitations: r.Limitations}
	rows := map[string]*htmlRow{}
	for _, e := range r.Entries {
		row := &htmlRow{ID: e.Route.ID, SourceType: "HTTP route", Signature: e.Route.Method + " " + e.Route.PathPattern, Location: e.Route.Registration, Gaps: e.Gaps}
		rows[e.Route.ID] = row
		v.Rows = append(v.Rows, row)
	}
	for i, u := range r.Uses {
		row := rows[u.RouteID]
		if row == nil {
			row = &htmlRow{ID: u.RouteID, SourceType: "Unresolved entry", Signature: "Route association unavailable"}
			rows[u.RouteID] = row
			v.Rows = append(v.Rows, row)
		}
		f := r.Findings[i]
		row.Items = append(row.Items, htmlItem{f, u})
		switch f.RecordType {
		case "finding":
			row.Findings++
			v.Findings++
		case "protection":
			row.Protected++
			v.Protected++
		default:
			row.Review++
			v.Review++
		}
	}
	for _, row := range v.Rows {
		row.Types = rowTypePills(row.Items)
	}
	sort.SliceStable(v.Rows, func(i, j int) bool {
		a, b := v.Rows[i], v.Rows[j]
		if (a.Findings > 0) != (b.Findings > 0) {
			return a.Findings > 0
		}
		if (a.Review > 0) != (b.Review > 0) {
			return a.Review > 0
		}
		return a.Signature < b.Signature
	})
	return htmlReport.Execute(w, v)
}
func displayLabel(value string) string {
	switch value {
	case "observed_pattern":
		return "Observed code pattern"
	case "supported_potential":
		return "Supported potential issue"
	case "needs_review":
		return "Analysis needs review"
	case "needs_evidence":
		return "Evidence incomplete"
	case "protected_for_claim":
		return "Protection supported for this claim"
	case "not_applicable":
		return "Not applicable to this claim"
	case "not_run", "":
		return "Not assessed"
	case "completed":
		return "Assessment completed"
	case "failed":
		return "Assessment failed"
	}
	return value
}
