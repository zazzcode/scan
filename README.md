# Zazzles

Zazzles is a Go static application security scanner (SAST) and dependency-audit CLI under active
development. It aims to help developers and coding agents understand how request inputs reach
security-sensitive operations,
what validation and protection already hold, and where code needs hardening. Reports should give
reviewers source evidence and explicit analysis gaps so they can judge each finding.

The approach combines Tree-sitter source extraction with code-flow analysis and optional decision-model
assessment of relevant evidence. Field validation and sink protection are separate questions: a valid
email address does not establish that a SQL query uses bound parameters. Unknown behavior stays
unknown rather than becoming a safety claim or a vulnerability by assumption.

## Explore sample findings

These reports show how a Zazzles scan works end to end: each is the unedited output of an offline scan of
an intentionally vulnerable application, with route and browser entries, source-to-sink evidence, per-path
protection, and Findings, Review and Protected results. Open one in your browser:

| Application | What the current sample shows | Output |
|---|---|---|
| **OWASP Juice Shop** | 152 route entries and 9 browser entries; DOM search and cross-tier track-order XSS, SQL/NoSQL query-construction findings, and a credential written to access logs | [View in browser](https://zazzcode.github.io/zazzles/docs/samples/juice-shop/findings.html) · [HTML file](docs/samples/juice-shop/findings.html) · [JSON](docs/samples/juice-shop/findings.json) · [SARIF](docs/samples/juice-shop/findings.sarif) |
| **Node.js Goof** | 19 route entries; reflected login XSS, stored Markdown XSS with separate create/update write paths, and two escaped outputs marked Protected | [View in browser](https://zazzcode.github.io/zazzles/docs/samples/nodejs-goof/findings.html) · [HTML file](docs/samples/nodejs-goof/findings.html) · [JSON](docs/samples/nodejs-goof/findings.json) · [SARIF](docs/samples/nodejs-goof/findings.sarif) |
| **Vulnerable TypeScript app** | 19 route entries; a reflected-XSS search response, plus Review records for a password logged in SQL text and other incomplete analysis | [View in browser](https://zazzcode.github.io/zazzles/docs/samples/vulnerable-typescript-app/findings.html) · [HTML file](docs/samples/vulnerable-typescript-app/findings.html) · [JSON](docs/samples/vulnerable-typescript-app/findings.json) · [SARIF](docs/samples/vulnerable-typescript-app/findings.sarif) |

**View in browser** opens the hosted copy published by GitHub Pages ([all samples](https://zazzcode.github.io/zazzles/docs/samples/)). GitHub's
file view shows HTML source rather than running it, so you can also download an HTML file and open
it locally. Findings appear first; expand rows for evidence and remediation, or click a category pill for its
glossary. No setup or model calls are needed. These [evolving samples](docs/samples/README.md) include
coverage limits and refresh instructions; they will be updated as Zazzles improves.

## Goals and current scope

The goal is to report vulnerabilities, weak security practices, and sufficient protection by route,
with useful explanations and remediation guidance. Longer term, Zazzles aims to match or exceed Bearer's
SAST and developer-workflow capabilities with fewer false positives, measured at comparable detection
coverage on independently labeled cases. Bearer and Semgrep Community Edition serve as local comparison
tools. That is a project objective; broad parity is not established.

Today, Zazzles analyzes local JavaScript and TypeScript source. It inventories Express-style routes
(including CommonJS exports and static router mounts) and browser entries, follows a bounded subset of
handlers, helpers, branches and promise continuations, and exports source evidence as JSON, standalone
HTML, or SARIF. Supported capabilities:

| Area | Current support |
|---|---|
| **Cross-site scripting (XSS)** | Reflected server output (Express `send`/`write`/`end`, EJS, Handlebars, Pug); template construction before compilation; stored XSS through Mongoose/Sequelize model fields, with each write path kept separate; DOM sinks, jQuery, React, Vue and Angular templates, including trust bypasses; cross-tier links from a browser output to the server route and JSON response field that produce it. Protection is context-specific (HTML text, attribute, URL, script, style) and must hold on every modeled path. |
| **Injection** | String-built SQL, NoSQL selectors and `$where` construction, Sequelize raw queries |
| **Input validation** | Request fields reaching sinks without modeled validation, kept separate from sink protection |
| **Logging** | Request data and credentials reaching verified loggers |
| **Outbound requests** | Request-influenced `fetch`/`axios` destinations |
| **Judgment (optional)** | Versioned OpenAI Decisions or Jev questions; answers never replace missing source evidence |

See the [XSS analysis guide](docs/analysis/xss.md) for the exact contracts and limits.

Full control-flow and whole-program interprocedural analysis, independent findings review, Python and Go
analysis, and dependency vulnerability scanning remain planned. Recognizing a file extension in the file list
does not establish analysis support. A clean scan does not certify that an application is secure.

## Roadmap: vulnerability types

Each type is delivered as its own bounded increment with source-reviewed answer keys, paired
protected/unknown controls, and fresh sample reports. The order follows the
[implementation plan](.zazz/project/implementation-plan.md#future-capability-branches): Python support
comes next, then the remaining types below for JavaScript, TypeScript and Python.

| Vulnerability type | Status | First acceptance target |
|---|---|---|
| Input validation | Delivered | Request fields reaching sinks without modeled validation |
| SQL injection (string-built queries) | Delivered | Juice Shop login/search query construction |
| NoSQL injection (selectors, `$where`) | Delivered | Juice Shop and Goof Mongo/Mongoose queries |
| Sensitive-data logging | Delivered | Credentials and request data reaching verified loggers and access logs |
| Outbound request destinations | Review records | Request-influenced `fetch`/`axios` targets; SSRF judgment is planned below |
| Cross-site scripting (XSS) | In progress | Reflected, stored, template, DOM and cross-tier cases across all three samples; Juice Shop stored/JWT/dialog origins remain |
| **Python support** | **Next** | Python adapter for the delivered types above on a vulnerable Python app, with the same reports and judgment |
| Hard-coded secrets | Planned | Private material in source; public verification keys stay negatives |
| Command injection | Planned | `exec` of request input (TypeScript app `/api/ping`) |
| Path traversal | Planned | `fs.readFile` of a request path (`/api/files`) |
| Code injection | Planned | `eval` of a request expression (`/api/calculate`) |
| XML external entities (XXE) | Planned | XML parse of a request body (`/api/parse-xml`) |
| Open redirect | Planned | `res.redirect` of a request URL |
| Weak randomness | Planned | `Math.random` used as a token |
| Cookie, session and security headers | Planned | `httpOnly`, `secure`, session secret, missing baseline headers |
| Server-side request forgery (SSRF) | Planned | Judge existing outbound destination candidates |
| Prototype pollution | Planned | `_.merge` of request objects |

After all of these OWASP attack types: dependency vulnerability scanning (SCA) for JavaScript,
TypeScript and Python libraries. Other languages such as Go and Java come after that, along with finding
baselines and CI gating. The
[Semgrep CE and Bearer comparisons](.zazz/evaluation/semgrep-results.md) track each type against
other scanners on the same labeled cases.

## Agent-first philosophy

Our goal is to help agents fix code correctly. Report established vulnerabilities **and poor
practice worth remediating or refactoring**, even when an exploit has not been demonstrated.
Preserve sufficient protection, explain uncertainty honestly, and provide evidence and targeted
remediation. Validation and destination-specific protection can be separate improvements at the
same operation. Read the [project philosophy](docs/philosophy.md).

## Documentation

Start with the [documentation index and introduction](docs/README.md):

- [Philosophy](docs/philosophy.md) — why and when an agent should fix code
- [Getting started](docs/getting-started.md) — prerequisites, compilation, and `.env`
- [Running scans](docs/running-scans.md) — modes, entry selection, output folders, and exports
- [Interpreting results](docs/interpreting-results.md) — Findings, Review, Protected, evidence, and remediation
- [Finding categories](docs/glossary.md) — definitions and current coverage
- [Decision-model setup](docs/decision-models.md) — optional assessment and provider credentials
- [CI/CD](docs/ci.md) — JSON, HTML, and SARIF artifacts
- [Contributing](docs/contributing.md) — development workflow and verification

## Setup

Install the Go toolchain declared in `go.mod` and a C compiler for the Tree-sitter cgo bindings.
Build from the checkout root:

```bash
go build -o bin/scan ./cmd/scan
```

For a new local setup, copy [.env.sample](.env.sample) to `.env`; merge settings if `.env` already
exists so you preserve your keys. Zazzles does not auto-load `.env`. Offline scanning needs no API key.
See the [setup guide](docs/getting-started.md) for full instructions.

### Optional decision-model setup

Use [decision-model setup](docs/decision-models.md) when explicitly enabling assessment. The scanner
reads the selected provider's key from its process environment. Model requests transmit source
evidence and need an explicit request cap; a key is not a reason to run paid assessment automatically.

## Run

```bash
export SCAN_FINDINGS_DIR="$HOME/Dev/zazzcode/zazzles/_test-findings"
./bin/scan --strings --entry server.ts /path/to/juice-shop
./bin/scan --strings --entry src/server.ts /path/to/vulnerable-typescript-app
```

Each run creates `<output-root>/<repo-name>/run-<id>/` containing `findings.json`, `findings.html`, and
`findings.sarif`. Open the printed HTML path to review expandable source rows and category definitions.
`--output-dir` overrides the environment variable; `--output-dir ""` disables bundle output.

You can regenerate views without rescanning or calling a model:

```bash
./bin/scan --export-html /path/to/run/findings.json > /path/to/run/findings.html
./bin/scan --export-sarif /path/to/run/findings.json > /path/to/run/findings.sarif
```

The [running guide](docs/running-scans.md) covers inventory modes, route/check selection, output
formats, and limitations. JSON retains evidence for agents, HTML supports human review, and SARIF
provides standard security-tool interchange. Keep source-containing reports outside Git by default.

## Interpret results

| Column | Meaning | Action |
|---|---|---|
| **Findings** | Established problems, including poor practice that merits correction without a demonstrated exploit | Remediate or refactor the named problem and verify the change |
| **Review** | Evidence does not yet establish a problem or adequate protection | Resolve the named uncertainty before deciding on a fix |
| **Protected** | Evidence supports protection for the named claim | Preserve it; assess other fields and uses independently |

**Poor practice** is a finding to fix, not a request to wait for an exploit. The current deterministic
rules report string-built SQL and executable NoSQL `$where` constructions as poor practice.
A missing modeled validator or document constraint alone remains Review. A model's potential
vulnerability answer is retained for investigation; it does not automatically establish a finding.

**Categories** are overlapping labels, not additional findings. Route location identifies a route's
declaration; Finding location identifies the operation in the expanded record. The current report
still groups several concerns within a use record; distinct claim-level counting is planned. Read the
[results guide](docs/interpreting-results.md) for statuses, counts, and remediation examples.

Exit `0` means the operation completed, not that the target is safe. Exit `1` signals execution or
requested-assessment failure; exit `2` signals invalid usage. A finding-severity CI gate is not yet
implemented.

## Contributing

Developing Zazzles: use the [contributor guide](docs/contributing.md), [AGENTS.md](AGENTS.md), and
[developer agent orientation](.zazz/docs/zazzles-agent-orientation.md) for workspace conventions, branches,
commands, and the internal documentation map. Run relevant scoped checks;
CI verifies formatting, vet, build, and race tests. Work stays on feature branches, draft PRs target
`stage`, and humans review and merge. Live model calls are separate from offline tests.

Internal engineering material lives under `.zazz/`: [architecture](.zazz/architecture.md),
[implementation plan](.zazz/project/implementation-plan.md), [contracts](.zazz/specifications/README.md),
and the [engineering documentation map](.zazz/docs/zazzles-agent-orientation.md#documentation-map).

Analysis contributors: read the [Zazzles analysis rules and guidelines](docs/analysis/README.md) for
SQL/NoSQL and XSS judgment, input-validation policy, and calibration controls. Comparison runs use the
[Bearer runbook](docs/samples/nodejs-goof/bearer-runbook.md) and the
[Semgrep CE cheat sheet](docs/semgrep-ce-cheat-sheet.md).
