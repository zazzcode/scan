# VulnPy detection guide and comparison checklist

[Sample reports](../../docs/samples/vulnpy/README.md) ·
[Expected cases](../../docs/samples/vulnpy/cases.json) ·
[ThreatByte guide](threatbyte-detection-guide.md) ·
[FastAPI answer guide](python-insecure-app-detection-guide.md)

## Purpose

Use this living guide to expand Zazzles coverage without changing expected answers to match current
output. SQL injection, input validation, logger leaks and XSS are current analysis families. Other
VulnPy trigger families are planned work for future branches. A supported vulnerability family does
not imply support for every framework, callback, dynamic registration or library pattern.

Reviewed source: [Contrast-Security-OSS/vulnpy](https://github.com/Contrast-Security-OSS/vulnpy),
MIT, revision `b81885e58283a3449b578d20c20557aff637e5ea`. README describes a vulnerable function
library with framework adapters, not just one fixed app. The examples cover Flask, FastAPI, Django,
Falcon, Pyramid, Bottle, WSGI, Quart and aiohttp variants. The README also announces that the project
is no longer maintained. Source, README and tests were read without installing or executing the app.

## SQL injection: three required positive cases

`src/vulnpy/trigger/sqli.py` contains three variants. `_db_reset` establishes an actual in-memory
sqlite3 connection; `_execute` supplies its cursor to a nested callback. The public trigger function
captures user_input and the callback formats it into SQL before execution. The wrapper catches
exceptions and returns an error string; catching errors does not repair query construction.

| Case | Source construction and sink | Expected answer |
| --- | --- | --- |
| VP-SQL-EXECUTE | do_sqlite3_execute → format at 36 → execute at 37 | SQL-injection positive; first inserted value becomes SQL text |
| VP-SQL-EXECUTEMANY | do_sqlite3_executemany → format at 47 → executemany at 48 | SQL-injection positive; the parameter operand binds a different value, not the interpolated first value |
| VP-SQL-EXECUTESCRIPT | do_sqlite3_executescript → format at 58 → executescript at 59 | SQL-injection positive; executable script contains request text and is not value-bound |

The map in `trigger/util.py` discovers trigger functions using imports/introspection; get_trigger
selects the named function. Framework view factories call it with request input. For Flask, input
comes from the User-Input header, query parameter or form via `_get_user_input`; for FastAPI, it is
the query parameter of a dynamically constructed view. SQL trigger names become paths such as
`/vulnpy/sqli/sqlite3-execute`, with Flask allowing GET/POST and optional trailing slash.

`tests/trigger/test_sqli.py` supplies concrete row-manipulation expectations for all three variants,
including a quote-containing invalid input. Those source tests are stronger expected-behavior
evidence than the mere word “sqli” in a filename. They were inspected, not executed in the target.

- [x] Resolve generated registration and imported trigger selection without executing project code.
- [x] Retain the input captured by the nested callback and the cursor passed by `_execute`.
- [x] Keep the SQL-text value separate from executemany's bound value.
- [x] Do not label constant schema initialization or SELECT_ALL as user-controlled SQL.
- [x] Add protected parameterized counterparts in scanner fixtures/remediation copies.

## XSS: required positive paths

| Case | Source trace | Expected answer |
| --- | --- | --- |
| VP-XSS-FLASK | _get_user_input → generated xss/raw view → concatenation at flask/blueprint.py:55 → returned HTML at 57 | Reflected XSS |
| VP-XSS-FASTAPI | user_input query parameter → generated xss/raw view → concatenation at fastapi/vulnerable.py:42 → HTMLResponse at 44 | Reflected XSS |
| VP-XSS-JINJA | apps/flask_app.py query input at 20 → template.render at 19 → apps/templates/jinja/xss.html | Raw template-output XSS; the environment has no autoescape setting and the template explicitly uses safe |

`trigger/xss.py:do_raw` returns an empty string; it is not the actual HTML sink. The vulnerable
markup is assembled by the framework adapter. The framework tests assert that input appears in
that output, but their ordinary marker is not itself a browser exploit test.

- [x] Follow the actual HTML response/renderer rather than inventing a sink in do_raw.
- [x] Resolve factory parameters name/trigger before declaring a raw branch feasible.
- [ ] Model the Environment/get_template/render binding and the concrete template trust bypass.
- [x] Keep raw versus escaped text counterparts and unknown custom configuration controls.
- [ ] Extend the same review to other framework adapters as their extraction support is implemented.

## Validation and logger-leak controls

The arbitrary user_input fields have no general business-format requirement established by their
names. Query protection and output encoding are the primary concrete constraints in these examples.
Do not prescribe arbitrary alphanumeric or length restrictions as an invented field-policy finding.
Derive any future field constraints from actual route contracts and enforcing implementation.

Constant diagnostic prints in example routes and trigger lookup errors do not automatically expose
password/email content. A path that logs such content needs independent source/classification and
emission evidence. No confirmed PII logger-leak answer is invented here merely to fill that family.

- [x] Keep missing business policy as Review, separately from concrete SQL/XSS defects.
- [ ] Require actual sensitive content and logger binding before promoting a leak.
- [ ] Test password/email leaks and protected logging controls in separate scanner fixtures.

## Historical baseline and known coverage work

Baselines used the same clean target revision and pinned local rules without credentials or model
calls. Zazzles ran on the original read-only checkout; comparison engines parsed copies of tracked
files. Raw artifacts are under `_test-findings/vulnpy/` and normalized results under the sample link.
Tool output scopes/exclusions, rule hashes, exit codes and source hashes are retained.

| Tool | Result | Meaning |
| --- | --- | --- |
| Zazzles | 11 static route entries, 6 unassessed records, 4 explicit registration/path warnings | Generated trigger routes remain unsupported; one separate Jinja output candidate is unresolved |
| Semgrep CE 1.183.0 | 8 results; no scan errors | Autoescape configuration warning plus hash, deserialization, shell and container findings; no three SQL trigger findings |
| Bearer e3753de, rules 82c335f | 21 results; workers completed | Command/code execution, deserialization, XML, traversal, weak randomness/hash; no three SQL trigger findings |

The first Bearer attempt returned zero after sandbox-blocked worker sockets. That report is invalid
and retained only as failed history. The successful second run, with local worker sockets allowed,
is the baseline. No comparison-rule packs were downloaded for this addition.

The six source-positive cases are retained as **known misses**, not changed into negatives. Four
additional checks require explicit warnings for generated FastAPI paths and Flask add_url_rule
registrations. The case evaluator reports six misses and four coverage-control passes and exits
nonzero. The original three-app 29-case checks remain separate. A framework/model limitation is
coverage work within existing families, not evidence the app is safe.

- [ ] Add bounded, source-backed factory/registration/dispatch/callback support in future hardening.
- [ ] Upgrade the Jinja Environment/template binding so the candidate becomes a supported raw claim.
- [ ] Recalibrate only once necessary static evidence exists; model confidence cannot fill these gaps.
- [ ] Record tool disagreements per source case; 8 versus 21 is not an accuracy ranking.

## Future vulnerability-family branches

Trigger modules also cover command injection, deserialization, weak hashing/randomness, parameter
pollution, traversal, ReDoS, SSRF, unsafe code execution, XXE and XPath injection. Preserve these as
future answer-key work. The comparison tools' broader results can guide source review, but do not
replace an independent expected result and protected counterpart for each new class.

- [ ] For each future family, review the actual operation, input contract, tests and framework path.
- [ ] Keep sink-level defects distinct from the many framework/method instances exposing them.
- [ ] Preserve target source integrity and avoid running its deliberately unsafe trigger functions.

## D5 confirmation and remaining gaps

The [source reference's confirmed-findings table](../../docs/samples/vulnpy/reference.md#findings-confirmed-by-scan-and-code-inspection--d5)
records the five SQL/raw-XSS source cases recovered by the full 2026-11-04 scan and then checked
against actual code. Main-app acceptance covers 12 method/route contexts; separate mounted instances
produce 28 contextual vulnerability findings across the full checkout. These are static source
confirmations, not runtime exploit executions. The 827 model assessments completed successfully;
model choices without necessary source/policy evidence remain Review.

Upstream supplies deliberately vulnerable functions and behavioral tests, rather than a formal
expected-findings manifest. The SQLite tests establish concrete input/row behavior and framework
XSS tests establish input inclusion; source review supplies the security reasoning. Evaluation labels
remain outside scanner/model inputs. Historical saved findings, raw answers, hashes and case labels
are unchanged. `VP-XSS-JINJA`, other framework adapters and additional families remain separate work.
Use the [D5 results](common-entry-dispatch-results.md) for checks, artifacts and supported-subset limits.
