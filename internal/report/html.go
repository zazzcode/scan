package report

import (
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"html/template"
	"io"
	"sort"

	"github.com/zazzcode/zazzles/internal/analysis"
)

//go:embed report.html
var htmlSource string

//go:embed report-sort.js
var sortSource string

var htmlReport = template.Must(template.New("report").Funcs(template.FuncMap{
	"label": displayLabel,
}).Parse(htmlSource))

type htmlItem struct {
	Finding analysis.Finding
	Use     analysis.CandidateUse
}
type htmlRow struct {
	ID, SourceType, Signature   string
	Path, Method                string
	Location                    analysis.Span
	Items                       []htmlItem
	Types                       []htmlTypePill
	Gaps                        []analysis.Gap
	Findings, Review, Protected int
}
type htmlView struct {
	SortScript                  template.JS
	ScriptHash                  string
	Glossary                    []analysis.FindingType
	Root                        string
	Rows                        []*htmlRow
	ClientRows                  []*htmlRow
	Gaps                        []analysis.Gap
	Limitations                 []string
	Findings, Review, Protected int
}

// HTML renders a standalone snapshot with native disclosure rows and no network resources.
func HTML(w io.Writer, r *analysis.AnalysisReport) error {
	Enrich(r)
	v := htmlView{Glossary: r.FindingTypes, Root: r.Root, Gaps: r.Gaps, Limitations: r.Limitations}
	hash := sha256.Sum256([]byte(sortSource))
	v.ScriptHash = base64.StdEncoding.EncodeToString(hash[:])
	// Trust only the compiled-in script. Report/source values stay escaped in data attributes.
	v.SortScript = template.JS(sortSource)
	rows := map[string]*htmlRow{}
	for _, e := range r.Entries {
		row := &htmlRow{ID: e.Route.ID, SourceType: "HTTP route", Path: e.Route.PathPattern, Method: e.Route.Method, Signature: e.Route.Method + " " + e.Route.PathPattern, Location: e.Route.Registration, Gaps: e.Gaps}
		rows[e.Route.ID] = row
		v.Rows = append(v.Rows, row)
	}
	for _, e := range r.ClientEntries {
		row := &htmlRow{ID: e.ID, SourceType: "Browser " + e.Kind, Path: e.Location.File, Signature: e.Name, Location: e.Location}
		rows[e.ID] = row
		v.Rows = append(v.Rows, row)
	}
	for i, u := range r.Uses {
		identity := u.RouteID
		if u.ClientID != "" {
			identity = u.ClientID
		}
		row := rows[identity]
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
		// Within a route, show records in the same precedence as the table: Findings, Review, Protected.
		sort.SliceStable(row.Items, func(i, j int) bool { return itemRank(row.Items[i].Finding) < itemRank(row.Items[j].Finding) })
		row.Types = rowTypePills(row.Items)
	}
	sort.SliceStable(v.Rows, func(i, j int) bool {
		a, b := v.Rows[i], v.Rows[j]
		if a.Findings != b.Findings {
			return a.Findings > b.Findings
		}
		if a.Review != b.Review {
			return a.Review > b.Review
		}
		if a.Protected != b.Protected {
			return a.Protected > b.Protected
		}
		if (a.Path == "") != (b.Path == "") {
			return a.Path != ""
		}
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		return a.Method < b.Method
	})
	serverRows := make([]*htmlRow, 0, len(v.Rows))
	for _, row := range v.Rows {
		if len(row.SourceType) >= 8 && row.SourceType[:8] == "Browser " {
			v.ClientRows = append(v.ClientRows, row)
		} else {
			serverRows = append(serverRows, row)
		}
	}
	v.Rows = serverRows
	return htmlReport.Execute(w, v)
}
func displayLabel(value string) string {
	switch value {
	case "confirmed_vulnerability":
		return "Confirmed vulnerability · static analysis"
	case "poor_practice":
		return "Poor practice · remediate"
	case "observed_pattern":
		return "Legacy observation · review evidence"
	case "supported_potential":
		return "Potential vulnerability · review"
	case "needs_review":
		return "Analysis needs review"
	case "protected_for_claim":
		return "Protection supported for this claim"
	case "not_applicable":
		return "Not applicable to this claim"
	case "unknown", "uncertain", "needs_evidence":
		return "Inconclusive"
	case "not_run", "":
		return "Not assessed"
	case "completed":
		return "Assessment completed"
	case "failed":
		return "Assessment failed"
	}
	return value
}

// itemRank orders a route's records: established vulnerabilities, other findings, Review, Protected.
func itemRank(f analysis.Finding) int {
	switch {
	case f.RecordType == "finding" && f.Status == "confirmed_vulnerability":
		return 0
	case f.RecordType == "finding":
		return 1
	case f.RecordType == "protection":
		return 3
	default:
		return 2
	}
}
