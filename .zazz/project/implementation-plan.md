# Implementation Plan

How we get from an empty repository to the final solution, in phases that each end in something
runnable and verifiable. Design context is in [`../architecture.md`](../architecture.md).

**Status: draft.** Phase boundaries are settled; the detail inside each phase is indicative until the
phase starts.

## Acceptance targets

Two real TypeScript/Node.js applications, plus a hand-built labeled corpus.

| Target | What it exercises | Why it is here |
|---|---|---|
| TS microservice over **MongoDB** | routes (sources), Mongo CRUD (`nosql` sinks), outbound API calls (`ssrf` sinks), config | NoSQL has no parameterization escape hatch, so it should push findings onto Track B |
| TS app over **SQL** | the same, with `sql` sinks | SQL bound parameters are near-decisive, so it should settle deterministically |
| **Labeled corpus** | deliberately vulnerable *and* deliberately safe variants of the same endpoints | Ground truth. Without it, precision and recall numbers are meaningless |

Same language, two databases, two sink classes — that pair is a deliberate experiment, not just extra
test data. It tests the claim in `architecture.md` § 2 that the deterministic tier carries the SQL case
and the LLM carries the NoSQL case. If the split does not come out that way, the track boundary is wrong
and we revisit it early rather than after building on it.

## Ground rules

1. **Deterministic before LLM.** Value must exist with no model configured. Track A ships before any
   broker work.
2. **Ground truth before detection.** The labeled corpus lands in Phase 1, before any rule is written.
3. **Resolution before rules.** Rules running over an unresolvable graph produce noise, not findings.
4. **Bundle before LLM.** The bundle is the model's entire world; validate it with no model in the loop.
5. **Every phase ends runnable**, with a test that would fail if the phase regressed.
6. **No speculative abstraction** — see [`../standards/code-structure.md`](../standards/code-structure.md).

## Phases

### Phase 0 — Skeleton, fixture access, CI

**Goal:** toolchain, fixtures, and a quality gate exist.

- Go module, `cmd/scan`, walker that lists files, detects language, and respects ignore rules.
- Both TS apps vendored as read-only fixtures under `testdata/` (or cloned by a script — never modified).
- CI: build, `go vet`, `go test -race`, `gofmt -l`.

**Deliverable:** `scan --list <repo>` prints a file/language manifest.
**Acceptance:** runs against both fixtures; test asserts file and language counts; CI green.
**Why first:** proves the toolchain, the fixtures, and the repo shape before any analysis.

### Phase 1 — Labeled corpus and calibration

**Goal:** ground truth exists, and we know what context makes the model reason well — before building the
machinery around it.

- Hand-build paired vulnerable/safe endpoint fixtures in TS covering: NoSQL operator injection, SQL
  injection, command injection, SSRF, path traversal, weak crypto, hardcoded credentials, missing
  authorization, sensitive data in logs, permissive CORS, insecure deserialization.
- Each case records its expected outcome: rule ID, expected track, expected verdict, and a written
  reason.
- **Calibration, no Go code.** Hand-write the bundle for each case as a YAML file (snippets, the guard,
  the sink, the mitigation, the sink class), send it to the model, and record the verdict against the
  expectation. Iterate on the prompt and the bundle contents.
- Add a handful of **Track C scopes** — a route file with no ownership check, a chain that authenticates
  but never authorizes — and calibrate scope review the same way.

**Deliverable:** `testdata/corpus/` with expectations, plus a calibration log and the first prompt and
response schema.
**Acceptance:** a human reviews every case and its stated reason; paired cases differ *only* in the
mitigation; the calibration log records where the model was right and wrong, and why.
**Why second:** this is the measuring instrument *and* the prompt. It de-risks the model work for the cost
of writing YAML, and its output is reused directly by the phases that follow. It is not a go/no-go gate —
it is calibration that has to happen anyway, done before anything depends on it.

### Phase 2 — Parse and structural graph (JS/TS)

**Goal:** a deterministic graph of the fixture code.

- Settle the Tree-sitter binding by benchmark (cgo vs. pure Go) — this decides the static-binary goal.
- Grammar, parser worker pool, declarations, imports, call sites.
- Node identity including deterministic synthetic names for anonymous arrow handlers.
- `graph.json` export.

