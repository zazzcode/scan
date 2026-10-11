# AGENTS.md — scan

Primary entry point for agents and developers working in this repo. Read
[`.zazz/docs/scan-agent-orientation.md`](.zazz/docs/scan-agent-orientation.md) for the full
orientation; this file declares the operating rules.

## Project

Scan — an OWASP Top 10 static analysis scanner. It builds a language-agnostic structural graph of a
target repository with Tree-sitter, isolates source-to-sink paths, and hands the pruned paths to an
LLM broker for vulnerability analysis.

- Product docs: `.zazz/project.md` (project overview) and `.zazz/project/` (functionality, scope, architecture)
- Stack: Go 1.27+ · Tree-sitter (multi-language parsing) · LLM broker over an OpenAI-compatible HTTP API
- Initial release languages: JavaScript (Node.js), TypeScript, Go, Python, Java
- Repository shape: **single Go module** (`github.com/zazzcode/scan`)

## Document model

- **DOCS_ROOT: `.zazz/`** (repo-relative)
- `.zazz/project.md` — the project overview and durable orientation document (read first)
- `.zazz/project/` — detailed product docs: functionality, scope, and `architecture.md`
- `.zazz/specifications/` — deliverable specifications (tracked in Git)
- `.zazz/standards/` — engineering standards; route through `.zazz/standards/index.yaml`
- `.zazz/docs/` — methodology guides; first-read: `scan-agent-orientation.md`
- `.zazz/ephemeral/` — gitignored working directory (run logs, QA notes, handoffs). Never commit it.

## Branch model

- **Integration branch: `stage`.** All PRs target `stage`. Keep the `stage/` worktree clean and
  read-only.
- **Release branch: `main`.** At release time a human merges `stage` → `main` and tags a version
  (`vX.Y.Z`) on `main`. Agents never merge or tag releases.
- **`main` receives changes only by release promotion** — never direct feature work, and no
  hotfix-to-`main` path: expedited fixes also flow through `stage`.
- Flat, hyphenated branch names only (`feature-taint-graph`, not `feature/taint-graph`) so branch
  name = sibling worktree directory name.
- Layout: `.bare/` container + sibling worktrees (`.zazz/docs/worktree-setup.md`,
  `.zazz/docs/wt-cheat-sheet.md`).

## Git identity

This tree lives under `/Users/michael/Dev/zazzcode`: commits are authored as **zazzcode** (via
`includeIf`) and `gh`/git HTTPS auth uses the zazzcode profile. Verify with `git config user.name`
before committing.

## Skills

Authoritative local skill set: `.agents/skills/` (vendored from zazz-skills + local). Load only the
skill matching the current task. Current set: `feature-doc-builder`, `spec-builder`, `spec-driven`,
`handoff`, `standard-builder`, `worktree`, `doc-check`, `pr-builder`, `pr-review`.

## Repository commands

- Build: `go build ./...`
- Test: `go test ./...`
- Vet: `go vet ./...`
- Format: `gofmt -l -w .`
- No required environment or secrets for the current bootstrap.

## Not in use

- **zazz-board** is not used in this project for now — no board mirroring or tracking; execution
  records stay in `.zazz/ephemeral/` or worktree-local files.

## Execution discipline

Follow `.zazz/docs/agent-execution-discipline.md`: work in isolated worktrees, verify before
assuming, scope edits to the branch's diff against `stage`, run checks on the task's file set, and
stop to ask when a destructive or ambiguous situation appears.
