# Vulnerable TypeScript app sample

[All samples](../README.md) · [HTML report](findings.html) · [JSON evidence](findings.json) · [SARIF](findings.sarif)

Generated 2026-11-01 from Ryan McMorrow's
[vulnerable TypeScript application](https://github.com/ryanmcmorrowsnyk/vulnerable-typescript-app),
revision `7f8b240a2f02beba099da3632d1ba679ff8bc32f`. Extraction/judgment code: scan `c21070d`;
presentation uses the renderer checked in alongside this sample. Schema version 1, presentation 3.

Entry: `src/server.ts`, with no route filter. The fresh offline scan contains **19 route entries**
and **14 Review records**. No semantic model calls were made. Zero reported Findings is not a safety
verdict: several intentionally vulnerable behaviors need detector families that are not implemented.
The sample demonstrates retained evidence, incomplete analysis, and routes without reported findings.

The application's SQL-looking string is constructed and logged, with no database call at that site.
Its array `find` is not a database API. The scanner does not invent a SQL/NoSQL injection finding from
those shapes. Sensitive-data logging and other README-listed weakness families remain follow-up work.

```bash
./bin/scan --strings --entry src/server.ts --output-dir "$SCAN_FINDINGS_DIR" /path/to/vulnerable-typescript-app
```

Source excerpts retain their upstream attribution. Local target paths are normalized to the public
project name. Follow the [refresh procedure](../README.md#refresh-as-scan-improves) as analysis improves.
