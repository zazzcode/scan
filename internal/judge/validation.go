package judge

import "github.com/zazzcode/zazzles/internal/analysis"

const validationQuestionVersion = "validation-use-v2"

var validationQuestions = []question{
	{"field_requirement", "Does the supplied source or explicit requirement establish a concrete policy for the named field? established requires actual policy evidence, not a field name or a preference for stricter validation. Client-side form constraints establish intended policy but do not prove server enforcement."},
	{"validation", "Assess ONLY the complete enumerated constraints in declared_requirement.validation and the named field_policy role. Do not expand the policy from other fields, security preferences, or unrelated handler code. A required-only field with a rejecting missing/empty guard is sufficient even if stronger validation would be desirable. On the supplied source-to-use path, does server-side enforcement meet this precise field policy? sufficient requires the condition/helper and enforced rejection on every relevant path; insufficient requires a supported path bypassing or lacking the established requirement; uncertain requires missing relevant code or semantics. Do not confuse parameter binding or HTML escaping with field validation."},
	{"enforcement_coverage", "Is enough relevant server implementation present and structurally linked to assess enforcement? complete requires the actual guards/helpers, branch/return behavior and downstream constraints that can enforce this requirement; incomplete means a required definition, dispatch, model constraint or path is missing. Source comments and benchmark descriptions are not evidence of an expected answer."},
}

func validationUse(u analysis.CandidateUse) bool {
	return u.Kind == "input_validation_candidate" && u.Completeness != nil && u.Completeness.ClaimKind == "validation"
}
func validationChoices(name string) []map[string]string {
	switch name {
	case "field_requirement":
		return []map[string]string{{"value": "established", "description": "A concrete field policy is source-backed or explicitly supplied."}, {"value": "unknown", "description": "No concrete applicable requirement is established."}, {"value": "not_applicable", "description": "This use has no applicable field requirement."}}
	case "validation":
		return choices
	case "enforcement_coverage":
		return []map[string]string{{"value": "complete", "description": "Required implementation and path evidence are supplied."}, {"value": "incomplete", "description": "A required definition or path fact is missing."}}
	}
	return nil
}
func applyValidationPolicy(u analysis.CandidateUse, a *analysis.Assessment) {
	a.FindingStatus, a.Protection = "needs_evidence", "uncertain"
	blockers := append([]analysis.Gap{}, u.Completeness.Blockers...)
	roles := append(append([]string{}, u.Completeness.Required...), "request_source", "field_policy", "route_handler", "control_flow")
	seen := map[string]bool{}
	for _, role := range roles {
		if seen[role] {
			continue
		}
		seen[role] = true
		if len(u.Completeness.Evidence[role]) == 0 {
			blockers = append(blockers, analysis.Gap{Kind: "validation_evidence_missing", Detail: role, At: u.Sink.Span})
		}
	}
	if a.QuestionVersion != validationQuestionVersion {
		blockers = append(blockers, analysis.Gap{Kind: "validation_protocol_missing"})
	}
	a.BlockingGaps = blockers
	if len(blockers) > 0 || a.Answers["enforcement_coverage"] != "complete" || a.Answers["field_requirement"] != "established" {
		a.Reason = "The field policy or relevant enforcement implementation is incomplete; supply the named evidence rather than inferring a validation defect."
		return
	}
	switch a.Answers["validation"] {
	case "insufficient":
		a.FindingStatus, a.Protection = "poor_practice", "insufficient"
		a.Reason = "A source-backed field requirement is not enforced on a supported server path; sink protection is a separate claim."
	case "sufficient":
		a.FindingStatus, a.Protection = "protected_for_claim", "sufficient"
		a.Reason = "The supplied policy, enforcing guards/helpers and path evidence support this exact field constraint."
	default:
		a.Reason = "The supplied implementation does not support a definite validation judgment."
	}
}
func protocolVersion(u analysis.CandidateUse) string {
	if validationUse(u) {
		return validationQuestionVersion
	}
	if logUse(u) {
		return logQuestionVersion
	}
	if u.Kind == "xss_candidate" {
		return xssQuestionVersion
	}
	return questionVersion
}
