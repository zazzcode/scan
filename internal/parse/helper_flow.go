package parse

import (
	"fmt"

	sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/zazzcode/scan/internal/graph"
)

// helperExpansion bounds one handler's helper-call traversal. The stack stops
// recursion along the current path; budget counts call expansions, not unique names.
type helperExpansion struct {
	budget int
	stack  map[string]bool
}

type boundParam struct {
	name string
	node *sitter.Node
}

// enterHelper analyzes a same-file or relative-imported function body.
// entered means sinks inside the body were sought. precise means the call result
// is the modeled return value rather than a may-dependency on every argument.
func (f *flow) enterHelper(call *sitter.Node, name string, args []value) (ret value, entered bool, precise bool) {
	if f.helpers == nil || callHasSpread(call) {
		return value{}, false, false
	}
	hm, helper := f.set.resolveCall(f.module, call)
	if hm == nil || helper == nil || !supportedCallee(helper) {
		return value{}, false, false
	}
	params, ok := functionParameters(helper, hm.src)
	if !ok || len(args) > len(params) {
		return value{}, false, false
	}
	key := fmt.Sprintf("%s:%d", hm.file, helper.StartByte())
	if f.helpers.stack[key] {
		g := f.gap("recursive_helper", "recursive helper call is not re-entered: "+name, call)
		f.gaps = append(f.gaps, g)
		return value{}, false, false
	}
	if f.helpers.budget <= 0 || f.depth >= 3 {
		f.gaps = append(f.gaps, f.gap("helper_call_limit", "helper call expansion reached its depth or count limit: "+name, call))
		return value{}, false, false
	}
	f.helpers.budget--
	f.helpers.stack[key] = true
	defer delete(f.helpers.stack, key)

	f.context = append(f.context, hm.evidence("helper_definition", helper))
	before := len(f.gaps)
	child := flow{
		module:  hm,
		set:     f.set,
		route:   f.route,
		env:     map[string]value{},
		context: append(append([]graph.Evidence{}, f.context...), f.module.evidence("helper_call", call)),
		gaps:    append([]graph.Gap{}, f.gaps...),
		depth:   f.depth + 1,
		helpers: f.helpers,
		facts:   f.facts,
	}
	unbound := false
	for i, p := range params {
		if i < len(args) {
			// Bind on the callee flow so the parameter span stays in the callee file.
			child.env[p.name] = child.transfer(args[i], "parameter_binding", p.node)
			continue
		}
		unbound = true
		child.env[p.name] = value{gaps: []graph.Gap{{Kind: "unknown_value", Detail: "helper parameter has no argument: " + p.name, At: spanOf(hm.file, p.node)}}}
	}
	child.statements(helper.ChildByFieldName("body"))
	f.uses = append(f.uses, child.uses...)
	for _, g := range child.gaps[before:] {
		switch g.Kind {
		case "work_limit", "helper_call_limit", "recursive_helper":
			f.gaps = append(f.gaps, g)
		}
	}
	precise = len(child.returns) > 0 && !unbound
	for _, g := range child.gaps[before:] {
		switch g.Kind {
		case "unsupported_control_flow", "work_limit", "unsupported_binding", "helper_call_limit", "recursive_helper":
			precise = false
		}
	}
	if !precise {
		return value{}, true, false
	}
	for _, r := range child.returns {
		ret = combine(ret, r)
	}
	if len(child.returns) > 1 {
		ret.gaps = append(ret.gaps, f.gap("conditional_merge", "multiple helper returns merged without path feasibility: "+name, call))
	}
	return ret, true, true
}

func supportedCallee(n *sitter.Node) bool {
	return n.Kind() == "function_declaration" || isFunction(n)
}

func functionParameters(fn *sitter.Node, src []byte) ([]boundParam, bool) {
	params := fn.ChildByFieldName("parameters")
	if params == nil {
		return nil, true
	}
	// `value => ...` stores the identifier directly. `(value) => ...` wraps it.
	if params.Kind() == "identifier" {
		return []boundParam{{name: params.Utf8Text(src), node: params}}, true
	}
	out := []boundParam{}
	for i := uint(0); i < params.NamedChildCount(); i++ {
		node := params.NamedChild(i)
		if node.Kind() == "required_parameter" || node.Kind() == "optional_parameter" {
			node = node.ChildByFieldName("pattern")
		}
		if node != nil && node.Kind() == "assignment_pattern" {
			node = node.ChildByFieldName("left")
		}
		if node == nil || node.Kind() != "identifier" {
			return nil, false
		}
		out = append(out, boundParam{name: node.Utf8Text(src), node: node})
	}
	return out, true
}

func callHasSpread(call *sitter.Node) bool {
	args := call.ChildByFieldName("arguments")
	if args == nil {
		return false
	}
	for i := uint(0); i < args.NamedChildCount(); i++ {
		if args.NamedChild(i).Kind() == "spread_element" {
			return true
		}
	}
	return false
}
