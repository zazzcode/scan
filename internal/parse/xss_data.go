package parse

import (
	"fmt"

	sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/zazzcode/zazzles/internal/graph"
)

func (f *flow) recordXSSData(node *sitter.Node, v value, context, kind string) {
	f.uses = append(f.uses, graph.StringUse{
		ID:      fmt.Sprintf("render-data:%s:%s:%d", f.route.ID, f.module.file, node.StartByte()),
		RouteID: f.route.ID, Kind: "xss_data_candidate", Sink: f.module.evidence("rendering_data", node),
		Dependencies: v.deps, Context: append([]graph.Evidence{}, f.context...), Gaps: append([]graph.Gap{}, v.gaps...),
		XSS: &graph.XSSFacts{Type: kind, EntryKind: "server", OutputContext: context, ContentType: "application/json",
			Protection: "unknown", ProtectionReason: "Data transfer alone is not XSS; a browser rendering sink and its protection must be linked."},
		Assessment: graph.Assessment{Status: "not_run"},
	})
}
