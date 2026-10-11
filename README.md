# scan

`scan` is an OWASP Top 10 static analysis scanner. It walks a target repository, parses each source
file with Tree-sitter, resolves a language-agnostic structural graph (imports, declarations, call
sites), identifies **sources** (request entry points) and **sinks** (database, shell, filesystem,
network, and eval operations), prunes the graph down to the source-to-sink paths that matter, and
hands those paths to an LLM broker for vulnerability analysis.

The design follows the same split Graphify uses — deterministic structural extraction, then an
inference/evaluation pass — but is purpose-built for security tracing rather than general codebase
understanding.

- **Language:** Go 1.27+
- **Parsing:** Tree-sitter (multi-language; normalized into one graph schema)
- **Initial release languages:** JavaScript (Node.js), TypeScript, Go, Python, Java
- **Analysis:** LLM broker over an OpenAI-compatible HTTP API
- **Status:** bootstrap — repository scaffolding only, no implementation yet

## Documentation

The docs root is `.zazz/`.

- [`.zazz/project.md`](.zazz/project.md) — project overview — **start here**
- [`.zazz/docs/scan-agent-orientation.md`](.zazz/docs/scan-agent-orientation.md) — agent orientation and routing
- [`.zazz/project/`](.zazz/project/01-functionality.md) — functionality, scope, and architecture
- [`.zazz/standards/index.yaml`](.zazz/standards/index.yaml) — engineering standards routing table
- [`.zazz/specifications/`](.zazz/specifications/README.md) — deliverable specifications

## Repository layout

This repository uses the Zazz bare-repo container model:

```text
scan/
├── .bare/    # shared Git directory
└── stage/    # integration worktree (branch: stage)
```

- **`stage` is the integration branch** and the GitHub default branch. All PRs target `stage`.
- **`main` is the release branch.** It advances only by release promotion (`stage` → `main` plus a
  version tag on `main`), which is human-operated. Never commit to `main` directly.
- Branch names are flat and hyphenated (`feature-taint-graph`, not `feature/taint-graph`) so the
  branch name maps directly to a sibling worktree directory.
- Create new worktrees as siblings of `stage/`. See
  [`.zazz/docs/worktree-setup.md`](.zazz/docs/worktree-setup.md).

## Commands

```bash
go build ./...     # build all packages
go test ./...      # run all tests
go vet ./...       # static checks
gofmt -l -w .      # format
```

No environment variables or secrets are required for the current bootstrap.
