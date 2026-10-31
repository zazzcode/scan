# Query structure and input validation comparison

**2026-10-31.** Offline scan after the query-structure and validation pass. No live judge call.
Assessments remain `not_run`. A gap is a candidate for an agent fix, not a confirmed exploit.
Training snippets and the Juice Shop frontend are outside this comparison. Bearer totals that
include them are not the denominator.

Targets: Juice Shop `f189f521266f42c6bd05884468292000cee852ed` with `--entry server.ts`, and
vulnerable-typescript-app `7f8b240a2f02beba099da3632d1ba679ff8bc32f` with `--entry src/server.ts`.
Bearer results are the pinned run already stored for those revisions. Bearer has no NoSQL rule
hit on either target.

## Juice Shop documented injection and validation

| Documented item | Source | Scan | Bearer |
|---|---|---|---|
| Login SQL, challenges loginAdmin / loginBender / loginJim | `routes/login.ts:34` template | SQL candidate, `query_text_interpolated`, email validation absent | SQL injection at line 34 |
| Search SQL, unionSqlInjection / dbSchema | `routes/search.ts:23` template after a length truncation | SQL candidate, `query_text_interpolated`, validation absent. The length truncation is not treated as validation | SQL injection at line 23 |
| NoSQL DoS, noSqlCommandChallenge | `routes/showProductReviews.ts:36` `$where` concatenation | NoSQL candidate, `query_text_interpolated`. `Number` / `trunc` is unresolved validation | No NoSQL finding |
| NoSQL exfiltration, noSqlOrdersChallenge | `routes/trackOrder.ts:18` `$where` template | NoSQL candidate, `query_text_interpolated`. `replace` / `trunc` is unresolved validation | No NoSQL finding |
| NoSQL manipulation, noSqlReviewsChallenge | `routes/updateProductReviews.ts:17` `update({_id: req.body.id}, {$set}, {multi: true})` | NoSQL candidate, `query_document_unvalidated`, validation absent | No NoSQL finding |
| Chat catalog search | `routes/chat.ts:117` and `:123` `Op.like` templates | SQL candidates, `query_text_interpolated`. The tool argument is not a direct request field | No SQL finding |
| Chat review lookup | `routes/chat.ts:149` `$where` concatenation | NoSQL candidate, `query_text_interpolated` | No NoSQL finding |
| Empty registration / repetitive registration / admin registration | `server.ts` middleware around the user create route | Not a new query sink. Field reads on that route are validation candidates when the handler reads them. The passwordRepeat bypass remains the existing account-creation check | No input-validation rule |

Constant `models.sequelize.query('SELECT sql FROM sqlite_master')` in search is not a candidate.
The login password is hashed inside the same interpolated string, so the structure defect remains.

## Vulnerable TypeScript app documented items

The README names `main.ts` and `app.controller.ts`. Those files are not in this revision. The
running source is `src/server.ts`.

| README claim | Actual source | Scan | Bearer |
|---|---|---|---|
| SQL injection | `src/server.ts:73` builds a SQL string and `console.log`s it. No database call | No SQL candidate. Username is an `input_validation_unresolved` field because `users.find` receives it and is not a validator | No SQL injection. Logger finding at the log line |
| Other README rows (command, file, XSS, eval, XXE, redirect, secrets, cookies) | Present in `src/server.ts` | Outside this branch. Routes are inventoried. Request fields that do not reach a query are validation candidates. The proxy remains the outbound candidate at line 125 | Those rules fire. See the earlier 28-finding run |

## What this branch does not claim

- Every request field in both repositories has been proven unvalidated. Unresolved middleware can
  still validate a field. `input_validation_unresolved` means a call received the field and its
  effect is unknown.
- Custom sanitizer quality is not decided here. An allow-list that really constrains a quoted SQL
  value can still be a sufficient structure control. That remains a judge question. The
  concatenation fact stays visible.
- Outbound helper calls, especially `lib/webhook.ts`, are repeated once per calling route. That is
  helper expansion, not a query-structure result.
- Logger leaks, secrets, command injection, path traversal, XSS, `eval`, XXE, open redirect, weak
  randomness, cookie defaults, and prototype pollution are the [future branches](../project/implementation-plan.md#future-capability-branches).