**Deliverable:** `scan --graph <repo>` emits `graph.json`.
**Acceptance:** golden-file test; determinism test (two runs byte-identical); every route handler in both
fixtures has a node.
**Risk:** anonymous-node naming stability under unrelated edits — test it explicitly.

### Phase 3 — Resolution and edges

**Goal:** a connected graph where a route reaches its sinks.

- Resolver: relative imports, barrel re-exports, `tsconfig` path aliases, monorepo workspaces,
  `node_modules` boundary.
- `CALLS`, `REFERENCES`, `IMPORTS`/`RE_EXPORTS`, `INHERITS`/`IMPLEMENTS` edges.
- Call-site argument bindings (the § 8 schema).

**Deliverable:** `scan path <A> <B>` prints a resolved chain, Graphify-style.
**Acceptance:** from a route handler in the Mongo fixture, the chain reaches the Mongo call, with
argument bindings on each edge. Same for the SQL fixture.
**Why here:** the highest-risk work in the project. Everything downstream is noise if resolution is wrong.

### Phase 4 — Classification

**Goal:** roles and sink classes assigned.

- Rule engine plus per-language query sets.
- Sources: Express/Fastify route registration and handler signatures.
- Sinks with sink class: `nosql`, `sql`, `ssrf`, `shell`, `path`, `log`, `html`, `deserialization`.
- Guards: validation, sanitization, authorization.
- Data-type classification (PII, PHI, credentials) as an impact signal.

**Deliverable:** `scan --classify` shows roles, sink classes, and data types.
**Acceptance:** every route in both fixtures is a source; every Mongo/SQL/HTTP call is a sink with the
correct class; the auth middleware is a guard.

### Phase 5 — Track A pattern rules → first useful tool

Runs in parallel with Phase 6: Track A needs the graph (§ 3 of the architecture) but not the path
tracker.

**Goal:** **a working deterministic scanner with zero LLM calls.**

- The OWASP pattern set from `architecture.md` § 4, weighted to the fixtures: request object as Mongo
  filter, `$where`, `$regex` from user input, SQL built by concatenation, `exec`/`spawn` with a template
  literal, SSRF from a request URL, path traversal, weak crypto, hardcoded connection strings and
  credentials, `rejectUnauthorized: false`, CORS `*` with credentials, sensitive data logged, missing
  security headers, unsafe deserialization.
- Finding record: rule ID, `detection: pattern`, OWASP category, CWE, severity, data classes, evidence.
- JSON and Markdown reports; CI-friendly exit codes.

**Deliverable:** `scan <repo>` produces findings with no API key configured.
**Acceptance:** run against the corpus; compare against expectations; write down every false positive and
false negative rather than tuning silently.
**Why here:** this is the milestone where the tool is genuinely useful, and it de-risks everything after
it. If Track A alone is not useful, the LLM stage will not save it.

### Phase 6 — Code path tracking

**Goal:** the spine. Candidate paths, pruned, pre-filtered, and prioritized.

- Reachability, dead-end pruning, minimal-subgraph reconstruction, deduplication.
- Pre-filter: drop paths where no call-site argument is source-derived or non-literal.
- Prioritization by sink severity, data-class impact, and node centrality when the set must be capped.

**Deliverable:** `scan --paths` lists candidate paths with their argument bindings.
**Acceptance:** expected paths present in both fixtures; constant-only paths dropped; candidate count
**and dropped count** reported so nothing disappears silently and LLM cost stays visible.

### Phase 7 — Context bundling

**Goal:** a bundle good enough that a human could adjudicate from it alone.

- Path summary, ordered snippets, argument bindings per edge, sink class and mitigation idiom, data
  types, and **off-path guard nodes** with their definitions.
- Track C scope bundles: the file or chain, its neighbours, and the middleware order around it.

**Deliverable:** `scan --explain <path-id>` dumps the bundle.
**Acceptance:** manual inspection. For a Mongo-injection candidate the bundle must contain the validation
middleware and the schema, not just the query. If a human cannot decide from the bundle, neither can the
model.
**Why before the LLM:** fully validatable with no model in the loop.

### Phase 8 — Path review (Track B)

**Goal:** verdicts for tracked paths.

