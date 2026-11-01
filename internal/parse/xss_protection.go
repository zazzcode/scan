package parse

import (
	"regexp"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

var htmlAttribute = regexp.MustCompile(`(?i)([a-z_:][a-z0-9_:.-]*)\s*=\s*(["']?)[^"']*$`)

func htmlOutputContext(prefix string) string {
	lower := strings.ToLower(prefix)
	if strings.LastIndex(lower, "<script") > strings.LastIndex(lower, "</script") && strings.LastIndex(lower, ">") > strings.LastIndex(lower, "<script") {
		return "script"
	}
	if strings.LastIndex(lower, "<style") > strings.LastIndex(lower, "</style") && strings.LastIndex(lower, ">") > strings.LastIndex(lower, "<style") {
		return "style"
	}
	open := strings.LastIndex(prefix, "<")
	if open > strings.LastIndex(prefix, ">") {
		match := htmlAttribute.FindStringSubmatch(prefix[open:])
		if len(match) == 0 {
			return "html_tag"
		}
		attribute := strings.ToLower(match[1])
		if strings.HasPrefix(attribute, "on") || attribute == "srcdoc" || attribute == "style" {
			return "active_attribute"
		}
		if attribute == "href" || attribute == "src" || attribute == "action" || attribute == "formaction" {
			return "url"
		}
		if match[2] == "" {
			return "unquoted_attribute"
		}
		return "html_attribute"
	}
	return "html_text"
}

func (f *flow) xssProtection(v value, context string, escaped bool) (string, string) {
	if escaped {
		if context == "html_text" || context == "html_attribute" {
			return "encoded", "The template HTML-escapes this value in HTML text or a quoted non-active attribute."
		}
		return "unknown", "Template HTML escaping does not establish protection in this output context."
	}
	allProtected := len(v.deps) > 0
	hasRaw, hasUnknown, hasIneffective := false, false, false
	for _, dependency := range v.deps {
		protected, unknown, ineffective := false, false, false
		for _, transfer := range dependency.Transfers {
			if transfer.Kind != "call_argument_dependency" {
				continue
			}
			m := f.set.modules[transfer.Span.File]
			if m == nil {
				unknown = true
				continue
			}
			call := nodeAt(m.tree.RootNode(), transfer.Span.StartByte, transfer.Span.EndByte)
			mode := verifiedOutputProtection(m, call)
			if mode == "markdown" && f.set.packageVersion("marked") == "0.3.5" {
				mode = "ineffective"
			}
			switch mode {
			case "encode":
				protected = context == "html_text" || context == "html_attribute"
				unknown = !protected
			case "sanitize":
				protected = context == "html_text"
				unknown = !protected
			case "ineffective":
				protected, ineffective = false, true
			default:
				unknown, protected = true, false
			}
		}
		allProtected = allProtected && protected && !unknown
		hasUnknown = hasUnknown || unknown
		if !protected && !unknown {
			hasRaw = true
			hasIneffective = hasIneffective || ineffective
		}
	}
	if allProtected {
		return "sanitized", "Every modeled source path passes through a verified encoder or sanitizer appropriate to this output context."
	}
	if hasRaw {
		if hasIneffective {
			return "ineffective", "A source path reaches raw HTML through marked 0.3.5, whose sanitize option has a documented entity/URL bypass (upstream marked PR 592)."
		}
		return "raw", "A source path reaches raw output without modeled context-appropriate encoding or sanitization; unresolved sibling paths remain in the evidence."
	}
	if hasUnknown {
		return "unknown", "A transformation or its adequacy for this context is unresolved; no protection is assumed."
	}
	return "unknown", "No complete source-to-output path was established."
}

func verifiedOutputProtection(m *sourceModule, call *sitter.Node) string {
	if call == nil || call.Kind() != "call_expression" {
		return ""
	}
	fn, args := call.ChildByFieldName("function"), call.ChildByFieldName("arguments")
	if fn == nil || args == nil || args.NamedChildCount() != 1 {
		return ""
	}
	parts := strings.Split(fn.Utf8Text(m.src), ".")
	if shadowed(call, parts[0], m.src) || bindingReassigned(m, parts[0]) {
		return ""
	}
	imp, ok := m.imports[parts[0]]
	if !ok {
		return ""
	}
	method := imp.name
	if len(parts) == 2 && (imp.name == "default" || imp.name == "*") {
		method = parts[1]
	}
	switch imp.module {
	case "marked":
		if len(parts) == 1 || method == "parse" {
			return "markdown"
		}
	case "entities", "he":
		if method == "encode" || method == "encodeHTML" || method == "escapeHTML" {
			return "encode"
		}
	case "escape-html":
		if len(parts) == 1 {
			return "encode"
		}
	case "sanitize-html":
		if len(parts) == 1 {
			return "sanitize"
		}
	case "dompurify", "isomorphic-dompurify":
		if method == "sanitize" {
			return "sanitize"
		}
	}
	return ""
}
