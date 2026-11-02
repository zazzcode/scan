# ThreatByte detection guide and expected findings

[Python implementation results](python-support-results.md) ·
[FastAPI answer guide](python-insecure-app-detection-guide.md) ·
[Existing comparison cases](../../docs/samples/ThreatByte/cases.json)

## Purpose and maintenance

This is a living answer guide and review checklist for Zazzles. It defines what a scan should
recognize, which controls it should preserve, and the evidence needed for each judgment. It is not
a frozen result-count target. Update expectations when reviewed application behavior changes;
retain independent rationale rather than making labels follow scanner output. These answers must
never become hints in production detector logic or decision-model requests.

SQL injection, input validation, logger leaks and supported XSS flavors are current analysis
families. Other advertised families remain development work for future branches, not abandoned
scope. Finding counts, source paths and coverage gaps must remain distinct.

Source references below were checked against ThreatByte revision
`ede28f27e7e5061866a67b19b6088df4dd11c569`. Function names and flow relations are the durable
anchors; line numbers aid navigation. Source files remained unchanged and the checkout was clean.
The manual review did not launch the app, import its modules, or access its on-disk database.

Source navigation: [routes.py](https://github.com/anotherik/ThreatByte/blob/ede28f27e7e5061866a67b19b6088df4dd11c569/server/routes.py),
[database binding](https://github.com/anotherik/ThreatByte/blob/ede28f27e7e5061866a67b19b6088df4dd11c569/db/initialize_db.py),
[README](https://github.com/anotherik/ThreatByte/blob/ede28f27e7e5061866a67b19b6088df4dd11c569/README.md).

## Author's documented answers

The author's [Part 1 walkthrough](https://medium.com/@anotherik/exploring-web-vulnerabilities-with-threatbyte-from-exploitation-to-defense-part-1-4a693b4fa3fc)
identifies dashboard file search as SQL-injectable: a quote produces a SQLite error, followed by
read-oriented injection examples. It locates input at `routes.py:126`, SQL construction at 137 and
execution at 138, and recommends driver binding. It separately demonstrates reflected XSS in the
same search field, locating the trust bypass at `dashboard.html:64` and recommending removal of
`safe`. Its credential-recovery chain also discusses MD5 password storage. Those are the article's
concrete examples; it is not an exhaustive route-by-route answer key. The profile/signup conclusions
below come from direct source review, not an invented article verdict for those routes.

A database syntax error is a useful signal, not sufficient proof by itself. The complete
source-to-execution path establishes the SQL defect. Broader impact claims or CVSS scores in an
article do not automatically establish every capability of a particular SQLite operation.

## GET /profile: SQL protection is correct

Expected SQL classification: **Protected for this query's SQL syntax**, not protected for the
entire endpoint. This is the HTML page route, distinct from the REST API's profile resource.

| Step | Actual source | Security meaning |
| --- | --- | --- |
| Entry | `profile`, routes 147–154 | If session username is absent, return the login redirect |
| Value | `session['user_id']`, line 151 | Passed as an argument to the helper, not appended to SQL |
| Helper | `get_user_details_from_database`, lines 406–413 | Opens a verified sqlite3 connection, makes a cursor and fetches one row |
| Connection | `db/initialize_db.py:6–9` | `sqlite3.connect(DATABASE_URI)` with `sqlite3.Row`; this is the real driver |
| SQL text | routes 410 | Constant `SELECT * FROM users WHERE id = ?` |
| Value slot | routes 411 | `cursor.execute(query, (user_id,))`; one unquoted placeholder, one bound value |
| Output | routes 152 | Passes the returned row to `profile.html`; rendering is a separate claim |

There is no f-string, percent formatting, concatenation or identifier substitution in this query.
The trailing comma makes `(user_id,)` a one-element tuple. Even a manipulated session value would
remain SQL data at this operation. Session authenticity, authorization, row ownership, missing users
and field validity require their own checks; none changes this parameter-binding conclusion.
The [sqlite3 binding contract](https://docs.python.org/3/library/sqlite3.html#how-to-use-placeholders-to-bind-values-in-sql-queries)
supports this distinction.

Reusable SQL checklist:

- [ ] Resolve the request/session value through the actual helper argument binding.
- [ ] Verify the driver, constant query text, placeholder count and separate operand.
- [ ] Preserve SQL protection for quote-bearing and expression-shaped bound values.
- [ ] Do not borrow the profile query's protection for dashboard search or HTML output.
- [ ] Keep authentication/authorization and input-validity questions separate.

### Profile output expectations

The same endpoint contains a stored-XSS defect. `country` comes from signup form input at routes 83,
is bound into the insert at 99, is read by the profile helper, and is emitted through
`{{ user.country | safe}}` at `client/templates/profile.html:41`. A second write path exists through
`update_profile` at routes 186 and 234. Binding stores the characters; it does not HTML-encode them.
The `safe` filter removes the normal escaping guarantee. **Expect stored XSS**, with writing and
reading entries and the HTML-text sink retained in the trace.

For SQLite scalar fields emitted without a trust bypass, the exact text-output claim should remain
protected under the verified default renderer. [Flask's documented setup](https://flask.palletsprojects.com/en/stable/templating/)
enables autoescaping for `.html` templates. The reviewed application contains no custom autoescape
or Markup conversion on this read/render path.

| Profile field | Template line | Expected output judgment |
| --- | --- | --- |
| username | 39 | Protected HTML text via autoescaping |
| email | 40 | Protected HTML text via autoescaping; not proof of valid email syntax |
| country | 41 | Stored-XSS finding: explicit `safe` bypass |
| created_at | 42 | Protected output encoding; exact stored-origin completeness is a separate question |
| role | 44 | Protected HTML text; authorization is a separate future-branch question |
| permissions | 45 | Protected HTML text; authorization is separate |
| team | 46 | Protected HTML text; not evidence of input validation |
| last_login | 47 | Protected output encoding; timestamp origin is separate |

- [ ] Remove `safe` for the country text output in a remediation copy; preserve ordinary escaping.
- [ ] Verify angle brackets, quotes and ampersands render as text in all protected expressions.
- [ ] Do not label stored-origin uncertainty as an encoding defect. If the scanner's claim requires
      a missing origin, retain Review and explain that narrower protection is directly established.

## POST /signup: SQL protection and validation have different answers

Source trace: form fields enter at routes 81–85. Lines 87–88 reject missing/empty values, and
89–90 reject unequal password/confirmation values. `custom_hash(password)` runs at 93; its actual
body at 403 computes MD5. The database helper resolves to sqlite3. Signup then executes:

| SQL operation | Text location / execute | Parameters | Expected SQL judgment |
| --- | --- | --- | --- |
| Insert user | 98 / 99 | username, email, country, hashed password; four `?` slots | Protected |
| Read created user | 105 / 106 | `(username,)`; one `?` slot | Protected |
| Set last_login | 113 | timestamp and username; two `?` slots | Protected |

All three keep values separate from static SQL. Quotes in a username do not change query structure.
Password hashing is unrelated to that guarantee. A sqlite3 connection's `with` block commits or
rolls back a transaction; it does not close the connection automatically. The statements after
that block are not automatically unreachable. See the
[connection context-manager contract](https://docs.python.org/3/library/sqlite3.html#how-to-use-the-connection-context-manager).

### Exact validation requirements

The linked signup form (`client/templates/signup.html:13–27`) establishes the following policy.
The table is deliberately per constraint, not a judgment of the entire account-creation flow.

| Requirement | Enforcing source | Expected answer |
| --- | --- | --- |
| username present/nonempty | routes 87–88 | Enforced |
| email present/nonempty | routes 87–88 | Enforced |
| country present/nonempty | routes 87–88 | Enforced |
| password present/nonempty | routes 87–88 | Enforced |
| confirmation present/nonempty | routes 87–88 | Enforced |
| password equals confirmation | routes 89–90 | Enforced |
| email syntax from form `type="email"` | No equivalent server check before persistence | Poor input-validation practice |

Database `NOT NULL` and `UNIQUE` declarations do not validate email syntax. Browser constraints can
be bypassed with direct requests. A returned error page before persistence enforces rejection even
if it returns HTTP 200. No source policy here establishes a country allow-list, password length,
complexity rule or username-as-email requirement; do not invent those defects.

- [ ] Keep all three SQL operations protected independently of the email-format defect.
- [ ] Identify email as the affected field and name format, rather than vaguely saying “validate input.”
- [ ] Verify direct requests with malformed nonempty email are rejected before insert in a repair copy.
- [ ] Preserve the five presence checks and password-matching guard as positive controls.
- [ ] Keep the stored-country-to-profile XSS trace visible even though its insert is bound.
- [ ] Track MD5 password storage for the future cryptography/password-security branch.

## Dashboard search: a required positive SQL-injection case

This is **GET /dashboard?search=...**, not the separate `/search?query=...` handler. The latter only
constructs a result list in the reviewed source and has no database execution.

Authenticated dashboard flow: routes 124 checks the session; 126 reads `search`; 131–134 performs a
bound username lookup and obtains `user_id`. If search is nonempty, 137 formats both values into
SQL and 138 executes it without a parameter operand:

```python
s_query = f"SELECT filename FROM files WHERE user_id = {user_id} AND filename LIKE '%{search_query}%'"
cursor.execute(s_query)
```

Expected: **SQL-injection finding**, not merely unknown validation or construction advice. The
initial bound username lookup cannot protect a later interpolated statement. A single quote in
search changes SQL parsing. The non-search alternative at 140, by contrast, binds `user_id` and
should retain SQL-syntax protection when its own provenance is established.

The source-specific repair pattern keeps both values out of SQL structure:

```python
s_query = "SELECT filename FROM files WHERE user_id = ? AND filename LIKE ?"
cursor.execute(s_query, (user_id, f"%{search_query}%"))
```

Here the f-string builds a **bound data value**, which is appropriate; it does not build SQL text.
LIKE `%` and `_` matching behavior is a separate search-policy question, not SQL injection after
binding. Do not recommend stripping apostrophes or arbitrary punctuation from legitimate searches.

The same source flows through routes 143 into the quoted input attribute at
`client/templates/dashboard.html:64`, where `| safe` permits reflected XSS. Expect **a separate
reflected-XSS claim**. A SQL fix alone does not fix that output; remove the trust bypass as well.

- [ ] Catch the search SQL sink at routes 138 with the input and construction assignments attached.
- [ ] Distinguish the protected lookup at 132 and protected no-search alternative at 140.
- [ ] Catch reflected XSS at dashboard template 64 with the surrounding attribute markup.
- [ ] Require independent SQL-binding and HTML-output repair/regression checks.

## Browser XSS and logger-leak expectations

At `dashboard.html:116–123`, a form submit handler reads a comment via jQuery `.val()`, then appends
an HTML string containing that value. This is an independent **DOM-XSS** path. It does not need an
invented API producer or database origin. A repair builds a list node and sets its text, for example
`$commentList.append($('<li>').text(comment))`. Native DOM `append(string)` and jQuery `.text()` are
text operations and must not inherit the HTML-parser sink classification.

The reviewed Python log calls are exception prints at routes 102/241 and `print(session)` at
`server/api/v1/api_v1.py:182`. No direct password/email log operand is established in signup or
profile. Exception text and a whole session need content/sensitivity tracing; the mere existence of
a logger, a container or a field name is not a confirmed personal-data leak.

- [ ] Keep the comment input-to-HTML-parser path as a positive DOM-XSS control.
- [ ] Trace emitted exception/session content before promoting a logger-leak finding.
- [ ] Do not invent password logging because the same route hashes or stores a password.
- [ ] In scanner fixtures/repair copies, pair actual password/email emissions with redaction,
      constant-event and non-content controls; do not inject fixture leaks into the benchmark app.

## Comparison guidance and remaining work

The source review agrees with the scanner's protected profile SQL and all three protected signup
SQL records. The scanner catches dashboard SQL injection, reflected search XSS, stored country XSS
and browser comment XSS. It separately reports the signup email-format practice defect.

Several lessons remain actionable rather than being hidden by totals:

- [ ] Resolve the dashboard no-search record's missing-request-source gate without inventing input.
      Its bound call shape is directly visible; a database-derived ID is not raw query syntax.
- [ ] Separate escaped-output adequacy from unnecessarily broad stored-origin uncertainty on
      created_at/last_login. Preserve uncertainty for any genuinely missing claim evidence.
- [ ] Verify protected SQL headings state the supported protection, rather than retaining generic review wording.
- [ ] Keep coarse SQL validation uncertainty separate from independently established presence guards.
      Operation-level model answers are not standalone field-policy findings.
- [ ] Show only final Protected claims in green; show concrete investigation for Review and concrete
      operation/constraint repairs for defects, consistently in HTML, JSON and SARIF.

An isolated SQLite check extracted only the reviewed SQL expressions from the source AST. With
SQLite 3.54.0, a quote-bearing username round-tripped through signup; expression-shaped profile and
username-lookup operands remained data; the timestamp update affected the intended row. An ASCII apostrophe in
the reconstructed dashboard f-string raised `OperationalError`; the proposed bound search handled
it as data. Markup round-tripped unchanged through bound storage, confirming that SQL binding is
not output encoding. These are operation checks, not HTTP exploitation or a deployment audit.

## Future-branch checklist

The README also advertises insecure configuration, weak cryptography, brute-force/session issues,
BOLA/BOPA/BFLA, CSRF, traversal, unrestricted upload and SSRF. Preserve all as planned coverage.
The [family inventory](../../docs/samples/python-documented-claims.md) and source cases provide
starting points; a family name alone is not a reviewed exploit or location-level answer.

- [ ] Add source-backed cases and protected counterparts for each future family before claiming coverage.
- [ ] Include the MD5 password path and relevant session/cookie configuration in that work.
- [ ] Keep article impact chains distinct from independently counted root defects.
- [ ] Recheck this guide against source and new scanner behavior as rules evolve; never train the
      production classifier to agree with app-specific expected labels.
