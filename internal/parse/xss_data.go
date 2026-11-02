package parse

import (
	"fmt"

	sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/zazzcode/zazzles/internal/analysis"
)

// jsonProducer is a server JSON response value, kept with its fields so a browser
// consumer of the same route and response field can be linked to its origin.
type jsonProducer struct {
	route analysis.Route
	value value
	sink  analysis.Evidence
}

func (f *flow) recordXSSData(node *sitter.Node, v value, context, kind string) {
	if f.set != nil && f.route.ID != "" {
		f.set.producers = append(f.set.producers, jsonProducer{f.route, v, f.module.evidence("rendering_data", node)})
	}
	f.uses = append(f.uses, analysis.CandidateUse{
		ID:      fmt.Sprintf("render-data:%s:%s:%d", f.route.ID, f.module.file, node.StartByte()),
		RouteID: f.route.ID, Kind: "xss_data_candidate", Sink: f.module.evidence("rendering_data", node),
		Dependencies: v.deps, Context: append([]analysis.Evidence{}, f.context...), Gaps: append([]analysis.Gap{}, v.gaps...),
		XSS: &analysis.XSSFacts{Type: kind, EntryKind: "server", OutputContext: context, ContentType: "application/json",
			Protection: "unknown", ProtectionReason: "Data transfer alone is not XSS; a browser rendering sink and its protection must be linked."},
		Assessment: analysis.Assessment{Status: "not_run"},
	})
}
