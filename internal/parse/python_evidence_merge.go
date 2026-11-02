package parse

import (
	"encoding/json"
	"github.com/zazzcode/zazzles/internal/graph"
)

func pythonMergeGaps(a, b []graph.Gap) []graph.Gap {
	out := []graph.Gap{}
	seen := map[graph.Gap]bool{}
	for _, g := range append(append([]graph.Gap{}, a...), b...) {
		if !seen[g] {
			seen[g] = true
			out = append(out, g)
		}
	}
	return out
}
func pythonMergeEvidence(a, b []graph.Evidence) []graph.Evidence {
	out := []graph.Evidence{}
	seen := map[string]bool{}
	for _, e := range append(append([]graph.Evidence{}, a...), b...) {
		raw, _ := json.Marshal(e)
		key := string(raw)
		if !seen[key] {
			seen[key] = true
			out = append(out, e)
		}
	}
	return out
}
