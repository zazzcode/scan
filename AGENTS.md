# AGENTS.md — scan

Read this file and the [agent orientation](.zazz/docs/scan-agent-orientation.md) at session start.
Then load only the supporting context required for the task; links are not automatic includes.

## Load context by task

| Task | Load |
|---|---|
| Product, scope, or architecture decisions | Orientation's [documentation map](.zazz/docs/scan-agent-orientation.md#documentation-map); select the relevant product/design doc |
| Author or refine a spec | Assigned phase/spec in the [plan](.zazz/project/implementation-plan.md), `spec-builder`, and [spec hygiene](.zazz/standards/spec-hygiene.md) |
| Implement, test, resume, or verify work | Active greenlit spec/run log, [execution discipline](.zazz/docs/agent-execution-discipline.md), `spec-driven`, and matched [standards](.zazz/standards/index.yaml) |
| Create/manage a worktree | [Worktree setup](.zazz/docs/worktree-setup.md); [Worktrunk cheat sheet](.zazz/docs/wt-cheat-sheet.md) only when using Worktrunk |
| Refine SDLC or ownership/process rules | [Zazz methodology](.zazz/docs/zazz-methodology.md) |
| Edit documentation or standards | Matched [standards](.zazz/standards/index.yaml) and `doc-check` |
| Prepare or review a PR | [PR process](.zazz/standards/pr-process.md) and the matching `pr-builder`/`pr-review` skill; [human-review guide](.zazz/docs/human-in-loop-pr-review-strategy.md) only for review policy |

Skills live in `.agents/skills/`; load only the matching skill. Keep this entry point small and place
procedures in their owning guides.

## Branch model

`stage` integrates reviewed work; `main` receives human-operated release promotion and version tags.
All fixes flow through `stage`. Agents MUST NOT merge, tag releases, or push directly to integration.
Current exception: assigned concept/architecture/methodology documentation may be edited in `stage/`;
this does not permit spec work or feature implementation there. Detailed layout, naming, and operations:
[worktree setup](.zazz/docs/worktree-setup.md).

## Definition and spec-driven delivery

Each spec MUST be authored, refined, and implemented in its own dedicated feature worktree from
`stage`, with its own draft PR targeting `stage`. Implementation requires human Deliverable Owner
greenlight. A proposed phase/backlog entry is not approval. Follow [execution discipline](.zazz/docs/agent-execution-discipline.md)
with the scan-specific branch/spec rules above taking precedence.

## Document model

**DOCS_ROOT: `.zazz/`**. Specs are Git-tracked in `.zazz/specifications/`; execution scratch belongs in
gitignored `.zazz/ephemeral/`, never committed. Zazz Board is not used. See the
[documentation map](.zazz/docs/scan-agent-orientation.md#documentation-map) for durable document owners.

## Repository commands

Run in the active worktree using its declared Go toolchain and the assigned spec's scoped checks:

- List inputs: `go run ./cmd/scan --list .` (add `--json` before the path for JSON).
- Build/test/vet: `go build ./...`, `go test ./...`, `go vet ./...`.
- Format check: `gofmt -l .`; format only changed Go files when needed.

Before committing/pushing, verify path/branch, `git config user.name` (**zazzcode**), and the zazzcode
Git/`gh` authentication profile. Test/evidence conventions live in [testing](.zazz/standards/testing.md).
