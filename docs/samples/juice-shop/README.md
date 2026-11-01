# OWASP Juice Shop sample

[All samples](../README.md) · [HTML report](findings.html) · [JSON evidence](findings.json) · [SARIF](findings.sarif)

Generated 2026-11-01 from [OWASP Juice Shop](https://github.com/juice-shop/juice-shop), v20.3.0,
revision `f189f521266f42c6bd05884468292000cee852ed`. Extraction: scan `logger-leak` branch with
sensitive-data logging analysis; model judgments come from the `c21070d` gpt-6-luna run (question
protocol `string-use-v8`) and are reused only where a use's evidence is byte-identical. Presentation
uses the renderer checked in alongside this sample. Schema version 1, presentation version 3.

This full entry inventory contains **152 route entries** and **160 use records**. Seven records are
Findings: three source-supported static injection judgments; one credential written to access logs
(`GET /rest/user/change-password` reads `current` and `new` passwords from the URL query, and
`morgan('combined')` records the full URL); and three query-construction findings. The other 153
records are Review, not proven vulnerabilities. Counts are per use/route; categories can overlap.
Independent findings review and runtime confirmation were not performed.

Entry: `server.ts`, with no route filter. Thirteen database uses keep their completed model
assessments. The login SQL, track-order NoSQL and chat SQL uses gained dependency evidence through
chained calls (for example the hashed password in the login query), so their earlier answers were
not reused; they await a new judgment and show as construction-practice findings meanwhile. Logging
and other non-query uses were not semantically assessed. This is full route inventory coverage within
the selected entry, not full detector or semantic coverage.

Offline inventory command:

```bash
./bin/scan --strings --entry server.ts --output-dir "$SCAN_FINDINGS_DIR" /path/to/juice-shop
```

The sample also retains the actual bounded model assessments described above. That offline command
alone does not recreate those answers. Viewing/exporting requires no key or live requests.
Follow the [refresh procedure](../README.md#refresh-as-scan-improves) when updating results.
Upstream source excerpts are attributed to Juice Shop contributors under the included
[upstream license](UPSTREAM-LICENSE.txt).
