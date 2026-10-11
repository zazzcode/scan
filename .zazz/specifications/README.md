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

## Deliverables

None yet. The first deliverable is the core pipeline skeleton, sequenced as step 1 in
[`../project/02-project-scope.md`](../project/02-project-scope.md) and designed in
[`../architecture.md`](../architecture.md).

Draft it with the `spec-builder` skill before implementation begins, then implement it in a sibling
worktree on a flat branch name (`deliverable-core-pipeline-skeleton`), opening a PR into `stage`.
