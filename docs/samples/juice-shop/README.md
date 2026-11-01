# OWASP Juice Shop sample

[All samples](../README.md) · [HTML report](findings.html) · [JSON evidence](findings.json) · [SARIF](findings.sarif)

Generated 2026-11-01 from [OWASP Juice Shop](https://github.com/juice-shop/juice-shop), v20.3.0,
revision `f189f521266f42c6bd05884468292000cee852ed`. Extraction/judgment checkpoint: scan `c21070d`;
presentation uses the renderer checked in alongside this sample. Schema version 1, presentation
version 3, OpenAI Decisions gpt-6-luna question protocol `string-use-v8`.

This full entry inventory contains **152 route entries** and **137 use records**. Six records are
Findings: five source-supported static injection judgments and one construction-practice finding.
The other 131 records are Review, not proven vulnerabilities. Counts are per use/route; categories
can overlap. Independent findings review and runtime confirmation were not performed.

Entry: `server.ts`, with no route filter. Sixteen of 17 database uses received completed model
assessments: five vulnerable and eleven Inconclusive. One chat-query bundle exceeded the 64 KiB
payload limit before transmission; the other 120 non-query uses were not semantically assessed.
This is full route inventory coverage within the selected entry, not full detector or semantic coverage.

Offline inventory command:

```bash
./bin/scan --strings --entry server.ts --output-dir "$SCAN_FINDINGS_DIR" /path/to/juice-shop
```

The sample also retains the actual bounded model assessments described above. That offline command
alone does not recreate those answers. Viewing/exporting requires no key or live requests.
Follow the [refresh procedure](../README.md#refresh-as-scan-improves) when updating results.
Upstream source excerpts are attributed to Juice Shop contributors under the included
[upstream license](UPSTREAM-LICENSE.txt).
