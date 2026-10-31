# Login email validation and SQL-use protection

**Phase acceptance case, 2026-10-31; Juice Shop `v20.3.0`.** Keep this independent from the
[registration case](07-registration-validation-case.md). Both are declared in the
[phase requirements file](../evaluation/juice-shop/validation-checks.json).

## Source-backed case

The local registration is `POST /rest/user/login` at `server.ts:624`. It calls the imported login
factory. `routes/login.ts:32–34` passes `req.body.email || ''` into the SQL-text template used by
`models.sequelize.query`. The password expression is hashed first; that is a different operand and
cannot establish protection for email. The query options set a return model and plain-result behavior;
they do not bind the email value as a SQL data parameter.

Relevant provenance includes the namespace import at `routes/login.ts:13`, the Sequelize import and
constructor/factory in `models/index.ts:28–44`, and its exported connection at `models/index.ts:77`.
Earlier middleware and the pre-login helper must be inspected because validation could occur before
the query or mutate the request in place. The checked bundle now supplies those source definitions
rather than relying on a function's name. Mutable bindings and dynamic extensions are not universally
resolved; supported findings are scoped to the documented models/analysis subset.

## Independent questions

- **Field requirement:** is the request email a string with acceptable email-address syntax before
  authentication, according to the declared application policy?
- **SQL-use requirement:** is email kept separate from executable SQL structure at this query?
- **Password handling:** password storage/hash adequacy is separate from email format and SQL use;
  it requires its own rule/claim and regression evidence.

The structural observation is raw email interpolation into SQL text. Injection protection is not
established by the existence of an email-format check or the hashing of a different value. A model's
insufficient-validation answer alone must not become an injection finding without source/sink and
query-context evidence. Conversely, parameterized SQL can protect a query even if a separate field
validation policy is violated.

References: [Sequelize query semantics](https://sequelize.org/docs/v6/core-concepts/raw-queries/),
[OWASP SQL injection prevention](https://cheatsheetseries.owasp.org/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.html),
and the user's [Juice Shop walkthrough](https://www.automat-it.com/blog/hack-the-owasp-juice-shop-application-and-protect-it-with-aws-waf-part-2/).
The pinned source is the case evidence; the walkthrough is benchmark context, not a scanner verdict.

## Complementary remediation expectations

A future route report should propose parameter binding for SQL operands, server-side email type/format
constraints, and authentication regression tests as separate measures. Suitable password storage and
other authentication controls are additional independent work. Do not recommend email validation
alone as a sufficient SQL-injection fix or claim any single change makes the route impervious.

## Verification for this phase

Checked mode must supply route/handler, relevant pre-use helpers, the email dependency, executable
query operand and library-binding evidence. Password-only transformation uncertainty and result-model
metadata do not automatically block an email SQL claim. An unresolved helper receiving email, a
missing sink binding, or incomplete required context still blocks the affected claim.

Verify vulnerable interpolation, safe parameter binding, custom validation, ignored validation results,
shadowed helpers, unrelated bookkeeping, and unknown guard/library cases. Snapshot/golden and policy
fixtures test evidence construction and status routing. A live Luna call on this known vulnerable
case tests integration; it does not establish held-out accuracy or runtime exploitability. Results:
[phase evaluation](../evaluation/juice-shop/results.md).
