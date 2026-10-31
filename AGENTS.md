# AGENTS.md — scan

Read this file and the [agent orientation](.zazz/docs/scan-agent-orientation.md) at session start.
Then load only the supporting context required for the task; links are not automatic includes.

## Load context by task

| Task | Load |
|---|---|
| Verify routes or plan string validation | Orientation's [current work](.zazz/docs/scan-agent-orientation.md#current-work) and [task reading paths](.zazz/docs/scan-agent-orientation.md#task-reading-paths) |
| Configure the GPT decision API | [OpenAI Decisions setup](.zazz/docs/openai-decisions-setup.md); no key is needed for offline checks |
| Graph structure / outbound service paths | [Repository graph design](.zazz/architecture/repository-graph.md), then relevant architecture sections |
| Product, scope, or architecture decisions | Orientation's [documentation map](.zazz/docs/scan-agent-orientation.md#documentation-map); select the relevant product/design doc |
| Author or refine a spec | Assigned phase/spec in the [plan](.zazz/project/implementation-plan.md), `spec-builder`, and [spec hygiene](.zazz/standards/spec-hygiene.md) |
| Implement, test, resume, or verify work | Active implementation contract/run log, [execution discipline](.zazz/docs/agent-execution-discipline.md), `spec-driven`, and matched [standards](.zazz/standards/index.yaml) |
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

## Authorization and delivery

An explicit user request to implement authorizes the requested work in the assigned isolated
feature worktree. Agents MUST NOT require a separate Deliverable Owner ceremony or repeated
confirmation. Record the bounded scope and acceptance evidence in a tracked implementation contract;
refine it in place as the user steers the work. A backlog entry alone is not an implementation request.
The human user owns scope and review decisions; this is the meaning of Deliverable Owner in older guides.

Use the assigned worktree for related architecture, documentation, implementation, and tests. Split
independent future deliverables when useful; do not create extra worktrees or PRs merely for the
steps of the active task. Existing human merge/release boundaries still apply. Repository-local rules
in this section take precedence over stricter greenlight/topology wording in the methodology,
execution discipline, and skills. Current authorization: implement `route-string-validation` here,
per the user's 2026-10-31 instructions.

## Document model

**DOCS_ROOT: `.zazz/`**. Specs are Git-tracked in `.zazz/specifications/`; execution scratch belongs in
gitignored `.zazz/ephemeral/`, never committed. Zazz Board is not used. See the
[documentation map](.zazz/docs/scan-agent-orientation.md#documentation-map) for durable document owners.

## Live decision-model usage for this branch

Current user restriction (2026-10-31): live judge calls MUST use OpenAI Decisions with
`gpt-6-luna` only. Keep request counts small and explicit; no automatic provider/model fallback.
The user initially funded $10; that is context, not authorization to exhaust the balance.
Alternative providers remain part of the design and may be tested with local HTTP stubs.
Never print the `.env` contents or credentials; pass the selected key only to the child process.

## Repository commands

Run in the active worktree using its declared Go toolchain and the assigned spec's scoped checks:

- List inputs: `go run ./cmd/scan --list .` (add `--json` before the path for JSON).
- Route inventory: `go run ./cmd/scan --routes --json <target-root>`; matches are syntactic candidates.
  Verification scope and limitations: [route assessment](.zazz/project/06-route-validation-assessment.md).
- Build/test/vet: `go build ./...`, `go test ./...`, `go vet ./...`.
- Format check: `gofmt -l .`; format only changed Go files when needed.

Before committing/pushing, verify path/branch, `git config user.name` (**zazzcode**), and the zazzcode
Git/`gh` authentication profile. Test/evidence conventions live in [testing](.zazz/standards/testing.md).
