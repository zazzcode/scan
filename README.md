# scan

An OWASP Top 10 code scanner built on a **graph-narrowed LLM review** model. It tracks code paths across
a repository with Tree-sitter and hands the paths that matter to a model for judgment.

> **Status: early.** Phase 0 — the skeleton, walker, language registry, and CI — is in place. No parsing,
> graph, rules, or model calls yet.

## The idea

Traditional SAST builds a source-to-sink dataflow engine: hand-coded taint propagation through
parameters, returns, fields, and aliases. That is years of engineering, it is brittle across frameworks,
and it still produces findings nobody can act on.

`scan` takes a different route. It builds a deterministic structural graph of the target codebase with
Tree-sitter — the Graphify model — and tracks the code paths through it. A deterministic rule tier handles
everything decidable by shape. Everything requiring judgment goes to the model.

## Premise

**A capable LLM reasons about code better than a static analyzer, and handles the cases that break static
analysis** — custom validators, framework indirection, dynamic dispatch, authorization and logic flaws,
and whether a mitigation is *sufficient*. Those are semantic questions, and they are where SAST is
structurally blind or noisy.

## Code path tracking is the spine

Everything depends on resolving what calls what across files. It is what makes the model stage tractable:

- **It reduces what the model must review.** A tracked path hands it exactly the code on that path, in
  order, with the guard and sink already identified — instead of a repository and a search problem.
- **It provides structure for analysis.** A path is an explicit claim: this entry point reaches this
  operation, through these functions, with these guards, passing these arguments. The model reasons over
  that structure, and both the model and a human reviewer can cite it.

## Three detection tracks

| Track | Detection | Deterministic? |
|---|---|---|
| **A — pattern rules** | A rule matches a known-bad shape; the rule *is* the verdict | Yes, end to end. No model call |
| **B — path rules** | A rule finds a source-to-sink path; the model judges whether the data is sufficiently validated | Discovery yes; verdict no |
| **C — scope review** | The model reviews a scope — a route file, a controller, a middleware chain — for authorization and logic gaps | Discovery yes; findings no |

Every finding records which track produced it, so a consumer can gate hard on Track A and treat
low-confidence Track B and C findings as advisory.

**`scan` is useful with no API key configured.** Track A and the full deterministic pipeline still work.

## Language scope

TypeScript/JavaScript first, then Python, then Go and Java. The core is language-neutral: adding a
language must mean adding a grammar, query sets, and a resolver — never changing the core.

Two TypeScript acceptance fixtures — one over MongoDB, one over SQL — exist to test a hypothesis: SQL
paths should mostly settle deterministically because bound parameters are near-decisive, while NoSQL paths
should mostly escalate to Track B because Mongo has no equivalent escape hatch.

## Positioning

A semantic review tool, complementary to a dataflow SAST rather than a replacement for it. A dataflow
analyzer over-reports inside the scope it models; `scan` reaches further semantically and under-reports
outside what its graph and scopes cover. Different error profiles — run both.

It is not a whole-program analyzer, so: a clean run is not a security guarantee. CI gating on new findings
is meaningful; "no findings" is not a certification.

## Documentation

- [`.zazz/project.md`](.zazz/project.md) — project overview — **start here**
- [`.zazz/architecture.md`](.zazz/architecture.md) — tech stack and technical approach
- [`.zazz/project/implementation-plan.md`](.zazz/project/implementation-plan.md) — phased plan to the final solution
- [`.zazz/project/01-functionality.md`](.zazz/project/01-functionality.md) — capabilities and coverage
- [`.zazz/project/02-project-scope.md`](.zazz/project/02-project-scope.md) — boundaries, sequencing, targets
- [`.zazz/standards/index.yaml`](.zazz/standards/index.yaml) — engineering standards routing table
- [`.zazz/specifications/README.md`](.zazz/specifications/README.md) — deliverable specifications

## Repository layout

This repository uses the Zazz bare-repo container model:

```text
scan/
├── .bare/    # shared Git directory
└── stage/    # integration worktree (branch: stage)
```

- **`stage`** is the integration branch and the GitHub default branch. PRs target `stage`.
- **`main`** is the release branch. It advances only by release promotion plus a version tag, which is
  human-operated. Never commit to `main` directly.
- Branch names are flat and hyphenated, so the branch name maps to a sibling worktree directory.

## Commands

```bash
go build ./...                       # build
go test ./...                        # test
go test -race ./...                  # test with the race detector
go vet ./...                         # static checks
gofmt -l -w .                        # format

go run ./cmd/scan --list .           # file and language manifest
go run ./cmd/scan --list --json .    # machine-readable manifest
```

Exit codes: `0` success, `1` error, `2` usage.
