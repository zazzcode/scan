# Scan — project overview

The durable orientation document for this repository. Read it before architecture or implementation
documents.

**Status: draft, under active design iteration.** We are settling the design before writing code.

## What this is

`scan` is an OWASP Top 10 code scanner built on a **graph-narrowed LLM review** model. It is
deliberately not a traditional SAST tool.

Traditional SAST builds a source-to-sink dataflow engine: hand-coded taint propagation through
parameters, return values, struct fields, and aliases. That is years of engineering, it is brittle
across frameworks, and it still produces findings nobody can act on.

`scan` takes a different route. It builds a deterministic structural graph of the target codebase with
Tree-sitter — the Graphify model — and tracks the code paths through it. A deterministic rule tier
handles everything decidable by shape. Everything requiring judgment goes to the model, which reasons
over a prepared path or scope rather than raw files.

## Premise

**A capable LLM reasons about code better than a static analyzer, and handles the cases that break static
analysis.** Custom validators, framework indirection, dynamic dispatch, authorization and logic flaws,
and the question of whether a mitigation is *sufficient* — those are semantic questions, and they are
where SAST is structurally blind or noisy.

We do not try to out-engineer SAST on its own terms. We do the deterministic work static analysis is
genuinely good at, and hand the semantic questions to the model.

## Code path tracking is the spine

Everything depends on being able to resolve what calls what across files. It is what makes the model
stage tractable, in two ways:

- **It reduces what the model must review.** Pointing a model at a repository spends its whole context
  window locating relevant code, with no guarantee it found the right code. A tracked path hands it
exactly the code on that path, in order, with the guard and sink already identified.
- **It provides structure for analysis.** A path is an explicit claim — this entry point reaches this
  operation, through these functions, with these guards, passing these arguments. The model reasons over
  that structure rather than over loose files, and both the model and a human reviewer can cite it.

It is also the only reason the tracks share anything: Track A needs call relationships, Track B is
source→sink paths, and Track C scopes are graph-derived neighbourhoods. One graph, three uses. See
[`architecture.md`](architecture.md) § 3.

## Three detection tracks

The core structural idea, and the reason the tool is not "solely relying on the LLM".

- **Track A — pattern rules.** Deterministic end to end, zero model calls. The rule *is* the verdict, so
  findings are byte-stable across runs. This track carries most of the OWASP Top 10. These findings gate
  CI by default.
- **Track B — path rules.** A rule finds a source→sink path; its *existence* is deterministic, but
  whether the data is sufficiently validated is the model's judgment. Identity comes from the path, never
  from the verdict.
- **Track C — scope review.** The model reviews a scope rather than a path: a route file, a controller, a
  middleware chain, an auth module. This is where authorization, logic, and design flaws live — the
  findings whole-program analyzers are structurally blind to, because they are not dataflow properties.

Every finding records which track produced it. That is what lets a consumer gate hard on Track A and
treat low-confidence Track B and C findings as advisory.

**Consequence: `scan` is useful with no API key configured.** Track A and the full deterministic
pipeline still work. The degraded mode is a product mode, not a fallback.

## Division of labour

| Stage | Who | What |
|---|---|---|
| Structure | Graph (deterministic) | What code exists, what calls what, what is a source, a sink, or a guard, what sink class, and what arguments are passed at each call site |
| Discovery | Graph (deterministic) | Which source-to-sink paths exist, pruned and pre-filtered to the ones worth examining |
| Pattern detection | Rules (deterministic) | Known-bad shapes: weak crypto, `eval`, unparameterized queries, hard-coded credentials, permissive CORS |
| Judgment | LLM (semantic) | Is this path real, is the guard sufficient, is the sink exploitable |

The rule of thumb: **if it can be decided by matching a known shape, it is a rule. If it requires judging
whether a human's mitigation is adequate, it is the LLM.**

## Why this instead of rules alone, or dataflow

