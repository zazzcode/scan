package judge

import (
	"encoding/json"
	"github.com/zazzcode/zazzles/internal/graph"
	"testing"
)

func TestLanguageNeutralEvidenceCoverageLedger(t *testing.T) {
	for _, file := range []string{"entry.py", "entry.ts", "entry.go"} {
		t.Run(file, func(t *testing.T) {
			e := graph.Evidence{Kind: "helper_definition", Text: "relevant implementation", Span: graph.Span{File: file, StartLine: 3}, SHA256: "snapshot"}
			u := graph.StringUse{Context: []graph.Evidence{e}, Completeness: &graph.Completeness{ClaimKind: "injection", Required: []string{"helper_definitions", "sink_binding"}, Evidence: map[string][]graph.Evidence{}}}
			raw, err := json.Marshal(wireBundle(u))
			if err != nil {
				t.Fatal(err)
			}
			var b struct {
				Coverage []evidenceCoverage `json:"evidence_coverage"`
				Evidence []graph.Evidence   `json:"evidence"`
			}
			if err := json.Unmarshal(raw, &b); err != nil {
				t.Fatal(err)
			}
			states := map[string]string{}
			for _, r := range b.Coverage {
				states[r.Role] = r.State
			}
			if states["helper_definitions"] != "supplied" || states["sink_binding"] != "missing" {
				t.Fatalf("wrong coverage %v", states)
			}
			if len(u.Completeness.Evidence) != 0 {
				t.Fatal("bundle indexing altered canonical completeness")
			}
		})
	}
}
