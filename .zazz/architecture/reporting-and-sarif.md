# Route-oriented findings and SARIF export

**Design and bounded implementation, 2026-10-31.** Use one versioned internal
model for route analysis, evidence and findings, with SARIF 2.1.0 as an implemented findings
export. Avoid separate detection/policy implementations for JSON and SARIF.

## Canonical report

The target repository is identified by snapshot and declared scope. Route records group by method,
path/receiver/mount identity and connect to handlers, input definitions, transfers, guards and sensitive
uses. Findings anchor to specific use/constraint pairs, so SQL injection and email-format validation
remain distinct even on the same route. Unknown/partial routes remain visible as coverage records.

Keep normalized findings separate from the full analysis graph. The current evidence JSON is a
bounded precursor, not a finalized public vulnerability schema. A later schema contract must settle
compatibility and edit-stable identity before baseline/CI workflows depend on it.

Proposed records:

| Record | Content |
|---|---|
| Target / scan scope | Source hashes, revision if available, selection, tool/model versions and exclusions |
| Route | Method/path, receiver/mount/handler provenance and applicable chain |
| Use / witness | Source and sink argument roles, transformations, control facts and assumptions |
| Finding | Rule/constraint, route/use reference, CWE/OWASP edition, severity policy and evidence status |
| Judge assessment | Provider/model/question versions, typed answers, usage and completion state |
| Completeness | Required facts, relevant blockers, unrelated coverage gaps and bounds |
| Review / runtime evidence | Separate independent review and controlled runtime confirmation states |
| Recommendations | Multiple complementary changes, scope/preconditions, evidence and regression expectations |
| Dependency finding | Package/version/advisory/dependency path and separate route-reachability status |

Recommendations should be organized around a route and its relevant graph slice. For login, an
appropriate set can include SQL parameter binding, server-side input type/format constraints, suitable
password storage, and authentication regression tests. Each addresses a different property. Do not
label email validation alone a sufficient SQL-injection fix or guarantee an impervious route.

## SARIF mapping

The exporter follows the
[OASIS specification plus errata](https://docs.oasis-open.org/sarif/sarif/v2.1.0/sarif-v2.1.0.html).
The current mapping is below; specific downstream consumer support has not been tested.

| Scan concept | SARIF representation |
|---|---|
| Detector / constraint definitions | Tool driver rules, rule identifiers and help |
| Source-backed finding | Result with primary location and related locations |
| Supported possible source-to-sink path | codeFlows / threadFlows with source locations |
| Route grouping | Logical locations plus namespaced scan properties carrying method/path/route ID |
| Stable identity | Snapshot-local IDs only; edit-stable fingerprints/baselines remain planned |
| Graph views | Local control edges in properties; full directed graph export remains planned |
| Concrete verified patch | fixes only when an actual edit is available; prose options remain recommendations/help |
| Judge/completeness/review details | Namespaced properties linked to the canonical evidence |
| Execution failure / partial scope | Invocation notifications and explicit completeness properties |

SARIF supports paths and directed graphs; a custom JSON format is not required merely to exchange
findings. The internal graph remains useful for analysis ownership, indexed traversal, caches and
rich route views. Export only justified code flows; do not turn an unordered dependency set or a
truncated candidate into a fabricated execution path. Preserve approximations and missing facts.

SARIF can carry source and dependency findings. SCA mapping must retain package/advisory identity
and manifest/lockfile locations without inventing a source-code sink. Route reachability is an
additional assessment, not the advisory's existence.

## Implemented interface and remediation

`--strings --sarif` exports a current scan. `--export-sarif <saved-report.json>` converts a canonical
JSON report offline with no new judge request. These flags are mutually exclusive with `--json`.
The exporter does not recompute or upgrade saved assessments.

Results map supported_potential to fail/warning, protected_for_claim to pass/none, needs_evidence
and incomplete assessments to review/note, and not_applicable to notApplicable/none. A pass applies
only to the declared claim. Validation and use-protection answers remain separate in messages and
properties: safe SQL parameterization can coexist with missing domain validation.

Rule help and per-result recommendations include the project tenet that findings cover confirmed
vulnerabilities and weak or poor practice, because an agent can harden both cheaply, plus
complementary actions and **illustrative, unverified
code templates** where the supported library model permits them. They are not patches and do not
appear in SARIF fixes. The Decisions choice API supplies typed assessments; it does not generate code.
A future optional Responses generation stage could produce tailored examples using separately capped
Luna requests, but must preserve assumptions and require checking before any applied fix. That stage
is not implemented or used in this evaluation.

## Delivery and verification

The current exporter was validated offline against the official SARIF 2.1.0 schema on the two Juice
Shop findings and vulnerable/protected/unknown comparison reports. Unit tests cover statuses, escaped
locations, preserved incomplete assessment, deduplicated related locations and absence of fabricated
registration value paths or unverified fixes. No GitHub/IDE ingestion test has been performed.

Phase 9 still owns the stable normalized finding contract, identity/baselines, CI policy and consumer
integration. This bounded adapter does not complete that phase. Neither JSON nor SARIF changes
finding eligibility, runtime-confirmation state or independent-review policy.
