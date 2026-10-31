# Route string validation — implementation contract

**Authorization:** explicit user request, 2026-10-31, to implement in the isolated
`route-string-validation` worktree. No further greenlight ceremony is required. The user also
requested provider-neutral judging and a graph design covering outbound service exchanges.

## Outcome

Add a bounded JS/TS route-to-string-use evidence mode alongside the unchanged route-v1 inventory.
Resolve direct named relative imports and handler factories for the selected Juice Shop login/search
cases. Track request-field dependencies through supported local expressions/assignments to candidate
SQL and outbound HTTP calls. Report source-backed evidence, transfer limits and uncertainty; never
label missing evidence protected. Make OpenAI judge access explicit and optional; retain Jev and
future adapters in the design without claiming unimplemented transports exist.

## Implementation boundary

- Tree-sitter identifies syntax; no regex security verdicts. A separate evidence schema preserves
  route, source, operation and candidate-use relationships and source hashes.
- The first slice supports local named ESM imports, direct/factory handlers and local expression
  dependencies. Unsupported branches, closures, wrappers, heap/cross-function transfers and external
  response continuations are explicit gaps, not completed analysis claims.
- SQL `.query` shapes are candidate sink models, not proven library identity or confirmed vulnerabilities.
  Outbound calls record request argument evidence; opaque response behavior stays unknown.
- `--strings [--json] [--entry <relative-file>] [--route <path-pattern>] <target-root>` runs offline by default. Optional
  `--judge openai` uses `OPENAI_API_KEY`; `--judge jev` uses `TYPESAFE_API_KEY`; each assessment preserves provider/model and failure state.
  No automatic provider fallback. Jev uses its verified provider-specific choice contract and requires an explicit model name.
- Models judge protection adequacy from bounded evidence. Offline output has `not_run` assessment;
  no protection or exploitability claim is inferred from missing validation alone.
- No live requests during tests or this session without credentials and a bounded evaluation scope.
  No target application execution or HTTP probing.

## Acceptance evidence

- Existing list/route tests remain unchanged and pass.
- Fixture routes connect named imported factories to request sources and candidate use spans.
- Supported assignment/concatenation/interpolation preserves dependencies; reassignment removes
  stale dependencies; unknown calls retain possible dependence with explicit uncertainty.
- HTTP-client route-shaped calls are not promoted to verified Express entry points in the new mode.
- Outbound argument evidence and response-flow limitations remain distinct.
- Offline runs are deterministic; parse/resolution/flow gaps and judge `not_run` are visible.
- Stubbed judge tests reject missing/duplicate/refused/invalid answers and non-success responses.
- CLI tests cover new mode, incompatible flags, bad entry paths and missing target roots.
- Manual Juice Shop verification records login/search candidates and remaining limitations.

## Design and operating references

[Graph design](../architecture/repository-graph.md), [architecture](../architecture.md),
[plan](../project/implementation-plan.md), [AGENTS.md](../../AGENTS.md),
[testing](../standards/testing.md), [code structure](../standards/code-structure.md).
Execution history belongs in `.zazz/ephemeral/route-string-validation-run-log.md`.

## Implemented evidence and known boundaries

The report is a route-linked collection of source-to-use dependency paths, not the complete canonical
IR/CFG described by the long-term graph design. Supported expressions yield local syntactic
relationships. Unknown transformations and branch merges retain possible dependencies with gaps;
those are not proven execution paths. Source size is capped at 2 MiB per loaded module, expression
work at 10,000 per handler, and dependency/transfer paths at 64 with explicit limit reasons.

Both OpenAI and Jev adapters share normalized validation/protection choice answers, model/provider,
question version and common usage counters. No model probability is used for policy. Legacy static gaps
force uncertainty; checked candidates use required roles and claim-relevant blockers. Independent review
and a vulnerability CI gate remain future work. The original connection test and the subsequent [two-case evaluation](../evaluation/juice-shop/results.md)
completed; no runtime target confirmation or held-out accuracy measurement was performed.

## Two-case implementation increment

Authorized user steering adds explicit method/path/field requirements via `--checks <json-file>`.
The tracked [phase cases](../evaluation/juice-shop/validation-checks.json) cover login email/SQL and
registration nonempty credentials. These are policy requirements, not expected judge answers.

- Attach applicable earlier middleware, same-route handlers, relevant input helper bodies and
  model/setter/hook context. Support the declared namespace/concise-factory and finale resource subset.
- Emit supported local branch/successor edges for registration middleware; preserve return/throw
  termination and source spans. Do not claim a complete interprocedural CFG.
- Deduplicate wire evidence and refer to it from dependencies, control edges and required fact roles.
- For checked candidates, require declared evidence roles and block on claim-relevant unknowns.
  Preserve unrelated coverage gaps. Legacy candidates remain conservative.
- Report supported_potential/protected_for_claim/needs_evidence separately from raw judge answers;
  no status implies runtime confirmation or independent review.
- Verify vulnerable, protected and unknown fixtures, missing bindings/guards, shadowed helpers,
  method-specific grouping, bounded relevant context and stable checked-evidence goldens.

The final two Luna calls returned insufficient validation/protection and supported_potential for both
cases. See [evidence design](../architecture/evidence-and-claims.md),
[login case](../project/08-login-validation-case.md), [registration case](../project/07-registration-validation-case.md)
and [results](../evaluation/juice-shop/results.md). No target code was modified. SARIF 2.1.0 export and illustrative remediation templates are implemented;
model-generated code, independent review and SCA remain planned work.

## Reporting and comparison increment

- `--sarif` exports strings output; `--export-sarif <saved-json>` converts offline without another API call.
- Preserve route identity, source locations, supported dependency paths, raw validation/protection
  facets, claim status, blockers, coverage and review/runtime state. Do not fabricate registration
  value flow or export unverified code examples as fixes.
- Verify official schema conformance and original vulnerable/protected/unknown controls with expected
  labels excluded from the submitted bundle. Preserve mismatches and evidence corrections in the
  [comparison results](../evaluation/corpus/README.md).
- Keep automatic novel-repository discovery as the next bounded Phase 3/5–8 extension in the
  [evidence design](../architecture/evidence-and-claims.md#discovery-on-novel-target-repositories).

## Implementation And Review Change Log

### 2026-10-31 — Route-focused reporting and original comparison controls

**Source:** user requests for iterative accuracy work, SARIF, remediation examples and novel targets.
**Changed sections:** reporting/comparison increment and implemented boundaries.
**Rationale:** distinguish actual graph evidence quality from judge confidence and known target labels.
**Summary:** add offline SARIF export/templates, original controls and verified import provenance in bundles.
**Verification impact:** schema checks, live capped Luna comparisons and offline export/status regressions.
