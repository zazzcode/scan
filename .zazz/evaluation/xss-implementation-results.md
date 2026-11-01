# XSS first-pass implementation and evaluation

**2026-11-01, `xss-response`.** Discovery/reporting is implemented and the user requested that
additional analysis/judging hardening continue with a different agent in the same worktree.
The [handoff](../ephemeral/xss-analysis-handoff.md) is the entry point for that session.

## Current samples

| Target | Server entries | Client entries | Use records | Findings | Review | Protected |
|---|---:|---:|---:|---:|---:|---:|
| [Goof](../../docs/samples/nodejs-goof/findings.html) | 19 | 0 | 33 | 2 | 29 | 2 |
| [Juice Shop](../../docs/samples/juice-shop/findings.html) | 152 | 9 | 197 | 7 | 190 | 0 |
| [Vulnerable TypeScript app](../../docs/samples/vulnerable-typescript-app/findings.html) | 19 | 0 | 19 | 1 | 18 | 0 |

These are use-level results, not confirmed-vulnerability totals. Categories overlap. Goof's login
has a source-complete static Luna XSS judgment; its stored Markdown output is an orange construction
finding while the model's exposure answer remains unknown. The two Protected records are the edit
template's content text/quoted-attribute outputs. Its account-details whole-object render options
remain Review, despite observed Handlebars escaping. Its `/about_new` query log is not labeled a
sensitive-data leak without data-class evidence. Input-validation and Mongoose NoSQL records and
the static `/users` mount are included; TypeORM object operations are not invented as raw SQL.

Juice Shop's client finding is the browser query → verified Angular trust bypass → search template
flow. Unresolved API producer fields, JWT/gallery/table/component/dialog links and trusted-wrapper
return types remain Review. The TypeScript sample's HTML search response is a reflected-XSS finding;
its SQL-looking login log remains distinct from an executed query. No current historical SQL model
answers were copied onto changed evidence/rules in the two offline targets.

## Decisions and source references

Five specifically approved `gpt-6-luna` requests completed: **23,436 input tokens, 0 output tokens**.
[Raw requests/responses and scored answers](../../docs/samples/nodejs-goof/decisions/20261101-five-case-01/README.md)
are saved. No more live calls were made to preserve response confidence/probabilities or render
reports. Exact matching requests were replayed offline. The API supplied no narrative reasons;
policy explanations are scanner-authored. Scores are uncalibrated audit metadata, not evidence.

Source-reviewed expected outcomes and hashes are maintained separately for
[Goof](../../docs/samples/nodejs-goof/xss-answer-key.md),
[Juice Shop](../../docs/samples/juice-shop/xss-answer-key.md), and
[TypeScript](../../docs/samples/vulnerable-typescript-app/xss-answer-key.md).
These references are not detector input, runtime confirmation or a held-out accuracy benchmark.
Original target checkouts remain clean and pinned. Remediation uses independent copies.

## Increment and validation record

Saved report snapshots under each sample's `zazzles/` directory cover reflected output, template
construction, stored rendering, judgment, browser discovery and the current combined analyzers.
Earlier protection labels are historical and superseded by the audited current canonical reports.
Metadata records scanner working-source hashes, binary SHA256, target revision and source hashes.
The initial Bearer comparison is preserved separately; no raw comparator outcomes were rewritten.

Validation: native Go 1.27.3 darwin/arm64; `go test -race ./...`, `go vet ./...`, build and changed-Go
format checks passed. Tests cover positive/protected/unknown contexts, CommonJS/router bindings,
stored field links, templates, DOM/Angular/React/Vue, protocol consistency, retained response scores,
identical-request replay and client identity/export headings. The remediation filter was exercised
offline against mixed and protected-only sources, including reset behavior. HTML/JSON/SARIF were
regenerated from canonical JSON. Visual browser QA was not performed because file URLs were blocked
by the browser tool; structural/export checks do not replace a visual review.

Further hardening priorities and exact disagreements are in the local, gitignored handoff. No broader language
or whole-program accuracy claim is made; Python extraction is the next planned adapter after XSS.
