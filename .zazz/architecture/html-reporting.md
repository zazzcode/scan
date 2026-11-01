# Standalone HTML reports for CLI runs

## Purpose and ownership

HTML is scan's core human-facing report for CLI runs. The canonical JSON is the richer agent-readable
artifact; SARIF remains the standard interchange output. See [reporting and SARIF](reporting-and-sarif.md)
for shared statuses and evidence ownership. HTML contains scan's own results only. Comparisons to
other tools belong in evaluation documents, never in the product report.

## Data path

```text
Source analysis → raw uses / assessments → human-readable finding records → canonical JSON
                                                                         ├→ standalone HTML
                                                                         └→ SARIF 2.1.0
```

`--strings --json` includes additive human-readable `findings` and `presentation_version` fields.
`--export-html <saved.json>` reads a schema-v1 string report and renders HTML offline; no target-source
read, key, model request, application execution, or server is needed. `--export-sarif` uses the same
saved report. `--strings --output-dir <base>` or `SCAN_FINDINGS_DIR` writes all three formats under
`<base>/<repo-name>/run-<id>/`. The explicit flag overrides the environment; an empty flag disables file
output. Repository names come from the target basename; clones with equal basenames share a repository
folder but never a run folder. Run directories use private permissions, and output files are created
exclusively with mode 0600. Existing runs are never overwritten. The output base is configurable,
without a machine-specific default. The local workspace convention is a sibling `_test-finding/`.

The CLI keeps stdout JSON/SARIF redirectable while reporting the bundle path on stderr. Text mode
prints bundle paths instead of the full text report when file output is active. Without file output,
existing stdout behavior remains available. Export and list modes ignore the output environment
variable; an explicit output-dir flag is valid only for strings scans. `.env` is not auto-loaded.
The user places regenerated JSON, HTML, and SARIF in the same run folder. Each generated view is a snapshot and must be regenerated after JSON changes.

The exporter rebuilds deterministic presentation from saved raw uses, supporting older schema-v1
reports without human fields. It preserves raw assessments and blockers; it does not invent a model
answer or reclassify a finding as confirmed. `internal/report/findings.go` is the common presentation
rule layer. Driver-specific remediation comes from `guidance.go`; illustrative examples are not patches.
The CLI rejects invalid JSON, unsupported schema versions, route-only reports, and incompatible flags.
Saved input is capped at 64 MiB by the shared JSON export reader.

## Entry grid and disclosure

Show one expandable entry per inventoried source, grouped by route ID rather than path alone. Columns:
source type, method/path signature, registration location, finding count, analysis-review count, and
protected-claim count. Findings appear first, then unresolved analysis and zero-finding entries;
within these groups, sort by signature with stable input order for ties. Distinct HTTP methods sharing a path remain separate. Today the source type is
HTTP route; a use without a matching entry appears under an explicit unresolved entry, never an invented
route. Other entry types require future detector contracts rather than labels inferred by the viewer.

Native `details` / `summary` elements use CSS Grid to align row summaries. This is a disclosure list
with column styling, not an ARIA interactive grid. Keyboard activation, visible focus, and native
expand/collapse behavior work without JavaScript. Narrow layouts stack the columns and label counts.
The expanded row contains:

- Named observed issues, incomplete analysis items, and scoped protection with distinct status labels.
- Plain-English behavior and separate input-validation / use-protection explanations.
- Source locations, operation excerpt, resolved inputs and transfer evidence, and additional context.
- Raw assessment state/provider/model, declared checks, gaps, blockers, and scope limitations.
- Remediation actions and optional explicitly unverified illustrative examples.

Zero-finding entries remain visible with “No findings reported within the analyzed scope,” plus their
analysis gaps. They do not receive a safety badge. Counts describe route/use records, not unique
vulnerabilities; shared helper observations can appear under several routes. No severity or exploit
confirmation is invented for offline observations. The review and runtime-confirmation stages remain
unimplemented; opening an HTML report does not perform either stage.

## Standalone behavior and source safety

Go `html/template` escapes report text, code, source paths, and locations. Source excerpts are untrusted
text and must never be cast to trusted HTML. CSS and markup are bundled with the scanner, with no
scripts, remote fonts, CDN dependencies, automatic JSON fetch, or analytics. A restrictive content
security policy disables scripts and network resources. The generated file can be opened directly
with `file://`; browser access to a neighboring JSON file is unnecessary.

Reports contain source excerpts and must be handled with the target source's access controls. The
current viewer does not link to or execute arbitrary source paths. Optional future filtering, sorting,
and expand-all controls would require a separate UX decision; they are not needed for native row
expansion. Standard browser find remains available, although matching inside collapsed disclosures is
browser-dependent. Printing follows the current disclosure state; expand required rows before printing.

## CI/CD retention

Point the output root at the runner's temporary storage. Upload the complete tree with an always-run
artifact step, retaining JSON for agents, HTML for human download, and SARIF for standard ingestion.
Keep the scan step's failure visible: bundle creation does not mask a requested-assessment failure.
Early setup/analysis failure may produce no bundle; artifact upload must handle that case explicitly.
SARIF ingestion is separate from artifact upload, and a finding-based CI gate is not implemented.
Use CI-managed secrets only for deliberately configured paid runs; offline export needs none.

## Verification

Tests cover source/method grouping, orphan uses, zero-finding entries, source-text escaping,
recommendations and evidence, deterministic output, old schema-v1 input, invalid input, and incompatible
flags. Detector baselines are compared against fresh offline runs on the two pinned target repositories;
presentation changes must not alter raw use or gap totals. Browser inspection checks real expanded
findings and remediation. JSON/HTML/SARIF pairs and verification notes belong in ignored local run
artifacts; reusable rendering and classification tests belong in Go test sources.
