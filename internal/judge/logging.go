package judge

import "github.com/zazzcode/zazzles/internal/analysis"

// logQuestionVersion is versioned separately so query assessments keep their protocol.
const logQuestionVersion = "log-use-v1"

var logQuestions = []question{
	{"data_class", "Classify the data carried by the request value(s) that this logging operation receives, as named in the evidence. credential = a secret that authenticates or authorizes (password, token, API key, session cookie, Authorization header value, one-time code). personal_data = information that identifies or describes a person (email, phone, government ID, payment card, birth date). other = neither. Use the request contract, how the route uses the value (hashing, comparison with a stored credential, storage in a credential attribute) and HTTP semantics. A field or variable name alone is a hint, not proof; choose unknown when the evidence does not establish the class."},
	{"content_exposure", "Does this emission write the content of that value, wholly or in part, into log output? yes = the value, a substring, or a text-preserving conversion reaches a logged argument or the logged URL; no = only derived facts (length, presence, equality result, type) or a verified redacted form are logged, or the value is not logged; unknown = an unmodeled transformation, serializer or format decides. For request-line logging, the logged URL includes its query string when the format contains :url or is a predefined morgan format."},
	{"redaction", "Is effective redaction or masking applied to that value before this emission on every relevant path? yes = a supplied transformation or logger configuration demonstrably removes or masks the content before emission; no = no redaction is shown; unknown = a call or configuration whose effect is not shown may redact; not_applicable = no content is exposed. A function named mask or redact whose body is not shown is unknown. Redaction after the emission does not count."},
	{"log_exposure", "Does this logging operation expose sensitive data content to the log? vulnerable = credential or personal-data content reaches the log without effective redaction on a supported path; protected = only non-revealing facts or effectively redacted content is logged; unknown = data class, content or redaction cannot be decided from the supplied evidence; not_applicable = the logged data is not sensitive. Constant messages, startup banners, ports, and error stacks without sensitive request data are not exposures. Do not judge who reads the log or its retention; that is outside this question."},
}

func logChoices(name string) []map[string]string {
	switch name {
	case "data_class":
		return []map[string]string{
			{"value": "credential", "description": "The value authenticates or authorizes a principal."},
			{"value": "personal_data", "description": "The value identifies or describes a person."},
			{"value": "other", "description": "The evidence establishes that the value is neither a credential nor personal data."},
			{"value": "unknown", "description": "The evidence does not establish the data class."},
		}
	case "content_exposure":
		return []map[string]string{
			{"value": "yes", "description": "The value's content, or part of it, is written to the log."},
			{"value": "no", "description": "Only derived non-content facts or a verified redacted form are written, or the value is not logged."},
			{"value": "unknown", "description": "An unmodeled transformation, serializer or format decides what is written."},
		}
	case "redaction":
		return []map[string]string{
			{"value": "yes", "description": "Effective redaction before emission is shown on every relevant path."},
			{"value": "no", "description": "No redaction is shown."},
			{"value": "unknown", "description": "A call or configuration whose effect is not shown may redact."},
			{"value": "not_applicable", "description": "No content of the value is exposed."},
		}
	case "log_exposure":
		return []map[string]string{
			{"value": "vulnerable", "description": "Sensitive content reaches the log without effective redaction on a supported path."},
			{"value": "protected", "description": "Only non-revealing facts or effectively redacted content are logged."},
			{"value": "unknown", "description": "Data class, content or redaction cannot be decided from the supplied evidence."},
			{"value": "not_applicable", "description": "The logged data is not sensitive."},
		}
	}
	return nil
}

func logUse(use analysis.CandidateUse) bool { return use.Kind == "log_sink_candidate" }

// applyLogPolicy requires binding and source evidence and answers that agree with
// each other and with the scanner's own content facts before a definitive status.
func applyLogPolicy(use analysis.CandidateUse, a *analysis.Assessment) {
	ans := a.Answers
	switch ans["log_exposure"] {
	case "vulnerable":
		a.Protection = "insufficient"
	case "protected":
		a.Protection = "sufficient"
	case "not_applicable":
		a.Protection = "not_applicable"
	default:
		a.Protection = "uncertain"
	}
	var blockers []analysis.Gap
	for _, g := range use.Gaps {
		switch g.Kind {
		case "log_literal_sensitivity_unresolved", "log_transform_unresolved", "input_helper_unresolved", "python_decorator_unresolved", "route_mount_unresolved", "work_limit", "dependency_limit", "transfer_limit":
			blockers = append(blockers, g)
		}
	}
	binding := false
	for _, e := range use.Context {
		binding = binding || e.Kind == "logger_binding"
	}
	if !binding {
		blockers = append(blockers, analysis.Gap{Kind: "logger_binding_missing", Detail: "verified logger or request-logging binding is required"})
	}
	source := false
	for _, d := range use.Dependencies {
		source = source || d.Source.Kind == "request_source" || d.Source.Kind == "literal_sensitive_source" || d.Source.Kind == "storage_source"
	}
	if !source {
		blockers = append(blockers, analysis.Gap{Kind: "log_source_missing", Detail: "request, literal personal-data or verified storage source evidence is required"})
	}
	a.BlockingGaps = blockers
	if len(blockers) > 0 {
		a.FindingStatus = "needs_evidence"
		a.Protection = "uncertain"
		a.Reason = "logging claim evidence is incomplete"
		return
	}
	sensitive := ans["data_class"] == "credential" || ans["data_class"] == "personal_data"
	nonContentOnly := hasLogGap(use, "log_non_content_only") && !hasLogGap(use, "log_protocol_credential", "log_sensitive_field_hint", "log_request_container", "log_transform_unresolved", "log_request_line_credential", "log_request_line_hint")
	switch ans["log_exposure"] {
	case "vulnerable":
		if !sensitive || ans["content_exposure"] != "yes" || (ans["redaction"] != "no" && ans["redaction"] != "not_applicable") || nonContentOnly {
			a.FindingStatus = "needs_evidence"
			a.Protection = "uncertain"
			a.Reason = "Exposure answer is not supported by consistent data-class, content and redaction answers, or contradicts the scanner's non-content facts."
			return
		}
		a.FindingStatus = "confirmed_vulnerability"
		a.Reason = "Verified logger binding and source evidence, with consistent answers that sensitive content reaches the log without redaction; this is a static judgment, not runtime confirmation."
	case "protected":
		if ans["content_exposure"] != "no" && ans["redaction"] != "yes" && ans["data_class"] != "other" {
			a.FindingStatus = "needs_evidence"
			a.Protection = "uncertain"
			a.Reason = "Protection answer is not supported by content, redaction or data-class answers."
			return
		}
		a.FindingStatus = "protected_for_claim"
		a.Reason = "Assessment supports that this emission does not expose sensitive content for the supplied claim."
	case "not_applicable":
		if ans["data_class"] != "other" {
			a.FindingStatus = "needs_evidence"
			a.Protection = "uncertain"
			a.Reason = "Not-applicable exposure requires a non-sensitive data class."
			return
		}
		a.FindingStatus = "not_applicable"
		a.Reason = "Assessment establishes that the logged data is not sensitive."
	default:
		a.FindingStatus = "needs_evidence"
		a.Reason = "Logging exposure assessment is missing or inconclusive."
	}
}

func hasLogGap(use analysis.CandidateUse, kinds ...string) bool {
	for _, g := range use.Gaps {
		for _, k := range kinds {
			if g.Kind == k {
				return true
			}
		}
	}
	return false
}
