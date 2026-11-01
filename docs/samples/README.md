# Sample findings reports

[Project introduction](../../README.md) · [User documentation](../README.md) · [Interpret results](../interpreting-results.md)

Explore actual scan output without installing the targets, configuring keys, or calling a model.
Open the [hosted samples](https://zazzcode.github.io/zazzles/) in your browser, or download a standalone HTML file and open it
locally. GitHub's file view displays HTML source; it does not run the report.

| Sample | HTML report | Agent-readable evidence | Standard interchange |
|---|---|---|---|
| [OWASP Juice Shop](juice-shop/README.md) | [findings.html](juice-shop/findings.html) | [findings.json](juice-shop/findings.json) | [findings.sarif](juice-shop/findings.sarif) |
| [Vulnerable TypeScript app](vulnerable-typescript-app/README.md) | [findings.html](vulnerable-typescript-app/findings.html) | [findings.json](vulnerable-typescript-app/findings.json) | [findings.sarif](vulnerable-typescript-app/findings.sarif) |
| [Node.js Goof](nodejs-goof/README.md) | [findings.html](nodejs-goof/findings.html) | [findings.json](nodejs-goof/findings.json) | [findings.sarif](nodejs-goof/findings.sarif) |

All three current reports include XSS discovery and the **Hide protected results** remediation
filter. Server/API routes and browser/client-side sources have separate sections and headings.
[Source-reviewed answer keys and remediation-copy policy](reference-policy.md) explain why specific
baseline flows are vulnerable, protected or unresolved, independently of detector output. Goof's
[Bearer baseline](nodejs-goof/bearer/20261101-goof-baseline-02/SUMMARY.md) and
[five-case Decisions evaluation](nodejs-goof/decisions/20261101-five-case-01/README.md) are saved
separately from the canonical zazzles results.

Reports start with the grouped Findings → Review → Protected sort, highest counts first, then route
path A–Z; blank rows appear last. Click that grouped header or Route signature to change the order.
Expand a row for evidence and remediation. Type pills open the glossary. HTML is standalone, with
no external resources or model calls. JSON preserves richer analysis context.

These are evolving examples, not a complete benchmark or independent accuracy claim. Sample READMEs
record target versions, actual assessment coverage, and current limitations. Inconclusive and
unassessed records remain visible. The target root is replaced with its public project name; relative
source locations, hashes, use IDs, and provider judgments are preserved. Source excerpts belong to
the upstream projects; vulnerable fixture values are training data, not scan provider credentials.

## Refresh as zazzles improves

Build scan, then regenerate both checked-in HTML/SARIF views from the saved JSON:

```bash
go build -o bin/scan ./cmd/scan
python3 docs/samples/refresh.py --scan ./bin/scan
```

This updates presentation only, offline. To publish improved analysis, run a fresh target scan into
your normal external findings directory, replace the corresponding sample `findings.json` with that
run's canonical JSON, update its README's target revision/date/assessment coverage, then run the
refresh command. The script normalizes the public root label and exports matching HTML/SARIF.

Saved model judgments may be carried into a fresh scan only for a use whose whole evidence is
byte-identical to the judged use. A use whose evidence changed is published as not assessed until it
is judged again; never copy an old answer onto changed evidence.
Do not copy old judgments onto changed source evidence or invent missing model assessments.
Review the diff and preserve upstream attribution before committing the refreshed samples.

Routine scans continue to write outside the worktree. Only these deliberately curated samples are
checked in. Keep all three sample applications current when analysis or reporting behavior changes.

## Browser-hosted previews

The **Publish sample findings** GitHub Actions workflow publishes this directory to GitHub Pages,
including an index linking all three reports. A repository maintainer first selects **GitHub Actions**
under Settings → Pages and allows `stage` in the `github-pages` environment. The
workflow then runs when sample files or the workflow change on `stage` (the default branch), and can
be run manually there. Feature branches publish when they merge into `stage`.
The hosted address is GitHub's project-site URL, https://zazzcode.github.io/zazzles/. It serves pages only after a maintainer
enables Pages for the repository; until then, download and open the HTML file.
Only the curated sample directory is published; no scan execution or provider keys are needed.
