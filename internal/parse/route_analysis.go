package parse

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/zazzcode/zazzles/internal/analysis"
)

func (s *moduleSet) analyzeRoute(out *analysis.AnalysisReport, m, owner *sourceModule, r analysis.Route, call *sitter.Node, receiver string, context []analysis.Evidence) {
	e := analysis.Entry{Route: r, Gaps: []analysis.Gap{}}
	if shadowed(call, receiver, m.src) && !parameterReceiver(m, call, receiver) {
		e.Gaps = append(e.Gaps, analysis.Gap{Kind: "ambiguous_receiver", Detail: "receiver is shadowed in enclosing scope", At: r.Registration})
		out.Entries = append(out.Entries, e)
		return
	}
	args := call.ChildByFieldName("arguments")
	var h *sitter.Node
	if args != nil && args.NamedChildCount() > 1 {
		h = args.NamedChild(args.NamedChildCount() - 1)
	}
	hm, body, req, gaps := s.handler(m, h)
	var parts map[string]string
	if body != nil && req == "" {
		parts, gaps = requestParts(hm, body, gaps)
	}
	e.Gaps = append(e.Gaps, gaps...)
	if body != nil {
		evidence := hm.evidence("handler", body)
		e.Handler = &evidence
		f := flow{module: hm, set: s, route: r, request: req, response: handlerParameter(body, 1, hm.src), registration: owner,
			env: map[string]value{}, context: append([]analysis.Evidence{m.evidence("registration", call), evidence}, context...),
			gaps: append([]analysis.Gap{}, e.Gaps...), helpers: &helperExpansion{budget: 24, stack: map[string]bool{}}, facts: newInputFacts(), parts: parts}
		f.gaps = append(f.gaps, analysis.Gap{Kind: "middleware_unresolved", Detail: "effective middleware/guard chain is not modeled", At: r.Registration})
		s.noteMiddlewareFields(&f, m, call)
		f.statements(body)
		f.finishValidation()
		s.requestLineLogs(&f, m, call, receiver)
		out.Uses = append(out.Uses, f.uses...)
		e.Gaps = append(e.Gaps, f.gaps[len(e.Gaps):]...)
	}
	out.Entries = append(out.Entries, e)
}
