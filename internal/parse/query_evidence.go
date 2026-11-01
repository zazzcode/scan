package parse

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/zazzcode/scan/internal/graph"
	"strings"
)

// refineQueryEvidence separates value/guard evidence assembled for this claim from
// unrelated limits inherited from discovery. Original gaps remain in the report.
func (s *moduleSet) refineQueryEvidence(use *graph.StringUse) {
	c := use.Completeness
	if c == nil || c.ClaimKind != "injection" {
		return
	}
	c.Required = append(c.Required, "control_flow")
	// Missing middleware/helper bodies remain explicit blockers from chain assembly.
	chainSupplied := true
	for _, g := range c.Blockers {
		switch g.Kind {
		case "applicable_middleware_unresolved", "middleware_mount_unresolved", "input_helper_unresolved", "helper_context_limit", "unresolved_factory", "unresolved_handler":
			chainSupplied = false
		}
	}
	if chainSupplied && len(s.contextLimits) == 0 {
		c.ModelFacts = []string{
			"Scope is direct Express registrations with syntactically matched earlier middleware and same-route handlers. Matching source bodies or explicit transport/authentication adapter contracts are supplied; unsupported/dynamic chains remain blockers.",
			"Known express.json/body-parser parsing adapters parse request values and may admit nested JSON objects. They do not enforce an id/email/name field schema or remove query operators from nested fields. A parser call alone is not a validator.",
			"For this supported registration scope, the discovery-stage generic middleware_unresolved warning is superseded by the supplied route/middleware context. Judge guard effects from those bodies and actual input transformations; do not invent additional undisclosed middleware.",
		}
		coverage := []graph.Gap{}
		for _, g := range c.Coverage {
			if g.Kind == "middleware_unresolved" {
				c.Resolved = append(c.Resolved, g)
			} else {
				coverage = append(coverage, g)
			}
		}
		c.Coverage = coverage
	}
	blockers := []graph.Gap{}
	for _, g := range c.Blockers {
		covered := false
		m := s.modules[g.At.File]
		var n *sitter.Node
		if m != nil {
			n = nodeAt(m.tree.RootNode(), g.At.StartByte, g.At.EndByte)
		}
		switch g.Kind {
		case "helper_call_limit":
			// Discovery's traversal differs from independently assembled input-helper context.
			// Missing relevant helper bodies still block via input_helper/context-limit gaps.
			covered = len(s.contextLimits) == 0 && len(c.Evidence["control_flow"]) > 0
		case "unknown_call_result", "unknown_transfer":
			covered = !gapTouchesSQLSource(g, *use, use.Check.Fields)
			if !covered && nativeStringOperation(m, n) {
				covered = true
				c.Evidence["transformation"] = append(c.Evidence["transformation"], m.evidence("transformation", n))
			}
		case "conditional_merge":
			if n != nil && n.Kind() == "ternary_expression" {
				covered = true
				c.Evidence["conditional_value"] = append(c.Evidence["conditional_value"], m.evidence("conditional_value", n))
			}
		case "unresolved_wrapper":
			covered = len(c.Evidence["wrapper_definition"]) > 0
		}
		if covered {
			c.Coverage = append(c.Coverage, g)
		} else {
			blockers = append(blockers, g)
		}
	}
	for _, e := range c.Evidence["transformation"] {
		module := s.modules[e.Span.File]
		if module == nil {
			continue
		}
		node := nodeAt(module.tree.RootNode(), e.Span.StartByte, e.Span.EndByte)
		if node == nil {
			continue
		}
		fn := node.ChildByFieldName("function")
		if fn != nil && fn.Utf8Text(module.src) == "Number" && nativeStringOperation(module, node) {
			c.ModelFacts = append(c.ModelFacts, "The exact transformation "+e.Text+" at "+e.Span.File+" is a verified unshadowed ECMAScript Number call. ONLY that call result is numeric (including NaN/Infinity) or throws; this guarantee does NOT apply to the original input or other conditional branches. Evaluate each alternative separately. String truncation and newline removal do not make other alternatives numeric and do not remove query-expression operators. Numeric business requirements remain separate.")
		}
	}
	c.Blockers = blockers
	s.conditionalHelpers(use)
}

func nativeStringOperation(m *sourceModule, n *sitter.Node) bool {
	if m == nil || n == nil || n.Kind() != "call_expression" {
		return false
	}
	fn := n.ChildByFieldName("function")
	if fn == nil {
		return false
	}
	text := fn.Utf8Text(m.src)
	if (text == "Number" || text == "String") && !shadowed(n, text, m.src) && m.funcs[text] == nil && topVariable(m, text) == nil {
		if _, ok := m.imports[text]; !ok {
			return true
		}
	}
	if fn.Kind() != "member_expression" {
		return false
	}
	prop := fn.ChildByFieldName("property")
	if prop == nil {
		return false
	}
	switch prop.Utf8Text(m.src) {
	case "substring", "substr", "slice", "replace", "trim":
		return true
	}
	return false
}

