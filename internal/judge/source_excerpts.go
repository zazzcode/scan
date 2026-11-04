package judge

import (
	"fmt"
	"sort"

	"github.com/zazzcode/zazzles/internal/analysis"
)

type wireEvidence struct {
	analysis.Evidence
	TextRef string `json:"text_ref,omitempty"`
}

// Source text is shared without merging evidence roles or dropping containing
// controls. Every referenced excerpt reconstructs the original bytes and span.
func sharedSourceExcerpts(table []analysis.Evidence) ([]wireEvidence, []analysis.Evidence) {
	order := make([]int, 0, len(table))
	for i, e := range table {
		if e.Text != "" && e.SHA256 != "" && e.Span.File != "" && len(e.Text) == e.Span.EndByte-e.Span.StartByte {
			order = append(order, i)
		}
	}
	sort.SliceStable(order, func(a, b int) bool {
		x, y := table[order[a]], table[order[b]]
		if x.Span.File != y.Span.File {
			return x.Span.File < y.Span.File
		}
		if x.SHA256 != y.SHA256 {
			return x.SHA256 < y.SHA256
		}
		if x.Span.StartByte != y.Span.StartByte {
			return x.Span.StartByte < y.Span.StartByte
		}
		return x.Span.EndByte > y.Span.EndByte
	})
	excerpts := []analysis.Evidence{}
	refs := map[int]int{}
	for _, i := range order {
		e := table[i]
		matched := false
		if len(excerpts) > 0 {
			j := len(excerpts) - 1
			p := &excerpts[j]
			if p.Span.File == e.Span.File && p.SHA256 == e.SHA256 && e.Span.StartByte >= p.Span.StartByte && e.Span.StartByte <= p.Span.EndByte {
				overlap := p.Span.EndByte - e.Span.StartByte
				if overlap > len(e.Text) {
					overlap = len(e.Text)
				}
				start := e.Span.StartByte - p.Span.StartByte
				if p.Text[start:start+overlap] == e.Text[:overlap] {
					if e.Span.EndByte > p.Span.EndByte {
						p.Text += e.Text[overlap:]
						p.Span.EndByte = e.Span.EndByte
						p.Span.EndLine = e.Span.EndLine
					}
					refs[i] = j
					matched = true
				}
			}
		}
		if !matched {
			e.Kind = "source_excerpt"
			refs[i] = len(excerpts)
			excerpts = append(excerpts, e)
		}
	}
	out := make([]wireEvidence, len(table))
	for i, e := range table {
		out[i].Evidence = e
		if j, ok := refs[i]; ok {
			out[i].Text = ""
			out[i].TextRef = fmt.Sprintf("s%d", j)
		}
	}
	return out, excerpts
}
