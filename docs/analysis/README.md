# Analysis rules and guidelines

[Documentation index](../README.md) · [Result interpretation](../interpreting-results.md)

These product documents explain how zazzles evaluates code and how contributors and coding agents
should improve its analysis. They belong outside `.zazz/`, alongside user documentation.

| Guide | Scope |
|---|---|
| [Database injection](database-injection.md) | SQL/NoSQL library contracts, source flow, branch constraints, and verdicts |
| [Input validation](input-validation.md) | Business requirements, enforcing checks, and separate destination protection |
| [Sensitive-data logging](sensitive-data-logging.md) | Logger bindings, data classes, content versus derived facts, and outcomes |
| [Cross-site scripting](xss.md) | Source-to-rendering evidence, HTML contexts, template escaping, sanitization and XSS outcomes |
| [Analysis changes and calibration](calibration.md) | General rules, positive/negative controls, and acceptance evidence |

The current scanner implements these rules in Go evidence assembly, question construction, and
classification. These Markdown files describe and govern that behavior; editing them alone does
not change a running binary or load a new model prompt. Keep implementation, tests, and these guides
aligned when changing a rule. Analysis decisions preserve the raw source and provider answers.

Sensitive-data logging analysis covers explicit logger calls in analyzed route handlers and morgan
request logging of URL query credentials, with a logging judgment protocol. Scoped live calibration and its limits are recorded in the
[Python evaluation](../../.zazz/evaluation/python-support-results.md).

`.zazz/` owns internal architecture, specifications, project planning, evaluation history, and SDLC
standards. User instructions, product philosophy, analysis rules, and zazzles operation guides live in
root `docs/`. See [contributing](../contributing.md) for repository workflow.

## Embedded Analysis Guidelines

The scanner also embeds versioned language **Analysis Guidelines** from
[`internal/guidelines/`](../../internal/guidelines/). These documents combine secure coding
recommendations with instructions for interpreting source evidence. The semantic judge receives
the full applicable document; reports retain its ID and version under `analysis_guidelines`.
They do not implement detection or override missing evidence, application requirements or Go policy
gates. Unlike the product guides on this page, their text is included in model requests.

“Poor practice” remains a finding classification for established defects. It is independent of the
former `Practice` name for guideline documents. See the [architecture terminology](../../.zazz/architecture.md#1-architecture-premise).

The [bounded Apollo GraphQL adapter](../running-scans.md#bounded-apollo-graphql-analysis) adds
resolver-specific argument sources to these existing analysis rules. Parent/context/info provenance,
SDL syntax coverage and application policy remain separate from SQL protection and claim status.

The [bounded Node.js gRPC adapter](../running-scans.md#bounded-nodejs-unary-grpc-analysis) joins
source-visible unary registrations to request-field inputs and existing SQL claims. Proto field
types/decoding do not establish validation or protection; metadata/context/callback roles remain
separate. Streaming and unresolved descriptor/handler bindings retain visible coverage gaps.
