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
request logging of URL query credentials, with a logging judgment protocol. Live calibration is pending.

`.zazz/` owns internal architecture, specifications, project planning, evaluation history, and SDLC
standards. User instructions, product philosophy, analysis rules, and zazzles operation guides live in
root `docs/`. See [contributing](../contributing.md) for repository workflow.
