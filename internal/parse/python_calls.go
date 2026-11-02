package parse

import (
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/zazzcode/zazzles/internal/graph"
)

func (f *pythonFlow) call(n *sitter.Node) pythonValue {
	m := f.module
	fn := n.ChildByFieldName("function")
	name := pyText(m, fn)
	q := f.qualified(name)
	if q == "flask.render_template" {
		f.renderTemplate(n)
		return pythonValue{rendered: true}
	}
	nodes, values, all := f.callArguments(n)

	object := pyField(fn, "object")
	method := pyText(m, pyField(fn, "attribute"))
	var receiver pythonValue
	if object != nil {
		receiver = f.eval(object)
	}
	if v, ok := f.frameworkCall(n, q, nodes, values, all); ok {
		return v
	}

	if receiver.binding == "db_connection" && method == "cursor" {
		return pythonValue{binding: "db_cursor", evidence: append(receiver.evidence, m.evidence("sink_binding", n))}
	}
	if (receiver.binding == "db_connection" || receiver.binding == "db_cursor") && (method == "execute" || method == "executemany" || method == "executescript") {
		result := f.sql(n, object, values, receiver)
		result.binding = "db_cursor"
		return result
	}
	if receiver.binding == "db_cursor" && (method == "fetchone" || method == "fetchall") {
		v := receiver
		v.binding = "db_row"
		v.evidence = append(v.evidence, m.evidence("storage_read", n))
		return v
	}
	if receiver.binding == "jinja_template" && method == "render" {
		v := pyCombine(receiver, all)
		v = pyTransfer(v, m.evidence("template_render", n))
		v.binding = ""
		return v
	}
	if method == "get" && len(values) > 0 {
		key, _ := pyLiteral(m, nodes[0])
		if len(receiver.deps) > 0 || receiver.fields != nil {
			v := f.projectValue(receiver, key, n)
			if len(values) > 1 {
				v = pyCombine(v, values[1])
			}
			return v
		}
	}
	if f.outbound(q) {
		return f.outboundCall(n, q, fn, values)
	}

	if f.logger(q, name, receiver, method, n) {
		for _, v := range values {
			all = pyCombine(all, f.logValue(n, v))
		}
		f.log(n, all)
		return pythonValue{}
	}
	if (name == "len" || name == "bool") && f.builtin(name) {
		return pyTransfer(pyProtect(all, "non_content"), m.evidence("non_content", n))
	}
	if (name == "int" || name == "float") && f.builtin(name) {
		return pyTransfer(pyProtect(all, "numeric"), m.evidence("numeric_conversion", n))
	}
	if name == "str" && f.builtin(name) {
		return pyTransfer(all, m.evidence("string_conversion", n))
	}
	if method == "format" || method == "join" {
		v := pyCombine(receiver, all)
		v.built = true
		return pyTransfer(v, m.evidence("string_construction", n))
	}
	if method == "strip" || method == "lower" || method == "upper" || method == "split" {
		return pyTransfer(pyCombine(receiver, all), m.evidence("string_transform", n))
	}
	if method == "getLogger" && strings.HasPrefix(q, "logging.") {
		return pythonValue{binding: "logger", evidence: []graph.Evidence{m.evidence("logger_binding", n)}}
	}
	if target, key := f.set.resolve(m, name, 0); target != nil && target.funcs[key] != nil {
		// Local assignments/parameters override helper declarations and imported helpers.
		if _, shadow := f.env[strings.Split(name, ".")[0]]; !shadow {
			return f.helper(n, target, target.funcs[key], values)
		}
	}
	v := pyCombine(receiver, all)
	if len(v.deps) > 0 {
		v = pyTransfer(pyProtect(v, "unknown"), m.evidence("call_argument_dependency", n))
		v.gaps = append(v.gaps, f.gap("python_call_unresolved", "unmodeled return/content effect: "+name, n))
	}
	return v
}
func (f *pythonFlow) builtin(name string) bool {
	_, local := f.env[name]
	return !local && f.module.vars[name] == nil && f.module.funcs[name] == nil && f.module.imports[name].node == nil
}
func (f *pythonFlow) helper(call *sitter.Node, m *pythonModule, fn *sitter.Node, values []pythonValue) pythonValue {
	key := m.file + ":" + pyText(m, fn.ChildByFieldName("name"))
	if f.depth >= 3 || *f.work >= 24 || f.stack[key] {
		v := pythonValue{}
		for _, a := range values {
			v = pyCombine(v, a)
		}
		v = pyProtect(v, "unknown")
		v.gaps = append(v.gaps, f.gap("work_limit", "Python helper expansion depth/cycle/work limit", call))
		return v
	}
	*f.work++
	f.stack[key] = true
	defer delete(f.stack, key)
	f.context = append(f.context, m.evidence("helper_definition", fn), f.module.evidence("call_binding", call))
	h := f.clone()
	h.module = m
	h.depth++
	h.env = map[string]pythonValue{}
	h.imports = map[string]pythonImport{}
	h.context = append(h.context, m.evidence("helper_definition", fn), f.module.evidence("call_binding", call))
	params := fn.ChildByFieldName("parameters")
	pos := 0
	for i := uint(0); params != nil && i < params.NamedChildCount(); i++ {
		p := params.NamedChild(i)
		name := p
		if p.Kind() != "identifier" {
			name = p.ChildByFieldName("name")
			if name == nil && p.Kind() == "typed_parameter" && p.NamedChildCount() > 0 {
				name = p.NamedChild(0)
			}
		}
		if name == nil {
			h.gaps = append(h.gaps, graph.Gap{Kind: "python_parameter_unresolved", Detail: "unsupported helper parameter", At: spanOf(m.file, p)})
			continue
		}
		v := pythonValue{}
		if pos < len(values) {
			v = pyTransfer(values[pos], m.evidence("argument_to_parameter", p))
		} else if keyword := pyKeywordIn(f.module, call, pyText(m, name)); keyword != nil {
			v = pyTransfer(f.eval(keyword), m.evidence("argument_to_parameter", p))
		} else {
			v = h.eval(p.ChildByFieldName("value"))
		}
		h.env[pyText(m, name)] = v
		pos++
	}
	h.block(fn.ChildByFieldName("body"))
	f.mergeGuards(h)
	return pyTransfer(h.returns, f.module.evidence("helper_return", call))
}
func (f *pythonFlow) outbound(q string) bool {
	for _, p := range []string{"requests.", "httpx."} {
		if strings.HasPrefix(q, p) {
			switch strings.TrimPrefix(q, p) {
			case "get", "post", "put", "patch", "delete", "head", "request":
				return true
			}
		}
	}
	return q == "urllib.request.urlopen"
}
func (f *pythonFlow) logger(q, name string, v pythonValue, method string, n *sitter.Node) bool {
	if name == "print" && f.builtin("print") {
		return true
	}
	if strings.HasPrefix(q, "logging.") || v.binding == "logger" {
		switch method {
		case "debug", "info", "warning", "warn", "error", "exception", "critical", "log":
			return true
		}
		switch strings.TrimPrefix(q, "logging.") {
		case "debug", "info", "warning", "error", "exception", "critical", "log":
			return true
		}
	}
	return false
}
func (f *pythonFlow) log(n *sitter.Node, v pythonValue) {
	var claims []graph.Gap
	for _, d := range v.deps {
		class, hint := logDataClass(d.container, d.key)
		if d.container == "session" {
			class = "container"
		}
		if class == "" {
			continue
		}
		kind, detail := "log_sensitive_field_hint", hint+" field hint: "+d.key
		switch {
		case d.protection == "non_content":
			kind, detail = "log_non_content_only", "only a derived non-content fact is emitted"
		case d.protection == "unknown":
			kind, detail = "log_transform_unresolved", "unmodeled content transformation"
		case class == "protocol_credential":
			kind, detail = "log_protocol_credential", "HTTP credential content is emitted"
		case class == "container":
			kind, detail = "log_request_container", "whole request/session container may include sensitive data"
		}
		claims = append(claims, f.gap(kind, detail, n))
		if d.dep.Source.Kind == "literal_sensitive_source" {
			claims = append(claims, f.gap("log_literal_sensitivity_unresolved", "Email syntax alone does not establish that a constant address is private or identifies a person", n))
		}
	}
	if len(claims) == 0 {
		return
	}
	u := f.addUse("log_sink_candidate", n, v)
	u.Context = append(u.Context, f.module.evidence("logger_binding", n.ChildByFieldName("function")))
	u.Gaps = append(u.Gaps, claims...)
	u.Gaps = append(u.Gaps, f.gap("log_destination_unresolved", "logger handlers/formatters/retention not modeled", n))
}

func (f *pythonFlow) callArguments(n *sitter.Node) ([]*sitter.Node, []pythonValue, pythonValue) {
	nodes := pyArgs(n)
	values := make([]pythonValue, len(nodes))
	var all pythonValue
	for i, a := range nodes {
		values[i] = f.eval(a)
		all = pyCombine(all, values[i])
	}
	args := n.ChildByFieldName("arguments")
	if args != nil {
		for i := uint(0); i < args.NamedChildCount(); i++ {
			a := args.NamedChild(i)
			if a.Kind() == "keyword_argument" {
				all = pyCombine(all, f.eval(a.ChildByFieldName("value")))
			}
		}
	}
	return nodes, values, all
}
