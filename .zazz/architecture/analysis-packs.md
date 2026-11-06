# Embedded analysis components

The [SQL-first contract](../specifications/analysis-pack-mvp.md) implements the
[analysis pack proposal](../project/analysis-pack-mvp.md) on merged common-entry/dispatch base
`f8b7dbdfe1521d9b914f271fa68b26b01fdec97b` (PR #8); the [NoSQL/validation contract](../specifications/components-01-nosql-validation.md)
made those families the second and third consumers, followed by the
[XSS/browser contract](../specifications/components-02-xss.md). This document owns component responsibilities,
source inventory, compatibility, migration status and acceptance evidence. The
[public analysis guide](../../docs/analysis/database-injection.md#query-component-selection) owns SQL
and NoSQL interpretation, and the [XSS guide](../../docs/analysis/xss.md) owns rendering
interpretation; [the implementation plan](../project/implementation-plan.md#analysis-pack-mvp-initiative)
owns future sequencing.

## Session 02 closeout and session 03 starting point

Runtime checkpoint `4e0be3e` on the existing `analysis-pack-mvp` feature branch includes the
session01 foundation (`7610193`), XSS component migration (`30fada9`) and authorized minor cleanup.
SQL, NoSQL, standalone validation and XSS are working component consumers. Their protocols are
`sql-components-v2`, `nosql-components-v1`, `validation-components-v1` and `xss-components-v1`;
report schema4, evidence-version3 and route-inventory schema1 are independent contracts.
Schemas1–4 still read/export with original historical answers; changed requests reject replay.

The post-push Juice Shop check used clean pinned revision
`f189f521266f42c6bd05884468292000cee852ed`, with all 27 reference hashes verified before/after.
Accepted pre-XSS `7610193` and current `30fada9` scans have exact canonical equality after removing
only the 23 new XSS selection records: 165 entries (152 HTTP, 13 browser), 208 uses, 8 construction
findings and 200 Review records. DOM search and cross-tier track-order remain the two XSS
construction findings. Stored/JWT/model-setter/collection/dialog origin gaps remain Review.
Both provider encoders construct all 23 rendering requests; localhost stubs verify the two
source-reviewed raw controls plus replay; all 20 JSON-only observations remain unassessed. Largest bodies are
62,738 bytes (OpenAI) and 62,108 (Jev), below 65,536. These are offline preservation checks, not
live-model accuracy or resolution of the known coverage misses. Raw report/verification artifacts:
`_test-findings/juice-shop/run-3290785101/`; accepted pre-XSS run: `run-3975911309`.
The small cleanup removes an unreachable XSS protocol fallback, a test-only production constant
and SQL-specific wording in the shared replay error. Full race, vet, build, formatting and docs
checks passed after cleanup; no source facts, protocols or successful judgments changed.

Next is [session03: outbound observations](../specifications/components-03-outbound.md).
It has not been implemented. Outbound still uses existing generic legacy assessment; explicit
observation-only `not_run`, zero model calls and Preview rejection are its planned eligibility
correction. Freeze that legacy behavior separately from source/report parity before changing it.
Preserve earlier-family effective requests and browser method/path/field provenance. Logging
migration and complete SSRF/exfiltration/browser-security assessment remain separate future work.
Resume the assigned feature tree at actual latest HEAD, including this documentation closeout,
and preserve uncommitted work. Commit/feature-push delivery was authorized by the user; no PR,
integration merge, new worktree/agent, paid call or target execution is implied for the next session.

## Boundary and ownership

`internal/analysispacks` owns scanner-authored component descriptors, dependency validation,
deterministic selection and bounded interpretation material. Components describe existing compiled
capabilities; they do not load or execute adapters. The package does not import parse, judge or report.
Shared serializable references live in `internal/analysis`; extraction remains in `internal/parse`,
transport and mandatory claim gates in `internal/judge`, and presentation in `internal/report`.

| Knowledge | Baseline owner | Component responsibility and remaining owner |
| --- | --- | --- |
| JS/TS syntax, Express, dynamic HTTP, GraphQL, unary gRPC, browser entry/flow | `internal/parse/` adapter and flow files | Reference supported compiled subsets and their limits; extraction remains in parse |
| Python Flask/Blueprint, RESTX, FastAPI, DB-API flow | `internal/parse/python_*.go` | Reference supported compiled subsets; no new Python framework/heap/async support |
| Language interpretation | `internal/guidelines/javascript.md`, `python.md`, `go.md` | Registered families select validated existing sections; unmigrated families retain full language documents; Go guidance does not establish Go support |
| Sequelize raw SQL | `internal/parse/checks.go` | Catalog owns the verified operation contract; parse still proves receiver/import/slot bindings |
| Python SQL drivers | `internal/parse/python_dbapi.go` | Catalog owns sqlite3, Psycopg 2/3, PyMySQL and MySQL Connector operation contracts; parse still proves connection/cursor provenance |
| SQL questions and choices | `internal/judge/client.go`, `injection.go` | SQL components own the selected question/answer material; judge encodes, decodes and applies Go policy |
| Required SQL evidence | `internal/parse/checks.go`, `python_sql.go`; `internal/judge/policy.go` | Components declare obligations; core input-source, sink-binding, invoked-callable and control-flow roles remain mandatory |
| NoSQL contracts | `internal/analysispacks/nosql.go` (catalog owner); parse proves collection/model bindings | Migrated: MarsDB collection and Mongoose query/casting contracts are catalog components; selector, `$where`, schema/version and all-path semantics preserved |
| XSS contracts/questions | `internal/analysispacks/xss.go` and `xss_selection.go`; validated guideline projections | Migrated: catalog owns primary rendering contracts, five questions and supplementary context/protection/storage/browser/response material; parse retains source facts and judge retains mandatory gates |
| Standalone validation | `internal/analysispacks/validation.go` (questions/choices); policy assembly stays in parse; gates in `internal/judge/validation.go` | Migrated: requirement/enforcement questions are catalog components; established/absent policies, enforced/bypassed/missing guards and API-less field claims preserved |
| Logging | Logger/source extraction in parse; `internal/judge/logging.go` | Inventoried, not migrated; retain sensitive content, emission and redaction gates |
| Outbound requests | Destination/source extraction in parse; Python HTTP-client API models; existing generic judge/report paths | Inventoried, not migrated; preserve current fetch/axios/Python destination facts; full SSRF assessment is a separate feature |
| Findings and exports | `internal/report/`, `cmd/scan/` | Existing presentation counts/statuses remain; canonical JSON carries detailed selection provenance |

A material item has one runtime owner. Family integration reuses selected shared language sections
and moves operation/question knowledge rather than creating alternative legacy policy paths.
Unmigrated families (logging, outbound) retain their existing behavior until their own migrations.

Catalog identities are `pack.<family>`, `claim.<family>`, `evidence.<family>`, `flow.local`,
`language.javascript`, `language.python`, `language.nosql.javascript`,
`language.validation.javascript`, `language.validation.python`, `entry.http`, `entry.graphql`,
`entry.grpc` and `api.<operation-contract-id>`. A normal SQL or NoSQL candidate selects seven
components; a standalone validation claim selects five. Language versions follow the canonical
JavaScript/TypeScript v5 and Python v3 guidance; the SQL projection keeps the accepted
`language.javascript`/`language.python` identities, and each family projection has distinct
component identities because its selected bytes differ. Operation IDs remain
`sequelize-raw-query-v1`, `python-<driver>-dbapi-v1` for the five existing Python drivers,
`marsdb-collection-v1`, `mongoose-model-query-v1` and the pinned `mongoose-model-query-v2`.
A language version identifies the underlying guidance; the digest also covers the actual selected text.

## Descriptor and selection contract

Each embedded component has a stable ID, explicit version, kind, domains, applicability,
compatibility, pinned dependencies, compiled/source owner, required roles, limits and test references.
Its SHA-256 covers content and metadata, including compiled-owner references, so changing bytes
without bumping a version still changes identity. Ordered dependency closure produces a stable selection digest. Pack identity is versioned
independently of a target application or dependency version.

Language, entry framework, verified library operation and claim compose independently. Selection
uses supported source facts and scanner configuration. Installed packages, helper names and target
comments do not enable trusted contracts or define policy. Unknown application library versions
retain the contract's bounded assumptions; a scanner contract version is not an installed version.

Static registered-family assembly records its selection and effective obligations on the candidate. Existing
`completeness.required` and source gaps are preserved; `analysis_components.required` records selected
obligations, and Go judgment validates their union with core/claim roles. Judgment
validates/consumes that same selection. Mandatory Go obligations are unioned with component roles;
component content cannot remove evidence gates, invent an edge or clear a blocker. Caller-specific
controls remain attached to invocation contexts, including shared callables exposed through several
entries. Protection cannot transfer from one caller, field, operand or branch to another.

Invalid embedded catalog configuration fails before model requests. Missing or ambiguous candidate
applicability retains the candidate and an explicit unresolved state, without generic fallback.
Selected-for-analysis provenance is distinct from a completed model assessment.

The model receives selected interpretation and the actual trace. Compiled-owner identifiers,
regression labels/test names and unrelated catalog entries are metadata, not model instructions.
No catalog-wide list or expected result is sent. Existing source deduplication and lossless compact
wire encoding are retained; neither removes controls or obligations.

## Bounds and compatibility

The selection closure is limited to 32 components and selected interpretation material to 16 KiB.
These bounds leave headroom for the small working SQL composition while preventing unrestricted
catalog expansion into a request. The existing 64 KiB final provider-body bound remains authoritative
for both provider encoders; source and selected material count toward it. Requests may use existing
lossless compaction. Excess material or a still-oversize body remains unresolved with no HTTP call;
required content is never truncated. Measurements are recorded under verification below.

Fresh analysis JSON uses schema 4. SQL requests use `sql-components-v2` (the NoSQL/validation
generalization changed component compatibility metadata and therefore SQL provenance bytes; guidance,
questions, choices and outcomes are unchanged), NoSQL uses `nosql-components-v1` and standalone
validation `validation-components-v1`. The accepted foundation used schema 3 with
`candidate-injection-v1` SQL/NoSQL requests and `validation-use-v2` validation requests; saved
reports keep those answers but cannot replay onto changed requests. XSS uses `xss-components-v1`,
logging `log-use-v1`, and generic candidates `candidate-use-v10`. Evidence-version 3 and
route-inventory schema 1 are independent contracts.

The candidate's `analysis_components` records pack/component IDs, versions, digests, effective
required roles, limits and optional `failure` coverage provenance, separately from original source
gaps. `assessment.request_digest` records the effective request for assessed registered-family claims. Decision traces carry
`analysis_components` and `request_digest` alongside the exact request/response bodies. Exact effective provider-body matching and protocol checks still gate Replay. Changing
selected content/metadata, evidence, instructions, questions/choices or model rejects an old response,
even if someone forgot to bump a material version. No automatic persistent decision cache is added.

Saved schema-1/2/3 reports remain exportable through the existing compatibility boundary. Export
preserves original answers and source-schema provenance; it does not invent a component selection or
claim that a saved result was freshly judged. JSON is canonical provenance. HTML/SARIF retain existing
issue counts and status semantics, including a construction-practice concern beside protected SQL.

## Migration ledger

| Family | Status | Next bounded task and acceptance controls |
| --- | --- | --- |
| SQL injection with query validation/construction facets | Migrated / offline verification passed | Catalog/selection, request composition, provenance and actual JS/TS/Python pipeline parity; raw, bound, constrained, unknown and context-isolation controls; `sql-components-v2` provenance after family generalization with frozen discovery hashes and outcomes |
| Command injection | Separate next new family | Add verified process/shell contracts and command-specific questions through accepted components; positive/protected/unknown controls and unchanged other families |
| XSS | Migrated / offline verification passed | Selected rendering/context/protection/storage/browser/response material; five unchanged named questions, mandatory origin/context/path gates, JSON-only observations, exact replay and separate bounded fetch-method correction |
| Standalone validation | Migrated / offline verification passed | Requirement/enforcement questions and field-policy guidance selected per claim; established/absent policies, enforced/bypassed/missing guards, client-only policy, unknown helpers, API-less claims, JS/Python parity and claim isolation; no invented standards |
| Logging | Inventoried / not migrated | Move sensitive-content/emission/redaction contracts and questions; paired leaks, metadata/constants, redaction and unknowns |
| JS/TS NoSQL | Migrated / offline verification passed | MarsDB/Mongoose contracts, selector/`$where` guidance and all-path/coverage questions selected per claim; raw selector, constrained scalar, conditional alternatives, Mongoose schema/version, unknown helper/API, shared-caller isolation and replay controls; Python NoSQL stays unsupported |
| Outbound requests | Spec prepared, session 03 / not migrated | Move existing destination-contract/guidance material; retain destination/body/response separation, binding/source gaps and current outcomes; do not imply complete SSRF assessment |

The user revised session grouping on 2026-11-05. Implementation order is
[01 NoSQL plus standalone validation](../specifications/components-01-nosql-validation.md),
[02 all XSS types including browser](../specifications/components-02-xss.md), then
[03 outbound observations](../specifications/components-03-outbound.md). Logging remains separate.
Sessions 01 and 02 are implemented; session 03 remains an unassigned implementation contract.
XSS uses coordinated server/template, stored and browser/cross-tier checkpoints. The plan tracks
actual future sequencing; command injection remains the next new vulnerability family.

### Extension risks recorded by source review

- Resolved by session 01: catalog compatibility, required roots, `Select`/`Attach`, material
  requirements and request-digest handling were SQL-specific. Family registrations now own candidate
  eligibility, question protocols, roots, roles, cardinality and policy owners, and required
  questions are validated by registered names and choices, not a universal count.
- Resolved by session 01: guidance projection relied on exact Markdown headings and silent literal
  phrase removal. Projections now declare required named sections; duplicate, renamed or missing
  sections and missing required phrase edits fail before scanning, and the SQL projection's effective
  bytes are frozen by test.
- Resolved by session 02: browser fetch uses actual literal methods/default GET; unresolved options,
  methods and Request objects keep a gap and cannot select a guessed server producer. The separate
  before/after table below records the only authorized detection correction. No runtime exploit or
  corpus-wide accuracy claim is made.
- [JS/TS outbound discovery](../../internal/parse/flow.go) uses fetch/axios-shaped names and retains
  unverified-binding/response gaps. It combines argument dependencies without the same destination
  precision as the existing Python adapter. Session 03 must retain these limits, never confer a trusted
  API contract by name, and separate URL/body/response and browser/server evidence.
- Generic SQL construction questions do not establish outbound destination security. Session 03
  makes observation-only model eligibility explicit; full SSRF requires its own later questions,
  evidence gates and evaluation. Client-side requests must not be labeled server SSRF.

These are maintenance/extension constraints and bounded correctness concerns, not evidence that the
SQL refactor regressed. Existing offline checks preserve the measured SQL baseline; they do not
establish general live accuracy or coverage of every browser vulnerability.

Scan, Inspect, Audit and Fix remain independent product goals. Existing validation and construction
practice can share components without establishing a general quality engine. A team standard requires
explicit trusted selection and claim-specific evidence. This repository's `.zazz/` development rules
are not automatic policy for targets. General standards ingestion, downloaded packs, plugins,
independent Findings Review and automated remediation remain outside this MVP.

## Verification and limits

Verification uses existing source fixtures and local provider stubs. It compares discovery, bindings,
source spans/hashes, invocation contexts, blockers, coverage, raw answers and final report outcomes
against the accepted foundation. Only explicitly versioned schema/provenance/envelope changes are
allowed. Existing unmigrated-family controls remain regression coverage. Expected labels remain in
test assertions/stubs, outside judgment input.

Offline tests demonstrate composition, transport, policy and reporting behavior. They do not measure
live-model accuracy, general precision/recall, custom-validator generalization or paid-call savings.
No live requests, target execution or sample-application edits are needed for this migration.

### SQL migration verification, 2026-11-05

On native Go 1.27.3 darwin/arm64 (module minimum 1.27.2), uncached scoped tests pass for parse,
judge, analysis, report and CLI. SQL-specific evidence includes:

- `TestSQLDiscoveryBaseline`: frozen accepted-foundation hashes across nine fixture sets, allowing
  only schema 4 and `analysis_components` fields to differ; source facts and existing gaps are preserved.
- `TestCatalogValidation`, `TestSelectionDeterminism`, `TestSQLApplicability` and
  `TestSQLMaterialAndSelectionLimits`: closure, dependency/version/applicability conflicts, byte changes
  without version bumps, additive obligations and explicit limits.
- `TestSQLComponentPayload`: selected SQL assessment, decoding and exact replay through both OpenAI
  and Jev localhost transports; unrelated guidance/catalog/test metadata excluded.
- `TestSQLComponentPipelineParity`: six SQL controls each in JavaScript/TypeScript GraphQL/gRPC,
  expected per-entry outcomes and JSON/HTML/SARIF count/status/source/raw-assessment parity.
- `TestSQLPythonComponentPipeline`: five drivers × raw, bound, dynamic-identifier and unresolved
  transformation controls (20 cases), with selected contracts, policy/replay and cross-format checks.
- `TestSQLHTTPComponentAnswers`, `TestSQLSelectionFailurePipeline`, `TestSQLComponentLimits` and
  `TestSQLComponentReplay`: independent facets, unknown/not-applicable answers, zero-call failures,
  oversize rejection and stale-request controls. Existing injection/context/hex/conditional tests remain.
- `TestComponentReportCompatibility` plus existing schema-1/2 tests: old answers/protocols stay
  historical; schema-3 export does not fabricate component or request provenance.

Measured selection material is 6,074 bytes for JS/TS and 5,293 bytes for Python sqlite3; each uses seven
components. Across 35 valid SQL fixture candidates, final provider bodies range from 22,668 to 48,318
bytes. The same baseline fixture set's largest body was 48,138 bytes: the observed maximum increases
by 180 bytes, remaining below 65,536. Material stays well below the 16,384-byte cap; no content is
truncated to fit. These fixture measurements are not whole-repository payload guarantees.

`BenchmarkSQLComponentSelection` on Apple M6 measured 2,493 ns/op, 2,765 B/op and 20 allocations/op.
The old path had no component selection, so this is added local selection work, not a scanner-wide
speedup. Timing is host-specific evidence, not a correctness threshold or live-model cost estimate.

Final independent checks after the required-root configuration fix all pass:

```bash
GOTOOLCHAIN=local GOPROXY=off GOFLAGS='-mod=readonly -buildvcs=false' \
  GOCACHE="$PWD/.zazz/ephemeral/go-cache" go test -race ./... -count=1
GOTOOLCHAIN=local GOPROXY=off GOFLAGS='-mod=readonly -buildvcs=false' \
  GOCACHE="$PWD/.zazz/ephemeral/go-cache" go vet ./...
GOTOOLCHAIN=local GOPROXY=off GOFLAGS='-mod=readonly -buildvcs=false' \
  GOCACHE="$PWD/.zazz/ephemeral/go-cache" go build ./...
```

All nine packages pass the race suite. Changed/new Go files pass the format gate; the diff passes
whitespace/conflict checks, and new/changed documentation links resolve. The obsolete parse-owned
Python DB-API contract file was removed after moving its sole owner to the catalog. No paid calls,
sample edits, feature commit/push or integration merge were performed as part of this refactor.
The user-requested local branch synchronization was a fast-forward to the already merged foundation.

### NoSQL and standalone validation migration verification, 2026-11-05

Second/third consumer migration under the [combined contract](../specifications/components-01-nosql-validation.md),
implemented on the accepted (uncommitted) SQL tree recorded in the session run log. On native Go
1.27.3 darwin/arm64, uncached scoped tests and the full race/vet/build suite pass. Evidence:

- `TestFamilyComponentSelection`, `TestFamilyCatalogValidation`, `TestFamilyGuidanceCompleteness`:
  independent family registrations; missing family roots, unknown protocols, drifted claim questions,
  contradictory language applicability and renamed/missing/duplicate guidance sections fail before
  scanning; API-less validation selects; Python NoSQL and unregistered kinds gain no selection.
- `TestNoSQLComponentPipelineParity`: raw selector, constrained scalar (protected and
  optimistic-vulnerable), `$where` numeric-only vs raw-sibling alternatives, Mongoose pinned
  4.2.4 schema/version (`mongoose-model-query-v2`), unknown helper, fake collection API and shared
  caller with per-invocation guard isolation; per-route outcomes, blockers, replay and
  JSON/HTML/SARIF parity through actual parsed fixtures.
- `TestValidationComponentPipelineParity`: free text, absent policy and ignored type-only guards
  (JS), required-only enforced, source-backed format unenforced, client-only policy and unknown
  helpers (Python); optimistic answers cannot bypass missing field-policy roles; no API contract
  in any validation payload.
- `TestNoSQLValidationClaimIsolation`: one Python field simultaneously SQL-protected and
  policy-invalid; one JS selector field injection-protected beside an unresolved standalone record;
  no duplicate issue counts across JSON/HTML/SARIF.
- `TestFamilyComponentPayloads`, `TestFamilyComponentFailures`, `TestFamilyComponentReplay`: both
  provider encoders carry only selected material with exact question sets and provenance; zero HTTP
  calls on selection failures; stale-request/protocol/evidence mutations reject replay.
- `TestComponentReportCompatibility` extended: schema-3 exports keep historical
  `candidate-injection-v1`/`validation-use-v2` answers and absent provenance (including display);
  schema-4 exports keep `pack.nosql`/`pack.validation` provenance and request digests.
- Existing NoSQL/validation regressions (`TestNoSQLAllPathProtectionGate`,
  `TestMongooseQueryProtectionControls`, `TestConditionalQueryOutputConstraints`,
  `TestAutomaticCollectionProvenance`, `TestValidationPolicyRequiresIndependentEvidence`,
  `TestValidationRespectsAdditionalRequiredRoles`, `TestPythonValidationPolicyClaims`,
  `TestPythonFormPolicyIsNotEnforcement`) and the frozen SQL discovery baseline all pass unchanged;
  `string-flow` goldens are byte-identical to the accepted SQL tree.

Measured on fixtures: NoSQL selects 7 components / 11,092 bytes of material with final bodies of
28,755 (jev) and 29,451 (openai) bytes; validation selects 5 components / 4,358 (JS) and 3,646
(Python) bytes with bodies of 16,031 and 16,265 bytes. All remain inside the 32-component,
16 KiB material and 64 KiB request bounds with no truncation. SQL bodies stay in the recorded
22,668–48,318 byte range. `BenchmarkSQLComponentSelection` measures 2,767 ns/op, 4,200 B/op and
34 allocations/op on Apple M6 (accepted SQL baseline: 2,493 ns/op, 2,765 B/op, 20 allocs/op);
host-specific timing evidence, not a correctness threshold.

No paid calls, sample edits, target execution, feature commit/push or integration merge were
performed. JS standalone field records now use the field-policy protocol instead of the generic
practice pair whose sink-protection question they could not answer; their outcomes remain
Review/needs_evidence, now for the honest missing-policy reason.

### XSS migration verification, 2026-11-05

Accepted predecessor: clean `7610193`, containing session-01 implementation `9c83a07` and its
cleanup. The ignored run log records the complete tracked-blob manifest and passing baseline.
One agent used the existing `analysis-pack-mvp` worktree; no predecessor work was discarded.

`pack.xss`, `claim.xss` and `evidence.xss` register `xss_candidate`. The catalog owns the three
existing Express, browser and Python rendering contracts and the unchanged five named questions
and choices. `render.*` is a singleton primary rendering role; `context.xss`, `protection.xss`,
`storage.xss`, `browser.xss` and `response.xss` are supplementary interpretation roles selected by
actual rendering/type facts. `flow.xss` references existing compiled extraction without changing
SQL/NoSQL/validation descriptors or requests. Template sinks retain their actual spans/files;
rendering contracts identify their source adapter for language selection. The browser contract is
rendering interpretation, not an invented server API. Data-only JSON records acquire no selection
and make no decision request. Existing compiled encoder/sanitizer binding/version/options logic,
independent writing paths and source/coverage limits remain authoritative.

The staged server/template, stored, browser/framework/event/inline-script and cross-tier checks pass:

- `TestXSSComponentSelection`, `TestXSSMaterialBoundsAndConflicts`, `TestXSSComponentInventory`:
  browser-only selection, required roots, incompatible primary contracts, stale bytes without a
  version bump, bounds and actual Python/browser fixture inventory.
- `TestXSSServerComponentParity`, `TestXSSStoredComponentParity`, `TestXSSBrowserComponentParity`,
  `TestXSSCrossTierComponentParity`: parsed source through both OpenAI/Jev localhost transports
  to JSON/HTML/SARIF, with independent expected raw/protected/unknown results; EJS/Handlebars/Pug,
  Jinja, Mongoose/Sequelize, DOM/jQuery, React/Vue/Angular and Python inline scripts.
- Existing response/template/compile/sanitizer, stored sibling-path, framework, event/limit and
  Python-inline regression matrices pass. Frozen SQL discovery projections and non-XSS suites pass.
- `TestXSSComponentPayloadsAndPolicy`, `TestXSSDataOnlyNotAssessed`, `TestXSSComponentFailures`,
  `TestXSSComponentReplay`: correct selected questions, no unrelated policy/test labels, optimistic
  answer rejection, required-role union, unresolved origins, zero-call selection/body failures,
  exact replay and changed method/evidence/material/model/protocol rejection.
- `TestXSSComponentReportParity`: six real event registrations, four independent use contexts,
  two raw, one protected and one unknown outcome; JSON/SARIF keep zero-use coverage and HTML
  keeps browser grouping/source locations. `TestXSSComponentHistory` retains original XSS answers
  and versions through schema1–4 read/export; fresh reports remain schema4.

The fetch exception is separate from organizational parity. The pre-fix GET-hardcoded control was
run against the same fixtures and failed; the corrected implementation passes. GET supplies raw
`req.query.get`, while POST supplies independently constrained `req.body.post` at the same `/item`.

| Browser request | Before | After |
| --- | --- | --- |
| Absent/undefined options, `{}`, explicit GET | GET origin; raw | Same GET origin/protection |
| Literal lower-case `post` | Wrong GET origin; raw | POST origin only; sanitized on its own path |
| Dynamic method/options or alias | Guessed GET origin; raw | Empty unresolved method; explicit method gap; no linked producer |
| Spread, duplicate or computed property | Guessed GET origin; raw | Method gap; no linked producer |
| Request object | Unsupported endpoint; no demonstrated origin | Explicit unresolved method; no guessed producer; rendering record retained |

`TestBrowserFetchMethodIsolation`, `TestBrowserFetchRequestObjectRetainsGap`, `TestPythonInlineFetchMethodGap` and existing
`TestCrossTierResponseFieldLinks` cover literal/default methods, wrong fields, ambiguous routes and
unresolved producers. New origin/method evidence cannot reuse a saved request. No other detector
behavior correction, framework/flow expansion or vulnerability family was introduced.

Measured selected material: reflected JS 8 components / 8,115 bytes; stored Python 9 / 7,394 bytes;
cross-tier JS/TS 10 / 8,836 bytes, below 32 / 16,384. The 88 measured pipeline/inventory provider bodies
range from 16,961 to 38,293 bytes, below 65,536; source evidence is never truncated to fit. Local selection benchmark:
6,437 B/op, 44 allocations/op (Apple M6, Go1.27.3); no timing threshold or scanner-wide speedup claim.
Test output lives outside the worktree under `_test-findings/components-xss/20261105-offline/`.

Full contributor checks pass: `go test -race ./... -count=1` (nine packages), `go vet ./...`,
`go build ./...`, `gofmt -l .` (empty), whitespace/conflict and scoped manual documentation/link
checks. X1–X8 are covered; no unmet criterion remains in this contract. Subsequent user-authorized
feature commits/pushes and minor cleanup are recorded in the checkpoint above. No paid calls,
sample edits, target execution, new worktrees/agents, PR creation or integration merge were performed. Live-model accuracy,
broader browser security and general Python API-to-browser linkage remain outside this contract.
