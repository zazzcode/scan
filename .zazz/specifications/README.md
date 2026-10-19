# Specifications

Deliverable specifications — the bounded execution contracts that agent sessions implement.

- authored with the `spec-builder` skill (see `.agents/skills/spec-builder/`)
- quality bar: `../standards/spec-hygiene.md`
- after greenlight, the implementation lifecycle runs per `.agents/skills/spec-driven/`

## Storage policy

Specifications are **tracked in Git** in this repo (Git-native audit trail). They are committed to the
deliverable's branch and reviewed through the same PR flow as code.

Working notes, run logs, QA findings, and handoffs do **not** live here — they go to
`.zazz/ephemeral/` (gitignored) or a worktree-local location.

## Specification index

Every specification is authored, refined, and implemented in its own feature worktree from `stage`,
with its own draft PR targeting `stage`. The current definition session does not create spec contracts.

- [Phase 0 skeleton](phase0-skeleton/README.md) — historical foundation contract, retained for context.
  The revised [Phase 0 plan](../project/implementation-plan.md#phase-0--foundation-audit-and-tooling-decision)
  audits existing artifacts before further work.
- Remaining specs are not authored yet. Their proposed slices live in the
  [implementation plan](../project/implementation-plan.md); each needs separate human greenlight.

See [repository state](../project.md#concept-stage-repository-state) for the current concept stage and
[specification-driven delivery](../docs/zazz-methodology.md#specification-driven-delivery-cycle) for
ownership, approval, execution, and review requirements. Earlier LLM-era phase numbering is superseded
by the current plan.