// wrapperContext includes verified local wrapper bodies for semantic assessment;
// names such as asyncHandler are never treated as protection by themselves.
func (s *moduleSet) wrapperContext(m *sourceModule, arg *sitter.Node) []graph.Evidence {
	if arg == nil || arg.Kind() != "call_expression" {
		return nil
	}
	args := arg.ChildByFieldName("arguments")
	if args == nil || args.NamedChildCount() != 1 {
		return nil
	}
	target, fn := s.resolveCall(m, arg)
	if target == nil || fn == nil {
		return nil
	}
	return []graph.Evidence{target.evidence("wrapper_definition", fn)}
}

func (s *moduleSet) authenticationAdapter(m *sourceModule, arg *sitter.Node) []graph.Evidence {
	if arg == nil || arg.Kind() != "call_expression" {
		return nil
	}
	target, fn := s.resolveCall(m, arg)
	if target == nil || fn == nil {
		return nil
	}
	body := fn.ChildByFieldName("body")
	if body == nil || body.Kind() != "call_expression" {
		return nil
	}
	callee := body.ChildByFieldName("function")
	args := body.ChildByFieldName("arguments")
	if callee == nil || args == nil || args.NamedChildCount() != 1 {
		return nil
	}
	imp, ok := target.imports[callee.Utf8Text(target.src)]
	if !ok || imp.module != "express-jwt" {
		return nil
	}
	config := args.NamedChild(0)
	for config != nil && (config.Kind() == "as_expression" || config.Kind() == "parenthesized_expression") {
		config = config.NamedChild(0)
	}
	if config == nil || config.Kind() != "object" {
		return nil
	}
	for i := uint(0); i < config.NamedChildCount(); i++ {
		pair := config.NamedChild(i)
		key := pair.ChildByFieldName("key")
		if key == nil {
			return nil
		}
		switch strings.Trim(key.Utf8Text(target.src), "'\"") {
		case "secret", "algorithms", "credentialsRequired":
		default:
			return nil
		}
	}
	return []graph.Evidence{target.evidence("authentication_adapter", fn), target.evidence("authentication_adapter", imp.declaration)}
}

func (s *moduleSet) conditionalHelpers(use *graph.StringUse) {
	c := use.Completeness
	seen := map[string]bool{}
	var include func(*sourceModule, *sitter.Node, int)
	include = func(m *sourceModule, n *sitter.Node, depth int) {
		if depth > 2 {
			return
		}
		walkCalls(n, func(call *sitter.Node) {
			fn := call.ChildByFieldName("function")
			if fn == nil {
				return
			}
			name := fn.Utf8Text(m.src)
			if nativeStringOperation(m, call) {
				return
			}
			target, helper := s.resolveCall(m, call)
			if target != nil && helper != nil {
				key := target.file + ":" + helper.Utf8Text(target.src)
				if seen[key] {
					return
				}
				seen[key] = true
				evidence := target.evidence("condition_helper", helper)
				use.Context = append(use.Context, evidence)
				c.Evidence["condition_helper"] = append(c.Evidence["condition_helper"], evidence)
				if body := helper.ChildByFieldName("body"); body != nil {
					include(target, body, depth+1)
				}
			} else if depth == 0 {
				c.Blockers = append(c.Blockers, graph.Gap{Kind: "condition_helper_unresolved", Detail: "condition callee lacks a supported definition: " + name, At: spanOf(m.file, call)})
			}
		})
	}
	for _, e := range c.Evidence["conditional_value"] {
		m := s.modules[e.Span.File]
		if m == nil {
			continue
		}
		n := nodeAt(m.tree.RootNode(), e.Span.StartByte, e.Span.EndByte)
		if n == nil {
			continue
		}
		if condition := n.ChildByFieldName("condition"); condition != nil {
			include(m, condition, 0)
		}
		conditionalConstraints(m, n, c)
	}
	if len(c.Evidence["conditional_value"]) > 0 {
		c.ModelFacts = append(c.ModelFacts, "Conditional value expressions and any directly resolved condition helper bodies are supplied. A static path may be judged under an explicit possible configuration precondition; do not require an observed runtime request/configuration. Reject a path if supplied source proves its precondition impossible. Do not report all runtime configurations as vulnerable merely because one supported branch is unsafe.")
	}
}
