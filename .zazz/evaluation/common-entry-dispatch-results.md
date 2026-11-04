# Common entry and dispatch implementation results

## D6 — bounded Apollo GraphQL resolvers, 2026-11-04

Implemented against D5 `7fa55e4d8e4055532993eac4073f9d2c25d59b80` in the existing feature
worktree. The [D6 contract](../specifications/common-entry-dispatch-06-graphql.md) owns acceptance;
D7 has not started. No paid calls, target application execution/modification, new worktree or PR.
The future GraphQL application remains deferred evaluation.

| Acceptance | Passing source-backed evidence |
| --- | --- |
| AC1: Apollo/schema/resolver binding | `internal/parse` TestGraphQLResolverBindings and variants: literal SDL, imported/local/inline/method resolvers, independent schema instances, object fields, unresolved/default/composed maps |
| AC2: independent inputs and SQL transfers | TestGraphQLInputRoles and transfer controls: args to verified Sequelize SQL; parent/context/info unknown; scalar names establish no business policy or protection; calls alone do not establish flow |
| AC3: real completeness and policy | parse/judge TestGraphQLProtectionContexts: raw/bound/unknown SQL, guarded/open shared sink, unknown applicable guard, contextual controls, local Decisions HTTP stub, exact replay and changed-coordinate rejection |
| AC4: actual exports, selectors and bounds | parse/report/CLI TestGraphQLPipelineAndBounds variants: real JS and typed TS fixtures, schema3 JSON, resolver-specific HTML/SARIF, zero-use gaps, exact source hashes/spans, typed selectors, root/symlink escapes, shadowed constructors and deterministic budgets |
| Regression/quality | Scoped suites, full `go test -race ./... -count=1`, `go vet ./...`, `go build ./...`, changed-file formatting and diff review |

Each pinned source-only JS/TS fixture yields nine entries, six SQL candidates and one separate
argument-validation observation. Offline reports contain five construction-practice findings and
two Review observations, with every assessment `not_run`. Under the local stub, three SQL contexts
are confirmed, bound and guarded SQL injection facets are Protected, and the unknown transform
remains evidence-limited despite a protected stub answer. Construction practice keeps the guarded
string-built query orange; only the bound record becomes a final green protection record. A separate
unknown-guard fixture also rejects the stub's protected answer. These tests establish policy and
transport behavior, not model accuracy or runtime exploitability.

A scanner built from an archive of the exact D5 commit produced identical JSON evidence to D6 for
six regression sets: JS/TS dynamic dispatch (four entries/four uses each), Python dynamic VulnPy
(38/31), existing string flow (one/one), and JS/TS browser events (six/four each). No historical report,
raw answer, golden or target source changed. Saved GraphQL exports preserve input hashes and retain
JSON/SARIF statuses, invocation identities, actual sink locations and related schema locations.

