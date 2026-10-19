# Scan — project overview

**Project definition: proposed direction, 2026-10-19.** This overview owns product purpose, direction,
and repository state. The [agent orientation](docs/scan-agent-orientation.md#documentation-map) routes
to the detailed documents and defines their responsibilities.

## Purpose

`scan` is a proposed Go static application security scanner for selected OWASP Top 10 risks. It aims
to give developers reviewable source evidence, security engineers useful triage context, and CI an
explicit policy over findings and analysis limits. A clean run is not a security certification.

## Product direction

The chosen direction combines **Tree-sitter extraction, Go control-flow/dataflow analysis, Jev
assessment of string validation/sanitization, then an LLM Findings Reviewer for identified issues**.
All scanner and code-analysis implementation is Go, including model clients and orchestration.
Jev's role is fast, lower-cost semantic decisions; the reviewer provides deeper findings review.
Experiments measure this speed/cost objective alongside how well Jev recognizes adequate
protection and how the reviewer corroborates or challenges its identified issues. Local-only operation
remains useful, with omitted assessment/review stages explicitly marked.

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

No security scanner has been built or validated. The tree contains bootstrap source artifacts:
language-extension registry, deterministic file walker, `scan --list [--json] <path>`, tests, and CI
configuration. Their health and fit have not been established in this definition session; they
correspond to the historical [Phase 0 specification](specifications/phase0-skeleton/README.md).

There is no parser, graph engine, security rule engine, Jev client, or security findings report yet.
The repository is a single Go module, `github.com/zazzcode/scan`, currently declaring Go `1.27.2`.
Recognizing a language extension does not establish security-analysis support for that language.

Current work refines product, architecture, and, as assigned, Zazz methodology documents. It does not
implement the roadmap or create new execution specifications. Workspace conventions and the `stage`
definition-session rules live in the [orientation](docs/scan-agent-orientation.md).

## Delivery and evidence

The [implementation plan](project/implementation-plan.md) owns phase decomposition, future spec slices,
dependencies, acceptance targets, and phase-specific evidence. Those phases are proposals, not approved
contracts or scheduled milestones. Specifications will be authored in separate sessions and greenlit
through the [Zazz delivery lifecycle](docs/zazz-methodology.md#specification-driven-delivery-cycle).

The plan includes the cloned OWASP Juice Shop reference and small paired fixtures. Its
[acceptance-target section](project/implementation-plan.md#acceptance-targets) owns their role and setup
work. Source handling, accuracy, latency, cost, and packaging remain hypotheses to measure rather than
established product promises.
