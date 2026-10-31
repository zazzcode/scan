id: python
version: 1
languages: python

# Python practice

This standard is applied when a Python file is assessed. The current route scanner does not yet emit Python sink candidates.

## Posture

Report a confirmed injection and weak or poor practice. An agent can harden both cheaply. Unproven exploitability does not remove the finding.

## Validation

Restrict each field as far as its business requirement allows. A missing check is insufficient. `isinstance(value, str)` or a length check is insufficient when a tighter check would still meet the business need. A name, a comment, or an unshown helper is not validation.

## Free text

Comments, messages, descriptions, reviews, feedback, and notes may contain any character. Do not require an allow-list. The injection control is the sink.

## Sinks

Do not build SQL or NoSQL text with f-strings, `%` formatting, `+`, or `format`. Pass values as DB-API parameters. A parameter is sufficient protection for statement structure even when validation is weak. A type or length check does not make a formatted statement safe.

## Do not claim

Do not invent a guard, a parameter binding, or a branch that the evidence does not show.
