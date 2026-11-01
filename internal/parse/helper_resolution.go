package parse

import (
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/zazzcode/scan/internal/graph"
)

func (s *moduleSet) resolveCall(m *sourceModule, call *sitter.Node) (*sourceModule, *sitter.Node) {
	fn := call.ChildByFieldName("function")
	if fn == nil {
		return nil, nil
	}
	name := fn.Utf8Text(m.src)
	if shadowed(call, strings.Split(name, ".")[0], m.src) {
		return nil, nil
	}
	if fn.Kind() != "identifier" {
		return s.resolve(m, name)
	}
	for p := call.Parent(); p != nil; p = p.Parent() {
		if p.Kind() != "function_declaration" {
			continue
		}
		body := p.ChildByFieldName("body")
		if body == nil {
			continue
		}
		for i := uint(0); i < body.NamedChildCount(); i++ {
			d := body.NamedChild(i)
			id := d.ChildByFieldName("name")
			if d.Kind() == "function_declaration" && id != nil && id.Utf8Text(m.src) == name {
				return m, d
			}
		}
	}
	return s.resolve(m, name)
}

func (s *moduleSet) inputHelperGaps(m *sourceModule, body *sitter.Node, req string, fields []string) []graph.Gap {
	gaps := []graph.Gap{}
	walkCalls(body, func(call *sitter.Node) {
		fn := call.ChildByFieldName("function")
		args := call.ChildByFieldName("arguments")
		if fn == nil || args == nil {
			return
		}
		name := fn.Utf8Text(m.src)
		if strings.HasSuffix(name, ".query") {
			return
		}
		relevant := false
		var inspect func(*sitter.Node)
		inspect = func(n *sitter.Node) {
			if n.Kind() == "identifier" && n.Utf8Text(m.src) == req {
				relevant = true
			}
			if n.Kind() == "member_expression" {
				text := n.Utf8Text(m.src)
				for _, f := range fields {
					if text == req+".body."+f || text == req+".query."+f {
						relevant = true
					}
				}
				if text == req+".body" {
					relevant = true
				}
				return
			}
			for i := uint(0); i < n.NamedChildCount(); i++ {
				inspect(n.NamedChild(i))
			}
		}
		inspect(args)
		if !relevant {
			return
		}
		if _, helper := s.resolveCall(m, call); helper != nil {
			return
		}
		gaps = append(gaps, graph.Gap{Kind: "input_helper_unresolved", Detail: "call receiving checked input lacks a definition: " + name, At: spanOf(m.file, call)})
	})
	return gaps
}
