package parse

import (
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	"github.com/zazzcode/zazzles/internal/analysis"
)

// conditionalConstraints describes output alternatives, rather than propagating
// a conversion's guarantee across sibling branches. It does not decide injection.
func conditionalConstraints(m *sourceModule, n *sitter.Node, c *analysis.Completeness) {
	if n == nil || n.Kind() != "ternary_expression" {
		return
	}
	var alternatives []*sitter.Node
	var collect func(*sitter.Node, int) bool
	collect = func(value *sitter.Node, depth int) bool {
		if value == nil || depth > 8 {
			return false
		}
		if value.Kind() == "parenthesized_expression" {
			return collect(value.NamedChild(0), depth+1)
		}
		if value.Kind() == "ternary_expression" {
			return collect(value.ChildByFieldName("consequence"), depth+1) && collect(value.ChildByFieldName("alternative"), depth+1)
		}
		alternatives = append(alternatives, value)
		return true
	}
	if !collect(n, 0) {
		return
	}
	allNumeric := true
	for _, value := range alternatives {
		fact := "No numeric output guarantee is established for this alternative; assess its actual constraints independently."
		if numericQueryValue(m, value) {
			fact = "This alternative produces only a number on normal completion (or throws); the original request text is not its output."
		} else {
			allNumeric = false
		}
		c.ModelFacts = append(c.ModelFacts, "Conditional output alternative "+value.Utf8Text(m.src)+" at "+m.file+": "+fact)
	}
	if allNumeric {
		c.ModelFacts = append(c.ModelFacts, "Every output alternative of the exact expression "+n.Utf8Text(m.src)+" at "+m.file+" is proven numeric on normal completion. There is no unconstrained text-output sibling branch in this expression. This is a value-type fact, not business validation or approval of executable-string construction; assess the actual database slot separately.")
	}
}

func numericQueryValue(m *sourceModule, n *sitter.Node) bool {
	if n == nil {
		return false
	}
	switch n.Kind() {
	case "number":
		return !strings.HasSuffix(n.Utf8Text(m.src), "n")
	case "parenthesized_expression":
		return numericQueryValue(m, n.NamedChild(0))
	case "call_expression":
		fn := n.ChildByFieldName("function")
		return fn != nil && fn.Utf8Text(m.src) == "Number" && nativeStringOperation(m, n)
	}
	return false
}
