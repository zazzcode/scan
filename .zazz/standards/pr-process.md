---
last_updated_at: 2026-10-11
---

# PR process

This standard governs how PRs are titled, scoped, structured, and merged in this repository. It covers
the bracketed scope-label format used in every PR title, the one-logical-change-per-PR rule, the
prohibition on transitional cruft that lands dead on merge, and the CVE-prefix convention for security
remediations.

## PR titles — bracketed scope labels

Every PR title begins with one or more bracketed scope labels, before any descriptive text. Labels map to
the package layout in [`../architecture.md`](../architecture.md) § 18, so a reader can tell which part of
the tool a PR touches before opening it.

| Label | Scope |
|---|---|
| `CLI` | `cmd/scan` — flags, output formats, exit codes |
| `PARSE` | `internal/parse`, `internal/lang` — Tree-sitter integration, grammars, parsing |
| `GRAPH` | `internal/graph` — node and edge model, resolver, argument bindings |
| `RULES` | `internal/rules`, `internal/patterns` — rule engine, pattern rules, classification |
| `TRACE` | `internal/trace` — code path tracking, pruning, pre-filter, prioritization |
| `SCOPES` | `internal/scopes` — Track C scope selection and bundles |
| `LLM` | `internal/llm`, `internal/context` — broker client, prompts, bundle construction |
| `FIND` | `internal/findings` — identity, ledger, baseline, diffing |
| `REPORT` | `internal/report` — JSON and Markdown emitters |
| `CI` | GitHub Actions, tooling, and format gates |
| `DOC` | Documentation, standards, and specifications |
| `AGENT` | Agent skills and agent configuration under `.agents/` |

Cross-cutting PRs combine labels with `+` and **no spaces** inside the brackets. The canonical form is
`[GRAPH+TRACE]`, never `[GRAPH + TRACE]` or `[GRAPH] + [TRACE]`.

A label that maps to a package does not mean every change in that package needs every label — use the
labels for the packages the PR actually changes.

### Desired ✅

```text
[PARSE] Add TypeScript grammar and extension registry
[GRAPH] Resolve barrel re-exports through index.ts
[TRACE] Bound path search depth and report dropped candidates
[RULES] Add NoSQL operator-injection pattern rule
[LLM] Add three-state verdict response schema
[FIND] Fingerprint findings from rule id and symbol chain
[CI] Add race-detector test job
[DOC] Add path-tracking standard and register in index
[GRAPH+TRACE] Add argument bindings and consume them in path search
[CLI+REPORT] Add --scopes and Markdown scope findings
```

### Not desired ❌

```text
Add TypeScript grammar
# missing scope label entirely

[GRAPH] + [TRACE] Add argument bindings
# split brackets; canonical form is [GRAPH+TRACE]

[GRAPH + TRACE] Add argument bindings
# spaces around +; canonical form is [GRAPH+TRACE]
```

## PR scope — one logical change per PR

A single PR scopes one logical change: one feature, one fix, one refactor, or one migration. Incidental
fixes discovered while working on a feature go in their own one-liner PR, or at minimum an isolated commit
on a separate PR. Drive-by edits ("while I'm here") go in their own PR.

Multi-scope PRs are reserved for changes that are genuinely cross-cutting — a graph change plus the
tracker code that consumes it — not for bundling unrelated work behind a multi-scope label. The label
declares cross-cutting *coupling*, not "I happened to touch all of these this week."

The rule extends to commits inside a PR: when a feature commit unavoidably includes an unrelated fix, the
fix lands in its own commit with a clear message so `git revert` and `git blame` work cleanly.

### Desired ✅

```text
[TRACE] Fix depth bound excluding sinks at the boundary
# one-liner extracted from a feature PR

[GRAPH+TRACE] Track paths across resolved barrel imports
# feature only, no bundled refactors or drive-by fixes
```

### Not desired ❌

```text
[GRAPH+TRACE] Track barrel-import paths + rename Node.Role + reformat package
# three logically distinct changes bundled into one PR
```

## No transitional cruft on merge

Transitional types, parameters, and interfaces added "for the migration" are removed in the same PR if
they offer no value after merge. If the only reason the code exists is "support for the previous
signature," update the callers in the same PR and delete the bridge.

Concrete patterns to remove before merge:

- An unused parameter retained "for compatibility" with a previous signature.
- An exported symbol kept only because one caller has not been updated yet.
- A deprecated struct field the new implementation never reads.
- A legacy adapter type whose only consumer is the code this PR replaces.

### Desired ✅

```go
// walk.Walk emits jobs on the channel; the consumer owns closing it.
func Walk(root string, jobs chan<- Job) error {
	// ...
}
```

### Not desired ❌

```go
// _ignorePaths is retained "for compatibility" with the previous signature.
// The new implementation never reads it; it becomes cruft the moment this
// PR merges.
func Walk(root string, _ignorePaths []string, jobs chan<- Job) error {
	// ...
}
```

```go
// LegacyNode is kept only so the old builder still compiles.
type LegacyNode struct {
	// shape matches the pre-refactor model, not the new one
}
```

## CVE remediation PR titles

PR titles that remediate a published CVE include the CVE identifier in the descriptive portion of the
title, after the scope label, in the form `[<scope>] <CVE-ID>: <human description>`. The identifier in the
title is what makes PR-list grep, dashboards, and post-incident audits work — burying it in the body
breaks all three.

A CVE in a Go dependency is a dependency bump, which is out of scope for this tool (see
[`../project/02-project-scope.md`](../project/02-project-scope.md) § 3). This convention applies to CVEs in
our own code or in the Tree-sitter bindings we vendor.

### Desired ✅

```text
[PARSE] CVE-2026-12345: bump tree-sitter runtime past parser overflow
[CI] CVE-2026-99999: pin actions/checkout to a patched release
```

### Not desired ❌

```text
[PARSE] Bump tree-sitter
# no CVE identifier; remediation context is invisible in PR lists
# and the changelog
```

## Draft-first, human merge

- Open PRs as drafts. A draft PR is the normal working state, and CI runs on it.
- PRs target `stage`. Never `main`.
- **Agents do not merge PRs.** A human reviewer approves and merges.
- Do not merge locally into `stage`; the PR is the integration path.
- Release promotion (`stage` → `main` plus a version tag) is human-operated.

## Related standards

- [`code-structure.md`](code-structure.md) — the package layout the scope labels map to.
- [`spec-hygiene.md`](spec-hygiene.md) — the quality bar for the specification a deliverable PR
  implements.
- [`docs-hygiene.md`](docs-hygiene.md) — `[DOC]`-scoped conventions for standards documents and READMEs.
