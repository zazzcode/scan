package judge

import "github.com/zazzcode/scan/internal/graph"

// applyPolicy uses claim-scoped completeness when supplied. Legacy candidates
// remain conservative until their dependency and binding evidence is assembled.
func applyPolicy(use graph.StringUse, a *graph.Assessment) {
	a.Protection = a.Answers["protection"]
	if use.Completeness == nil {
		if len(use.Gaps) > 0 {
			a.Protection = "uncertain"
			a.FindingStatus = "needs_evidence"
			a.BlockingGaps = use.Gaps
			a.Reason = "legacy candidate has no claim-scoped completeness evidence"
		}
		return
	}
	c := use.Completeness
	blockers := append([]graph.Gap{}, c.Blockers...)
	for _, required := range c.Required {
		if len(c.Evidence[required]) == 0 {
			blockers = append(blockers, graph.Gap{Kind: "required_evidence_missing", Detail: required, At: use.Sink.Span})
		}
	}
	a.BlockingGaps = blockers
	if len(blockers) > 0 {
		a.Protection = "uncertain"
		a.FindingStatus = "needs_evidence"
		a.Reason = "claim-relevant evidence is incomplete; model answers and unrelated coverage gaps are retained"
		return
	}
	switch a.Protection {
	case "insufficient":
		a.FindingStatus = "supported_potential"
		a.Reason = "source-backed claim has no declared blocking gaps; this is not runtime confirmation"
	case "sufficient":
		a.FindingStatus = "protected_for_claim"
		a.Reason = "assessment applies only to the stated constraint and supplied analysis scope"
	case "not_applicable":
		a.FindingStatus = "not_applicable"
	default:
		a.FindingStatus = "needs_evidence"
	}
}
