# OWASP Juice Shop sample

[All samples](../README.md) · [HTML](findings.html) · [JSON](findings.json) · [SARIF](findings.sarif) ·
[Source-reviewed reference](xss-answer-key.md) · [Provenance](metadata.json)

Refreshed 2026-11-01 from the [upstream application](https://github.com/juice-shop/juice-shop), revision
`f189f521266f42c6bd05884468292000cee852ed`, using the `xss-response` working sources and binary identified in metadata.
The original checkout remains unchanged. Remediation exercises use an independent copy.

The current report contains 152 server/API entries, 9 browser/client entries, and 197 use records: 7 Findings and 190 Review records. One XSS construction finding is the browser query → Angular trust bypass → search template flow. Other client/API-origin and collection/dialog links remain Review. No Protected records are asserted by this extraction snapshot.

All implemented analyzer families remain visible; this is not an XSS-only route filter. HTML has
separate Server/API and Browser/client-side sections and a Hide protected results control. The
filter changes presentation only; JSON/SARIF retain all evidence. Counts are per use, categories
overlap, and no-findings does not mean an application is safe.

This refresh is offline. Historical model judgments were not carried onto changed evidence/rules;
their earlier versions remain in Git history. No live model calls were made for this target in this
session. Source-reviewed expected properties are in the reference above, not inferred from finding
counts or comments. See the [static baseline/copy policy](../reference-policy.md).

Upstream source excerpts are attributed under the [included license](UPSTREAM-LICENSE.txt).