- Broker client over an OpenAI-compatible endpoint; three-state verdict; mandatory cited evidence.
- Stub for tests; degraded mode when no endpoint is configured.
- Model ID and prompt version recorded per scan.

**Deliverable:** `scan` produces `validated` / `not_validated` / `uncertain` with citations.
**Acceptance:** stubbed tests including malformed responses; then a manual run where the deliberately
vulnerable Mongo filter is `not_validated` and its schema-validated pair is `validated`. **This is where
the paired corpus from Phase 1 earns its keep.**
**Hypothesis to record:** how many SQL vs. NoSQL paths reached Track B, and how many verdicts were
`uncertain`. A high `uncertain` rate means the bundle is missing context, not that the model is bad.

### Phase 8b — Scope review (Track C)

**Goal:** authorization, logic, and design findings.

- Scope selection: which route files, controllers, middleware chains, and auth modules are worth
  reviewing, and how that scales with repository size.
- Scope bundles and the scope-review prompt, calibrated in Phase 1.
- Findings attached to an anchor symbol, with severity and evidence rather than a path verdict.

**Deliverable:** `scan --scopes` and Track C findings in the report.
**Acceptance:** the corpus's deliberate authz gap is found; the scope-selection count stays bounded on both
fixtures so cost is predictable.
**Why here:** it shares the bundler and the broker with Phase 8, so it is cheap once those exist — and it
is the mode where semantic review most clearly exceeds static analysis.

### Phase 9 — Reporting, identity, findings tracking

**Goal:** usable in CI.

- Unified report across all three tracks, with `detection` on every finding.
- Fingerprint, previous-findings log, Git-tracked baseline, rule-level skips.
- CI gates on new findings only.

**Deliverable:** run twice → zero new findings; introduce a vulnerability → exactly one new finding.
**Acceptance:** Track A findings byte-identical across runs; a changed LLM verdict does **not** create a
new finding.

### Phase 10 — Rescan cache

**Goal:** second run is fast.

- Parse cache keyed by content hash; classification keyed by hash plus ruleset version.
- Path invalidation per `architecture.md` § 13; unchanged Track B paths keep their verdict.

**Deliverable:** second run on an unchanged fixture is materially faster and spends zero tokens.
**Acceptance:** cached and uncached runs produce identical output.

### Phase 11 — Python adapter

**Goal:** prove the abstractions hold.

- Decorators as sources, `self`, relative imports, `__init__.py` re-exports.

**Deliverable:** a Python fixture scanned end to end with no core changes.
**Acceptance:** the only new code is a language adapter and query sets. If core files changed, the
`architecture.md` § 6 interfaces were wrong — record what broke.

### Phase 12 — Go and Java adapters

Receiver methods, interface dispatch, annotations, and overloads. Java overloads are the test of the node
ID scheme.

## Measurement

- Run the corpus after Phases 5 and 8; record precision and recall per rule and per track.
- Track C is measured against the deliberately-gapped scopes in the corpus, plus a review of the
  false-positive rate on scopes that are actually fine — the failure mode there is noise, not silence.
- Record the **dropped-path count** from the pre-filter alongside findings, so cost tuning cannot hide
  lost findings.
- Differential testing against Bearer on both fixtures — a cheap external check on the Track A rules.
- Track candidate-path counts per fixture to keep LLM cost honest and visible.
- Track the `uncertain` rate as a first-class metric, not a footnote.

## Known gaps to design for

- **Finding grouping.** The same rule and sink reached by eight routes must present as one finding with
  eight entry points, or the report becomes noise.
- **Framework facts.** A curated per-library store of what actually mitigates each sink class, so the
  model reasons from facts rather than guessing at library behaviour.
- **Local binding table.** Extending argument binding with per-function local assignments, still
  syntactically. Closes the most common flow gap in real Express code.
- **Gating policy.** Gate on new fingerprints; surface verdict changes on unchanged code as
  informational, not as new findings.

## Deliberately deferred

- Inter-procedural dataflow — the whole point of the design is not to build it.
- LLM-driven traversal (the "deep dive" mode).
- A06 dependency scanning — a different tool.
- Bearer-style privacy report output; the `data_classes` signal already exists to build on.
- Graph database storage (LadybugDB) — revisit only if the graph outgrows memory.
