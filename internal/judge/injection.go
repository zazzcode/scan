package judge

import "github.com/zazzcode/zazzles/internal/graph"

func questionsFor(use graph.StringUse) []question {
	if logUse(use) {
		return logQuestions
	}
	if use.Kind != "sql_query_candidate" && use.Kind != "nosql_query_candidate" {
		return questions
	}
	out := append([]question(nil), questions...)
	out = append(out, question{"injection_protection", "Can request-controlled input change executable database query syntax or selector operators on a supported path? Choose vulnerable only when the supplied source, verified library contract and relevant guard/path evidence establish that property. Choose protected for verified separation (bound SQL data, constrained scalar selectors, or effective context-specific conversions/filters); poor construction practice alone is not injection. Choose unknown for missing source, API, guard or necessary path facts. Choose not_applicable only if this is not an injection-sensitive use. For protected, establish protection on ALL feasible input paths, not merely one branch. Evaluate every output alternative: a verified Number result or numeric literal cannot carry arbitrary input text into an expression; if all alternatives have that constraint, the conditional does not introduce a raw-text bypass. A Number conversion in one ternary branch does not constrain an unconverted sibling branch. Truncation or newline removal alone does not exclude query operators/calls; assess other alternatives without borrowing the numeric guarantee. Evaluate the actual operation and slot, not category labels or source comments."})
	if use.Kind == "nosql_query_candidate" {
		out = append(out,
			question{"all_paths_protected", "Does positive evidence establish injection protection on ALL relevant feasible branches? yes=every source-to-operation alternative has adequate syntax/operator protection; no=at least one supported branch leaves attacker syntax/operator control; unknown=an alternative's protection cannot be decided. A numeric guarantee applies only to that conversion's output, never automatically to sibling branches. Business/format validation is separate."},
			question{"source_coverage", "For the injection question, is a relevant externally controlled source and its transfer into the actual selector/expression slot established? yes=the request source and transfer are established; no=positive evidence the relevant slot is not externally controlled; unknown=source or transfer facts are missing."},
			question{"library_coverage", "For the injection question, do the verified import/constructor evidence and supplied API contract establish the actual collection and selector/expression behavior? yes=the actual collection/API slot is identified with its contract; no=positive evidence this is a different API; unknown=library or slot semantics are missing."},
			question{"guard_coverage", "For the injection question, are relevant middleware and value constraints described well enough to assess their effect? yes=relevant guard/check effects are available (including an absence of guards in the supplied complete slice); no=an applicable check body is missing; unknown=coverage is ambiguous. Effective and ineffective guards can both have adequate coverage."},
			question{"path_coverage", "For the injection question, is at least one relevant path's feasibility established? yes=a relevant source-to-operation path is supported under stated possible conditions; no=the operation is proven unreachable for that input; unknown=necessary branch/callback/return facts are missing. A reported static path does not require executing the application."})
	}
	return out
}

func questionChoices(q question) []map[string]string {
	if c := logChoices(q.name); c != nil {
		return c
	}
	if q.name == "injection_protection" {
		return []map[string]string{
			{"value": "vulnerable", "description": "The supplied source and verified operation model support input changing executable query syntax or selector operators on a relevant path."},
			{"value": "protected", "description": "The supplied evidence establishes that the relevant input remains data or is constrained so it cannot change syntax/operators for this operation."},
			{"value": "unknown", "description": "Required source, library, guard, transformation, or path evidence is missing or ambiguous."},
			{"value": "not_applicable", "description": "This operation has no applicable injection question."},
		}
	}
	switch q.name {
	case "source_coverage", "library_coverage", "guard_coverage", "path_coverage", "all_paths_protected":
		return []map[string]string{
			{"value": "yes", "description": "The fact asked about by this question is established by the supplied evidence."},
			{"value": "no", "description": "The supplied evidence establishes the contrary fact or identifies the missing body, as defined in this question."},
			{"value": "unknown", "description": "The evidence does not establish the requested fact or its contrary."},
		}
	}
	return choices
}
