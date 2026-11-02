package judge

import (
	"encoding/json"
	"testing"

	"github.com/zazzcode/zazzles/internal/analysis"
)

func TestClaimScopedPolicy(t *testing.T) {
	for _, tt := range []struct {
		name, answer     string
		blocker, missing bool
		want             string
	}{
		{"positive with unrelated coverage", "insufficient", false, false, "supported_potential"},
		{"negative with evidence", "sufficient", false, false, "protected_for_claim"},
		{"critical guard missing", "insufficient", true, false, "needs_evidence"},
		{"required binding missing", "insufficient", false, true, "needs_evidence"},
		{"model uncertain", "uncertain", false, false, "needs_evidence"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c := &analysis.Completeness{Required: []string{"sink_binding"}, Evidence: map[string][]analysis.Evidence{"sink_binding": {{Kind: "sink_binding", Text: "new Sequelize()"}}}, Coverage: []analysis.Gap{{Kind: "external_response_unresolved"}}}
			if tt.blocker {
				c.Blockers = []analysis.Gap{{Kind: "unknown_guard"}}
			}
			if tt.missing {
				delete(c.Evidence, "sink_binding")
			}
			use := analysis.CandidateUse{Completeness: c, Gaps: []analysis.Gap{{Kind: "external_response_unresolved"}}}
			a := analysis.Assessment{Status: "completed", Answers: map[string]string{"protection": tt.answer}}
			applyPolicy(use, &a)
			if a.FindingStatus != tt.want {
				t.Fatalf("assessment=%+v want %s", a, tt.want)
			}
		})
	}
}

func TestWireEvidenceDeduplication(t *testing.T) {
	e := analysis.Evidence{Kind: "helper", Text: "function validate(value){return value}", SHA256: "snapshot", Span: analysis.Span{File: "x.ts", StartByte: 1, EndByte: 40}}
	use := analysis.CandidateUse{Sink: e, Context: []analysis.Evidence{e, e}, Completeness: &analysis.Completeness{Evidence: map[string][]analysis.Evidence{"helper": {e}}}}
	data, err := json.Marshal(wireBundle(use))
	if err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		Evidence []analysis.Evidence `json:"evidence"`
	}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded.Evidence) != 1 {
		t.Fatalf("duplicated evidence: %s", data)
	}
}
