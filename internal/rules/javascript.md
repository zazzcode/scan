id: javascript
version: 1
languages: javascript, typescript

# JavaScript and TypeScript practice

TypeScript types are not runtime validation. The same rules apply to `.js`, `.jsx`, `.ts`, and `.tsx`.

## Posture

Report a confirmed injection and weak or poor practice. An agent can harden both cheaply. Do not treat a finding as out of scope only because exploitability is unproven.

## Validation

Restrict each field as far as its business requirement allows. A missing check is insufficient. A check that only shows the value is a string, a number, or shorter than a limit is insufficient when a tighter check would still be valid input. Middleware counts only when its body shows the check. A function name, a comment, or a call whose body is not shown is not validation.

## Free text

Comments, messages, descriptions, reviews, feedback, notes, and similar fields may contain any character. Do not require an allow-list for them. The injection control is the sink.

## Sinks

Do not build a query, update, or insert by concatenation, template interpolation, or a `$where` string. Bind the value with the driver's parameter API. A bound value is sufficient protection for statement structure even when validation is weak. A type or length check does not make an interpolated statement safe. Logging a SQL-looking string is not a query. `Array.find` is not a NoSQL query.

## Do not claim

Do not invent a guard, a binding, or a branch that the evidence does not show. Unknown middleware and unknown custom checks stay uncertain.
