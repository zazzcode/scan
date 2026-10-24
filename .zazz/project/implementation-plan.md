# Proposed phased implementation plan

**Concept-stage planning, 2026-10-19. First-build pointer updated 2026-10-24.** Phases 0–13 below are
not approved for execution. This is a decomposition for later spec sessions, not a specification,
calendar schedule, or claim of completed functionality.

## First build

The [route-identification spec](../specifications/tree-sitter-route-identification-mvp.md) is greenlit
and implemented. `scan --routes` emits a JSON route index of Express-style registrations and same-file
handler names. It is not a resolved execution trace. The OpenAI Decisions API (`gpt-6-luna` on
`POST /v1/decisions`) remains the later assessor and is not part of this slice. Juice Shop manual
checks use tag `v20.3.0`. Phases 0–13 stay proposed. The next spec is a Go verdict for each external
string on a route: covered, unvalidated, unsanitized, or uncertain. That verdict is not this slice.
Context: [architecture](../architecture.md), [scope](02-project-scope.md),
[design proposal](03-jev-design-proposal.md), [Jev research](04-jev-validation-research.md), and
[Zazz methodology](../docs/zazz-methodology.md).

## Starting point and definition gate

Use the [repository-state record](../project.md#concept-stage-repository-state) for existing artifacts.
Phase 0 below audits their fitness; their presence does not establish an accepted implementation.

Every phase uses Go for scan-authored scanner, code-analysis, integration, and evaluation-harness code,
following the [implementation boundary](../architecture.md#1-architecture-premise). Target-language
fixtures and later Python/Go/Java adapters refer to languages being analyzed. Reference tools such
as IRIS inform the design without introducing their runtimes into scan.

Before the next implementation spec, select its JS/TS/framework slice, OWASP edition, analysis limits,
source-transmission policy, and evaluation criteria from the
[open decision register](03-jev-design-proposal.md#open-decision-register). Resolve only what the next
contract needs. Credential-dependent work is identified per phase; pending key status and setup
references live in the [research note](04-jev-validation-research.md#what-can-happen-before-a-key-exists).
Library candidates and selection ownership live in the
[proposal's library register](03-jev-design-proposal.md#library-choices). Owning specs record exact
versions, compatibility/build constraints, licensing, and verification before adding dependencies;
the pasted library examples do not constitute approved dependency selections.

## How a phase becomes agent-manageable work

Every spec is authored, refined, and implemented in its own sibling feature worktree, with its own
PR into `stage`. The definition-only `stage/` exception does not apply to specs. See the
[orientation](../docs/scan-agent-orientation.md#branches-and-worktrees).

- Suggested spec names below describe future contracts; they are not files created or greenlit now.
- Each spec has one reviewable outcome and its own scope, non-goals, inputs/outputs, dependencies,
  acceptance criteria, mapped tests/evidence, standards, and halt conditions.
- A phase containing multiple specs is implemented one contract at a time. Further split any contract
  that mixes independent outcomes or cannot produce a compact, verifiable handoff.
- Long-running execution resumes from a gitignored run log with current branch/spec, completed
  acceptance evidence, changed artifacts, failed checks, unresolved decisions, and next bounded step.
  Persist checkpoints before context loss; no chat-only progress state.
- Downstream work starts from accepted dependency contracts. Shared IR/schema changes require owner
  steering and updates to dependents, not silent drift.
- Ownership, greenlight, review, and storage follow the
  [methodology](../docs/zazz-methodology.md#specification-driven-delivery-cycle); branch/worktree policy
  lives in [AGENTS.md](../../AGENTS.md#branch-model).

## Acceptance targets

| Target | Role | Setup/limits |
|---|---|---|
| Small paired JS/TS corpus | Exact expected flows, guards, findings, and classification labels | Human-reviewed safe/vulnerable pairs; tuning and held-out partitions |
| OWASP Juice Shop | Realistic JS/TS repository and selected documented challenges | Clone supplied; pin the accepted revision and map selected challenges |
| SQL and MongoDB microfixtures | Compare sink/API-specific mitigation behavior | Isolated examples, not assumptions that SQL is safe or NoSQL unsafe |
| Hand-built IR/graph fixtures | Verify analysis independently of parsing/model output | Joins, cycles, guard bypass, unknown calls, call/return transfer, resource bounds |

Juice Shop is a [deliberately insecure Node.js/Express/Angular application](https://owasp.org/projects/juice-shop?tab=overview).
Map selected challenges to reachable code and expected findings at a pinned revision. Not every
challenge is statically detectable, and dependency/deployment/client/runtime-only issues may be outside
the chosen slice. Do not call total challenge count a scanner recall denominator. Keep expected labels
and challenge explanations out of model input. Paired cases live in our separate corpus rather
than modifying the reference checkout; selected built-in coding-challenge fixes can seed that corpus. Fixture access, licensing, inclusion/exclusion, and versioning
belong in a future spec. Local clone lookup and handling follow the
[orientation workspace conventions](../docs/scan-agent-orientation.md#portable-workspace-and-dev_root).

### Juice Shop case selection and paired fixes

Observed checkout revision: `399fba76194305a312693eee44f4df64bdf0c3ce`; the future corpus spec must explicitly
select its accepted revision and map each case to the documented challenge, root cause, source/value,
sink, relevant helpers, and expected result. Use the documented vulnerable case as the positive test,
then a reviewed protected counterpart and inadequate-fix variants as separate tests. Judge safety per
vulnerability/use; a SQL-safe fix is not proof of correct authentication or the absence of other flaws.

The clone includes coding-challenge snippets and explanation metadata under `data/static/codefixes/`.
For example, the [ignored replacement variant](https://github.com/juice-shop/juice-shop/blob/399fba76194305a312693eee44f4df64bdf0c3ce/data/static/codefixes/unionSqlInjectionChallenge_1.ts)
calls string replacement without using its returned value; the
[parameter-bound login fix](https://github.com/juice-shop/juice-shop/blob/399fba76194305a312693eee44f4df64bdf0c3ce/data/static/codefixes/loginAdminChallenge_4_correct.ts)
binds query values. Treat these as candidate cases, not automatically accepted labels. Verify code,
helper behavior, intended functionality, and sink-library semantics against their metadata.

Hide label-bearing filenames such as `_correct`, challenge IDs, vulnerability-marker comments,
solution explanations, and expected outcomes from model input while retaining a provenance mapping
outside the payload. Preserve real executable semantics and separately test misleading ordinary
comments. Partition by root cause/helper family, not challenge ID: several challenges may share one
underlying flaw. Include held-out original variants because public examples can be familiar to models.

The [prior-art study](05-prior-art-and-feasibility.md) informs parser/adapter/flow/rule/report specs:
Bearer for Go implementation patterns and CodeQL's open queries/libraries for semantic models and
regression cases. Keep these studies scoped to the next contract. Phase 1 may compare prepared
evidence with pinned Bearer results; treat human labels as the
reference, not Bearer output. This study does not change the chosen Jev Assessment → Findings Review roles.

## Dependency overview

| Phase | Depends on | Proposed outcome |
|---|---|---|
| 0 | Definition gate | Audit/retain bootstrap foundation |
| 1 | Definition gate | Corpus, question design, optional live feasibility evidence |
| 2 | 0; corpus shape from 1 | Parser choice, normalized IR, deterministic extraction |
| 3 | 2 | Resolution and security-role/library models |
| 4 | 3; labels from 1 | Useful deterministic rules and minimal finding output |
| 5 | 2–3; flow fixtures from 1 | Local CFG/dataflow and mitigation-use facts |
| 6 | 5 + resolution from 3 | Bounded cross-function flow and witnesses |
| 7 | 6; evidence needs from 1 | Complete, minimized evidence bundles |
| 8 | 7 + live feasibility from 1 | Jev Assessment, LLM Findings Reviewer, and reconciliation policy |
| 9 | 4, 6–8 | Unified reporting, identity, baseline, CI behavior |
| 10 | 9 | Correct incremental caches and resource/performance evidence |
| 11 | 10 + release-scope decision | Python adapter with its declared analysis coverage |
| 12 | 11 + separate language decisions | Go, then Java adapters in separate contracts |
| 13 | Separate proposal + relevant foundation | Optional narrow scope checks / coverage expansion |

Phases 4 and 5 are independent after their prerequisites. Offline phases 0–7 can progress after their
own greenlights without a Jev key. Phase 1 live feasibility and Phase 8 need credentials; absence of a
key does not justify inventing benchmark results or an automatic provider substitution.

## Ordered specification backlog

**These specs are planned, not authored or greenlit.** IDs below are planning references, not files.
Follow this recommended sequence, but start a spec only after its listed prerequisites are accepted.
Each row gets its own feature worktree and PR into `stage`, including investigation/experiment specs.
A successor branches from `stage` after its prerequisites merge; independent rows can proceed in
separate worktrees when authorized. Every implementation waits for its own spec greenlight.

| Order | Phase | Future spec | Prerequisites | Reviewable outcome |
|---|---|---|---|---|
| S01 | 0 | `foundation-audit` | Definition gate | Verified current behavior and gaps |
| S02 | 0 | `foundation-adjustments` | S01; only if needed | Bounded fixes to accepted foundation/toolchain |
| S03 | 1 | `corpus-and-juice-shop-map` | Definition gate | Pinned fixture scope and human labels |
| S04 | 1 | `semantic-question-design` | S03 | Jev validation/sanitization questions and reviewer rubric |
| S05 | 1 | `semantic-feasibility-eval` | S04; credentials/budgets | Measured Jev Assessment and Findings Review on prepared evidence |
| S06 | 2 | `parser-binding-evaluation` | S01, S03; S02 if needed | Binding/grammar/build decision |
| S07 | 2 | `analysis-ir-contract` | S06 | Normalized operations, typed adjacency contract, identity, provenance, capability limits |
| S08 | 2 | `jsts-cst-extraction` | S07 | Deterministic syntax extraction |
| S09 | 3 | `jsts-module-call-resolution` | S08 | Supported import/call/parameter/return resolution |
| S10 | 3 | `security-model-catalog` | S09 | Source/sink/guard and library models |
| S11 | 4 | `pattern-rule-contract` | S10, S03 | Rule metadata and evidence schema |
| S12 | 4 | `initial-pattern-rules-and-output` | S11 | Offline rule findings, marked not yet reviewed |
| S13 | 5 | `local-cfg` | S08, S10 | Control-flow blocks and branch/loop behavior |
| S14 | 5 | `local-value-flow` | S13 | Value transfer, reassignment, joins, and unknowns |
| S15 | 5 | `guard-use-and-mitigation-facts` | S14 | Guard/value/sink relationship and bypass evidence |
| S16 | 6 | `function-dependency-summaries` | S09, S15 | Bounded argument/return dependency summaries |
| S17 | 6 | `cross-function-witnesses` | S16 | Multi-file traces, recursion limits, gaps |
| S18 | 7 | `candidate-evidence-contract` | S17, S04 | Versioned assessment/review evidence contract |
| S19 | 7 | `bounded-bundle-builder` | S18 | Complete bounded bundles and inspector output |
| S20 | 8 | `jev-typed-adapter` | S19, S05 | Production Jev Assessment adapter and failure states |
| S21 | 8 | `llm-findings-reviewer` | S20 | Independent Findings Review and bounded reconciliation |
| S22 | 8 | `assessment-review-policy-and-evaluation` | S21, S12 | Integrated pipeline, human-label evaluation, audit/completeness policy |
| S23 | 9 | `finding-identity-and-grouping` | S22 | Stable identity across detector/reviewer assessments |
| S24 | 9 | `report-and-scan-comparison` | S23 | Unified findings/review/coverage and compatible-scan comparison |
| S25 | 9 | `baseline-and-ci-policy` | S24 | Reviewed baseline rules and incomplete-review CI behavior |
| S26 | 10 | `dependency-aware-cache` | S25 | Cache correctness across source/model/policy changes |
| S27 | 10 | `resource-and-performance-validation` | S26 | Measured latency, resource, spend, and cancellation limits |
| S28 | 10 | `release-packaging-and-readiness` | S27; release-scope decision | Reproducible selected-platform release evidence |
| S29 | 11 | `python-extraction-and-resolution` | S28; language-scope decision | Python adapter extraction/resolution |
| S30 | 11 | `python-flow-and-rules` | S29 | Python flow, rules, and assessment/review evaluation |
| S31 | 12 | `go-resolution-and-lowering` | S30; Go-scope decision | Go receiver/interface/lowering subset |
| S32 | 12 | `go-flow-and-rules` | S31 | Go flow, rules, and assessment/review evaluation |
| S33 | 12 | `java-resolution-and-lowering` | S32; Java-scope decision | Java types, annotations, overloads, lowering |
| S34 | 12 | `java-flow-and-rules` | S33 | Java flow, rules, and assessment/review evaluation |

S05's live experiment can remain pending while S06–S19 progress after their own greenlights. S12 and
S13 start independent rule/flow lanes after shared prerequisites. The initial release is evaluated at
S28; later language rows are a proposed expansion sequence, not initial-release requirements. SARIF
and Phase 13 capabilities need separately scoped backlog entries when approved.

## Phase 0 — Foundation audit and tooling decision

**Future specs:** `foundation-audit`; `foundation-adjustments` only if the audit identifies needed work.
Inspect the existing CLI/walker/registry, module/toolchain declaration, and CI against the refined scope.
Verify behavior and keep useful source. Decide any missing ignore/selection behavior explicitly.

**Evidence/exit:** reproducible manifest and code-check results, gap list, retained behavior, and
owner-approved toolchain/platform assumptions. Existing manifest exit codes are documented.
**Exclude:** parser, analysis, model integration, and speculative package scaffolding.
**Halt/replan:** toolchain unavailable or current behavior conflicts with a necessary contract.

## Phase 1 — Ground truth and semantic feasibility

**Future specs:** `corpus-and-juice-shop-map`; `semantic-question-design`; `semantic-feasibility-eval`.
These are separate offline and live deliverables. Start with a few named SQL/command/path cases; expand
NoSQL/SSRF and other families only as agreed. Include safe, vulnerable, and intentionally unknown cases,
ignored validator results, bypass branches, custom helpers, and misleading comments. Label validation,
context-specific protection, and actual guard use separately.

**Evidence/exit:** human-reviewed expectations, pinned reference revision and scope, held-out partition,
hand-prepared evidence bundles, versioned questions, and a declared evaluation protocol. Live spec sets
error/uncertainty/cost budgets before execution; report false-safe and false-alarm rates by family.
Experiment with Jev Assessment of validation/sanitization, then Findings Review of identified issues.
Use diagnostic deterministic/Jev-only/LLM-only baselines on the same evidence to measure contribution.
Audit Jev-protected cases for missed issues; reviewer agreement is not a human label.
**Exclude:** production scanner implementation and API requests in the offline specs.
**Halt/replan:** missing key for live work, contested labels, poor mitigation discrimination, or unclear
question semantics. Preserve offline artifacts; narrow the proposal rather than forcing a model result.

## Phase 2 — Parser decision, IR contract, and extraction

**Future specs:** `parser-binding-evaluation`; `analysis-ir-contract`; `jsts-cst-extraction`.
Evaluate selected grammars/bindings on real syntax and the intended build matrix. Define normalized
operations, values, blocks, source spans, provenance, identities, and serialized artifact versions.
Define scan-owned typed adjacency and forward/reverse relation indexes using the
[graph decision](../architecture.md#graph-storage-decision). Then extract declarations, imports, calls,
assignments, branches, and anonymous handlers. Verify the chosen parser's compatibility rather than
reopening library choices without a concrete blocker.

**Evidence/exit:** decision record, extraction goldens, malformed/partial-source diagnostics, native
resource lifecycle evidence, stable serial/parallel artifacts, and identity collision/edit cases.
Graph fixtures preserve parallel relation kinds, self-edges/cycles, and source provenance; initial
resource measurements establish a baseline for the later performance phase.
**Exclude:** resolved execution claims, taint propagation, and model calls.
**Halt/replan:** binding portability is incompatible with release needs or the IR cannot represent the
agreed local-flow subset; revise the contract before more adapters depend on it.

## Phase 3 — Resolution and source/sink/guard models

**Future specs:** `jsts-module-call-resolution`; `security-model-catalog`.
Resolve the declared aliases/re-exports/call patterns, registration-by-reference, parameters, and
returns. Classify source/sink/guard candidates with versioned query/rule/library data. Unknown targets
remain visible; imported package names alone do not prove API identity.

**Evidence/exit:** fixture source→call→sink relationships with exact binding evidence, unknown-edge
accounting, shadowing/alias negatives, and reviewed library-model preconditions. Map the selected
Juice Shop slice without claiming all application paths are resolved.
**Exclude:** a universal dynamic-language resolver and automatic sanitizer acceptance.
**Halt/replan:** false target resolution or framework assumptions invalidating the selected corpus.

## Phase 4 — Deterministic rules and minimal reports

**Future specs:** `pattern-rule-contract`; `initial-pattern-rules-and-output`.
Define reviewed rule metadata, versioned OWASP/CWE mapping, severity, and evidence requirements.
Implement a small rule set with safe/vulnerable pairs and minimal stable JSON findings. Rules needing
flow wait for the relevant supported-flow phase; do not fake it with string matching.

**Evidence/exit:** useful no-key/offline results, per-rule precision/recall on the declared corpus,
repeated-run equality, and evidence-backed positive/negative rule cases.
**Exclude:** broad OWASP completion claims, final baseline lifecycle, and semantic inference.
**Halt/replan:** a “pattern” verdict requires missing flow or runtime facts; narrow or defer that rule.

## Phase 5 — Local control flow and dataflow

**Future specs:** `local-cfg`; `local-value-flow`; `guard-use-and-mitigation-facts`.
Lower the agreed JS/TS subset into blocks and transfers. Define assignment/reassignment, expression
uses, branch joins, bounded loop handling, and unknown-operation behavior. Determine whether the
relevant validated/transformed value actually reaches a sink under the required guard outcome.

**Evidence/exit:** hand-built IR and parsed-source cases agree; guard-after-sink, bypass branch,
ignored return, changed value, cycles, and bounds are covered. Emit source-to-sink witnesses and gaps.
**Exclude:** broad alias/heap analysis and unsupported asynchronous semantics.
**Halt/replan:** proposed mitigation clears taint without value/control justification or loops fail to
converge/terminate under the declared bound.

## Phase 6 — Bounded interprocedural flow

**Future specs:** `function-dependency-summaries`; `cross-function-witnesses`.
Connect supported argument/parameter and return transfers. Define summary context, recursion/call
bounds, unknown/external calls, cancellation, and representative witness extraction/deduplication.
Report approximations and limits instead of enumerating every possible execution path.

**Evidence/exit:** multi-file vulnerable/safe fixtures, recursive/cyclic cases, helper sanitizer use,
unknown dispatch, and bounded runtime/memory. Every witness edge identifies its justification.
**Exclude:** exhaustive interprocedural alias solving and unstated framework behavior.
**Halt/replan:** limits silently lose expected findings or summaries conflate incompatible contexts.

## Phase 7 — Evidence bundles

**Future specs:** `candidate-evidence-contract`; `bounded-bundle-builder`.
Build minimal candidate bundles with snapshot/hash-bound spans, witness operations, branch/guard
facts, off-path helper definitions, sink semantics, library versions, and explicit missing evidence.
Add payload budgets and inspector output before any production semantic calls.

**Evidence/exit:** human review can reconstruct the case using only the bundle; missing/oversize
context is explicit; stable bundle hashes and allowed-source/redaction behavior are verified.
**Exclude:** classifier transport, free-form repository uploads, and generated citations.
**Halt/replan:** necessary evidence cannot fit or the source-transmission policy is unresolved.

## Phase 8 — Jev Assessment and Findings Review

**Future specs:** `jev-typed-adapter`; `llm-findings-reviewer`; `assessment-review-policy-and-evaluation`.
Choose Jev transport/model and validate typed IDs/options/ranges with bounded requests, deadlines,
retries, cancellation, budgets, and no-key/failure states. Add the LLM Findings Reviewer as a required
stage in the full pipeline. It independently inspects source evidence for every identified issue and
uncertain candidate before comparison with the Jev/rule assessment. Go retains separate results,
validates evidence IDs, and handles disagreement or missing context explicitly.

**Evidence/exit:** local HTTP stubs cover both adapters and partial failures; held-out human-labeled
evaluation measures validation/sanitization recognition, reviewer corrections and incorrect agreement,
false-safe risk, uncertainty, citation validity, review completeness, latency, and spend. Audit a sample
of Jev-protected candidates. Diagnostic baselines measure the roles; they do not turn the reviewer into
an optional production tier. No issue is labeled reviewed unless its review completed.
**Exclude:** autonomous graph exploration, generated fixes, and automatic safety acceptance/suppression.
**Halt/replan:** credentials unavailable, model changes, budget exceeded, evidence missing, or the quality
bar is not met. Keep incomplete stages and unresolved cases visible; reopen scope with the owner.

## Phase 9 — Unified reporting, identity, baseline, and CI

**Future specs:** `finding-identity-and-grouping`; `report-and-scan-comparison`;
`baseline-and-ci-policy`. Add SARIF only under a separately scoped export spec.
Unify rule findings, flow candidates, semantic assessments, and coverage. Specify collision/move
behavior, grouping, compatible-scan comparisons, model/policy changes, and incomplete-scan exit policy.
Keep baseline ownership/reason review separate from model assessment.

**Evidence/exit:** repeat scan has no structural identity churn; assessment changes remain visible;
new sink/rule is reported correctly; incomplete scans never silently resolve missing findings;
no-key and mixed-failure reports expose incomplete Findings Review; CI behavior matches reviewed policy.
**Exclude:** a vulnerability lifecycle service and model-driven suppression.
**Halt/replan:** fingerprints collide or unchanged identity hides a material changed assessment.

## Phase 10 — Incremental correctness, performance, release readiness

**Future specs:** `dependency-aware-cache`; `resource-and-performance-validation`;
`release-packaging-and-readiness` for the explicitly selected release scope.
Cache parses/IR/summaries/bundles/raw model responses with complete dependency/version keys. Reapply
policy independently. Benchmark cold/warm runs and limits on pinned fixtures; choose supported builds.

**Evidence/exit:** cached/uncached equivalence after changed validators, imports, callers/callees,
model/question/rule/config versions; no live requests for unchanged eligible bundles; measured CPU,
RSS, latency and spend; documented rule/language coverage and limitations; reproducible packaging.
**Exclude:** graph databases and promises of complete Top 10 coverage.
**Halt/replan:** stale caches alter findings, resource bounds fail, or the chosen release scope lacks
acceptance evidence. Humans alone promote `stage` to `main` and tag a release.

## Phase 11 — Python adapter

**Future specs:** `python-extraction-and-resolution`; `python-flow-and-rules`.
Model the declared decorators, imports, `self`, and framework patterns. Reuse the core where semantics
fit; record IR gaps and amend the governing contract before changing shared semantics.

**Evidence/exit:** Python safe/vulnerable and unknown cases pass through the agreed pipeline; coverage
matrix and bundle quality are reviewed; re-evaluate model questions on Python examples.
**Exclude:** automatic parity with all JS/TS constructs or a promise of zero core changes.
**Halt/replan:** the adapter would silently reinterpret core operations or claim unsupported dispatch.

## Phase 12 — Go and Java adapters

**Future specs:** `go-resolution-and-lowering`; `go-flow-and-rules`;
`java-resolution-and-lowering`; `java-flow-and-rules`. Each language is a separate decision and sequence.
Go stresses receivers/interfaces; Java stresses overloads, annotations, class/type boundaries.

**Evidence/exit:** language-specific resolution/flow fixtures, signature identity tests, explicit
unknowns, rules and semantic-question evaluation. Release support is declared per proven capability.
**Exclude:** complete interface/reflection resolution and blanket cross-language parity.
**Halt/replan:** ambiguous target selection, unsupported semantics, or language-specific evaluation gaps.

## Phase 13 — Optional scope checks and broader coverage

**Future specs:** one proposal and corpus contract per capability, followed by bounded implementation
specs. Revisit authorization/ownership, middleware ordering, logic constraints, additional framework
rules, or other languages based on measured user value.

**Evidence/exit:** enumerated questions, explicit required business/framework facts, stable anchors,
positive/negative corpus, uncertainty routing, and evidence of benefit beyond path detection.
**Exclude:** open-ended “review the application for security” as an unbounded model task.
**Halt/replan:** an expected answer requires business requirements or runtime facts absent from input.

## Progress and decision records

The phase list contains no checked-off security work. Update completion only from accepted spec
outcomes and evidence. Durable conclusions go into the proposal/architecture/scope docs; run logs and
resumption notes stay in `.zazz/ephemeral/`. Measure filtered/unexamined candidates alongside findings,
so cost reduction cannot conceal lost coverage. Revisit scope after each accepted capability instead
of treating this long-term plan as fixed.
