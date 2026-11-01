# Input validation analysis

[Analysis index](README.md) · [Database injection](database-injection.md)

A missing validator call is evidence about the analyzed slice, not proof that every execution path
accepts invalid input. Analysis MUST distinguish an absent observed check, an unresolved check, and
an established failure to enforce a declared field requirement.

## Establish the field policy

Use requirements demonstrated by the application or an explicit scan check: required presence,
scalar type, allowed values, format, length, range, or cross-field business constraints. Field names
alone do not define policy. Free-text punctuation is not intrinsically invalid.

A **Poor Input Validation** finding requires both an applicable requirement and evidence that it is
not enforced on a supported path. When business policy is missing, report the precise question as
Review / Inconclusive instead of inventing a restriction.

## Establish enforcement

Analysis MUST inspect what a validator does, whether its result is enforced, and whether execution
can reach the sensitive operation after failure. Include applicable middleware in execution order,
helper bodies, returned/reassigned values, and relevant branch/return facts. A function name or a
call receiving the field does not prove validation. An unresolved check is not a proven bypass.

Native Number conversion can constrain injection syntax without checking range, finiteness, or the
field's business validity. JSON parsing can accept nested objects without validating their schema.
SQL binding does not validate email format; an email-format check does not replace query binding.

## Report for a correct fix

Name the field, requirement, relevant path, and evidence of the enforcement defect. Keep destination
protection separate. Recommend changes that preserve valid application behavior, with normal,
rejected, boundary, protected, and inconclusive cases where relevant.

Current validation answers are facets of use records. Separate established per-field defect counting
is incomplete; do not add overlapping category pills together as independent findings.
