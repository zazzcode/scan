package parse

import (
	"fmt"

	"github.com/zazzcode/zazzles/internal/graph"
)

func assembleXSSClaims(report *graph.StringReport) {
	for i := range report.Uses {
		u := &report.Uses[i]
		if u.Kind != "xss_candidate" || u.XSS == nil {
			continue
		}
		c := &graph.Completeness{ClaimKind: "xss", Required: []string{"render_source", "render_binding", "render_context"},
			Evidence: map[string][]graph.Evidence{}, Scope: "Bounded source-to-rendering flow; output context and protections are scoped to this sink.",
			Coverage: append([]graph.Gap{}, u.Gaps...), Blockers: []graph.Gap{}}
		for _, dep := range u.Dependencies {
			c.Evidence["render_source"] = append(c.Evidence["render_source"], dep.Source)
		}
		for _, e := range u.Context {
			if e.Kind == "render_binding" {
				c.Evidence["render_binding"] = append(c.Evidence["render_binding"], e)
			}
		}
		c.Evidence["render_context"] = []graph.Evidence{u.Sink}
		for _, e := range u.Context {
			if e.Kind == "render_markup_context" {
				c.Evidence["render_context"] = append(c.Evidence["render_context"], e)
			}
		}
		if u.XSS.Protection == "unknown" || u.XSS.OutputContext == "unknown" || u.XSS.OutputContext == "mixed" {
			c.Blockers = append(c.Blockers, graph.Gap{Kind: "render_protection_unresolved", Detail: u.XSS.ProtectionReason, At: u.Sink.Span})
		}
		if (u.XSS.Type == "stored" || u.XSS.Type == "cross_tier") && !u.XSS.OriginResolved {
			c.Blockers = append(c.Blockers, graph.Gap{Kind: "storage_origin_unresolved", Detail: "No complete source-to-storage write link is established.", At: u.Sink.Span})
		}
		for _, gap := range u.Gaps {
			switch gap.Kind {
			case "work_limit", "dependency_limit", "transfer_limit", "gap_limit", "template_parse_unresolved":
				c.Blockers = append(c.Blockers, gap)
			}
		}
		c.ModelFacts = []string{"output_context=" + u.XSS.OutputContext, "protection=" + u.XSS.Protection, fmt.Sprintf("origin_resolved=%t", u.XSS.OriginResolved)}
		u.Completeness = c
	}
}
