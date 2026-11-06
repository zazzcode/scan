package judge

import (
	"github.com/zazzcode/zazzles/internal/analysis"
	"github.com/zazzcode/zazzles/internal/analysispacks"
)

// validationUse covers standalone field-policy records: Python claims with
// assembled policy evidence and JS/TS field records that only carry gaps.
func validationUse(u analysis.CandidateUse) bool {
	return u.Kind == "input_validation_candidate"
}

func applyValidationPolicy(u analysis.CandidateUse, a *analysis.Assessment) {
	a.FindingStatus, a.Protection = "needs_evidence", "uncertain"
	blockers := []analysis.Gap{}
	roles := []string{"input_source", "field_policy", "invoked_callable", "control_flow"}
	if u.Completeness != nil {
		blockers = append(blockers, u.Completeness.Blockers...)
		roles = append(append([]string{}, u.Completeness.Required...), roles...)
	}
	if u.Components != nil {
		// Selected obligations extend, never replace, the adapter's roles.
		roles = append(roles, u.Components.Required...)
	}
	seen := map[string]bool{}
	for _, role := range roles {
		if seen[role] {
			continue
		}
		seen[role] = true
		if u.Completeness == nil || len(u.Completeness.Evidence[role]) == 0 {
			blockers = append(blockers, analysis.Gap{Kind: "validation_evidence_missing", Detail: role, At: u.Sink.Span})
		}
	}
	if a.QuestionVersion != analysispacks.ValidationProtocol {
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
	if version := analysispacks.ProtocolFor(u.Kind); version != "" {
		return version
	}
	if logUse(u) {
		return logQuestionVersion
	}
	return questionVersion
}
