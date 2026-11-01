# Scan — agent orientation

Compact session context and document routing. Read [AGENTS.md](../../AGENTS.md) for binding rules and
when to load supporting guides. Keep this document under 400 lines; link to details rather than copy them.

## Table of contents

- [Current work](#current-work)
- [Portable workspace and DEV_ROOT](#portable-workspace-and-dev_root)
- [Branches and worktrees](#branches-and-worktrees)
- [Task reading paths](#task-reading-paths)
- [Documentation map](#documentation-map)
- [Working protocol](#working-protocol)
- [Methodology refinement](#methodology-refinement)
- [Commands and credentials](#commands-and-credentials)

## Current work

**Current worktree, 2026-11-01: `logger-leak`.** Verify the actual path, branch, and changes with
`pwd` and `git status --short --branch`. The user requested stable agent guidance, README setup/usage,
clearer validation and database results, fresh offline target scans, and JSON-driven standalone HTML with expandable rows and remediation.
See the [reporting contract](../specifications/human-readable-reports.md). Logger-leak detection follows
this reporting increment and is not implemented by it. The historical definition-only `stage/`
exception is not this assignment.

- **Implemented:** `--list` and a Tree-sitter JS/TS route inventory (`--routes`).
- **Verified:** Juice Shop `v20.3.0` login/search registrations; repeatable syntactic output.
  Whole-tree matches include HTTP clients and training snippets. Read the
  [assessment](../project/06-route-validation-assessment.md#route-verification) before claiming coverage.
- **Next objective:** extend automatic claim-scoped evidence assembly to novel targets, then independently
  evaluate findings; selected route/check bundles already connect inputs, helpers, guards and uses.
- **Long-term objective:** equal-or-better SAST detection/code analysis than Bearer and equivalent
  capabilities across CLI/Pro with fewer false positives. For each route, report sufficient protection,
  lacking protection, and poor code-practice findings; see [objectives](../project.md#long-term-objectives). This is an ambition,
  not current parity. Start with the local CLI comparison and build a versioned capability inventory.
- **Model direction:** GPT Decisions or Jev provides the semantic judge, with future adapters allowed.
  OpenAI is first. The independent LLM Findings Reviewer remains a separate role; provider schemas,
  keys, quality and calibration are separate. See the [graph design](../architecture/repository-graph.md)
  for inbound API routes, internal paths, outbound requests and response continuations.
- **Implemented in this worktree:** `--strings` emits bounded source-to-use dependency witnesses,
  resolves selected relative named imports and handler factories, and records candidate SQL/outbound
  HTTP uses. Optional `--judge openai|jev` provides separate validation/protection answers and explicit
  failures. See the [contract](../specifications/route-string-validation.md) and [setup](openai-decisions-setup.md).
- **Still incomplete:** general CFG/interprocedural/heap analysis, effective middleware, external
  response continuations and independent Findings Review. The [two-case Luna run](../evaluation/juice-shop/results.md) produced supported potential findings;
  [original controls](../evaluation/corpus/README.md) exercise vulnerable/protected/unknown cases;
  held-out repository accuracy and runtime confirmation remain unmeasured.
  Current user restriction: live calls use `gpt-6-luna` only, with small explicit request caps.
- **Authorization:** the earlier route/string implementation request is historical context, not
  authorization for a new scanner slice. Explicit user requests are sufficient under
  [AGENTS.md](../../AGENTS.md#authorization-and-delivery); no separate greenlight is required.

**Local session handoff:** when resuming this checkout, read
[SQL/NoSQL judgment handoff](../ephemeral/sql-nosql-judgment-handoff.md) first. It records current
calibration failures and the logger-leak gate. The earlier
[route-string-validation handoff](../ephemeral/route-string-validation-handoff.md) is historical context. This is an
ignored local execution artifact; durable design/contracts/results remain indexed below.
The initial [Bearer comparison](../evaluation/bearer-results.md) is complete: both samples and
thirteen controls, no new model calls. Supported same-file and relative-imported helper calls now
discover SQL sinks through argument, parameter and direct-return evidence. The next slice is
non-mutating guards versus reassignment, then a fresh capped Luna run on the changed bundles.
The [Juice Shop index](../evaluation/juice-shop/README.md) includes all 116 documented challenges and
a 671-finding review ledger; three findings are source-inspected false positives, 664 remain pending.

The [near-term plan](../project/implementation-plan.md#route-string-validation-delivery-sequence)
owns proposed delivery order. The [architecture](../architecture.md#current-provider-decision) owns
provider responsibilities. The [project overview](../project.md#concept-stage-repository-state) owns
repository capability state.

## Portable workspace and DEV_ROOT

**`<DEV_ROOT>` means this computer's development directory**, usually `Dev/` under the user's home.
It is documentation notation, not an implemented environment variable. Usual locations:

- Repositories: `<DEV_ROOT>/<account>/<project>`; scan's container is `<DEV_ROOT>/zazzcode/scan`.
- Reference clones: `<DEV_ROOT>/juice-shop` (application), `<DEV_ROOT>/bearer`, and `<DEV_ROOT>/codeql`
  (SAST design references).

Resolve/verify roots each session; durable docs/configuration MUST NOT hardcode usernames or absolute
checkout prefixes. Other documents should link here for workspace conventions.

## Branches and worktrees

Scan uses a shared bare Git directory with sibling worktrees. `stage` integrates work and `main` is
release-only. Use the assigned feature worktree for related design, implementation, and tests.
Independent future deliverables may use separate feature worktrees and PRs into `stage`.
Load [worktree setup](worktree-setup.md) for the layout, reference-clone handling, naming, creation,
and recovery details; [AGENTS.md](../../AGENTS.md#branch-model) declares the current `stage/` exception.

## Task reading paths

After `AGENTS.md` and this orientation, choose the smallest applicable row. Load linked sections first;
load whole design documents only when the task crosses their boundaries. A link is not an include.

| Task / question | Read next | Load more when |
|---|---|---|
| Run scan or confirm route extraction | [README](../../README.md#run), [verification and limits](../project/06-route-validation-assessment.md#route-verification) | Read the [route spec](../specifications/tree-sitter-route-identification-mvp.md) and `internal/parse` when changing the route contract |
| Plan this branch's input-validation work | [near-term sequence](../project/implementation-plan.md#route-string-validation-delivery-sequence), [assessment](../project/06-route-validation-assessment.md#architecture-and-plan-assessment) | Read architecture §§3–7 for resolution, value/guard semantics, or evidence schema changes |
| Design graph or outbound service analysis | [Repository graph design](../architecture/repository-graph.md) | Load architecture §§4–7 for shared IR, flow and evidence policy |
| Design model assessment / review | [current provider decision](../architecture.md#current-provider-decision), architecture §8, [API setup](openai-decisions-setup.md) | Read historical [Jev research](../project/04-jev-validation-research.md) only for evaluation ideas or a provider comparison |
| Set up a key or diagnose access | [API setup](openai-decisions-setup.md) | Follow its official API links for current schema/access; never load or print secrets |
| Author the next specification | Near-term plan, [spec index](../specifications/README.md), `spec-builder`, [spec hygiene](../standards/spec-hygiene.md) | Resolve only decisions required by that contract; the old route spec does not authorize this branch's implementation |
| Implement or resume | Current authorized spec, its local run log, [execution discipline](agent-execution-discipline.md), `spec-driven` | Select relevant standards from the standards index; explicit user implementation requests are sufficient authorization under AGENTS.md |
| Reassess scope or long-term architecture | Project overview, architecture, scope, then proposal | Load functionality for capability details; prior-art/research docs for evidence behind a disputed choice |

## Product documentation versus SDLC documentation

Root `docs/` owns public philosophy, user/contributor setup, CLI operation, result interpretation,
and [scan analysis rules and guidelines](../../docs/analysis/README.md). These are product documents
used by people running scans and agents refining detection. They MUST NOT be placed under `.zazz/`.
Load the [database injection](../../docs/analysis/database-injection.md),
[input validation](../../docs/analysis/input-validation.md), and
[calibration](../../docs/analysis/calibration.md) guides when changing those analyses.

`.zazz/` owns SDLC/process standards, project plans, implementation contracts, internal architecture,
and evaluation history. `.zazz/ephemeral/` is ignored execution scratch and session handoff only.
Public rules describe the implemented analysis; current Markdown files are not dynamically loaded
by the binary. Changes MUST keep those rules, Go evidence/judgment code, and tests aligned.

## Documentation map

Select the document that owns the topic; do not load every row. Research/design context is preserved
in these files independently of chat attachments.

| Document | Primary content |
|---|---|
| [User documentation](../../docs/README.md) | Product introduction, agent-first philosophy, setup, use, result interpretation and contributor guides |
| [Project overview](../project.md) | Purpose, [public philosophy](../../docs/philosophy.md), current state, direction |
| [Architecture](../architecture.md) | Flow diagrams, judge/provider contracts, analysis semantics, package boundaries |
| [Repository graph design](../architecture/repository-graph.md) | Target-repository nodes/edges, inbound/outbound API boundaries, response flow, identity and provenance |
| [SQL/NoSQL judgment design](../architecture/query-judgment.md) | Automatic query claims, model contracts, conditional output facts and policy gates |
| [Evidence and claims](../architecture/evidence-and-claims.md) | Source-backed route slices, relevant functions, checked requirements and finding eligibility |
| [Reporting and SARIF](../architecture/reporting-and-sarif.md) | Route-oriented output, remediation options, internal model and implemented bounded SARIF export |
| [Standalone HTML](../architecture/html-reporting.md) | Core human-facing CLI report, native disclosure rows and JSON/SARIF relationship |
| [SCA design](../architecture/sca-design.md) | Follow-on ecosystem audit adapters and reachability correlation |
| [Functionality](../project/01-functionality.md) | Capabilities, detection modes, rule families, language sequence, interfaces |
| [Scope](../project/02-project-scope.md) | Boundaries, quality hypotheses, risks |
| [Design proposal](../project/03-jev-design-proposal.md) | Rationale, supplied concepts, graph/library choices, decision register |
| [Jev research](../project/04-jev-validation-research.md) | Capability evidence, limitations, evaluation protocol |
| [Prior art and feasibility](../project/05-prior-art-and-feasibility.md) | Novelty assessment, pinned Bearer study, engineering risks |
| [Implementation plan](../project/implementation-plan.md) | Near-term route/string sequence, then long-term phases and backlog |
| [Route validation assessment](../project/06-route-validation-assessment.md) | Verified Juice Shop evidence, gaps, architecture review, bounded next scope |
| [Registration validation case](../project/07-registration-validation-case.md) | Next field-constraint case: middleware/guard bypass, generated persistence, evidence and judge questions |
| [Login SQL/validation case](../project/08-login-validation-case.md) | Separate email constraints and query semantics; phase regression expectations |
| [Two-case results](../evaluation/juice-shop/results.md) | Known-target Luna outcomes, usage, offline checks and remaining limits |
| [Bearer comparison plan](../evaluation/bearer-comparison.md) | Repeat/extension protocol; links to pinned results and capability inventory |
| [Original comparison controls](../evaluation/corpus/README.md) | Label-isolated corpus, Luna outcomes, provenance correction and limitations |
| [Decision-model user guide](../../docs/decision-models.md) | Provider credentials, environment loading, capped scanner assessment |
| [Integration evidence](openai-decisions-setup.md) | Historical provider setup, transport and access smoke-check evidence |
| [Zazz methodology](zazz-methodology.md) | SDLC, ownership, greenlight and review lifecycle |
| [Specifications](../specifications/README.md) | Execution contracts and index |
| [Standards index](../standards/index.yaml) | Task-specific engineering standards routing |

## Working protocol

Load [Agent Execution Discipline](agent-execution-discipline.md) for execution, testing, verification,
resumption, or handoff; [AGENTS.md](../../AGENTS.md) supplies repo-specific overrides. For planning,
use the [phase-slicing guide](../project/implementation-plan.md#how-a-phase-becomes-agent-manageable-work).

## Methodology refinement

The Project Owner may assign Zazz refinements during definition. Load the
[methodology](zazz-methodology.md) for that work; record whether a change is local or proposed upstream
and preserve human authority. Assigned process changes do not authorize scanner implementation.

## Commands and credentials

Commands and contributor checks: [README](../../README.md#contributing). Current live-model setup:
[OpenAI Decisions setup](openai-decisions-setup.md). Offline extraction and fixture work need no key.
The decision-model key is already configured in the local ignored `.env`; agents MUST NOT print it.
The scanner does not auto-load that file; see the [README](../../README.md#optional-decision-model-setup).
API setup is separate from permission to transmit target source or execute a paid evaluation.
Live requests remain restricted to OpenAI Decisions with `gpt-6-luna`, small explicit request caps,
and no provider/model fallback. The previously reported $10 funding is context, not a spending authorization.