The [public subset and limits](../../docs/running-scans.md#bounded-apollo-graphql-analysis) describe
literal named Apollo wiring, supported SDL syntax, independent unknown parent/context/info roles,
and explicit unsupported composition/default/custom-scalar/wrapper/parent-result behavior. The
adapter assumes Apollo Server 4 compatibility; installed versions and transport exposure are not
verified. Reached source/SDL/registration/record/helper limits remain visible gaps, never protection.

Local ignored evidence is under `.zazz/ephemeral/`: the append-only
`common-entry-dispatch-06-graphql-run-log.md`, `d6-{scoped,race,vet,build}.log`,
`d6-report-manifest.json`, `d6-regressions/`, and `d6-reports/`. Report verification covers generated
HTML source, disclosures/sort/filter markup and cross-format evidence/status parity. Browser policy
blocked the local-file visual preview; no visual browser QA is claimed here. Broader milestone visual
acceptance remains separate from these D6 adapter/export tests.

## D6 real-sample follow-up — 2026-11-04

At the user's request, fresh full source scans ran against all seven pinned application checkouts,
followed by a source-evidence audit, saved exports and exact historical decision replay. Every target
matched its recorded revision, reference hashes and clean Git state, including ignored files, before
and after scanning. No target was executed, modified or installed into; no new paid request was made.

| Fresh full scan | Entries | Uses | Offline finding / Review / Protected records |
| --- | ---: | ---: | --- |
| Juice Shop | 167 | 208 | 8 / 200 / 0 |
| Node.js Goof | 19 | 33 | 2 / 29 / 2 |
| Vulnerable TypeScript app | 19 | 20 | 1 / 19 / 0 |
| ThreatByte | 24 | 56 | 5 / 40 / 11 |
| Python insecure app | 1 | 2 | 1 / 1 / 0 |
| Vulnerable Flask App | 19 | 24 | 4 / 20 / 0 |
| VulnPy | 564 | 827 | 28 / 799 / 0 |

These fresh runs contain no model assessments. Their finding counts describe established
construction-practice observations; scoped protection can follow supported structural encoding facts.
They are not fresh vulnerability confirmations or proof of complete coverage.

The initial scan output matched the exact D5 source build for all seven apps. A stricter byte-span
audit then found an inherited logger-evidence defect in Goof and the TypeScript app: verified global
console bindings recorded descriptive text `global console.log` at a span containing `console.log`.
The fix uses the standard source-evidence constructor, preserving original text/hash/span and masked
model text. A failing-then-passing parsed-source regression covers five console methods/forms,
including an inline comment. Final reports retain identical entries, contexts, operations and
classifications; five apps remain completely identical to D5, while the other two differ only in
these two corrected witnesses and their two evidence-derived candidate/finding IDs. Old judgments
were not copied onto the changed evidence.

The final audit verifies 47,937 source-location/hash/line witnesses, including 46,035 byte-exact raw
excerpts. The remainder comprises 1,882 labeled field-projection annotations and 20 excerpt-free
browser scope markers; their original source hashes/spans are verified without relabeling them as
raw source text. JSON/HTML/SARIF entry/context/status/sink-location parity, zero-use entry visibility,
and identical generated/saved-export output pass. All seven curated saved JSON inputs remain
unchanged after offline export.

Current-source exact-request replay accepts all 827 VulnPy and 188 Juice Shop saved D5 decisions,
with zero new calls and complete assessed-report parity. VulnPy retains 28 contextual confirmed
findings and 799 evidence-limited results; the twelve main-app generated route/method contexts all
retain their source-backed SQL/raw-HTML findings. Juice Shop retains seven confirmed findings,
three construction-practice findings, 178 evidence-limited results and twenty unassessable API-data
observations. These are historical decisions applied only to identical current requests.

Original Python saved-report acceptance remains 29/29. Fresh source evidence passes all 22 structural
case checks; seven ThreatByte validation cases remain unassessed under the no-paid-call restriction,
with their policy candidates retained. The old VulnPy baseline intentionally still evaluates 4/10,
including its six documented positive misses; neither labels nor historical results were rewritten.
The separate Jinja coverage gap and other documented framework/flow limits remain. Node source
controls preserve Goof's raw/stored create+update and escaped/unknown outputs, TypeScript's raw
search output with no fabricated SQL sink, and Juice Shop's browser/cross-tier origins. Application
level GraphQL evaluation still awaits the future sample app.

Verification after the repair: focused logging tests, full `go test -race ./... -count=1`, vet/build,
formatting and diff review pass. Local ignored evidence:
`d6-real-sample-manifest.json`, `d6-real-sample-audit.json`, `d6-saved-python-evaluation.json`,
`d6-real-samples/`, `d6-real-logger-{red,green}.log`, and `d6-sample-fix-{race,vet,build}.log`.
The append-only D6 run log records commands, expected deltas, source hashes and publication checks.
No additional scanner regressions were found in this bounded sample verification; known unsupported
paths remain limits rather than protection claims. D7 has not started.

## D5 — bounded dynamic HTTP registration, 2026-11-04

Implemented in `common-entry-and-dispatch` against predecessor `fe11151bcb0a11102592a805e7d160fa57c303a9`.
The [D5 contract](../specifications/common-entry-dispatch-05-dynamic-http.md) owns scope and acceptance.
At the D5 checkpoint, D6 and additional framework/vulnerability families had not started.

Finite Express literal tables/for-of loops and direct local literal-argument factories now produce
source-backed registrations and separate closure/control contexts. Bounded Flask/FastAPI generation
resolves source module/member inventories, supported name transformations, local registry selection,
factory captures and independent application mounts without importing or executing target code.
The value trace retains request input, captured callbacks, actual callback arguments, helper-returned
SQLite connections and cursor/SQL operand provenance. Unsupported alternatives remain explicit gaps.

| Acceptance | Passing evidence |
| --- | --- |
| AC1: finite registration/provenance | `internal/parse` TestDynamicHTTPRegistrations and mounted variant; JavaScript/TypeScript parsed fixtures |
| AC2: context isolation | TestDynamicHTTPContextIsolation variants, Python repeated-path regression, raw/encoded/unknown controls, report assertions |
| AC3: bounded/conservative behavior | TestDynamicHTTPBounds variants: mutation, shadowing, unresolved/decorated callbacks, wrappers, recursion, limits and deterministic repeat scans |
| AC4: actual Python route-to-use traces | TestDynamicPythonVulnPyTraces in parse/judge, witness assertions, full/scoped pinned checkout scans and source/format audit |
| Real Decision-mode acceptance | 827 VulnPy and 188 Juice Shop actual completed requests; final exact-request replay of all 1,015 responses, no additional calls |
| Regression/quality | `go test -race ./... -count=1`, `go vet ./...`, `go build ./...`, scoped parse/report/CLI suite and `git diff --check` |

## Actual source-only sample results

Targets remained clean, including ignored files, at their pinned revisions before and after scans:
VulnPy `b81885e58283a3449b578d20c20557aff637e5ea`; Juice Shop
`f189f521266f42c6bd05884468292000cee852ed`. All VulnPy reference hashes were checked.
No target application was executed, installed into, modified or remediated.

| Offline scan | Entries | Candidate uses | Classification |
| --- | ---: | ---: | --- |
| Predecessor VulnPy Flask registration file | 0 | 0 | Explicit unsupported registration warnings |
| Predecessor VulnPy FastAPI registration file | 0 | 0 | Explicit unsupported registration warnings |
| D5 full VulnPy | 564 | 827 | 28 poor practice; 799 Review |
| D5 scoped Flask registration file | 474 | 756 | 24 poor practice; 732 Review |
| D5 scoped FastAPI registration file | 79 | 65 | 4 poor practice; 61 Review |
| D5 full Juice Shop | 167 | 208 | 8 poor practice; 200 Review; predecessor inventory/evidence/status parity preserved |

The original full VulnPy baseline's 11 static entries and six uses remain present. Its historical
six-positive-miss/four-warning evaluation, saved reports, answer keys and hashes are unchanged.
Resolved fresh registrations no longer carry obsolete blanket unsupported-registration warnings.
JSON/HTML/SARIF status, entry/context and primary sink-location parity was checked across all fresh
baseline/current/control/assessed bundles. The 12 required main-app positive contexts have 704
verified original source/hash/span evidence witnesses.

## Live Decision-mode findings and limits

User-authorized evaluations used OpenAI Decisions `gpt-6-luna`, with no retry or provider/model
fallback. All 1,015 actual requests completed. These are static model assessments checked against
source, not runtime exploit validation, a complete vulnerability inventory or an accuracy percentage.

**VulnPy:** 827 completed assessments produced 28 contextual confirmed vulnerabilities (21 SQL and
seven raw XSS), plus 799 evidence-limited results. All 12 required main Flask/FastAPI method/route
contexts were confirmed. Independent Flask application mounts explain the additional exposures:
there are five source cases, not 28 distinct defects. See the
[scan-plus-code-inspection table](../../docs/samples/vulnpy/reference.md#findings-confirmed-by-scan-and-code-inspection--d5).
The expectations derive from operations and tests; upstream's reviewed setup documentation does not
provide a formal per-route expected-findings manifest. The separate Jinja Environment/template
binding remains a known source-positive coverage gap. Observations lacking business-validation
policy remain evidence-limited rather than becoming invented policy violations.

**Juice Shop:** 188 model-eligible assessments completed across 208 uses. Twenty API JSON XSS-data
candidates correctly remained unassessed because no rendering claim was linked. Final findings:
seven confirmed vulnerabilities, three poor practice, 178 evidence-limited and 20 Review.
Source-inspected positives include product-search SQL (`routes/search.ts:23`), order/review NoSQL
(`trackOrder.ts:18`, `showProductReviews.ts:36`, `updateProductReviews.ts:17`), two sensitive request-line
logging contexts (`server.ts:366`) and cross-tier raw HTML (`track-result.component.html:9`, reviewed
JS-XSS-02). Other documented expected paths remain coverage/judgment limits. Login's SQL construction
remains poor practice because unmodeled password-hash transfers block a complete injection claim;
a vulnerable model answer does not override that blocker. The browser-search trust bypass retains
poor-practice status with exposure unknown.

Juice Shop's 14-row `autoModels`/Finale loop is an existing specialized persistence reference.
Four check-driven User-persistence contexts retain identical row/loop mapping, model constraints
and required evidence roles. D5 does not claim newly enumerated CRUD endpoints for every table row;
existing zero-argument login/search factories are regression references. Focused parsed fixtures
establish the new generic finite-table/literal-argument capability.

## Issues resolved and hardening

Overlapping enclosing source and repeated metadata, rather than uncontrolled route recursion,
caused oversized model requests. Lossless excerpt sharing plus an explicit compact evidence/gap
column schema retain roles, ordering, source spans/hashes, controls and references within the existing
64 KiB bound for both apps (preflight maxima: 47,517 bytes VulnPy; 65,245 Juice Shop). No record is
silently truncated; a remaining oversized request stays unassessed. Evidence envelope version is 6.

AST-derived model projections mask comments and genuine Python docstrings while local audit text
and hashes remain original. Executable f-strings, bytes literals, runtime strings and template code
are preserved. Regression tests cover reconstruction, comment filtering, executable expressions,
request size, replay matching and changed-closure rejection. Final replay matches every saved live
request exactly, including after the executable-string correction.

Additional fixes preserve identity for repeated-path factories with different captures, reject
unknown decorated callbacks, enforce mount-projection budgets, remove stale resolved-registration
warnings and give raw Python HTML responses encoding guidance. Existing golden changes affect
only the capability-limitation sentence. Unknown/mutated bindings cannot acquire fabricated protection.
Browser QA covered the generated route/sink/provenance, disclosures, keyboard sorting/focus, protected
filtering, visible zero-use gaps and a 390-pixel layout without horizontal overflow.

## Reproducible verification and local evidence

Use the declared Go toolchain; offline tests require no provider credentials:

```sh
GOTOOLCHAIN=local GOPROXY=off GOCACHE="$PWD/.zazz/ephemeral/go-cache" go test ./internal/parse ./internal/report ./cmd/scan
GOTOOLCHAIN=local GOPROXY=off GOCACHE="$PWD/.zazz/ephemeral/go-cache" go test -race ./... -count=1
GOTOOLCHAIN=local GOPROXY=off GOCACHE="$PWD/.zazz/ephemeral/go-cache" go vet ./...
GOTOOLCHAIN=local GOPROXY=off GOCACHE="$PWD/.zazz/ephemeral/go-cache" go build ./...
git diff --check
```

The append-only execution log and large/raw artifacts are local ignored evidence, not published
sample baselines. Paths below are relative to the assigned worktree's `.zazz/ephemeral/`:

- `common-entry-dispatch-05-dynamic-http-run-log.md`: commands, RED/fix evidence, issues and final audit.
- `d5-publish-{race,vet,build}.log`, `d5-final-projection-scoped.log`: final code checks.
- `d5-live/{vulnpy-01,juice-shop-01}/`: immutable original reports, exact request/raw-response traces,
  usage, run summaries and before/after target-integrity checks.
- `d5-publish-replayed/{vulnpy,juice-shop}/replay-summary.json`: final-source bundles with 827/188
  exact offline replays and zero new requests; 20 unassessable API records preserved.
- `d5-final-sample-manifest.json`: fresh full/scoped/baseline/control bundles, commands and output hashes;
  `d5-publish-offline.log`: counts and export checks, including saved JSON immutability.
- `d5-completion-audit.json`: inventory/evidence/format parity, 12-context matrix and source witnesses;
  `d5-node-persistence-parity.json`: actual Finale persistence parity.
- `d5-publish-scope-manifest.json`: exact implementation/documentation scope and SHA256 hashes.
- `d5-visual/`: generated-route and unresolved-control desktop/narrow screenshot evidence.
