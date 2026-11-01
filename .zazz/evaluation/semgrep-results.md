# Semgrep CE baseline and XSS comparison

**2026-11-01.** First Semgrep Community Edition baseline on the three pinned sample apps, compared with
Zazzles (hardening 10, commit `83aa7bd`) and the existing [Bearer results](bearer-results.md) on the
source-reviewed XSS answer keys. How to repeat it: [Semgrep CE cheat sheet](../../docs/semgrep-ce-cheat-sheet.md).

## Run

Semgrep 1.183.0, `--metrics=off`, offline against a rule snapshot downloaded once from the registry:
`p/javascript`, `p/typescript`, `p/xss` and `p/owasp-top-ten` (563 rules after merging, all valid).
Exclusions match the Zazzles source scope (`node_modules`, `.github`, Goof `exploits/` and
`public/js/bundle.js`, Juice Shop `three.js`, build output). Targets were unchanged before and after
(`git status --porcelain --ignored` empty).

| Target | Revision | Files | Results | Errors | Baseline |
|---|---|---:|---:|---:|---|
| Node.js Goof | `add14ba` | 30 | 6 | 0 | [JSON/SARIF/text](../../docs/samples/nodejs-goof/semgrep/20261101-baseline-01/metadata.json) |
| Vulnerable TypeScript app | `7f8b240` | 7 | 7 | 0 | [JSON/SARIF/text](../../docs/samples/vulnerable-typescript-app/semgrep/20261101-baseline-01/metadata.json) |
| OWASP Juice Shop | `f189f52` | 1,168 | 30 | 32 | [JSON/SARIF/text](../../docs/samples/juice-shop/semgrep/20261101-baseline-01/metadata.json) |

Juice Shop's 32 errors are syntax/partial-parse warnings in `data/static/codefixes/*` challenge snippets,
which are not complete TypeScript files. Curated copies make target paths relative and drop the local
rules-directory prefix from rule IDs; raw outputs remain under `_test-findings/<app>/semgrep-20261101-1/`.

## XSS answer-key comparison

Positive = expected vulnerable flow found as a finding. "Review" = Zazzles kept a record with a named
blocker instead of claiming a finding. Bearer columns use the saved Goof baseline and the earlier
TypeScript comparison. Bearer has not yet been run on Juice Shop for XSS.

| Case | Expected | Zazzles | Semgrep CE | Bearer |
|---|---|---|---|---|
| `GOOF-XSS-01` reflected login (EJS raw attribute) | positive | finding (red, replayed judgment) | missed | missed |
| `GOOF-XSS-02` stored Markdown, create | positive | finding, separate supported path | missed | missed |
| `GOOF-XSS-03` stored Markdown, update | positive | finding, separate supported path | missed | missed |
| `GOOF-XSS-N01`/`N02` escaped edit outputs | protected | Protected | no result | no result |
| `GOOF-XSS-N03` render options | Review | Review | no result | no result |
| `GOOF-XSS-N04` fixed analytics script | no finding | no finding | no finding | **false positive** (`dangerous_insert_html`) |
| `TS-XSS-01` reflected search | positive | finding | found (`direct-response-write`, `raw-html-format`) | found |
| `TS-XSS-N01` JSON output | no finding | no finding | no finding | — |
| `JS-XSS-01` DOM search trust bypass | positive | finding | missed | not run |
| `JS-XSS-02` track order (cross-tier) | positive | finding, branch paths | missed | not run |
| `JS-XSS-03` saved login IP via JWT | positive | Review (token origin) | missed | not run |
| `JS-XSS-04` admin email | positive | Review (stored origin) | missed | not run |
| `JS-XSS-05` feedback | positive | Review | missed | not run |
| `JS-XSS-06` product description dialog | positive | Review (dialog sink unlinked) | missed | not run |
| `JS-XSS-07` profile template construction | positive | Review (stored origin) | missed | not run |
| `JS-XSS-08` subtitles | Review | Review | no result | not run |
| `JS-XSS-N04` CAPTCHA trust bypass | no finding | Review, no finding | no finding | not run |

Totals over expected positives: Zazzles **6 of 11** found, with the other 5 kept as Review;
Semgrep CE **1 of 11**; Bearer **1 of 4** run (with one Goof false positive).

Semgrep's Angular `bypassSecurityTrust*` rule is taint-mode with intra-file sources, so it does not reach
Juice Shop's route-parameter, API-response or stored sources. CE has no cross-file or stored-flow
tracking (that is Semgrep Pro), and the snapshot packs have no rule for raw EJS `<%-` output or `marked`.
That is a capability difference on these labeled cases, not a general accuracy claim.

## Other Semgrep results (not XSS)

- **Goof:** cookie-session options and a hard-coded session secret (`app.js:42-43`), Dockerfile user.
- **TypeScript:** hard-coded session and JWT secrets, path traversal `:103`, string-built code `:137`,
  open redirect `:239`. Zazzles has no rule for secrets, path traversal or open redirect; these are
  candidate capability gaps, not verified cases.
- **Juice Shop:** six Sequelize injection results, `res.sendFile`, directory listing, a JWT secret,
  an open redirect, and Terraform/npm configuration rules.

Next: run Bearer on Juice Shop with the pinned rules; write answer-key cases for the non-XSS Semgrep
results before counting them; repeat with each Zazzles hardening increment against the same rule snapshot.
