# Contributing to scan

[Documentation index](README.md) · [Project philosophy](philosophy.md)

## Contributing

Use the same Go and C-toolchain prerequisites as user setup. Development uses a shared bare Git
repository with sibling feature worktrees; follow the [worktree setup guide](../.zazz/docs/worktree-setup.md)
for a new checkout or conversion. Work in the assigned feature worktree and verify the location,
branch, and existing changes with `pwd` and `git status --short --branch` before editing.
`stage` is the integration branch; `main` is reserved for human-operated releases.

Read [AGENTS.md](../AGENTS.md) for agent rules and the
[orientation](../.zazz/docs/scan-agent-orientation.md) for documentation routing. Agree on bounded scope
and acceptance criteria before implementation, keeping the contract in `.zazz/specifications/`.
Put transient run logs and handoffs in ignored `.zazz/ephemeral/`. Keep changes focused and preserve
unrelated work already in the checkout.

Run the relevant scoped checks while developing. The Go CI checks are:

```bash
gofmt -l .
go vet ./...
go build ./...
go test -race ./...
```

`gofmt -l .` should print nothing; format only changed Go files. Use `go test ./...` for a quick local
pass. Tests use local fixtures and HTTP stubs, so no decision-model key is required. For documentation
changes, use the repository's [doc-check workflow](../.agents/skills/doc-check/SKILL.md).

Submit focused draft PRs into `stage`, with scoped titles and validation evidence following the
[PR process](../.zazz/standards/pr-process.md). Human reviewers approve and merge; agents do not merge
or promote releases. The maintained checkout uses the `zazzcode` Git identity and authentication
profile; agents verify them before committing or pushing.


## Documentation ownership

The repository README introduces the product. Root `docs/` owns user and contributor guides,
remediation philosophy, and report interpretation. `.zazz/` owns engineering design, project plans,
contracts, evaluation records, standards, and SDLC processes. New user instructions belong here;
link to internal engineering evidence only when it helps contributors. Keep procedural details in
their owning guide and the README a useful standalone introduction.
