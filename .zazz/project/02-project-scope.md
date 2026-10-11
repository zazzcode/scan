# 02 — Project Scope

Boundaries, objectives, delivery sequence, and non-functional targets. Phased detail is in
[`implementation-plan.md`](implementation-plan.md); design is in [`../architecture.md`](../architecture.md).

## 1. Objective

Build a fast, multi-language, OWASP Top 10 code scanner whose findings are grounded in explicit
source-to-sink call paths or explicit known-bad patterns, so both humans and CI can act on them without
re-deriving the reasoning — and whose deterministic tier is large enough to be useful with no model
configured.

## 2. In scope

- Concurrent repository walking and Tree-sitter parsing.
- A normalized, language-agnostic graph with source, sink, and guard classification, sink classes, and
  data-type classification.
- Call-site argument binding.
- Deterministic pattern rules across the OWASP Top 10, with CWE mapping and severity.
- Source-to-sink path discovery with pruning and pre-filtering.
- Token-bounded context bundling including off-path guard definitions.
- Path review with a three-state verdict and cited evidence.
- Scope review for authorization, logic, and design gaps that are not dataflow properties.
- JSON and Markdown reporting.
- Stable finding identity and scan-to-scan tracking (new / persisting / resolved).
- Incremental rescans: reuse unchanged parse results and previously recorded paths.
- A Git-tracked findings baseline plus rule-level skips, so CI gates only on new findings.
- Grammars, query sets, and resolvers for TypeScript/JavaScript, Python, Go, and Java.

## 3. Out of scope

- **Inter-procedural dataflow analysis.** No taint propagation, no aliasing, no value tracking. This is
  the central scope decision, not an omission.
- **A06 — vulnerable and outdated components.** Dependency scanning is a different tool.
- Runtime instrumentation, dynamic analysis, or test-coverage-driven tracing.
- Automatic code remediation or patching.
- Hosting, multi-tenancy, or a persistent service — this is a CLI tool at this stage.
- Full runtime resolution of dynamically registered routes. An inference pass may flag these; they are
  not guaranteed to resolve.
- Secrets-scanning beyond hard-coded credentials found in source.
- Graph database storage. Revisit only if the graph outgrows memory.

## 4. Delivery sequence

Full detail in [`implementation-plan.md`](implementation-plan.md). The shape:

1. Skeleton, fixture access, CI.
2. Labeled corpus and model calibration — ground truth and prompt before any detection logic.
3. Parse and structural graph (TS/JS).
4. Resolution and edges — the highest-risk work, and what path tracking depends on.
5. Classification.
6. **Track A pattern rules — the first genuinely useful tool, with zero LLM calls.**
7. Code path tracking — discovery, pruning, pre-filter, prioritization.
8. Context bundling.
9. Path review (Track B) — three-state verdict with cited evidence.
10. Scope review (Track C) — authorization, logic, and design findings.
11. Reporting, identity, findings tracking.
12. Rescan cache.
13. Python adapter — the test of whether the abstractions hold.
14. Go and Java adapters.

Each step is a separate deliverable with its own specification and its own PR into `stage`.

## 5. Non-functional targets

- **Throughput:** parse and graph a ~1,000-file repository in low single-digit seconds on a laptop.
- **Parallelism:** saturate available cores for I/O and AST traversal; no global serialization point.
- **Memory:** bounded by the largest single file plus the graph; no full-repository text retention.
- **Determinism:** the structural graph, the path set, and **every Track A finding** must be identical
  across runs for the same input, independent of goroutine scheduling. Track B verdicts are explicitly
  not deterministic, which is why finding identity never depends on them.
- **Zero-token operation:** Track A and the full deterministic pipeline must work with no API key and
  make no network calls.
- **Token discipline:** the LLM must never receive a whole repository, only pruned, pre-filtered path
  bundles.
- **Portability:** a single static binary per platform; CGO must not be a hard requirement. The
  Tree-sitter binding decision is made by benchmark, and this target is what is at stake.
- **Local-first:** no source code leaves the machine except the specific path bundles sent to a
  configured endpoint.

## 6. Risks

- **Path explosion.** Reachability between sources and sinks is combinatorially nasty; pruning and
  pre-filtering are the core algorithmic risk and directly determine LLM cost.
- **Pre-filter over-aggression.** A pre-filter tuned for cost will silently drop real findings. It needs
  measurement against the corpus, and dropped paths should be reportable rather than invisible.
- **Resolution correctness.** Barrels, path aliases, and monorepo layouts in TS/JS. If resolution is
  wrong, every downstream finding is noise.
- **Anonymous node identity.** Inline arrow handlers are the most common real handler shape and the
  easiest to name unstably. Unstable naming means unstable fingerprints means scan-to-scan churn.
- **LLM non-determinism and cost.** Mitigated by keeping identity structural, caching verdicts for
  unchanged paths, and recording model and prompt versions per scan.
- **`uncertain` rate.** A high rate means the bundle is missing context, not that the model is weak. It
  is a bundle-quality metric and should be tracked as one.
- **False positives.** Both tracks can produce them. Every finding must carry its evidence so a human
  adjudicates in seconds.
- **Rule catalog maintenance.** The rule set is data and will grow; without a documented catalog and
  per-rule tests it will rot. Bearer's catalog is the maintenance reference.
- **Grammar and runtime drift.** Five languages means five independent Tree-sitter grammar upgrade
  paths.
- **Track boundary.** The claim that SQL settles deterministically and NoSQL escalates is a hypothesis,
  tested by the two TypeScript fixtures. If it does not hold, the boundary needs revisiting early.

## 7. Open questions

- Which Tree-sitter binding: cgo or a pure-Go runtime, and does either forfeit the static-binary target?
- Is the LLM broker provider-agnostic from day one, or OpenAI-compatible only initially?
- Should findings be schema-versioned from the first release?
- Where does the rescan cache live: a content-addressed directory or SQLite?
- How much rename churn is acceptable before a "likely moved" heuristic is warranted?
- Does the pre-filter threshold need to be configurable per repository, or is one default defensible?
