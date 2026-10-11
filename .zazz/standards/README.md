---
last_updated_at: 2026-10-11
---

# Scan Standards

Engineering standards for this repository. Agents route through
[`index.yaml`](index.yaml): match the task's paths and activities, then load only the matched documents.

## Curation policy

Two tiers exist (see [`contextual-split.md`](contextual-split.md) for the model):

1. **Vendored methodology standards** — stack-agnostic, adapted from
   [zazz-skills](https://github.com/zazzcode/zazz-skills) `docs/standards/`. Divergence is deliberate and
   recorded below.
2. **Scan stack standards** — authored **for this repo** with the `standard-builder` skill against real
   code. `code-structure.md` and `testing.md` landed with the bootstrap; the rest of the backlog is listed
   at the bottom of `index.yaml`.

The upstream stack-specific baselines (Python/Flask, SQL Server stored procedures, React/Next, AWS Lambda,
Drizzle/Neon) are deliberately **not** vendored — they describe a stack this repository does not use. This
is a single Go module with no web frontend, no SQL Server, and no serverless deployment.

## Local divergence from upstream

Recorded so a future upstream sync does not blindly overwrite them:

- **`contextual-split.md`** — rewritten. Upstream's copy is a catalog of the upstream repo's own standards
  (HTTP layer, stored procedures, data layer, service layer, reports, frontend, Python testing, tooling).
  That catalog describes another repository entirely, so the file was rewritten around scan's standards and
  backlog.
- **`pr-process.md`** — rewritten. Upstream's copy is roughly two-thirds specific to a monorepo with a
  frontend and a backend: `[FE]`/`[BE]`/`[DB]` scope labels, a legacy-screen-screenshot PR convention, and
  a root-only `.envrc`/`yamlfmt`/pre-commit section. The portable rules (bracketed scope labels, one
  logical change per PR, no transitional cruft, CVE identifiers in titles, draft-first human merge) were
  kept, and the labels now map to scan's package layout.
- **`docs-hygiene.md`** — one upstream Python example remains in the "working-doc identifiers" section; the
  substance is stack-agnostic and the example is illustrative. Scheduled to be genericized.
- **`docs-hygiene-reference-structure.md`** — the citation and section-ordering rules are kept; some
  upstream example paths remain and are scheduled to be genericized.

## Adding a standard

1. Run the `standard-builder` skill against actual code in the repo — standards describe real patterns, not
   aspirations.
2. Add the file to this directory.
3. Register it in `index.yaml` with `applies_to` paths/activities and a purpose.
4. Keep each standard scoped to one concern; split when it exceeds the size thresholds in
   [`code-structure.md`](code-structure.md).
