package report

import "github.com/zazzcode/scan/internal/graph"

// Guidance is evidence-scoped help, not a generated or verified patch.
// Every result includes the hardening posture: report proven vulnerabilities and
// weak or poor practice, because an agent can fix both cheaply.
type Guidance = graph.Recommendation

func guidance(use graph.StringUse) Guidance {
	kind := use.Kind
	if use.Check != nil {
		kind = use.Check.Kind
	}
	var g Guidance
	switch kind {
	case "sql_query", "sql_query_candidate":
		g = Guidance{Title: "Separate input data from SQL structure", Actions: []string{
			"Pass untrusted values through the verified driver's parameter-binding API. Keep SQL query, update, and insert structure separate from data. Review NoSQL APIs separately.",
			"Restrict each field as far as its business requirement allows. A string or length check is not injection protection. A comment or other free-text field may contain any character, so the sink construction is the injection control.",
			"Add regression cases for normal input, rejected malformed input and query-structure preservation; resolve inconclusive guard and binding evidence before applying a fix.",
		}, ExampleStatus: "not_generated; driver-specific adaptation and tests required"}
		if use.Completeness != nil && len(use.Completeness.Evidence["sink_binding"]) > 0 {
			g.Example = "// Illustrative Sequelize binding pattern; adapt schema and full authentication logic.\nconst rows = await sequelize.query(\n  'SELECT id FROM accounts WHERE email = $email',\n  { bind: { email: req.body.email } }\n);"
			g.ExampleStatus = "illustrative_template; not a verified application patch"
		}
	case "nosql_query", "nosql_query_candidate":
		g = Guidance{Title: "Keep request data out of NoSQL query structure", Actions: []string{
			"Reject objects and operators in fields that should be scalars before they reach a query, update, or insert document.",
			"Do not concatenate or interpolate request values into $where strings.",
			"Inspect numeric conversions and character filters on each actual path before judging their effect. Prefer data-only query APIs over executable $where strings.",
		}, ExampleStatus: "not_generated; driver-specific adaptation and tests required"}
	case "input_validation", "input_validation_candidate":
		g = Guidance{Title: "Validate and sanitize each request input", Actions: []string{
			"Define the server-side field policy, then enforce applicable type, presence, length, range, and format constraints before use. Free text must remain valid for its business purpose.",
			"Use destination-specific protection independently of field validation: SQL parameter binding or a verified NoSQL document API with constrained keys and scalar types. Do not apply a generic sanitizer as a substitute.",
			"Preserve an explicit Inconclusive result when a call receives the field but its effect has not been established.",
		}, ExampleStatus: "not_generated; field policy is route-specific"}
	case "account_creation", "account_creation_candidate":
		g = Guidance{Title: "Enforce required fields before account creation", Actions: []string{
			"Validate required fields independently of optional confirmation fields and before the persistence operation.",
			"Check runtime types, normalize according to policy, then enforce nonempty/format constraints; return immediately on rejection.",
			"Add appropriate model/database constraints as defense in depth and test actual persisted state, not only an HTTP response or challenge flag.",
		}, Example: "// Illustrative middleware; add the route's remaining schema requirements.\nconst { email, password } = req.body;\nif (typeof email !== 'string' || typeof password !== 'string' ||\n    email.trim().length === 0 || password.trim().length === 0) {\n  return res.status(400).json({ error: 'Invalid required fields' });\n}\n// Continue only after applying the remaining field policy.\nreturn next();", ExampleStatus: "illustrative_template; not a verified application patch"}
	default:
		g = Guidance{Title: "Review protection for the named use", Actions: []string{"Resolve the source-to-use relationship and applicable guards.", "Apply the destination/API-specific protection and the route's declared field requirements.", "Add vulnerable, protected and inconclusive regression cases."}, ExampleStatus: "no code example for an unverified library model"}
	}
	// Project tenet: an agent can harden weak practice cheaply, so every result says so.
	g.Actions = append([]string{"Remediate established vulnerabilities and poor practice. For Review items, resolve the named evidence before deciding whether or how behavior needs to change."}, g.Actions...)
	return g
}
