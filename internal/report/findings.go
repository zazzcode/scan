package report

import (
	"strings"

	"github.com/zazzcode/zazzles/internal/graph"
)

// Enrich adds deterministic display records while preserving raw evidence and assessments.
// Rebuilding from uses also lets older saved reports use the current presentation rules.
func Enrich(r *graph.StringReport) {
	r.PresentationVersion = 4
	r.FindingTypes = append([]graph.FindingType(nil), findingTypes...)
	r.Findings = make([]graph.Finding, 0, len(r.Uses))
	for _, u := range r.Uses {
		r.Findings = append(r.Findings, describe(u))
	}
}

func describe(u graph.StringUse) graph.Finding {
	f := graph.Finding{ID: "finding:" + u.ID, TypeIDs: findingTypeIDs(u), UseID: u.ID, RouteID: u.RouteID,
		ClientID:   u.ClientID,
		RecordType: "analysis", Status: "needs_review", Title: "Sensitive operation needs review",
		Summary:    "Scan recorded a possible sensitive operation. Its security behavior has not been established.",
		Validation: "Not established from the supplied evidence.", Protection: "Not established from the supplied evidence.",
		Uncertainty: "Review the evidence for the named claim. Independent findings review and runtime confirmation have not been performed.", Recommendations: guidance(u)}
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
			f.Status = "poor_practice"
			f.Title = "A value is built into SQL query text"
			f.Summary = "Refactor the string-built query to keep values separate from SQL structure using the verified driver API. This construction is a poor-practice finding even when an exploitable injection has not been demonstrated. Field validation is a separate concern."
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
			f.Status = "poor_practice"
			f.Title = "A value is built into a NoSQL $where expression"
			f.Summary = "Replace the concatenated or interpolated $where expression with an appropriate data-only query. This executable-string construction is a poor-practice finding even when conversions or filters may prevent exploitation. Review value types and preserve the intended query behavior."
		} else if hasGap(u, "query_document_unvalidated") {
			f.RecordType = "analysis"
			f.Status = "needs_review"
			f.Title = "NoSQL document type constraints need review"
			f.Summary = "A request value reaches a query, update, or insert document. Scan has not established the scalar type or permitted keys/operators for that slot. A structured document is not automatically unsafe; verify the driver and field constraints."
		}
	case "input_validation_candidate":
		f.Title = "Request-field validation needs review"
		f.Summary = "Scan identified a request field used by the handler. Review it against the route's actual business requirements, separately from database protection."
		if hasGap(u, "input_validation_absent") {
			f.RecordType = "analysis"
			f.Status = "needs_review"
			f.Title = "Request-field validation evidence is incomplete"
		}
	case "account_creation_candidate":
		f.Title = "Account-creation requirements need review"
		f.Summary = "Check required fields and enforcing guards before persistence. A modeled account-creation operation is not proof of a validation bypass."
	case "log_sink_candidate":
		describeLog(u, &f)
	case "xss_candidate", "xss_data_candidate":
		describeXSS(u, &f)
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
			f.Validation += " Model validation answer: " + displayLabel(answer) + "."
		}
		if u.Kind == "input_validation_candidate" && u.Completeness != nil && u.Completeness.ClaimKind == "validation" {
			f.Protection = "Field-policy enforcement: " + displayLabel(a.Protection) + ". " + a.Reason
		} else if u.Kind == "xss_candidate" {
			f.Protection = "Decision XSS exposure: " + displayLabel(a.Answers["xss_exposure"]) + ". Scanner policy: " + a.Reason
		} else if u.Kind == "log_sink_candidate" {
			f.Protection += " Assessed logging exposure: " + displayLabel(a.Answers["log_exposure"]) + " (data class " + displayLabel(a.Answers["data_class"]) + ", content " + displayLabel(a.Answers["content_exposure"]) + ", redaction " + displayLabel(a.Answers["redaction"]) + "). " + a.Reason
		} else {
			f.Protection = "Assessed injection/use protection: " + displayLabel(a.Protection) + ". " + a.Reason
		}
		switch a.FindingStatus {
		case "confirmed_vulnerability":
			if u.Kind == "xss_candidate" && a.QuestionVersion == "xss-use-v1" && len(a.BlockingGaps) == 0 && a.Answers["xss_exposure"] == "vulnerable" && u.XSS != nil && u.XSS.Protection != "unknown" {
				f.RecordType, f.Status = "finding", "confirmed_vulnerability"
				f.Title = "Cross-site scripting supported by source analysis"
				f.Summary = "A complete rendering claim and consistent XSS-specific assessment establish an unsafe source-to-sink path. This is a static judgment; independent runtime confirmation remains separate."
			}
			if u.Kind == "log_sink_candidate" && len(a.BlockingGaps) == 0 && a.Answers["log_exposure"] == "vulnerable" {
				f.RecordType = "finding"
				f.Status = "confirmed_vulnerability"
				f.Title = "Sensitive data written to logs, supported by source analysis"
				f.Summary = "The supplied source and logging-specific assessment establish that credential or personal-data content reaches log output without effective redaction. Remove it from the log or log a non-revealing fact instead."
			}
			if u.Completeness != nil && u.Completeness.ClaimKind == "injection" && len(a.BlockingGaps) == 0 && (a.Answers["injection_protection"] == "vulnerable" || a.Answers["injection_protection"] == "insufficient") {
				f.RecordType = "finding"
				f.Status = "confirmed_vulnerability"
				f.Title = "SQL injection supported by source analysis"
				if u.Kind == "nosql_query_candidate" {
					f.Title = "NoSQL injection supported by source analysis"
				}
				f.Summary = "The supplied source and injection-specific assessment establish that request input can affect query syntax or selector structure in the stated analysis scope. Remediate the injection using the verified database API; the underlying construction-practice concern is included in this one finding."
			}
			if u.Completeness != nil && len(u.Completeness.Evidence["conditional_value"]) > 0 {
				f.Uncertainty += " Conditional value paths are present. This static vulnerability judgment is scoped to an unsafe supported path under its configuration/branch precondition, not every runtime configuration; inspect the conditional-expression evidence."
			}
		case "poor_practice":
			if u.Kind == "input_validation_candidate" && a.QuestionVersion == "validation-use-v2" && len(a.BlockingGaps) == 0 && a.Answers["field_requirement"] == "established" && a.Answers["enforcement_coverage"] == "complete" && a.Answers["validation"] == "insufficient" {
				f.RecordType, f.Status = "finding", "poor_practice"
				f.Title = "Declared field validation is not enforced"
				f.Summary = "A source-backed requirement for the named field is not enforced on a supported server path. This is a validation-practice defect; exploitability at a downstream operation is assessed separately."
			}
		case "supported_potential":
			if f.Status != "poor_practice" {
				f.RecordType = "analysis"
				f.Status = "supported_potential"
			}
		case "protected_for_claim":
			if u.Kind == "input_validation_candidate" && u.Completeness != nil && u.Completeness.ClaimKind == "validation" {
				f.Title = "Declared field policy is enforced"
				f.Summary = "The supplied implementation and validation assessment support enforcement of the exact named field constraints. Downstream sink protection remains a separate claim."
			}
			if f.Status != "poor_practice" {
				f.RecordType = "protection"
				f.Status = "protected_for_claim"
			}
		case "not_applicable":
			if f.Status != "poor_practice" {
				f.RecordType = "analysis"
				f.Status = "not_applicable"
			}
		case "needs_evidence":
			if f.Status != "poor_practice" {
				f.Status = "needs_evidence"
			}
		}
	}
	if f.RecordType == "protection" && (a.Status == "failed" || protectionBlocked(u)) {
		f.Status = "needs_evidence"
	}
	if f.RecordType == "protection" && f.Status != "protected_for_claim" {
		f.RecordType = "analysis"
		f.Title = "Protection evidence is inconclusive"
		f.Summary = "Some protective behavior was observed, but the complete claim remains unresolved. Review the named source, path and assessment gaps before relying on this output."
	}
	if f.Status == "poor_practice" && u.Kind == "log_sink_candidate" {
		f.Uncertainty += " Who can read the log destination and how long it is retained have not been assessed; that does not remove this remediation task."
	} else if f.Status == "poor_practice" && u.Kind != "input_validation_candidate" {
		f.Uncertainty += " Exploitability is separate from the established construction practice; a missing or inconclusive semantic assessment does not remove this remediation task."
	}
	f.Recommendations = scopedGuidance(u, f)
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

// Display must not inherit a safe label from an old assessment whose source claim
// still declares blockers. The raw model answer remains available for audit.
func protectionBlocked(u graph.StringUse) bool {
	if len(u.Assessment.BlockingGaps) > 0 {
		return true
	}
	if c := u.Completeness; c != nil {
		if len(c.Blockers) > 0 {
			return true
		}
		for _, role := range c.Required {
			if len(c.Evidence[role]) == 0 {
				return true
			}
		}
	}
	return false
}
