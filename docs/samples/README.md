# Sample findings reports

[Project introduction](../../README.md) · [User documentation](../README.md) · [Interpret results](../interpreting-results.md)

Explore actual scan output without installing the targets, configuring keys, or calling a model.
Download a standalone HTML file and open it in your browser. GitHub's file view displays HTML source;
it does not run the report. After cloning this repository, open the local file directly.

| Sample | HTML report | Agent-readable evidence | Standard interchange |
|---|---|---|---|
| [OWASP Juice Shop](juice-shop/README.md) | [findings.html](juice-shop/findings.html) | [findings.json](juice-shop/findings.json) | [findings.sarif](juice-shop/findings.sarif) |
| [Vulnerable TypeScript app](vulnerable-typescript-app/README.md) | [findings.html](vulnerable-typescript-app/findings.html) | [findings.json](vulnerable-typescript-app/findings.json) | [findings.sarif](vulnerable-typescript-app/findings.sarif) |

Reports start with the grouped Findings → Review → Protected sort, highest counts first, then route
path A–Z; blank rows appear last. Click that grouped header or Route signature to change the order.
Expand a row for evidence and remediation. Type pills open the glossary. HTML is standalone, with
no external resources or model calls. JSON preserves richer analysis context.

These are evolving examples, not a complete benchmark or independent accuracy claim. Sample READMEs
record target versions, actual assessment coverage, and current limitations. Inconclusive and
unassessed records remain visible. The target root is replaced with its public project name; relative
source locations, hashes, use IDs, and provider judgments are preserved. Source excerpts belong to
the upstream projects; vulnerable fixture values are training data, not scan provider credentials.

## Refresh as scan improves

Build scan, then regenerate both checked-in HTML/SARIF views from the saved JSON:

```bash
go build -o bin/scan ./cmd/scan
python3 docs/samples/refresh.py --scan ./bin/scan
```

This updates presentation only, offline. To publish improved analysis, run a fresh target scan into
your normal external findings directory, replace the corresponding sample `findings.json` with that
run's canonical JSON, update its README's target revision/date/assessment coverage, then run the
refresh command. The script normalizes the public root label and exports matching HTML/SARIF.
Do not copy old judgments onto changed source evidence or invent missing model assessments.
Review the diff and preserve upstream attribution before committing the refreshed samples.

Routine scans continue to write outside the worktree. Only these deliberately curated samples are
checked in. Keep both sample applications current when analysis or reporting behavior changes.

## Browser-hosted previews

The manual **Publish sample findings** GitHub Actions workflow publishes this directory to GitHub
Pages, including an index linking both reports. A repository maintainer first selects **GitHub
Actions** under Settings → Pages, then runs that workflow from the branch containing these samples.
Use the deployment's actual published URL in README links once the first deployment succeeds.
This setup is committed; hosting is not yet activated or verified. Until then, download/open the HTML.
Only the curated sample directory is published; no scan execution or provider keys are needed.
