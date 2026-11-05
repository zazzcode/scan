# Analysis pack MVP proposal

**Planning decision, 2026-11-05.** `analysis-pack-mvp` is a separate initiative, with its own
future worktree/branch, after `common-entry-and-dispatch` completes acceptance and merges into
`stage`. It precedes command injection and the remaining vulnerability-family expansion.
This proposal authorizes documentation only; implementation needs its own bounded contract and
explicit assignment. It is not another deliverable in the common-entry/dispatch milestone.

The [implementation plan](implementation-plan.md#analysis-pack-mvp-initiative) owns sequencing.
The [technical approach](../../docs/technical-approach.md) owns the public explanation of semantic
judgment. This proposal owns the proposed pack boundary, MVP scope and evaluation goals.

## Purpose and architectural distinction

> Zazzles establishes supported code relationships through static analysis, then uses a decision
> model to evaluate whether the demonstrated controls satisfy a specific security requirement.
> Explicit evidence gates constrain the resulting judgment.

The decision model evaluates the supplied implementation, operand roles and relevant paths. Within
supported analysis boundaries, that can avoid needing a named sanitizer rule for every custom
validation function. A supplied function body and actual transfers/conditions let the judge assess
whether the control is adequate for this operation. A function name alone cannot establish protection.

Semantic judgment reduces the need to enumerate every acceptable implementation pattern; it does
not eliminate the need for analysis models. Sources, handler bindings, transfers, sink argument roles
and library behavior still need correct static support. Missing discovery means the judge never sees
the candidate. Missing helper bodies, unresolved effects and reached limits remain explicit gaps.

Established SAST includes structured program models, framework knowledge and dataflow, not just
text matching. For example, [Semgrep taint rules](https://docs.semgrep.dev/writing-rules/data-flow/taint-mode/overview)
define sources, propagation, sanitizers and sinks and can produce traces. Zazzles's design places a
decision model within claim judgment, constrained by static evidence and Go policy. This is an
architectural distinction, not proof of uniqueness, broader coverage or superior accuracy.

## Current foundation and proposed change

The current scanner already contains the components below. They ship in the binary and are not
user-downloadable analysis packs.

| Component | Current owner | Proposed pack role |
| --- | --- | --- |
| Language/framework discovery and bounded flow | `internal/parse/` | Reference implemented capabilities and their supported subsets; executable adapters stay compiled |
| Verified library/API contracts | Evidence assembly and candidate API models | Version applicable contract content and associate it with verified bindings |
| Versioned language Analysis Guidelines | `internal/guidelines/` | Identify embedded guidance and select applicable components explicitly |
| Claim questions and answer protocols | `internal/judge/` | Reference the question set and interpretation guidance for the named claim |
| Required evidence and consistency gates | Parse completeness and Go judge policy | Declare obligations while retaining mandatory core gates |
| Positive/protected/unknown controls | Existing fixtures and tests | Identify regression coverage for each declared capability |

The MVP organizes and versions existing knowledge. It does not implement new languages, frameworks
or vulnerability families merely by declaring them in a manifest. A `--checks` file remains an
application-requirement/selection input; it is not a language adapter or general rule-pack loader.
Product analysis documentation remains distinct from embedded model guidance: editing a public guide
alone does not change detection or the running judge's instructions.

## Two consumers, one explicit selection

| Consumer | Material it uses | Boundary |
| --- | --- | --- |
| Static engine | Applicability metadata, references to compiled source/sink/binding/flow models and evidence obligations | Establish source relationships before security judgment; pack declarations cannot invent edges |
| Decision Model Judge | Selected guidelines, applicable API contracts, claim questions and interpretation guidance alongside the trace | Assess the supplied implementation for a named requirement; no missing evidence filled by confidence |
| Reporter and replay validation | Selected component IDs, versions/digests, supported scope and judgment provenance | Keep results auditable and prevent changed requests from inheriting historical answers |

A pack is not itself the trace. The trace contains application evidence; the selected pack material
supplies scanner-authored interpretation context. Target code, comments, manifests and repository
text MUST remain untrusted evidence, not instructions that select or redefine security policy.

For example, a JS/TS gRPC SQL candidate could compose language guidance, the verified unary gRPC
input capability, a Sequelize operation contract and the SQL-injection claim definition. Only the
applicable judgment material accompanies that candidate's source/transfer/control evidence.
The scanner does not send executable adapters, test labels, unrelated vulnerability guidance or an
entire pack catalog to the judge.

Language, entry framework, library and vulnerability family SHOULD remain independently composable.
The design SHOULD avoid requiring a separate monolithic pack for every combination. Actual
selection MUST depend on supported source facts or explicit scanner configuration, not framework
names or installed-package declarations alone. Unknown versions/applicability MUST retain limits.

## Bounded MVP scope

- Define an internal embedded manifest with stable pack/component IDs, versions, compatibility
  references, applicability and test-coverage references. Keep scope and known limits discoverable.
- Inventory existing language guidance, verified API contracts, claim questions and required evidence;
  give each an explicit owner and reference without rewriting the analysis engine.
- Select the relevant components deterministically for a candidate. Missing or conflicting required
  components MUST remain explicit configuration/coverage failures; no permissive silent fallback.
- Assemble a bounded judge request containing the trace plus only the applicable interpretation
  material. Mandatory evidence and consistency gates MUST remain in Go and cannot be weakened by
  pack content. A pack MUST NOT mark a custom helper safe merely by its name.
- Record selected versions/digests in report/judgment provenance. Exact effective request content
  MUST participate in replay validation. Changed guidance, contracts, questions, evidence or request
  protocols MUST NOT reuse old answers. A representation change that alters an envelope needs an
  explicit version/compatibility decision, even if intended security behavior is unchanged.
- Reuse existing fixtures and offline stubs to demonstrate preserved findings, scoped protection,
  unknowns, evidence and output semantics. Use command injection as the next consumer of the new
  organization after the MVP is accepted, not as hidden work inside this initiative.

The implementation contract will determine the minimal file/schema changes after a current-code
inventory. This proposal does not fix a public pack format, CLI flag, package API or new runtime
plugin interface in advance of that evidence.

## Deferred scope

Downloadable registries, third-party or target-supplied policy packs, executable plugins, arbitrary
runtime imports, a new rule language, Semgrep YAML compatibility, automatic remote updates and
unrestricted user-authored safety contracts are outside the MVP. External distribution needs its own
trust, compatibility, update and validation design once real consumers establish the extension points.

A whole-program IR, broader helper/heap/async coverage, independent Findings Review, automatic
judgment caching and new vulnerability/language support remain separate initiatives. No paid
calibration, target execution or sample-application changes are required for this packaging MVP.

## Proposed acceptance and evaluation goals

- Existing JS/TS/Python cases retain their source bindings, claim outcomes, blockers and coverage
  through actual adapter-to-judge-to-report pipelines; zero-use entries and unsupported paths remain
  visible. Allow only explicitly reviewed provenance/envelope changes in baseline comparisons.
- Applicability controls cover supported and unsupported combinations, verified versus fake APIs,
  missing/conflicting components and shared-callable contexts without borrowed protection.
- Raw, protected and unknown controls demonstrate that the judge receives the right guidance and
  contracts and that optimistic answers cannot override missing evidence. Expected labels remain
  outside model requests.
- Deterministic selection and exact-request replay controls reject changed pack/guidance/contract
  content. Saved reports and raw historical responses remain intact; version changes cannot relabel
  those responses as fresh decisions.
- JSON/HTML/SARIF retain result parity and distinguish offline discovery, stub assessment and live
  model assessment. Request-size/selection limits remain explicit instead of truncating obligations.
- Measure selection/payload size and local overhead. Broader precision/recall, reduced false findings
  or false protection, custom-validator generalization and live model cost remain separate evaluated
  hypotheses requiring independently labeled held-out cases and explicit call authorization.

The first recovery boundary should be an inventory/manifest with unchanged behavior, followed by
selection and request/report provenance, then full regression and compatibility evidence. No broad
engine rewrite or unrelated cleanup is part of those checkpoints.

## Delivery boundary and next step

Finish and merge the eight-deliverable common-entry/dispatch milestone first. Then assign a separate
`analysis-pack-mvp` contract and create its own worktree/branch from the merged `stage` baseline using
the repository worktree workflow. Scope that contract around the current implementation, acceptance
controls, request/report compatibility and recovery points before changing scanner behavior.

The documentation prepared on the common-entry/dispatch branch records this future initiative only.
It neither starts its implementation nor authorizes creating the new worktree now. Subsequent
vulnerability-family work should use the accepted pack structure; the rest of the roadmap keeps
its existing order. External pack distribution can be reconsidered after several families validate
that the internal boundaries are useful.
