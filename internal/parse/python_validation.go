package parse

import (
	"path/filepath"
	"regexp"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/zazzcode/zazzles/internal/graph"
)

var pythonForm = regexp.MustCompile(`(?is)<form\b([^>]*)>(.*?)</form\s*>`)
var pythonInput = regexp.MustCompile(`(?is)<input\b[^>]*>`)
var pythonHTMLAttr = regexp.MustCompile(`(?i)([a-z][a-z0-9_-]*)\s*=\s*["']([^"']*)["']`)

// formPolicies records client-declared requirements only for a form that posts to
// this handler. Browser validation is policy evidence, never server enforcement.
func (s *pythonSet) formPolicies(r pythonRoute) []graph.Evidence {
	if r.receiver.kind != "flask" || r.receiver.templates == "" {
		return nil
	}
	var result []graph.Evidence
	seen := map[string]bool{}
	pyWalk(r.fn.ChildByFieldName("body"), func(n *sitter.Node) {
		if n.Kind() != "call" || s.qualified(r.module, pyText(r.module, n.ChildByFieldName("function"))) != "flask.render_template" {
			return
		}
		args := pyArgs(n)
		if len(args) == 0 {
			return
		}
		name, ok := pyLiteral(r.module, args[0])
		if !ok {
			return
		}
		file, ok := s.text.readText(filepath.ToSlash(filepath.Join(r.receiver.templates, name)))
		if !ok {
			return
		}
		for _, form := range pythonForm.FindAllStringSubmatchIndex(file.Text, -1) {
			tag := file.Text[form[2]:form[3]]
			handler := pyText(r.module, r.fn.ChildByFieldName("name"))
			actionMatches := strings.Contains(tag, "url_for('"+handler+"')") || strings.Contains(tag, `url_for("`+handler+`")`)
			if !actionMatches {
				continue
			}
			for _, input := range pythonInput.FindAllStringIndex(file.Text[form[4]:form[5]], -1) {
				start, end := form[4]+input[0], form[4]+input[1]
				text := file.Text[start:end]
				attributes := map[string]string{}
				for _, a := range pythonHTMLAttr.FindAllStringSubmatch(text, -1) {
					attributes[strings.ToLower(a[1])] = a[2]
				}
				if attributes["name"] == "" {
					continue
				}
				policy := attributes["type"] == "email" || strings.Contains(strings.ToLower(text), "required")
				for _, key := range []string{"pattern", "minlength", "maxlength", "min", "max"} {
					policy = policy || attributes[key] != ""
				}
				if !policy {
					continue
				}
				e := textRange(file, start, end, "form_field_requirement")
				key := e.Span.File + ":" + text
				if !seen[key] {
					seen[key] = true
					result = append(result, e)
				}
			}
		}
	})
	return result
}
