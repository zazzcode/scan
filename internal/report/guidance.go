package report

import "github.com/zazzcode/scan/internal/graph"

// Guidance is evidence-scoped help, not a generated or verified patch.
// Every result includes the hardening posture: report proven vulnerabilities and
// weak or poor practice, because an agent can fix both cheaply.
type Guidance struct {
	Title         string   `json:"title"`
	Actions       []string `json:"actions"`
	Example       string   `json:"illustrative_example,omitempty"`
	ExampleStatus string   `json:"example_status"`
}

func guidance(use graph.StringUse) Guidance {
	kind := use.Kind
	if use.Check != nil {
		kind = use.Check.Kind
	}
	var g Guidance
	switch kind {
	case "sql_query", "sql_query_candidate":
		g = Guidance{Title: "Separate input data from SQL structure", Actions: []string{
			"Pass untrusted values through the verified driver's parameter-binding API. Do not build executable SQL, NoSQL, update, or insert statements from the string.",
			"Restrict each field as far as its business requirement allows. A string or length check is not injection protection. A comment or other free-text field may contain any character, so the sink construction is the injection control.",
			"Add regression cases for normal input, rejected malformed input and query-structure preservation; inspect unknown guards and binding facts before applying a fix.",
		}, ExampleStatus: "not_generated; driver-specific adaptation and tests required"}
		if use.Completeness != nil && len(use.Completeness.Evidence["sink_binding"]) > 0 {
			g.Example = "// Illustrative Sequelize binding pattern; adapt schema and full authentication logic.\nconst rows = await sequelize.query(\n  'SELECT id FROM accounts WHERE email = $email',\n  { bind: { email: req.body.email } }\n);"
			g.ExampleStatus = "illustrative_template; not a verified application patch"
		}
	case "nosql_query", "nosql_query_candidate":
		g = Guidance{Title: "Keep request data out of NoSQL query structure", Actions: []string{
			"Reject objects and operators in fields that should be scalars before they reach a query, update, or insert document.",
			"Do not concatenate or interpolate request values into $where strings.",
			"Treat a length limit or partial character filter as unresolved until it is shown to dominate the sink.",
		}, ExampleStatus: "not_generated; driver-specific adaptation and tests required"}
	case "input_validation", "input_validation_candidate":
		g = Guidance{Title: "Validate and sanitize each request input", Actions: []string{
			"Enforce type, presence, length, and format on the server for every request field the route reads.",
			"Sanitize for the destination only after validation, and do not treat that sanitizer as proof the query structure is safe.",
			"Preserve an explicit result when a call receives the field but its effect is unknown.",
		}, ExampleStatus: "not_generated; field policy is route-specific"}
	case "account_creation", "account_creation_candidate":
		g = Guidance{Title: "Enforce required fields before account creation", Actions: []string{
			"Validate required fields independently of optional confirmation fields and before the persistence operation.",
			"Check runtime types, normalize according to policy, then enforce nonempty/format constraints; return immediately on rejection.",
			"Add appropriate model/database constraints as defense in depth and test actual persisted state, not only an HTTP response or challenge flag.",
		}, Example: "// Illustrative middleware; add the route's remaining schema requirements.\nconst { email, password } = req.body;\nif (typeof email !== 'string' || typeof password !== 'string' ||\n    email.trim().length === 0 || password.trim().length === 0) {\n  return res.status(400).json({ error: 'Invalid required fields' });\n}\n// Continue only after applying the remaining field policy.\nreturn next();", ExampleStatus: "illustrative_template; not a verified application patch"}
	default:
		g = Guidance{Title: "Review protection for the named use", Actions: []string{"Resolve the source-to-use relationship and applicable guards.", "Apply the destination/API-specific protection and the route's declared field requirements.", "Add vulnerable, protected and unknown regression cases."}, ExampleStatus: "no code example for an unverified library model"}
	}
	// Project tenet: an agent can harden weak practice cheaply, so every result says so.
	g.Actions = append([]string{"Report this whether it is a confirmed vulnerability or weak or poor practice. An agent can harden the route cheaply; do not drop the finding only because exploitability is unproven."}, g.Actions...)
	return g
}
