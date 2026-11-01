# Human-readable scan reports

## Outcome and authorization

The user requested clearer input-validation and database-sink results, fresh offline scans of Juice
Shop and the vulnerable TypeScript app, and an expandable standalone HTML report generated from saved
JSON. Implement in the assigned `logger-leak` worktree before logger-leak detection.

## Scope and decisions

- Add optional `findings` and report-presentation version fields to schema-v1 evidence JSON. Preserve
  raw entries, uses, gaps, and assessments. Human findings reference the originating use and route IDs.
- Findings name observed behavior in plain English. A query-shape match alone is an analysis item,
  not a confirmed injection. Separate input-validation evidence from query protection and distinguish
  observations, incomplete analysis, supported potential findings, and scoped protection.
- Preserve the existing detector and judge eligibility; do not promote findings based on display text.
- SARIF remains the standard interchange output and preserves human context in namespaced properties.
- The HTML contains only scan results; tool comparisons stay in evaluation documents.
- Use NoSQL consistently; comparator cells say “NoSQL: no finding reported.”
- Provide `--export-html <saved-json>` using Go `html/template`, embedded CSS, and native details/summary
  rows. No JavaScript, CDN, browser JSON fetch, or server. The HTML is regenerated from saved JSON and
  lives beside it; regeneration must require no target-source access or model call.
- Rows group by source entry, showing source type, method/path signature, location, counts, and analysis
  state. Expand for human findings, source excerpts, input/transfer/sink evidence, assessment state,
  limitations, and remediation. Retain entries without findings and unmatched uses without inventing routes.
- `--strings --output-dir <base>` writes a fresh private `<base>/<repo-name>/run-<id>/` folder with all three
  artifacts. `SCAN_FINDINGS_DIR` supplies the default base, overridden by the flag; an explicitly empty
  flag disables file output. Export/list modes ignore the environment variable. Explicit output-dir
  is accepted only for strings scans. Existing stdout JSON/SARIF modes remain machine-readable;
  text mode prints bundle paths. Requested-assessment failures retain the bundle and exit 1.
- Existing saved schema-v1 reports without presentation fields remain exportable. Invalid/schema-wrong
  input is rejected. Unknown properties remain compatible. Export must not mutate saved assessments.
- Update README compile/run/results guidance, comparison prose, and embedded JS/TS practice rules.
  No new detectors, paid calls, target execution, or target edits in this increment.

## Acceptance and verification

- Tests distinguish interpolated text, query documents, unresolved calls, validation observations,
  protected claims, and failed/not-run assessments without inventing confirmed vulnerabilities.
- HTML escapes hostile source/path text, preserves provenance and remediation, groups by route ID,
  retains zero-finding rows and orphan uses, and has no scripts or network resources.
- CLI tests cover valid saved JSON, invalid input, output modes, and conflicting flags.
- Test output-dir precedence, disabled file output, separate successive runs, JSON/SARIF validity,
  private file permissions, unavailable paths, and the retained exit-1 bundle after assessment failure.
- Rerun both pinned targets offline; compare detector use/gap totals against baseline, validate finding
  references, and export JSON/HTML pairs. Inspect representative expanded rows in a browser.
- Run Go tests with race detection, vet, build, format, and documentation checks. Use scoped tests
  during development. Evidence and generated reports live under ignored `.zazz/ephemeral/`.

## Execution rules

Use [AGENTS.md](../../AGENTS.md), [testing](../standards/testing.md),
[code structure](../standards/code-structure.md), and [docs hygiene](../standards/docs-hygiene.md).
Work serially in this checkout, preserve existing edits, and record evidence in
`.zazz/ephemeral/human-reports-run-log.md`. Stop for unresolved evidence or scope changes; do not
fabricate security judgments to complete the presentation. Human UX feedback guides the next iteration.
