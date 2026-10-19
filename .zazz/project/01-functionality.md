# 01 — Proposed functionality

**Concept-stage proposal, 2026-10-19.** These are desired capabilities, not built or accepted behavior.
Technical context: [architecture](../architecture.md). Delivery boundaries:
[scope](02-project-scope.md) and [implementation plan](implementation-plan.md).

## Purpose and users

Given a source repository, produce a ranked list of potential security issues with source evidence,
analysis limits, and a clear account of how each result was reached. Developers inspect the evidence;
security engineers triage it; CI applies an explicitly configured policy.

## Detection modes

| Mode | Proposed behavior | Evidence and limits |
|---|---|---|
| Track A: deterministic rules | Match explicit known-risk conditions without Jev | Rule ID, supported constructs, source spans; reproducible does not mean certainly exploitable |
| Track B: flow plus Jev | Trace supported string dependencies to a sink; assess validation/sanitization adequacy | Witnesses, guard facts, separate validation/sanitization facets, typed assessments |
| Findings Review | LLM reviews every identified rule/Jev issue and uncertain candidate against source evidence | Supported, disputed, or needs more evidence; separate review status and cited evidence IDs |
| Structural-only operation | Run local analysis without credentials or external calls | Pattern findings and path candidates; semantic evaluation explicitly `not_run` |
| Track C: future bounded scope checks | Assess explicitly enumerated authorization/logic questions | Deferred; separate evidence contract, corpus, and specs required |

Jev is the primary semantic assessor; the LLM Findings Reviewer is the planned review stage for
identified issues. No mode claims complete OWASP Top 10 coverage. Reports preserve source evidence,
Jev Assessments, and reviewer conclusions separately. A skipped/failed review stays explicit.

## Functional domains

| Domain | Proposed capability |
|---|---|
| FD-1: input selection | Deterministic manifest, language/path selection, ignore policy, skipped-file accounting |
| FD-2: extraction | Tree-sitter CSTs and versioned query captures with exact source locations |
| FD-3: normalization | Language adapters emit operations, symbols, blocks, values, and provenance |
| FD-4: resolution | Resolve supported imports, call targets, handlers, parameters, returns, and framework registrations |
| FD-5: classification | Identify sources, sinks, and candidate guards through rules/library models |
| FD-6: local flow | Track declared assignments, expression dependencies, branch joins, and bounded loops |
| FD-7: cross-function flow | Propagate supported argument/parameter and return dependencies through bounded summaries |
| FD-8: deterministic rules | Detect tested patterns and rule-specific mitigation facts without inference calls |
| FD-9: path candidates | Produce bounded source-to-sink witnesses; deduplicate and account for unexamined candidates |
| FD-10: evidence preparation | Bundle snippets, spans, bindings, guards, sink semantics, and completeness markers |
| FD-11: Jev Assessment | Experiment with versioned Noul/Choice questions about string validation and sanitization; validate typed responses |
| FD-12: findings review and policy | LLM reviews identified issues and uncertain candidates; Go reconciles separate assessments, exposes disputes/incomplete review, and applies policy |
| FD-13: reporting | Versioned JSON and readable Markdown; possible later SARIF export |
| FD-14: finding identity | Stable identity independent of model predictions; group related source traces |
| FD-15: comparison and baseline | New/persisting/changed assessments, reviewable suppressions, explicit incomplete-scan semantics |
| FD-16: incremental rescan | Reuse unchanged artifacts with dependency- and version-aware invalidation |

Optional name/pattern-based sensitive-data tags may inform prioritization later. They are heuristic
signals and are not proof of sensitive runtime data.

## Initial rule families and coverage

The first feasibility slice focuses on a few JS/TS source-to-sink families. Expand only after measuring
safe/vulnerable pairs and unsupported constructs.

| Family | Candidate evidence | Mitigation distinctions to test |
|---|---|---|
| SQL injection | Source-derived query fragments, query API, parameter positions | Bound values versus dynamic identifiers or SQL fragments |
| NoSQL operator injection | Untrusted values/objects entering database filters | Scalar coercion, schema enforcement, operator/object allow-lists |
| Command injection | Shell command construction and execution mode | Shell interpolation versus argument-vector execution; executable/option control |
| Path traversal | Untrusted path parts reaching filesystem operations | Normalization versus enforced base-directory containment |
| SSRF | Untrusted input affecting an outbound request target | Host/scheme restrictions, redirects, unresolved runtime/network behavior |
| Unsafe deserialization / HTML injection | Input reaching a relevant API or output context | Type- and context-specific protection |
| Pattern-only configuration / crypto | Explicit risky API/configuration use | Precisely documented rule preconditions and counterexamples |

Not every occurrence of `exec.Command`, SQL construction, a regex, or a validation API is a
vulnerability or a mitigation. Rules declare what they establish and when the case needs more evidence.
Missing authorization, insecure design, and deployment-dependent configuration cannot be inferred
reliably from an arbitrary code pattern alone; broad coverage of those remains future work.

OWASP mappings include an explicitly selected edition in the future catalog. CWE mappings and severity
come from reviewed rule metadata. The old unversioned A01–A10 list is retired rather than carried
forward as a release commitment. Dependency vulnerability scanning is excluded.

## Language sequence

JS/TS first, then Python, then Go and Java. The sequence and first release set remain proposed.
Each adapter includes extraction, resolution, lowering/flow semantics, library models, and a coverage
matrix. Recognizing a file extension is not analysis support. Additional languages may expose core IR
changes that need a reviewed contract revision.

## Interfaces and report semantics

The only source-level CLI foundation currently present is `scan --list [--json] <path>`. All analysis
interfaces are future spec work. Inspection of graph, classifications, witnesses, and explanation bundles
is desirable, but exact flags/subcommands are not decided by this document.

Security reports should include detection mode, source/sink evidence, rule/CWE/OWASP metadata,
severity policy, assessment status, model/question versions when used, and coverage gaps. Evaluation
states distinguish `not_run`, `completed`, and `failed`; a completed semantic assessment can be
`uncertain`. Neither a request failure nor absent credentials produces a safe verdict.

Stable candidate identity is distinct from current risk assessment. A changed assessment remains
visible even when the fingerprint is unchanged. A missing result from an incomplete scan cannot be
reported as resolved. CI policy and future security exit codes need their own spec.

## Decisions before specs

The [proposal register](03-jev-design-proposal.md#open-decision-register) tracks first-slice selection,
semantic-policy requirements, CI behavior, and allowed source transmission. The
[evaluation protocol](04-jev-validation-research.md#proposed-corpus-and-evaluation-protocol) defines
how to assess the resulting capability choices.
