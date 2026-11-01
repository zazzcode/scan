package report

import (
	"strings"

	"github.com/zazzcode/scan/internal/graph"
)

// Enrich adds deterministic display records while preserving raw evidence and assessments.
// Rebuilding from uses also lets older saved reports use the current presentation rules.
func Enrich(r *graph.StringReport) {
	r.PresentationVersion = 1
	r.Findings = make([]graph.Finding, 0, len(r.Uses))
	for _, u := range r.Uses {
		r.Findings = append(r.Findings, describe(u))
	}
}

func describe(u graph.StringUse) graph.Finding {
	f := graph.Finding{ID: "finding:" + u.ID, UseID: u.ID, RouteID: u.RouteID,
		RecordType: "analysis", Status: "needs_review", Title: "Sensitive operation needs review",
		Summary:    "Scan recorded a possible sensitive operation. Its security behavior has not been established.",
		Validation: "Not established from the supplied evidence.", Protection: "Not established from the supplied evidence.",
		Uncertainty: "This is bounded static analysis, not an execution trace or a confirmed exploit. Missing library behavior, control flow, and middleware effects remain unresolved.", Recommendations: guidance(u)}
	validationDetails := []string{}
	seenValidation := map[string]bool{}
	for _, g := range u.Gaps {
		detail := ""
		switch g.Kind {
		case "input_validation_absent":
			detail = "No validator call was identified for the named field in the analyzed handler and modeled middleware. This does not prove that every execution path lacks validation."
		case "input_validation_unresolved", "input_validation_middleware":
			detail = "The value reaches a call or middleware whose validation effect has not been established. Inspect the actual check and whether it runs before this use."
		case "validation_insufficient_for_sink", "validation_type_only":
			detail = "A type, conversion, or length restriction was observed. Whether it meets the field's business requirements still needs review; it does not establish query protection."
		case "sink_must_bind_free_text", "validation_free_text_bound":
			detail = "The field name suggests free text. Confirm that business policy; allowing punctuation is not by itself a validation defect."
		}
		if detail != "" && !seenValidation[g.Detail] {
			validationDetails = append(validationDetails, detail+" Evidence: "+g.Detail)
			seenValidation[g.Detail] = true
		}
	}
	if len(validationDetails) > 0 {
		f.Validation = strings.Join(validationDetails, " ")
	}
	switch u.Kind {
	case "sql_query_candidate":
		f.Title = "SQL database call needs binding review"
		f.Summary = "Scan identified a database-call pattern. Verify the driver and how it separates query structure from values."
		if hasGap(u, "query_text_interpolated") {
			f.RecordType = "finding"
			f.Status = "observed_pattern"
			f.Title = "A value is built into SQL query text"
			f.Summary = "The analyzed code builds database query text from a value rather than showing a separate bound parameter. This construction needs review independently of input validation."
			if strings.Contains(u.Sink.Text, "Op.like") {
				f.RecordType = "analysis"
				f.Status = "needs_review"
				f.Title = "Interpolated database search value needs library review"
				f.Summary = "A value is inserted into a database search expression. Scan has not established its external source or whether the database library safely separates the value from SQL structure. This is not a confirmed SQL-injection finding."
			}
		}
	case "nosql_query_candidate":
		f.Title = "NoSQL database operation needs review"
		if hasGap(u, "query_text_interpolated") {
			f.RecordType = "finding"
			f.Status = "observed_pattern"
			f.Title = "A value is built into a NoSQL $where expression"
			f.Summary = "The code constructs a $where expression using concatenation or interpolation. Review the value's origin and every applicable transformation; this observation alone does not establish an exploitable injection."
		} else if hasGap(u, "query_document_unvalidated") {
			f.RecordType = "finding"
			f.Status = "observed_pattern"
			f.Title = "Request data enters a NoSQL document without a modeled type constraint"
			f.Summary = "A request value reaches a query, update, or insert document. Scan has not established the scalar type or permitted keys/operators for that slot. A structured document is not automatically unsafe; verify the driver and field constraints."
		}
	case "input_validation_candidate":
		f.Title = "Request-field validation needs review"
		f.Summary = "Scan identified a request field used by the handler. Review it against the route's actual business requirements, separately from database protection."
		if hasGap(u, "input_validation_absent") {
			f.RecordType = "finding"
			f.Status = "observed_pattern"
			f.Title = "No validator call identified for a request field in the analyzed code"
		}
	case "account_creation_candidate":
		f.Title = "Account-creation requirements need review"
		f.Summary = "Check required fields and enforcing guards before persistence. A modeled account-creation operation is not proof of a validation bypass."
	case "outbound_http_candidate":
		f.Title = "Outbound request needs destination review"
		f.Summary = "Scan recorded a possible outbound HTTP use. Verify the destination, request-value handling, and applicable restrictions; this is not a confirmed SSRF finding."
	}
	if (u.Kind == "sql_query_candidate" || u.Kind == "nosql_query_candidate") && len(u.Dependencies) == 0 {
		f.Uncertainty += " No direct request-field dependency was resolved for this use; external control is not established."
	}
	a := u.Assessment
	if a.Status == "failed" {
		f.Uncertainty += " Semantic assessment failed: " + a.Reason
	} else if a.Status == "not_run" || a.Status == "" {
		f.Uncertainty += " Semantic assessment was not run."
	}
	if a.Status == "completed" {
		if answer := a.Answers["validation"]; answer != "" {
			f.Validation += " Model validation answer: " + answer + "."
		}
		f.Protection = "Assessed use protection: " + a.Protection + ". " + a.Reason
		switch a.FindingStatus {
		case "supported_potential":
			f.RecordType = "finding"
			f.Status = "supported_potential"
		case "protected_for_claim":
			f.RecordType = "protection"
			f.Status = "protected_for_claim"
		case "not_applicable":
			f.RecordType = "analysis"
			f.Status = "not_applicable"
		case "needs_evidence":
			f.Status = "needs_evidence"
		}
	}
	return f
}
func hasGap(u graph.StringUse, kind string) bool {
	for _, g := range u.Gaps {
		if g.Kind == kind {
			return true
		}
	}
	return false
}