- **Sufficiency of validation is a judgment call.** Custom validators, regex quality, ORM
  parameterization, template escaping, framework allow-lists. Encoding all of that as rules is a losing
  game; judging it is what language models are good at.
- **The expensive part of SAST is dataflow, and this split removes the need for it.** No fixpoint
  analysis, no aliasing, no inter-procedural propagation.
- **Findings stay explainable.** Every verdict is attached to a concrete call path with `file:line`
  evidence.

## Prior art

- **[Graphify](https://github.com/Graphify-Labs/graphify)** — Apache-2.0. Parses code with tree-sitter
  AST, explicitly *"deterministic, no LLM"*, with graph build measured at **zero LLM credits**. It is the
  strongest available evidence that the deterministic tier carries the load. We adopt its
  `EXTRACTED`/`INFERRED` edge tagging, its `graph.json` artifact, its type-hierarchy edges, its
  god-node centrality, and its local-first posture. Its path-tracking model is the basis for
  [`architecture.md`](architecture.md) § 3.
- **[Bearer](https://github.com/bearer/bearer)** — open-source SAST covering the OWASP Top 10 and CWE
  Top 25, with the dataflow analysis we are deliberately not building. That makes it the best reference
  for *where rules beat a model*: its published rule catalog is close to a specification for our Track A,
  and we adopt its CWE mapping, severity tiers, rule-level suppression, and data-type classification. It
  is a complement rather than a competitor: they model flows across files, we reach further semantically,
  and the two error profiles differ.

## Who it serves

- **Application security engineers** triaging a codebase for OWASP Top 10 exposure.
- **Developers** who want a path-level explanation of why something was flagged, not a category label.
- **CI pipelines** gating on new findings.

## Language scope

- **First implementation: JavaScript and TypeScript (Node.js).** Chosen first because they are the
  hardest resolution case — barrel re-exports, `tsconfig` path aliases, inline anonymous handlers,
  middleware registered by reference.
- **Two initial acceptance targets**, both TypeScript: one service over **MongoDB**, one over **SQL**.
  Same language, two sink classes, and a measurable hypothesis — SQL should mostly settle
  deterministically because bound parameters are near-decisive, while NoSQL should mostly escalate to
  Track B because Mongo has no equivalent escape hatch.
- **Second: Python.** Stresses a different axis — decorator-based routing, `self`, no static types.
- **Also in the initial release set: Go and Java.** Go for receiver methods and interface dispatch; Java
  for annotations and method overloads.
- The core is language-neutral. **Adding a language must mean adding a grammar, query sets, and a
  resolver. Never changing the core.**

## Major capabilities

- Concurrent repository ingestion and multi-language Tree-sitter parsing.
- **Code path tracking** across files — the spine that everything else depends on.
- A normalized, language-agnostic graph with source, sink, and guard classification, sink classes, and
  data-type classification.
- Call-site argument annotation, so the model's flow judgment is checkable rather than guessable.
- Deterministic pattern rules covering most of the OWASP Top 10, with CWE mapping and severity.
- Source-to-sink path discovery with pruning, pre-filtering, and prioritization.
- Token-bounded context bundling that includes off-path guard definitions.
- Path review with a three-state verdict and cited evidence.
- Scope review for authorization, logic, and design gaps that path-based analysis cannot reach.
- JSON and Markdown reporting; scan-to-scan findings tracking.

## Positioning

A semantic review tool, complementary to a dataflow SAST rather than a replacement for it. A dataflow
analyzer over-reports inside the scope it models; `scan` reaches further semantically and under-reports
outside what its graph and scopes cover. Different error profiles — run both.

It is not a whole-program analyzer, so one consequence is worth stating once: a clean run is not a
security guarantee. CI gating on new findings is meaningful; "no findings" is not a certification.

What it deliberately does not do is build a dataflow engine — no inter-procedural taint propagation, no
value tracking through statements, assignments, or branches. What it does record is call-site argument
annotation: at each call site, which argument expressions are passed and which parameter position each
lands in. That is extraction at the call site, not tracking through the program, and it is the flow
evidence the model relies on to check rather than guess.

## Findings tracking

Know what is **new**, what **persists**, and what is **resolved** since the last scan. That is set
membership over a log of previous findings — not a vulnerability lifecycle state machine.

- The tracked entity is the **deterministic structural path**: rule + the qualified symbol chain + guard
  status.
- The **LLM verdict is an annotation, not part of identity.** A verdict change never creates or resolves
  a finding.
- Human decisions — accepted, suppressed, false positive — live in a **Git-tracked baseline** reviewed in
  a PR, plus rule-level skips, not as rows in a database.
- `fixed`, `first_seen`, and `last_seen` are derived from the log, not stored state.

## Decisions log

- **Approach (2026-10-11):** graph-narrowed LLM review, not traditional SAST. No hand-coded
  source-to-sink string tracking.
- **Premise (2026-10-11):** a capable LLM reasons about code better than a static analyzer and handles
  the cases that break static analysis. Stated up front as the design's foundation.
- **Path tracking is the spine (2026-10-11):** tracking code paths across files is the core capability.
  It reduces what the model must review and gives it structure to reason over. All three tracks consume
  the same graph.
- **Three tracks (2026-10-11):** Track A pattern rules are deterministic and gate CI; Track B path rules
  get model judgment; Track C reviews a scope (route file, controller, middleware chain, auth module) to
  reach authorization and logic flaws that path analysis cannot. Every finding records its track.
- **Positioning (2026-10-11):** complementary to dataflow SAST with a different error profile, not a
  replacement. The caveat that a clean run is not a certification is stated once, not repeated.
- **Deterministic-first (2026-10-11):** keep the deterministic tier as large as it can honestly be.
  Track A plus the full pipeline works with no API key.
- **Flow grounding (2026-10-11):** the graph carries call-site argument annotations. Syntactic
  extraction, not dataflow.
- **Verdict shape (2026-10-11):** three states — validated, not-validated, uncertain. Not binary.
- **Traversal (2026-10-11):** we traverse deterministically and the LLM judges each prepared path.
  LLM-driven traversal is deferred as an opt-in deep-dive.
- **Finding identity (2026-10-11):** the structural path is the tracked identity; the LLM verdict is an
  annotation.
- **Findings tracking (2026-10-11):** previous-findings log plus diff, with a Git-tracked baseline and
  rule-level skips. No lifecycle state machine, no findings database.
- **Language order (2026-10-11):** TypeScript/JavaScript first (two fixtures: MongoDB and SQL), Python
  second, designed for Go and Java.
- **Sink classes (2026-10-11):** `nosql`, `sql`, `shell`, `path`, `ssrf`, `html`, `deserialization`,
  `log`. NoSQL escalates to Track B more often than SQL because Mongo has no parameterization escape
  hatch.
- **Prior art (2026-10-11):** Graphify validates the deterministic split; Bearer is the reference for
  what must stay deterministic and for our honest positioning.
- **Persistence (2026-10-11):** in-memory graph for execution; JSON for the artifact; a local cache for
  rescan speed. The graph store is never the execution engine.
- **Scope at bootstrap:** repository scaffolding and documentation only; no Go packages yet.

## Documentation set

- [`architecture.md`](architecture.md) — tech stack and technical approach
- [`project/implementation-plan.md`](project/implementation-plan.md) — phased plan to the final solution
- [`project/01-functionality.md`](project/01-functionality.md) — capabilities and coverage
- [`project/02-project-scope.md`](project/02-project-scope.md) — boundaries, sequencing, non-functional targets
- [`docs/scan-agent-orientation.md`](docs/scan-agent-orientation.md) — agent orientation and routing
- [`standards/index.yaml`](standards/index.yaml) — engineering standards routing table
- [`specifications/README.md`](specifications/README.md) — deliverable specifications
