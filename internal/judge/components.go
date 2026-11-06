package judge

import (
	"crypto/sha256"
	"fmt"
	"github.com/zazzcode/zazzles/internal/analysis"
	"github.com/zazzcode/zazzles/internal/analysispacks"
)

func componentQuestions(items []analysispacks.Question) []question {
	out := make([]question, 0, len(items))
	for _, q := range items {
		out = append(out, question{q.Name, q.Instructions})
	}
	return out
}

func requestDigest(use analysis.CandidateUse, payload []byte) string {
	if analysispacks.ProtocolFor(use.Kind) == "" {
		return ""
	}
	return fmt.Sprintf("%x", sha256.Sum256(payload))
}

func componentSelection(m *analysispacks.Material) *analysis.ComponentSelection {
	if m == nil {
		return nil
	}
	return m.Selection
}

func cloneComponents(s *analysis.ComponentSelection) *analysis.ComponentSelection {
	if s == nil {
		return nil
	}
	c := *s
	c.Components = append([]analysis.ComponentRef(nil), s.Components...)
	c.Required = append([]string(nil), s.Required...)
	c.Limits = append([]string(nil), s.Limits...)
	return &c
}
