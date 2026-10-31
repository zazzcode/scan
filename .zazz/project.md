# Scan — project overview

**Project definition: proposed direction, 2026-10-19.** This overview owns product purpose, direction,
and repository state. The [agent orientation](docs/scan-agent-orientation.md#documentation-map) routes
to the detailed documents and defines their responsibilities.

## Purpose

`scan` is a proposed Go static application security scanner for selected OWASP Top 10 risks. It aims
to give developers reviewable source evidence, security engineers useful triage context, and CI an
explicit policy over findings and analysis limits. A clean run is not a security certification.

## Long-term objectives

**Objective update, 2026-10-31:** build a source-code security scanner that is **as good as or better
than Bearer at SAST detection and code analysis**, with the long-term goal of supporting all equivalent
analysis and developer-workflow capabilities. This includes capabilities available across Bearer CLI
and Pro; the free CLI is the first local comparison baseline, not the ceiling for scan's ambition.

Track parity through a versioned capability inventory: supported languages/frameworks and weakness
families; source-to-sink taint and control/dataflow analysis; interprocedural and cross-file analysis;
sensitive-data classification and security/privacy impact; finding evidence and remediation guidance;
and reporting, configuration, incremental/baseline and CI workflows. Verify each comparator capability
and edition rather than assuming a marketing description establishes behavior. The detailed inventory
and phased delivery contracts remain to be developed.

Measure equal-or-better detection on matched, pinned source scopes and independently labeled,
held-out vulnerable/protected/unknown cases. Compare recall, precision, false safe claims, abstentions,
coverage and evidence quality; record latency, resource use and model cost alongside accuracy. More
findings alone do not establish superiority. Demonstrated parity applies only to the capabilities and
corpus actually evaluated; broad parity has not been achieved by the current implementation.

The intended advantage is graph-grounded semantic assessment of custom validation and sink protection,
with route-oriented explanations and complementary remediation. Test that advantage against traditional
SAST rather than assume that using AI improves detection. The
[Bearer comparison plan](evaluation/bearer-comparison.md) starts this evaluation; the
[phased plan](project/implementation-plan.md) owns delivery order. SCA remains a complementary objective
after the source-code analysis foundation.

## Product direction

The chosen direction combines **Tree-sitter extraction, Go control-flow/dataflow analysis, OpenAI Decision
assessment of string validation/sanitization, then an LLM Findings Reviewer for identified issues**.
All scanner and code-analysis implementation is Go, including model clients and orchestration.
The decision model's role is bounded semantic assessment; the reviewer provides deeper findings review.
Experiments must measure accuracy, latency, and total cost, including whether the model recognizes
adequate protection and whether the reviewer corroborates or challenges its identified issues. Local-only operation
remains useful, with omitted assessment/review stages explicitly marked.

**Direction update, 2026-10-24.** The first slice is greenlit and implemented:
[route identification](specifications/tree-sitter-route-identification-mvp.md). `scan --routes`
writes a JSON route index of Express-style JS/TS registrations and the handler names at the call.
It does not follow imports or decide whether a string is validated. The OpenAI Decisions API
(`gpt-6-luna`) is the later assessor of validation and sanitization. It is not called in this slice.
The [architecture provider decision](architecture.md#current-provider-decision) keeps GPT Decisions
and Jev as alternative semantic judge providers, with OpenAI first and future adapters allowed.
The [graph design](architecture/repository-graph.md) covers the target repository: API entry points,
internal paths, outbound service calls and their strings, responses, and subsequent value flow.
The 2026-10-31 [assessment](project/06-route-validation-assessment.md) records route verification
and the next branch scope. The implementation update below supersedes its initial pre-code status;
later [two-case Luna results](evaluation/juice-shop/results.md) supersede that initial no-live-test status.

Start with a small JS/TS slice, then expand from measured results. The
[capability catalog](project/01-functionality.md) defines detection modes, rule families, language
sequence, and interfaces. The [scope](project/02-project-scope.md) defines boundaries, quality
hypotheses, and risks. First-release coverage remains undecided.

Technical responsibilities live in the [architecture](architecture.md). The
[design proposal](project/03-jev-design-proposal.md) owns the change from generative call-path review
to supported Go dataflow plus bounded semantic classification, along with the rationale, retained
concepts, and open decisions. The [Jev research](project/04-jev-validation-research.md) records published
evidence and the evaluation needed to measure the chosen assessment-and-review pipeline.

## Concept-stage repository state

The module is `github.com/zazzcode/scan`, declaring Go `1.27.2`. `scan --list` walks a tree.
`scan --routes` parses JS/TS with Tree-sitter and emits the route index from the
[route-identification spec](specifications/tree-sitter-route-identification-mvp.md). There is no
complete CFG/interprocedural analysis, independent Findings Reviewer, or reviewed findings report yet.

**Implemented 2026-10-31:** `--strings` adds bounded route/handler resolution and request-value
use witnesses with source hashes, candidate SQL/outbound HTTP call models, and explicit gaps.
Optional `--judge openai|jev` sends bounded evidence through provider-specific typed adapters;
bounded Luna integration/control evaluations are recorded in [results](evaluation/juice-shop/results.md). See the [current contract](specifications/route-string-validation.md)
and [API setup](docs/openai-decisions-setup.md).

Phase 0's walker, language registry, and CI remain. Their historical record is the
[Phase 0 specification](specifications/phase0-skeleton/README.md). Recognizing a language extension
does not establish security-analysis support for that language.

Workspace conventions and the `stage` rules live in the
[orientation](docs/scan-agent-orientation.md).

## Route reporting and dependency direction

The product should explain findings by inbound route and source-to-use graph, separating field
validation from sink/API semantics and proposing multiple complementary improvements. A valid email
format does not establish SQL parameter separation. Keep the canonical internal evidence model and
provide the bounded SARIF 2.1.0 export; see [report design](architecture/reporting-and-sarif.md).
After the code-path foundation, add [SCA dependency auditing](architecture/sca-design.md) using
language/package-manager tools, with advisory evidence separate from reachability.

## Delivery and evidence

The [implementation plan](project/implementation-plan.md) owns phase decomposition, future spec slices,
dependencies, acceptance targets, and phase-specific evidence. Those phases are proposals, not approved
contracts or scheduled milestones. Specifications will be authored in separate sessions and greenlit
through the [Zazz delivery lifecycle](docs/zazz-methodology.md#specification-driven-delivery-cycle).

The plan includes the cloned OWASP Juice Shop reference and small paired fixtures. Its
[acceptance-target section](project/implementation-plan.md#acceptance-targets) owns their role and setup
work. Source handling, accuracy, latency, cost, and packaging remain hypotheses to measure rather than
established product promises.
