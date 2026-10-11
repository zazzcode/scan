---
last_updated_at: 2026-10-11
---

# Standards contextual split

This note explains how the standards are split and how the team maintains that split. It is written for
humans. `index.yaml` is the machine-readable source agents use to decide which standards to load.

## Why split standards

Standards have to be usable during implementation and review. When one file covers too many topics,
agents load rules unrelated to the task, humans scan too much text, and edits to unrelated rules collide
in the same document.

The goal is contextual loading: an agent working on path pruning should load the path-tracking rules, not
the LLM broker's prompt schema or the reporter's output format. A reviewer should be able to point to the
specific standard for the work under review.

Line-count thresholds are a reviewability signal:

- Under 400 lines is preferred. A file this size usually fits in agent context alongside the code, tests,
  specs, and command output for the task.
- Up to 500 lines is acceptable when the topic is still cohesive and another split would make discovery
  worse.
- Over 600 lines should block review unless the file is split or an explicit exception is accepted. At
  that size, the file is probably mixing contexts.

Splits follow work areas, not arbitrary halves. Do not create `part-1` / `part-2` files. Use names that
describe the context: parsing, graph model, rules, tracing, bundling, broker, reporting.

## Provenance requirement

Splitting a standard must not remove hard references. Each normative standard keeps a citation to the
evidence that justifies the rule — a sibling standard, a repo-relative source example, a public issue or
PR, or a short "review precedent" note when the originating discussion is not public.

This file is a map of the split. The rule text and its citations remain in the individual standards
listed below.

## Curation policy

Two tiers exist:

1. **Vendored methodology standards** — stack-agnostic, adapted from
   [zazz-skills](https://github.com/zazzcode/zazz-skills) `docs/standards/`. Local divergence is
   deliberate and recorded in [README.md](README.md).
2. **Scan stack standards** — authored for this repository with the `standard-builder` skill against real
   code.

The upstream stack-specific baselines (Python/Flask, SQL Server stored procedures, React/Next, AWS Lambda,
Drizzle/Neon) are deliberately **not** vendored. This repository is a single Go module with no web
frontend, no SQL Server, and no serverless deployment.

## Scan standards

### Structure and process

- `code-structure.md`: Go package layout (`cmd/`, `internal/`), file-size thresholds, contextual
  splitting, module cohesion, anti-slop rules, compute-once guidance, and determinism requirements. Use
  when creating or reviewing Go packages, or any authored file that may grow too large or duplicate
  structure.
- `testing.md`: table-driven Go tests, golden-file parsing fixtures, graph-fixture tests for path
  tracking, a stubbed broker, no network or wall-clock dependence, and the evidence expected on a PR. Use
  when writing or reviewing any test.
- `pr-process.md`: PR title scope labels, one logical change per PR, no transitional cruft on merge, CVE
  identifiers in titles, and the draft-first human-merge workflow. Use when opening, titling, or scoping a
  PR.

### Documentation

- `docs-hygiene.md`: voice, Desired/Not-desired examples, cross-document links, banned transient
  identifiers, and cleanup clauses. Use when writing standards docs or agent guides.
- `docs-hygiene-reference-structure.md`: standard-doc exclusions, citation formats, section ordering, and
  a worked structure. Use when shaping standards documentation.
- `spec-hygiene.md`: bounded scope, acceptance-criteria quality, test-plan linkage, halt conditions, and
  avoiding spec drift. Use when authoring or reviewing a deliverable specification.

### Agent skills and routing

- `agent-skills.md`: Agent Skills directory shape, portable front matter, metadata placement, description
  quality, progressive disclosure, and validation. Use when creating or modifying a skill under
  `.agents/skills/`.
- `contextual-split.md`: this file — how the standards are split, and which document to load for a task.

## To be authored for the scan stack

Backlog. Author each with the `standard-builder` skill against real code once the relevant package
exists — these describe real patterns, not aspirations.

- `path-tracking.md`: the code path model. What a path is, what the graph guarantees, what the tool does
  and does not claim, how reachability differs from dataflow, and how findings are phrased so the output
  never overclaims. This is the standard that keeps the tool honest about its own semantics.
- `parse-layer.md`: Tree-sitter grammar registry, parser worker pool, error handling, and the binding
  abstraction.
- `rule-authoring.md`: `.scm` query conventions, rule naming, sink classes, and how a rule declares its
  OWASP category, CWE, and severity.
- `graph-model.md`: node and edge schema evolution, node identity including anonymous callables, and
  graph export compatibility.
- `arg-binding.md`: call-site argument annotation conventions and the local binding table.
- `llm-broker.md`: prompt shape, structured response schema, verdict semantics, retry and cost controls,
  and model and prompt versioning.
- `finding-record.md`: the finding schema, fingerprint construction, previous-findings log, baseline
  format, and gating policy.
- `cli-surface.md`: flags, output formats, exit codes, and CI usage.
- `ci-workflows.md`: GitHub Actions triggers, Go build/vet/test/race jobs, and format gates.
- `observability.md`: structured logging and progress reporting for long scans.

## Maintenance rules

When adding a rule, put it in the file whose context an agent would naturally load for the task. If a
rule spans contexts, keep the normative rule in the most specific file and add a short cross-reference
only where it prevents missed discovery.

When a file approaches 400 lines, look for a contextual boundary. Split by task area, artifact type, or
review question. Do not split by page count alone.

After adding, moving, or renaming a standards file, update `index.yaml` with:

- the file name
- the paths or activities that should trigger the file
- a concise purpose naming the decisions the file governs
