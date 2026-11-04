# Common entry and dispatch implementation results

## D5 — bounded dynamic HTTP registration, 2026-11-04

Implemented in `common-entry-and-dispatch` against predecessor `fe11151bcb0a11102592a805e7d160fa57c303a9`.
The [D5 contract](../specifications/common-entry-dispatch-05-dynamic-http.md) owns scope and acceptance.
D6 and additional framework/vulnerability families have not started.

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
