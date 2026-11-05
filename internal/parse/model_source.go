package parse

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	"strings"
)

type sourceComment struct{ start, end uint }

// Comments and standalone Python docstrings are presentation data, not operation
// witnesses. Keep their byte positions while excluding their prose from model input.
func sourceComments(root *sitter.Node, src []byte, python bool) []sourceComment {
	var out []sourceComment
	var visit func(*sitter.Node)
	visit = func(n *sitter.Node) {
		if n.Kind() == "comment" {
			out = append(out, sourceComment{n.StartByte(), n.EndByte()})
			return
		}
		if python && n.Kind() == "expression_statement" && n.NamedChildCount() == 1 && n.NamedChild(0).Kind() == "string" {
			parent := n.Parent()
			owner := parent.Parent()
			doc := parent.Kind() == "module" || parent.Kind() == "block" && owner != nil && (owner.Kind() == "function_definition" || owner.Kind() == "class_definition")
			literal := strings.ToLower(n.NamedChild(0).Utf8Text(src))
			quote := strings.IndexAny(literal, "\"'")
			if quote < 0 || strings.ContainsAny(literal[:quote], "fb") {
				doc = false
			}
			if doc {
				for i := uint(0); i < parent.NamedChildCount(); i++ {
					c := parent.NamedChild(i)
					if c.Kind() == "comment" {
						continue
					}
					if c.Id() == n.Id() {
						out = append(out, sourceComment{n.StartByte(), n.EndByte()})
						return
					}
					break
				}
			}
		}
		for i := uint(0); i < n.NamedChildCount(); i++ {
			visit(n.NamedChild(i))
		}
	}
	visit(root)
	return out
}
func modelSource(src []byte, start, end uint, comments []sourceComment) string {
	var masked []byte
	for _, c := range comments {
		a, b := c.start, c.end
		if a < start {
			a = start
		}
		if b > end {
			b = end
		}
		if a >= b {
			continue
		}
		if masked == nil {
			masked = append([]byte{}, src[start:end]...)
		}
		for i := a - start; i < b-start; i++ {
			if masked[i] != '\n' && masked[i] != '\r' && masked[i] != '\t' {
				masked[i] = ' '
			}
		}
	}
	if masked == nil {
		return ""
	}
	return string(masked)
}
