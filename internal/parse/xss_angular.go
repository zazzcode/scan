package parse

import (
	"path/filepath"
	"regexp"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/zazzcode/zazzles/internal/graph"
)

var angularBinding = regexp.MustCompile(`(?i)\[(innerHTML|href|src)\]\s*=\s*(["'])([^"']+)["']`)

func (s *moduleSet) clientClass(m *sourceModule, class *sitter.Node) *clientContext {
	name := class.ChildByFieldName("name")
	context := &clientContext{entry: graph.ClientEntry{ID: clientID(m.file, name.Utf8Text(m.src), class.StartByte()), Kind: "component", Name: name.Utf8Text(m.src), Location: spanOf(m.file, class)},
		injected: map[string]clientInjection{}, methods: map[string]*sitter.Node{}, fields: map[string]*sitter.Node{}}
	body := class.ChildByFieldName("body")
	if body == nil {
		return context
	}
	for i := uint(0); i < body.NamedChildCount(); i++ {
		n := body.NamedChild(i)
		key := n.ChildByFieldName("name")
		if key == nil {
			continue
		}
		field := key.Utf8Text(m.src)
		if n.Kind() == "method_definition" {
			context.methods[field] = n
			if field == "constructor" {
				params := n.ChildByFieldName("parameters")
				if params != nil {
					for j := uint(0); j < params.NamedChildCount(); j++ {
						p := params.NamedChild(j)
						id, typ := p.ChildByFieldName("pattern"), p.ChildByFieldName("type")
						if id != nil && typ != nil {
							label := strings.TrimSpace(strings.TrimPrefix(typ.Utf8Text(m.src), ":"))
							if imp, ok := m.imports[label]; ok {
								context.injected[id.Utf8Text(m.src)] = clientInjection{imp.module, imp.name, p}
							}
						}
					}
				}
			}
			continue
		}
		v := n.ChildByFieldName("value")
		if v == nil {
			continue
		}
		context.fields[field] = v
		if v.Kind() != "call_expression" {
			continue
		}
		fn, args := v.ChildByFieldName("function"), v.ChildByFieldName("arguments")
		if fn == nil || args == nil || args.NamedChildCount() != 1 {
			continue
		}
		factory := m.imports[fn.Utf8Text(m.src)]
		if factory.module != "@angular/core" || factory.name != "inject" {
			continue
		}
		if imp, ok := m.imports[args.NamedChild(0).Utf8Text(m.src)]; ok {
			context.injected[field] = clientInjection{imp.module, imp.name, n}
			context.bindings = append(context.bindings, m.evidence("render_binding", n), m.evidence("render_binding", imp.declaration))
		}
	}
	walkCalls(m.tree.RootNode(), func(call *sitter.Node) {
		fn, args := call.ChildByFieldName("function"), call.ChildByFieldName("arguments")
		if fn == nil || args == nil || args.NamedChildCount() != 1 || call.Parent() == nil || call.Parent().Kind() != "decorator" || call.EndByte() > body.StartByte() {
			return
		}
		imp := m.imports[fn.Utf8Text(m.src)]
		if imp.module != "@angular/core" || imp.name != "Component" {
			return
		}
		config := args.NamedChild(0)
		path, ok := stringLiteral(property(config, "templateUrl", m.src), m.src)
		if !ok {
			return
		}
		if template, ok := s.readText(filepath.ToSlash(filepath.Join(filepath.Dir(m.file), path))); ok {
			context.template = &template
			context.angular, context.entry.Framework = true, "angular"
			context.bindings = append(context.bindings, m.evidence("render_binding", call), m.evidence("render_binding", imp.declaration))
		}
	})
	return context
}

func (s *moduleSet) angularOutputs(f *flow) {
	if f.client == nil || f.client.template == nil || !f.client.angular {
		return
	}
	file := *f.client.template
	matchedBypass := map[string]bool{}
	for _, match := range angularBinding.FindAllStringSubmatchIndex(file.Text, -1) {
		if htmlCommentAt(file.Text, match[0]) {
			continue
		}
		kind, expression := file.Text[match[2]:match[3]], strings.TrimSpace(file.Text[match[6]:match[7]])
		if !simpleTemplateField.MatchString(expression) {
			continue
		}
		v, ok := f.env[clientKey("this."+expression)]
		if !ok || len(v.deps) == 0 {
			continue
		}
		context := "html_text"
		if !strings.EqualFold(kind, "innerHTML") {
			context = "url"
		}
		before := len(f.uses)
		sink := textRange(file, match[0], match[1], "angular_output")
		f.recordXSS(sink, v, context, "text/html", v.trustedContext == "", f.client.bindings)
		if len(f.uses) == before {
			continue
		}
		u := &f.uses[len(f.uses)-1]
		if v.trustedContext == "" {
			if plainFrameworkInput(v) {
				u.XSS.Protection, u.XSS.ProtectionReason = "sanitized", "A plain browser/request value reaches Angular's one-way binding with no unknown transformation or trust bypass; its default sanitization applies."
			} else {
				u.XSS.Protection, u.XSS.ProtectionReason = "unknown", "Default Angular sanitization is expected for plain values, but the return type/transformations of this value are not established; a trusted wrapper cannot be excluded."
			}
		} else if v.trustedContext != context {
			u.XSS.Protection, u.XSS.ProtectionReason = "unknown", "The trust bypass and rendering context do not match."
		}
		if v.apiPath != "" {
			u.XSS.Type, u.XSS.OriginResolved = "cross_tier", false
			u.XSS.ResponsePath, u.XSS.ResponseMethod = v.apiPath, v.apiMethod
		}
		for _, d := range v.deps {
			for _, e := range d.Transfers {
				if e.Kind == "trust_bypass" {
					matchedBypass[e.Span.File+":"+e.Text] = true
				}
			}
		}
	}
	var uses []graph.StringUse
	for _, u := range f.uses {
		if u.Sink.Kind != "trust_bypass" || !matchedBypass[u.Sink.Span.File+":"+u.Sink.Text] {
			uses = append(uses, u)
		}
	}
	f.uses = uses
}

func plainFrameworkInput(v value) bool {
	if len(v.deps) == 0 {
		return false
	}
	for _, dep := range v.deps {
		if dep.Source.Kind != "browser_source" && dep.Source.Kind != "request_source" {
			return false
		}
		for _, transfer := range dep.Transfers {
			if transfer.Kind == "call_argument_dependency" || transfer.Kind == "trust_bypass" {
				return false
			}
		}
	}
	return true
}
