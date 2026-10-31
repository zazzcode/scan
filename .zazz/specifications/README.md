# Specifications

Deliverable specifications — the bounded execution contracts that agent sessions implement.

- authored with the `spec-builder` skill (see `.agents/skills/spec-builder/`)
- quality bar: `../standards/spec-hygiene.md`
- explicit user implementation requests authorize work; use `.agents/skills/spec-driven/` for evidence
  and review discipline, with the local authorization rules taking precedence

## Storage policy

Specifications are **tracked in Git** in this repo (Git-native audit trail). They are committed to the
deliverable's branch and reviewed through the same PR flow as code.

Working notes, run logs, QA findings, and handoffs do **not** live here — they go to
`.zazz/ephemeral/` (gitignored) or a worktree-local location.

## Specification index

Use the assigned feature worktree for related contract, design, implementation and tests. Explicit
user implementation requests authorize work under [AGENTS.md](../../AGENTS.md#authorization-and-delivery).
Independent future deliverables may use separate feature worktrees/PRs targeting `stage`.

- [Phase 0 skeleton](phase0-skeleton/README.md) — historical foundation contract, retained for context.
  The revised [Phase 0 plan](../project/implementation-plan.md#phase-0--foundation-audit-and-tooling-decision)
  audits existing artifacts before further work.
- [Tree-sitter route identification MVP](tree-sitter-route-identification-mvp.md) — greenlit
  2026-10-24 and implemented; syntax inventory only. Reverified on Juice Shop in the
  [route assessment](../project/06-route-validation-assessment.md#route-verification).
- [Route string validation](route-string-validation.md) — implementation explicitly authorized by
  the user on 2026-10-31 in the assigned feature worktree.
  Start from the [near-term sequence](../project/implementation-plan.md#route-string-validation-delivery-sequence).
- Remaining specs are not authored yet. Their proposed slices live in the
  [implementation plan](../project/implementation-plan.md); each needs a bounded contract and user authorization.

See [repository state](../project.md#concept-stage-repository-state) for the current concept stage and
[specification-driven delivery](../docs/zazz-methodology.md#specification-driven-delivery-cycle) for
ownership, approval, execution, and review requirements. Earlier LLM-era phase numbering is superseded
by the current plan.
