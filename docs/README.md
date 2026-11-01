# zazzles documentation

`zazzles` helps developers and coding agents harden application code from source-backed security
results. Start with a local JavaScript/TypeScript application, produce JSON/HTML/SARIF reports, and
review vulnerabilities, poor practice, unresolved questions, and sufficient protection separately.
The scanner is under active development with bounded analysis; see the limits in the running guide.

The agent-first objective is practical: **fix established problems correctly**, including poor
practice that deserves refactoring without proof of an exploit. Keep the evidence, business policy,
and intended behavior visible so an agent can apply and verify an appropriate change.

## Table of contents

| Guide | What you will learn |
|---|---|
| [Sample reports](samples/README.md) | Juice Shop and TypeScript reports; Node.js Goof evaluation setup and Bearer baseline runbook |
| [Philosophy](philosophy.md) | Remediation-first decisions, supported findings, uncertainty, and preserving protection |
| [Getting started](getting-started.md) | Go/C prerequisites, compilation, `.env.sample`, and environment loading |
| [Running scans](running-scans.md) | CLI modes, route scope, report folders, JSON/HTML/SARIF, and current limitations |
| [Interpreting results](interpreting-results.md) | Findings, Review, Protected, Categories, counts, and agent remediation |
| [Analysis rules and guidelines](analysis/README.md) | Source evidence, SQL/NoSQL judgment, input validation, and calibration |
| [Category glossary](glossary.md) | SQL/NoSQL, validation, outbound requests, and planned categories |
| [Decision-model setup](decision-models.md) | Optional provider keys, model selection, capped assessment, and failures |
| [CI/CD artifacts](ci.md) | Saving and consuming reports in build pipelines |
| [Semgrep CE cheat sheet](semgrep-ce-cheat-sheet.md) | Running Semgrep Community Edition locally as a comparison scanner |
| [Contributing](contributing.md) | Development checks, feature branches, review, and documentation ownership |

## First steps

[Compile scan](getting-started.md#compile), [run an offline scan](running-scans.md#run), then open the
printed HTML file and expand a route. Agents can read the matching JSON for source evidence and
remediation context. The files are snapshots; regenerate views after changing the JSON or renderer.

This root `docs/` directory is for users and contributors. `.zazz/` contains internal engineering
design, project planning, evaluation records, implementation contracts, and SDLC standards.
