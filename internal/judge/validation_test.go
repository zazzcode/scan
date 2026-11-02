package judge

import (
	"github.com/zazzcode/zazzles/internal/analysis"
	"testing"
)

func TestValidationPolicyRequiresIndependentEvidence(t *testing.T) {
	for _, tt := range []struct {
		name, validation, coverage, requirement, version, missing, want string
		blocker                                                         bool
	}{
		{name: "defect", validation: "insufficient", coverage: "complete", requirement: "established", version: validationQuestionVersion, want: "poor_practice"},
		{name: "enforced", validation: "sufficient", coverage: "complete", requirement: "established", version: validationQuestionVersion, want: "protected_for_claim"},
		{name: "unknown helper", validation: "insufficient", coverage: "complete", requirement: "established", version: validationQuestionVersion, blocker: true, want: "needs_evidence"},
		{name: "no policy", validation: "insufficient", coverage: "complete", requirement: "unknown", version: validationQuestionVersion, want: "needs_evidence"},
		{name: "missing role", validation: "sufficient", coverage: "complete", requirement: "established", version: validationQuestionVersion, missing: "control_flow", want: "needs_evidence"},
		{name: "incomplete", validation: "insufficient", coverage: "incomplete", requirement: "established", version: validationQuestionVersion, want: "needs_evidence"},
		{name: "old protocol", validation: "insufficient", coverage: "complete", requirement: "established", version: questionVersion, want: "needs_evidence"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c := &analysis.Completeness{ClaimKind: "validation", Evidence: map[string][]analysis.Evidence{}}
			for _, role := range []string{"request_source", "field_policy", "route_handler", "control_flow"} {
				if role != tt.missing {
					c.Evidence[role] = []analysis.Evidence{{Kind: role, Text: "source evidence"}}
				}
			}
			if tt.blocker {
				c.Blockers = []analysis.Gap{{Kind: "input_helper_unresolved"}}
			}
			u := analysis.CandidateUse{Kind: "input_validation_candidate", Completeness: c}
			a := analysis.Assessment{Status: "completed", QuestionVersion: tt.version, Answers: map[string]string{"validation": tt.validation, "enforcement_coverage": tt.coverage, "field_requirement": tt.requirement}}
			applyPolicy(u, &a)
			if a.FindingStatus != tt.want {
				t.Fatalf("got %+v want %s", a, tt.want)
			}
			if protocolVersion(u) != validationQuestionVersion || len(questionsFor(u)) != 3 {
				t.Fatal("wrong protocol")
			}
		})
	}
}

func TestValidationRespectsAdditionalRequiredRoles(t *testing.T) {
	c := &analysis.Completeness{ClaimKind: "validation", Required: []string{"schema_constraints"}, Evidence: map[string][]analysis.Evidence{}}
	for _, role := range []string{"request_source", "field_policy", "route_handler", "control_flow"} {
		c.Evidence[role] = []analysis.Evidence{{Kind: role}}
	}
	a := analysis.Assessment{QuestionVersion: validationQuestionVersion, Answers: map[string]string{"field_requirement": "established", "validation": "sufficient", "enforcement_coverage": "complete"}}
	applyPolicy(analysis.CandidateUse{Kind: "input_validation_candidate", Completeness: c}, &a)
	if a.FindingStatus != "needs_evidence" || len(a.BlockingGaps) != 1 {
		t.Fatal("adapter's additional required role was ignored")
	}
}
