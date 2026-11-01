# SQL/NoSQL judgment calibration

## Scope and evidence

Pinned Juice Shop target: v20.3.0, `f189f521266f42c6bd05884468292000cee852ed`.
These source-only experiments use OpenAI Decisions gpt-6-luna. No target execution or exploitation
was performed. Expected challenge answers and control labels were excluded from model bundles.
Reports are saved outside the checkout under `<DEV_ROOT>/zazzcode/zazzles/_test-findings/`.

Seven bounded batches made 28 completed requests: 235,102 input tokens and zero reported output
tokens. No provider errors or automatic retries/fallback. Local manifests retain exact source report
paths, typed answers, versions, and usage. A run can exit 1 for other unassessed uses after its
one-request cap even when the selected query completed.

## Calibration progression

- Initial SQL login/search judgments supported injection. NoSQL was inconclusive despite recognizing
  executable construction: collection provenance, guard scope and API semantics needed improvement.
- Verified MarsDB imports/constructors, question-specific choices, and stronger middleware evidence
  resolved small raw-expression/operator-selector cases. Numeric conversion initially remained inconclusive.
- Numeric output facts improved protection judgments, but v6 incorrectly protected Juice Shop reviews
  by borrowing a conversion guarantee from another branch. Condition helpers and branch scoping were added.
- V7 detected reviews and an unsafe conditional, but falsely flagged a numeric/constant conditional.
  Consistent injection/path answers did not prevent this false positive.
- V8 adds AST-derived facts for each output alternative, including all-numeric conditional outputs,
  and balanced semantic instructions. It does not force a protected/vulnerable label from a route name.

## Latest focused live checks: protocol string-use-v8

| Source case | Expected property | Actual injection judgment | Report folder |
|---|---|---|---|
| Juice Shop GET /rest/products/:id/reviews | Unsafe text branch under a configuration precondition | Vulnerable | juice-shop/run-3200022220 |
| Native Number conversion or raw request text | One unsafe alternative | Vulnerable | conditional-numeric/run-1239887494 |
| Native Number conversion or constant 0 | Every output alternative numeric | Protected | conditional-safe/run-64797390 |

The protected control still has an orange executable-string construction-practice finding; its
injection facet is Protected. That is not an injection false positive. Field business-validation
requirements remain separate. Latest focused batch: three requests, 24,637 input tokens.

Offline tests cover nested numeric/raw siblings and shadowed Number bindings. The full
`go test -race ./...` suite and `go vet ./...` passed after the code changes. Build passed; the sandbox
could not write an optional module stat-cache file. No browser interaction acceptance is claimed.

## Remaining acceptance work

The focused false positive is corrected, but this is not a general accuracy measurement. Recheck
raw/bound SQL, numeric/scalar/operator NoSQL and unknown-guard controls on v8; then refresh all five
Juice Shop target query routes and the vulnerable TypeScript app. Compare its README to actual source:
constructed/logged SQL without a database call must not become a database injection finding.

Create final consistent JSON/HTML/SARIF reports with actual matching-snapshot assessments. Do not
present a selected-route calibration report as a full assessed application. Distinct per-field
validation defect counting, independent Findings Review, and held-out validation remain incomplete.
Logger-leak implementation follows this work in a new session.
