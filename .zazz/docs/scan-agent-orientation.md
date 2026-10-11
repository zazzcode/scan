# Scan — agent orientation

First-read orientation for agents and developers working in this repository. It summarizes the repo
layout, the branch model, and where each kind of context lives. The operating rules are declared in
[`AGENTS.md`](../../AGENTS.md).

## What this project is

`scan` is an OWASP Top 10 static analysis scanner written in Go. It builds a structural graph of a
target repository with Tree-sitter, isolates the source-to-sink paths that carry untrusted input into
dangerous operations, and hands those pruned paths to an LLM broker for vulnerability analysis.

It is deliberately not a general codebase-understanding tool: every stage of the pipeline exists to
produce and narrow a security-relevant call path. See
[`../project/architecture.md`](../project/architecture.md) for the full design.

Initial release languages: **JavaScript (Node.js), TypeScript, Go, Python, Java**.

## Where things live

| Location | Contents |
|---|---|
| `.zazz/project.md` | Project overview — the durable orientation document |
| `.zazz/project/` | Functionality, scope, and architecture (the implementation overview) |
| `.zazz/standards/` | Engineering standards; route through `index.yaml` |
| `.zazz/specifications/` | Deliverable specifications (tracked in Git) |
| `.zazz/docs/` | Methodology guides (this directory) |
| `.zazz/ephemeral/` | Gitignored working notes, run logs, QA evidence, handoffs |
| `.agents/skills/` | Vendored Zazz skills; load only the one matching the task |

## Repo layout

`.bare/` is the shared Git directory; `stage/` is the integration worktree. Every other branch gets
its own sibling worktree. See [`worktree-setup.md`](worktree-setup.md) and
[`wt-cheat-sheet.md`](wt-cheat-sheet.md).

## Branch model

- **`stage`** — integration branch, and the GitHub default branch. All PRs target `stage`.
- **`main`** — release branch. Advances only by human-operated release promotion plus a version tag
  on `main`.
- Branch names are flat and hyphenated so branch name = worktree directory name.

## Working in this repo

1. Read `AGENTS.md` and this file.
2. For code work, load only the standards matched by `.zazz/standards/index.yaml`.
3. For product context, start with `.zazz/project.md`, then the relevant document under
   `.zazz/project/`.
4. For a tracked deliverable, use its specification under `.zazz/specifications/` as the execution
   contract.
5. Work in a sibling worktree on a flat branch name; open a PR into `stage`. Never merge locally.

## Repository commands

```bash
go build ./...
go test ./...
go vet ./...
gofmt -l -w .
```

## Current status

Bootstrap. The repository contains documentation, standards, and skills only — no Go packages yet.
The first deliverable is the core scanner pipeline described in
[`../project/architecture.md`](../project/architecture.md) and sequenced in
[`../project/02-project-scope.md`](../project/02-project-scope.md).
