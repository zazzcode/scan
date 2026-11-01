# scan

`scan` is a Go static application security scanner under active development. It aims to help
developers and coding agents understand how request inputs reach security-sensitive operations,
what validation and protection already hold, and where code needs hardening. Reports should give
reviewers source evidence and explicit analysis gaps so they can judge each finding.

The approach combines Tree-sitter source extraction with code-flow analysis and optional decision-model
assessment of relevant evidence. Field validation and sink protection are separate questions: a valid
email address does not establish that a SQL query uses bound parameters. Unknown behavior stays
unknown rather than becoming a safety claim or a vulnerability by assumption.

## Goals and current scope

The goal is to report vulnerabilities, weak security practices, and sufficient protection by route,
with useful explanations and remediation guidance. Longer term, scan aims to match or exceed Bearer's
SAST and developer-workflow capabilities with fewer false positives, measured at comparable detection
coverage on independently labeled cases. That is a project objective; broad parity is not established.

Today, scan analyzes local JavaScript and TypeScript source. It inventories Express-style route
registrations, follows a bounded subset of handlers and request-value dependencies into candidate SQL
and outbound HTTP uses, and exports source evidence as JSON, standalone HTML, or SARIF. Optional OpenAI Decisions or Jev
judging assesses validation and protection from that evidence.

Full control-flow and interprocedural analysis, independent findings review, Python and Go analysis,
and dependency vulnerability scanning remain planned. Recognizing a file extension in the file list
does not establish analysis support. A clean scan does not certify that an application is secure.

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
exists so you preserve your keys. Scan does not auto-load `.env`. Offline scanning needs no API key.
See the [setup guide](docs/getting-started.md) for full instructions.

### Optional decision-model setup

Use [decision-model setup](docs/decision-models.md) when explicitly enabling assessment. The scanner
reads the selected provider's key from its process environment. Model requests transmit source
evidence and need an explicit request cap; a key is not a reason to run paid assessment automatically.

## Run

```bash
export SCAN_FINDINGS_DIR="$HOME/Dev/zazzcode/scan/_test-finding"
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

Use the [contributor guide](docs/contributing.md) and [AGENTS.md](AGENTS.md). Run relevant scoped checks;
CI verifies formatting, vet, build, and race tests. Work stays on feature branches, draft PRs target
`stage`, and humans review and merge. Live model calls are separate from offline tests.

Internal engineering material lives under `.zazz/`: [architecture](.zazz/architecture.md),
[implementation plan](.zazz/project/implementation-plan.md), [contracts](.zazz/specifications/README.md),
and the [engineering documentation map](.zazz/docs/scan-agent-orientation.md#documentation-map).

Analysis contributors: read the [scan analysis rules and guidelines](docs/analysis/README.md) for
SQL/NoSQL judgment, input-validation policy, and calibration controls.
