# Vulnerable TypeScript app: source-reviewed reference

[Sample](README.md) · [Pinned file hashes](../../../.zazz/evaluation/sample-reference-manifests.json) ·
[Baseline and remediation policy](../reference-policy.md)

Reviewed 2026-11-01 at `7f8b240a2f02beba099da3632d1ba679ff8bc32f`.
The checked-in `src/` contains `server.ts`. The upstream README lists older NestJS/controller paths
and package-vulnerability totals; those labels do not establish current source locations, scanner
coverage or a fresh dependency audit. Use the actual source below as the reference.

| Case | Expected source-based classification | Why |
|---|---|---|
| `TS-XSS-01`, `GET /api/search` | Reflected server-XSS construction: positive | `src/server.ts:114-117` reads `req.query.query`, assigns it to `query`, and interpolates it into `<h1>` in `res.send`. Express string output is HTML; no output encoder/sanitizer is applied. The TypeScript `as string` cast does not validate or encode runtime input. |
| `TS-XSS-N01`, JSON output | Not an HTML-XSS sink | `res.json` serializes API data; a later browser rendering sink would need a separate linked claim. Do not turn proxy/API data into a server XSS finding solely because it contains a request dependency. |
| `TS-SQL-N01`, `POST /api/login` | No executed SQL-injection sink at this site | `src/server.ts:69-76` constructs SQL-looking text and logs it. It then calls the in-memory array's `users.find`, not a database. The source comment explicitly calls it a simulation. No database query receives the string. |
| `TS-LOG-01`, same login | Logged request content is established; sensitivity requires a stated contract | The query text includes the request field named `password` and reaches `console.log` at `:74`. That is a direct revealing content path if the supplied field is a credential. The name/endpoint suggests sensitivity, but the scanner must retain its classification evidence rather than infer an executed SQL query or invent a password comparison. |
| `TS-OUTBOUND-01`, `GET /api/proxy` | Unrestricted outbound destination: source-supported review target | `src/server.ts:121-128` forwards a request URL into `axios.default.get` without a local destination restriction. That is separate from XSS and requires client binding/destination analysis for a definitive SSRF claim. |
| `TS-CLIENT-N01` | No positive browser DOM-XSS example established | This pinned source has the Express server and no checked-in browser component/template application. Do not manufacture client routes or browser findings from server API paths. |

The existing app has no identified escaped counterpart for the search route; no protected XSS
output is asserted merely to create a balanced count. A remediated search counterpart belongs in
an independent copy. Scanner results, saved model answers and source-reviewed expectations are
separate artifacts. No new live model requests are required for unchanged-baseline comparisons.

Upstream source attribution is retained in the [included license](UPSTREAM-LICENSE.txt).
