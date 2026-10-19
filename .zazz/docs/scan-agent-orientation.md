# Scan — agent orientation

Compact session context and document routing. Read [AGENTS.md](../../AGENTS.md) for binding rules and
when to load supporting guides. Keep this document under 400 lines; link to details rather than copy them.

## Table of contents

- [Current work](#current-work)
- [Portable workspace and DEV_ROOT](#portable-workspace-and-dev_root)
- [Branches and worktrees](#branches-and-worktrees)
- [Documentation map](#documentation-map)
- [Working protocol](#working-protocol)
- [Methodology refinement](#methodology-refinement)
- [Commands and credentials](#commands-and-credentials)

## Current work

Concept/architecture definition is assigned to `stage/`. See the
[project overview](../project.md#concept-stage-repository-state) for status and the
[architecture](../architecture.md) for Jev Assessment → Findings Review.

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
release-only. Each spec has its own feature worktree and PR into `stage`, including spec authoring.
Load [worktree setup](worktree-setup.md) for the layout, reference-clone handling, naming, creation,
and recovery details; [AGENTS.md](../../AGENTS.md#branch-model) declares the current `stage/` exception.

## Documentation map

Select the document that owns the topic; do not load every row. Research/design context is preserved
in these files independently of chat attachments.

| Document | Primary content |
|---|---|
| [Project overview](../project.md) | Purpose, users, current state, direction |
| [Architecture](../architecture.md) | Flow diagrams, analysis semantics, stage contracts, package boundaries |
| [Functionality](../project/01-functionality.md) | Capabilities, detection modes, rule families, language sequence, interfaces |
| [Scope](../project/02-project-scope.md) | Boundaries, quality hypotheses, risks |
| [Design proposal](../project/03-jev-design-proposal.md) | Rationale, supplied concepts, graph/library choices, decision register |
| [Jev research](../project/04-jev-validation-research.md) | Capability evidence, limitations, evaluation protocol |
| [Prior art and feasibility](../project/05-prior-art-and-feasibility.md) | Novelty assessment, pinned Bearer study, engineering risks |
| [Implementation plan](../project/implementation-plan.md) | Phases, ordered spec backlog, prerequisites, acceptance targets |
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

Commands: [AGENTS.md](../../AGENTS.md#repository-commands). Live-model setup/prerequisites:
[Jev research](../project/04-jev-validation-research.md#what-can-happen-before-a-key-exists).
