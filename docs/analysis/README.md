# Analysis rules and guidelines

[Documentation index](../README.md) · [Result interpretation](../interpreting-results.md)

These product documents explain how scan evaluates code and how contributors and coding agents
should improve its analysis. They belong outside `.zazz/`, alongside user documentation.

| Guide | Scope |
|---|---|
| [Database injection](database-injection.md) | SQL/NoSQL library contracts, source flow, branch constraints, and verdicts |
| [Input validation](input-validation.md) | Business requirements, enforcing checks, and separate destination protection |
| [Analysis changes and calibration](calibration.md) | General rules, positive/negative controls, and acceptance evidence |

The current scanner implements these rules in Go evidence assembly, question construction, and
classification. These Markdown files describe and govern that behavior; editing them alone does
not change a running binary or load a new model prompt. Keep implementation, tests, and these guides
aligned when changing a rule. Analysis decisions preserve the raw source and provider answers.

Sensitive-data logging analysis is the next planned rule family. Its implementation and calibration
are pending; a glossary entry is not detector coverage.

`.zazz/` owns internal architecture, specifications, project planning, evaluation history, and SDLC
standards. User instructions, product philosophy, analysis rules, and scan operation guides live in
root `docs/`. See [contributing](../contributing.md) for repository workflow.
