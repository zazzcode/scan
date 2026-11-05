package parse

import (
	"github.com/zazzcode/zazzles/internal/analysis"
	"strings"
)

// queryConstraintReview preserves potential contradictions for the judgment gate.
// These witnesses never establish protection by themselves.
func (s *moduleSet) queryConstraintReview(use *analysis.CandidateUse) {
	c := use.Completeness
	if use.Kind == "sql_query_candidate" {
		allHex, count := true, 0
		for _, d := range use.Dependencies {
			if d.SinkSlot != "sql_text" {
				continue
			}
			count++
			hex := false
			for _, e := range d.Transfers {
				if e.Kind != "native_hash_hex" {
					continue
				}
				m := s.modules[e.Span.File]
				if m == nil {
					continue
				}
				if input, _ := nativeHashHex(m, nodeAt(m.tree.RootNode(), e.Span.StartByte, e.Span.EndByte)); input != nil {
					hex = true
				}
			}
			allHex = allHex && hex
		}
		if allHex && count > 0 {
			c.Evidence["sql_all_inputs_hex"] = []analysis.Evidence{use.Sink}
		}
	}
	if use.Kind != "nosql_query_candidate" {
		return
	}
	for _, e := range c.Evidence["control_statement"] {
		if e.Span.File != use.Sink.Span.File || e.Span.StartByte >= use.Sink.Span.StartByte || !strings.Contains(e.Text, "typeof") {
			continue
		}
		for _, d := range use.Dependencies {
			if strings.Contains(e.Text, d.Source.Text) {
				c.Evidence["selector_type_guard"] = append(c.Evidence["selector_type_guard"], e)
				break
			}
		}
	}
}
