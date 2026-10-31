package parse

import (
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/zazzcode/scan/internal/graph"
)

// middlewareContext preserves applicable earlier USE registrations. Built-in
// transport/parser adapters are declared assumptions, not custom validator rules.
func (s *moduleSet) middlewareContext(m *sourceModule, route graph.Route) ([]graph.Evidence, []graph.Gap) {
	out := []graph.Evidence{}
	gaps := []graph.Gap{}
	seen := map[string]bool{}
	walkCalls(m.tree.RootNode(), func(call *sitter.Node) {
		if int(call.StartByte()) >= route.Registration.StartByte {
			return
		}
		fn := call.ChildByFieldName("function")
		if fn == nil || fn.Kind() != "member_expression" {
			return
		}
		obj := fn.ChildByFieldName("object")
		prop := fn.ChildByFieldName("property")
		if obj == nil || prop == nil || prop.Utf8Text(m.src) != "use" || !m.receivers[obj.Utf8Text(m.src)] {
			return
		}
		args := call.ChildByFieldName("arguments")
		if args == nil || args.NamedChildCount() == 0 {
			return
		}
		first := uint(0)
		if args.NamedChild(0).Kind() == "array" {
			paths, dynamic, _ := pathArgs(args.NamedChild(0), m.src)
			if dynamic {
				gaps = append(gaps, graph.Gap{Kind: "middleware_mount_unresolved", Detail: "dynamic middleware path array", At: spanOf(m.file, call)})
				return
			}
			applies := false
			for _, path := range paths {
				if path == "/" || route.PathPattern == path || strings.HasPrefix(route.PathPattern, path+"/") {
					applies = true
				}
			}
			if !applies {
				return
			}
			first = 1
		}
		if path, ok := stringLiteral(args.NamedChild(0), m.src); ok {
			if strings.ContainsAny(path, "*:?()+[]") {
				prefix := path
				for i, c := range path {
					if strings.ContainsRune("*:?()+[]", c) {
						prefix = path[:i]
						break
					}
				}
				if prefix != "" && !strings.HasPrefix(route.PathPattern, prefix) && route.PathPattern != strings.TrimSuffix(prefix, "/") {
					return
				}
				gaps = append(gaps, graph.Gap{Kind: "middleware_mount_unresolved", Detail: "pattern middleware applicability requires a framework path model", At: spanOf(m.file, call)})
			} else if path != "/" && route.PathPattern != path && !strings.HasPrefix(route.PathPattern, path+"/") {
				return
			}
			first = 1
		}
		out = append(out, m.evidence("middleware_registration", call))
		for i := first; i < args.NamedChildCount(); i++ {
			arg := args.NamedChild(i)
			if transportMiddleware(m, arg) {
				continue
			}
			hm, body, _, more := s.handler(m, arg)
			if body == nil {
				for _, g := range more {
					g.Kind = "applicable_middleware_unresolved"
					gaps = append(gaps, g)
				}
				continue
			}
			// The full source is usable for judging even when the dependency engine
			// cannot lower a destructured URL-only request parameter.
			for _, g := range more {
				if g.Kind != "request_projection_unresolved" {
					gaps = append(gaps, g)
				}
			}
			out = append(out, hm.evidence("middleware_handler", body))
			s.includeHelpers(hm, body, &out, seen, 0)
		}
	})
	return out, gaps
}

func transportMiddleware(m *sourceModule, n *sitter.Node) bool {
	if n.Kind() == "call_expression" {
		n = n.ChildByFieldName("function")
	}
	if n == nil {
		return false
	}
	name := n.Utf8Text(m.src)
	if parts := strings.Split(name, "."); len(parts) > 1 {
		name = parts[0]
	}
	imp, ok := m.imports[name]
	if !ok {
		return false
	}
	switch imp.module {
	case "express", "body-parser", "cookie-parser", "compression", "morgan", "cors", "helmet", "i18n", "feature-policy", "express-robots-txt":
		return true
	}
	return false
}
