# Node.js Goof capability inventory — version 1

**Source-inspected 2026-11-01; initial baseline complete.** Target: Snyk Labs `nodejs-goof`, package
`1.0.1`, revision `add14ba59e98240d9e00a235dd7d42cd61ae9912`, clean local checkout at
`<DEV_ROOT>/nodejs-goof`. This extends the [comparator capability inventory](bearer-capabilities.md)
with target-specific cases; it does not amend the historical two-target results.

Sources inspected: target `README.md`, `package.json`, `app.js`, `routes/index.js`, and `views/`.
The README describes dependency vulnerabilities, Docker scanning, runtime alerts, and code risks.
Those descriptions identify evaluation leads. The matrix records actual source shapes, which can
differ from README examples (notably session configuration). No runtime exploit was attempted.
Run and curate the comparator using the [Bearer runbook](../../docs/samples/nodejs-goof/bearer-runbook.md).

## Evaluation and implementation contract

The initial documentation increment added target provenance, case inventory, a reproducible runbook
and curated result storage. The user subsequently authorized Bearer/current-zazzles comparison,
a scanner-independent answer key, and sequential XSS implementation with reports after each type.
The [XSS implementation contract](../specifications/xss-response.md) now owns implementation scope;
the [initial results](nodejs-goof-results.md) preserve pre-change observations. Documentation checks
and the 22-file snapshot verification passed; the valid Bearer run and current scanner outputs are
saved under `docs/samples/nodejs-goof`. Expected labels and hashes are in the
[answer key](../../docs/samples/nodejs-goof/xss-answer-key.md).

## Case matrix

Observations below refer only to the [initial pinned baseline](nodejs-goof-results.md). No
finding is not an explicit safety judgment. Zazzles acceptance describes required implementation,
not behavior achieved by the initial scan, whose CommonJS evidence analysis returned zero entries.

| Case / capability | Source-inspected target evidence | Observed Bearer | Zazzles acceptance / boundary |
|---|---|---|---|
| Reflected XSS, EJS attribute | `app.js:54` → `routes/index.js:80-85` passes `req.query.redirectPage` to `views/admin.ejs:17`, raw `<%- redirectPage %>` in a quoted `value` attribute | Not detected | XSS increment 1: route → render locals → template raw output and attribute context; no sanitizer established on this path |
| Stored Markdown HTML | `/create` saves `item` at `routes/index.js:168-188`; `/update/:id` writes `req.body.content` at `:234`; `/` reads todos at `:35-46` → `views/index.ejs:20` emits raw `marked(new String(todo.content))` | Not detected | Stored/template increments 1–3: persistence and render linkage; inspect `parse` on create separately from update; `marked` 0.3.5 with `sanitize: true` at `app.js:75` is not verified effective protection |
| Escaped EJS output | `views/edit.ejs:20,23` uses `<%= todo.content %>` in quoted attribute / HTML text contexts | No corresponding XSS finding; no protected assertion | Negative controls for context-appropriate escaped output; raw HTML elsewhere does not make every todo render an XSS sink |
| Handlebars profile rendering | `routes/index.js:104-120` passes the request object as render locals after selected field validation; `views/account.hbs` uses `{{firstname}}` and other escaped fields | No corresponding XSS or render-options finding | Separate escaped field rendering from request-controlled render options/layout; whole-object forwarding is not itself proof of raw-field XSS |
| Open redirect | `req.body.redirectPage` at `routes/index.js:54` → `adminLoginSuccess` → `res.redirect(redirectPage)` at `:74` | Not detected | Destination policy claim, distinct from the GET login template's XSS claim |
| NoSQL injection | `routes/index.js:52` forwards `req.body.username` / `password` into `User.find` fields | `javascript_express_nosql_injection` at :52 (also :185/:204/:232; unadjudicated) | Existing database analysis: establish Mongoose binding, accepted object shapes and operator interpretation; not an SQL string-construction case |
| Command execution | `routes/index.js:168-174` extracts a URL from request content and concatenates it into imported `child_process.exec` | Not detected | Command sink and shell-context analysis; regex matching alone does not establish shell escaping |
| Session configuration | Actual `app.js:42-46` sets a literal session secret and `cookie: { path: '/' }` | Secret rules at app.js:42/:43; cookie configuration at :45 | Cookie/session/header and hardcoded-value claims; verify dependency defaults rather than inherit the README's `secure`/`httpOnly` assertions |
| Logging fixture | `app.js:83-84` sends a token-shaped literal to `console.log` | No finding at token log :84; 21 logger records elsewhere | Sensitive-data logging: establish value and sink; label as an intentional training fixture |
| HTTP resource | `views/layout.ejs:7`, `views/layout.hbs:7` load a script over HTTP | No corresponding HTTP-resource finding | Insecure transport/resource claim, distinct from tainted HTML rendering |
| Test findings | `tests/authentication.component.spec.js` is included in the baseline scope | 5 hardcoded-secret records | Keep test/fixture locations separate in review and counts |
| Dependency vulnerabilities | Manifest pins e.g. `marked` 0.3.5, `st` 0.2.4, `mongoose` 4.2.4; README lists associated exploit demos | Not exercised as SCA | Package/advisory matching and reachability require a separate SCA contract; a SAST finding is not a dependency audit |
| Docker / runtime monitoring | README describes image scans and Snyk runtime reporting | Not exercised | Outside this static-source baseline and current XSS increments |

## Comparison rules

Keep source labels outside scanned inputs. The runbook snapshots application/templates/tests and
manifests while excluding exploit fixtures, documentation and prior SARIF. Templates are required
inputs even when the comparator cannot analyze their syntax. Record that limitation rather than
removing templates to make scopes appear equivalent. This snapshot differs from the earlier
JS/TS/JSON-only experiment; broad finding totals are not directly comparable.

Bearer absence is not a protected judgment. An old package's `sanitize: true` option is not proof of
effective sanitization, and package version alone is not proof of an exploitable render path.
Measure disagreements at the claim/source/sink/context level. Keep unresolved guards, sanitizer
behavior, template support and runtime reachability explicit. Independent labels, held-out accuracy
and runtime confirmation remain unmeasured.

The [XSS scope](../project/xss-detection-scope.md) owns delivery order. The
[sample README](../../docs/samples/nodejs-goof/README.md) owns saved baseline links and provenance.
