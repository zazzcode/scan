package report

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/zazzcode/zazzles/internal/graph"
)

func scopedGuidance(u graph.StringUse, f graph.Finding) Guidance {
	at := fmt.Sprintf("%s:%d", u.Sink.Span.File, u.Sink.Span.StartLine)
	if f.RecordType == "protection" {
		if u.Kind == "input_validation_candidate" && u.Check != nil {
			return Guidance{Title: "Preserve the enforced field constraints", Actions: []string{"Preserve the rejecting guard and return behavior for " + strings.Join(u.Check.Fields, ", ") + " at " + at + ". Exact policy: " + u.Check.Validation, "Test direct HTTP requests for each declared constraint and assert rejection occurs before a sensitive side effect. Do not infer SQL/XSS protection or stronger unstated requirements from this field-policy result."}, ExampleStatus: "no repair proposed for a protected field-policy claim"}
		}
		action := "Preserve the demonstrated protection at " + at + " for this exact operation: " + u.Sink.Text
		if u.XSS != nil {
			action = "Preserve escaping for " + u.Sink.Text + " at " + at + " in " + u.XSS.OutputContext + ". Do not add a trust bypass or disable the verified renderer's autoescape."
		}
		if u.Kind == "sql_query_candidate" {
			action = "Preserve the SQL text/data separation at " + at + ": " + u.Sink.Text + ". Keep values in the driver's separate parameter operand; do not move them into interpolated SQL."
		}
		return Guidance{Title: "Preserve the verified protection", Actions: []string{action, "Keep a regression test with HTML-significant or SQL-significant input, as appropriate to this claim. Input-validation requirements are assessed separately."}, ExampleStatus: "no repair proposed for a protected claim"}
	}
	if f.RecordType != "finding" {
		return investigationGuidance(u, at)
	}
	if u.Kind == "input_validation_candidate" && u.Check != nil {
		return Guidance{Title: "Enforce the declared field policy on the server", Actions: []string{
			"At " + at + ", validate " + strings.Join(u.Check.Fields, ", ") + " before using or persisting it. Enforce these exact source-declared constraints: " + fieldPolicyText(u),
			"Reject a direct HTTP request that violates each declared constraint and return before any side effect. Reuse an existing applicable server validator when available; browser attributes alone do not enforce this policy.",
			"Test direct requests bypassing the form, including missing, empty, malformed and boundary values relevant to the declared policy, plus a valid request. Assert rejected requests cause no persistence or other sensitive side effect.",
		}, ExampleStatus: "no automatic patch; preserve application-specific response and validation semantics"}
	}
	if u.Kind == "sql_query_candidate" && u.APIModel != nil && u.APIModel.Library == "sqlite3" {
		return sqliteGuidance(u, at)
	}
	if u.XSS != nil {
		return renderingGuidance(u, at)
	}
	g := guidance(u)
	g.Actions = append([]string{"Change the observed operation at " + at + ": " + u.Sink.Text}, g.Actions...)
	if filepath.Ext(u.Sink.Span.File) == ".py" {
		g.Example = ""
		g.ExampleStatus = "no verified source-derived replacement available"
		if u.Kind == "log_sink_candidate" {
			g.Actions = []string{"At " + at + ", remove the identified credential/personal-data operand from " + u.Sink.Text + ". Retain a constant event name and an allow-list of non-sensitive fields.", "Capture the emitted log in a test and assert that neither the original value nor any substring of it is present; test every branch that reaches this call."}
		}
	}
	return g
}
func investigationGuidance(u graph.StringUse, at string) Guidance {
	g := Guidance{Title: "Resolve the specific claim blockers", ExampleStatus: "no patch proposed while the claim is inconclusive", Actions: []string{"Investigate the operation at " + at + ": " + u.Sink.Text + ". Its current Review status does not establish a code defect."}}
	gaps := append([]graph.Gap{}, u.Assessment.BlockingGaps...)
	if u.Completeness != nil {
		gaps = append(gaps, u.Completeness.Blockers...)
	}
	seen := map[string]bool{}
	for _, gap := range gaps {
		key := gap.Kind + gap.Detail
		if seen[key] {
			continue
		}
		seen[key] = true
		g.Actions = append(g.Actions, fmt.Sprintf("Resolve %s at %s:%d: %s", gap.Kind, gap.At.File, gap.At.StartLine, gap.Detail))
		if len(g.Actions) >= 6 {
			break
		}
	}
	if u.XSS != nil && u.XSS.StorageKey != "" && !u.XSS.OriginResolved {
		g.Actions = append(g.Actions, "Trace the exact persisted field "+u.XSS.StorageKey+" from its writing entry to this read. Supply the write, read and intervening transformation/renderer evidence before labeling this stored-output claim protected.")
	}
	if u.Assessment.Status == "completed" && u.Assessment.FindingStatus == "needs_evidence" {
		detail := "Reconcile the semantic assessment with this exact source claim: " + u.Assessment.Reason
		if u.XSS != nil {
			detail += fmt.Sprintf(" Answers: context=%s; encoding=%s; sanitizer=%s; all paths protected=%s; exposure=%s. Static context=%s, protection=%s. Do not infer a code repair from contradictory answers.", u.Assessment.Answers["output_context"], u.Assessment.Answers["encoding_adequacy"], u.Assessment.Answers["sanitizer_effectiveness"], u.Assessment.Answers["xss_all_paths_protected"], u.Assessment.Answers["xss_exposure"], u.XSS.OutputContext, u.XSS.Protection)
		}
		g.Actions = append(g.Actions, detail)
	}
	if len(g.Actions) == 1 {
		g.Actions = append(g.Actions, "Identify the actual field requirement and the enforced guard before this operation, including rejection/return behavior and any helper body. Assess that requirement separately from the destination's protection.")
	}
	return g
}
func renderingGuidance(u graph.StringUse, at string) Guidance {
	g := Guidance{Title: "Repair this HTML output", ExampleStatus: "source-derived candidate; review intended output semantics and test before applying"}
	if strings.Contains(u.Sink.Text, "{{") && strings.Contains(u.Sink.Text, "|") && u.XSS.OutputContext != "script" && u.XSS.OutputContext != "url" {
		replacement := jinjaSafeFilter.ReplaceAllString(u.Sink.Text, "")
		if replacement != u.Sink.Text {
			g.Actions = []string{"At " + at + ", replace the raw template expression " + u.Sink.Text + " with " + replacement + " so the verified template engine escapes this value in " + u.XSS.OutputContext + ".", "Add a rendering test with quotes, angle brackets and ampersands; assert that the output remains text and cannot close its current attribute or element."}
			g.Example = replacement
			return g
		}
	}
	if match := jqueryListAppend.FindStringSubmatch(u.Sink.Text); len(match) > 0 {
		g.Actions = []string{"At " + at + ", build the list item as a DOM node and set its text from " + match[2] + ". The current string concatenation feeds that input to jQuery's HTML parser.", "Keep the existing list container and event flow. Test that an input containing an HTML tag is displayed literally and does not create another element."}
		g.Example = match[1] + ".append($('<li>').text(" + match[2] + "));"
		return g
	}
	g.Actions = []string{"At " + at + ", stop placing the source value into the HTML/template source consumed by " + u.Sink.Text + ".", "Use a constant template with the source value as a separate variable and explicitly enable context-appropriate autoescaping. Merely moving a value to Template.render does not enable Jinja autoescape.", "Test the same source-to-output path with HTML-significant input; retain a separate test for template-expression input so fixing template-source injection is not mistaken for fixing XSS."}
	return g
}

func fieldPolicyText(u graph.StringUse) string {
	var text []string
	if u.Completeness != nil {
		for _, e := range u.Completeness.Evidence["field_policy"] {
			text = append(text, e.Text)
		}
	}
	return strings.Join(text, "; ")
}
