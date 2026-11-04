# Node.js unary gRPC handler adapter — deliverable contract

Worktree / branch: `common-entry-and-dispatch`. Feature: common entry/dispatch.
Milestone deliverable: 7 of 8. Project and Deliverable Owner: Michael.
Drafted: 2026-11-03. Execution tracking: deliverable-specific ignored local run log only.
Coordination: one Contributor Agent working in the assigned worktree; no delegation required.
Proposed review shape: one milestone draft PR into `stage`, human review/merge.
The user requested design/specification now and will request implementation later.
This document does not itself authorize implementation, publishing, paid calls or merging.

## 0. Capability

Bind bounded grpc-js unary service/method registrations to JS/TS handlers and request inputs, reuse existing sink/claim analysis, and produce gRPC-specific reporting. Inventory unsupported streaming methods with explicit gaps.

## 1. Required reading

- [AGENTS](../../AGENTS.md#authorization-and-delivery) and
  [orientation](../docs/zazzles-agent-orientation.md#documentation-map).
- [Milestone index](common-entry-dispatch.md#delivery-order) and [Design §§4–8](../architecture/common-entry-dispatch.md#4-adapter-to-flow-boundary), [trace policy](../architecture/trace-judgment.md#separate-judgments-and-gates).
- Dependency: 06-graphql; reuse common records without depending on GraphQL semantics. [Previous deliverable](common-entry-dispatch-06-graphql.md)
- [Execution discipline](../docs/agent-execution-discipline.md#verification-discipline)
  and repository `.agents/skills/spec-driven/SKILL.md` for implementation.
- [Code structure §§1–6](../standards/code-structure.md#1-package-layout),
  [testing §§2–9](../standards/testing.md#2-test-kinds),
  [docs hygiene](../standards/docs-hygiene.md#cross-document-linking),
  [spec hygiene](../standards/spec-hygiene.md#acceptance-criteria-are-testable-from-the-criterion-alone),
  [PR process](../standards/pr-process.md) at publication.
  Recheck [standards index](../standards/index.yaml) against actual changed paths.
- Read the source files in §3 and adjacent existing tests before editing.
  Preserve current [analysis semantics](../../docs/analysis/README.md).

## 2. Invariants

- Existing supported HTTP and browser analyses MUST retain evidence, claim status and coverage semantics.
- Entry kind, dispatch mechanism, transport and implementation language MUST remain independent.
- Controls and input bindings MUST NOT cross invocation contexts; calls alone MUST NOT establish value flow.
- Unknown bindings and reached limits MUST remain explicit; a missing candidate MUST NOT imply protection.
- Source snapshots and historical raw answers MUST remain intact; changed requests MUST NOT reuse old answers.
- Only final scoped protection receives green styling. No entry-wide safety or runtime exploit claim is implied.

## 3. Scope and review shape

All edits MUST remain inside the assigned feature worktree. Allowed implementation paths
for this slice are listed below; focused adjacent consumers in the same named packages
may be updated where necessary for the same contract, with exact paths recorded in the
run log. New packages, framework subsets or public behavior beyond this contract require
an in-place contract revision. Unrelated cleanup is outside scope.

| Path | New / modified | Purpose |
| --- | --- | --- |
| `internal/parse/grpc_entries.go, grpc_proto.go, grpc_inputs.go, grpc_test.go` | new | Server/descriptor binding, proto3 subset and request roles |
| `internal/parse/analysis.go, modules.go, flow.go` | modified | gRPC discovery and language invocation seeding |
| `internal/parse/testdata/javascript/grpc/; internal/parse/testdata/typescript/grpc/` | new | Source-only proto/service/handler controls |
| `internal/report/entries.go; internal/judge/*_test.go; cmd/scan/analysis_test.go` | modified | Real adapter policy/report controls |
| `docs/running-scans.md; docs/analysis/README.md` | modified | Unary coverage and explicit generated/streaming limitations |
| `.zazz/specifications/common-entry-dispatch-07-grpc.md` | modified if needed | Contract decisions and evidence links |
| `.zazz/ephemeral/common-entry-dispatch-07-grpc-run-log.md` | ignored local | Append-only execution evidence |

One milestone PR preserves a single canonical migration and removes transitional state
before human integration. Each deliverable remains a distinct review/evidence unit in
that PR. Do not create additional worktrees or branches merely to represent checkpoints.
This proposed review shape remains reviewable by the user before implementation.

Exclusions: new vulnerability families, authorization engine, queues/jobs/CLI entry adapters,
new target languages, full IR, graph database, automatic cache, independent LLM reviewer,
paid calibration, application execution, sample remediation and unrelated formatting.

## 4. Decisions

**Decision 1.** Use @grpc/grpc-js plus @grpc/proto-loader literal loadSync/loadPackageDefinition and addService maps because they provide bounded source-visible descriptor and implementation bindings.

**Decision 2.** Support unary flow first because streaming requires separate event/value/control semantics; streaming declarations remain visible with a coverage gap.

**Decision 3.** Keep request fields, metadata and callback/context roles distinct because message types or a callable parameter named request do not establish source or protection semantics.

## 5. Agent implementation rules

Use the current assigned worktree as command cwd; verify path/branch/status before edits.
Explicit later implementation instructions authorize the requested slice under AGENTS;
do not request a separate ownership ceremony. Follow the milestone usage-window policy.
Record exact modified paths and compare the slice diff against its recorded start revision.

Keep coherent green recovery points. Commit/push only when authorized by the active user
request; verify the `zazzcode` identity/authentication before publication. Never merge,
push to `stage`/`main`, or create a release. Draft PRs target `stage` when publication is requested.
Do not use live keys or execute/install dependencies in target applications.

Halt dependent work for missing required reference data, a conflicting user edit, unresolved
scope/public-contract choice, destructive recovery, or an invariant violation. Investigate
routine failures autonomously; after three repeated failures without progress record the
blocker and surface it. Do not weaken tests or drop mandatory coverage to fit a window.

## 6. Acceptance criteria

- **AC1 — TestGRPCServiceBindings.** Verified construction/descriptor chain maps proto3 package/service/rpc declarations and literal implementation keys to actual handlers, accounting for supported method-name mapping with source provenance; separate server instances remain distinct. Verified by `TestGRPCServiceBindings` in the packages named in §7; new test names are planned contracts, not existing results.
- **AC2 — TestGRPCInputRolesAndControls.** call.request fields reach an existing SQL sink with supported transfers and separate metadata/context/callback roles; raw/bound/unknown controls exercise actual completeness and stub policy, including shared-callable distinct contexts. Verified by `TestGRPCInputRolesAndControls` in the packages named in §7; new test names are planned contracts, not existing results.
- **AC3 — TestGRPCUnknownAndStreaming.** Streaming, dynamic descriptors, missing/generated stubs, unknown interceptors, imports/include paths, method mismatches, fake addService and parser/expansion limits retain explicit located gaps without invented handlers/protection. Verified by `TestGRPCUnknownAndStreaming` in the packages named in §7; new test names are planned contracts, not existing results.
- **AC4 — TestGRPCPipeline.** Actual JS/TS unary adapter output reaches JSON/HTML/SARIF and gRPC selectors without method/path/legacy IDs; zero-use and streaming entries remain visible and repeated scans are deterministic. Verified by `TestGRPCPipeline` in the packages named in §7; new test names are planned contracts, not existing results.
- **Scope/checks.** Scoped tests pass and the slice diff contains only relevant approved paths;
  verify with `git diff --check`, `git status --short` and slice-diff inspection.

## 7. Test strategy

- `TestGRPCServiceBindings`: Verified construction/descriptor chain maps proto3 package/service/rpc declarations and literal implementation keys to actual handlers, accounting for supported method-name mapping with source provenance; separate server instances remain distinct.
- `TestGRPCInputRolesAndControls`: call.request fields reach an existing SQL sink with supported transfers and separate metadata/context/callback roles; raw/bound/unknown controls exercise actual completeness and stub policy, including shared-callable distinct contexts.
- `TestGRPCUnknownAndStreaming`: Streaming, dynamic descriptors, missing/generated stubs, unknown interceptors, imports/include paths, method mismatches, fake addService and parser/expansion limits retain explicit located gaps without invented handlers/protection.
- `TestGRPCPipeline`: Actual JS/TS unary adapter output reaches JSON/HTML/SARIF and gRPC selectors without method/path/legacy IDs; zero-use and streaming entries remain visible and repeated scans are deterministic.

Use compact table-driven controls and actual parsed fixtures where adapters are involved.
Hand-built records prove model/presenter boundaries only; they cannot substitute for actual
adapter-to-candidate-to-claim/report acceptance. Reuse adjacent existing tests for unchanged
behavior and inspect golden diffs rather than adding tests that mirror helper calls.

Scoped command (declared Go toolchain, no live network/model credentials):

```sh
GOTOOLCHAIN=local GOPROXY=off GOCACHE="$PWD/.zazz/ephemeral/go-cache" go test ./internal/parse ./internal/judge ./internal/report ./cmd/scan
git diff --check
```

At the deliverable's completed boundary run race tests for affected packages. Full
race/vet/build and sample/report validation are mandatory for the final milestone; the
index defines the commands and each later code change invalidates affected prior evidence.
Local loopback stubs test policy/transport, not model accuracy. A sandbox bind denial is
an environment limitation to resolve, not a passing test or a reason to use a paid endpoint.

## 8. TDD entry point and execution sequence

First failing test: `TestGRPCServiceBindings` with the observable contract in §6.

1. Start with TestGRPCServiceBindings using a local proto3 unary Search RPC and a literal handler map.
2. Implement bounded proto3 package/service/rpc/message reader and verified loader/server bindings, preserving schema hashes/spans and scoped path rules. Resolve supported implementation method naming explicitly.
3. Seed call.request separately from metadata, context and response callback; reuse existing SQL/value semantics and invocation completeness.
4. Run raw/bound/unknown and streaming/descriptor negative controls through real exports. Review upstream API and fixture manifest versions; do not execute stubs or code generators.

Illustrative implementation skeleton; names/bodies may adapt to verified existing code:

```go
// grpc_entries.go: verified server + loaded service descriptor -> RPC/handler dispatch.
// grpc_proto.go: bounded proto3 service/message syntax -> records or located gaps.
// grpc_inputs.go: call.request/metadata/context/callback roles -> invocation seeds.
```

Recovery checkpoints: Descriptor/service binding; then unary source flow; then streaming/unknown/report controls.
A checkpoint MUST leave a coherent runnable state before being called green; otherwise
record exact partial changes and failing checks, without claiming this deliverable complete.

## 9. Definition of done

- [ ] Required source/standards read and current scope verified.
- [ ] All acceptance criteria have passing, cited evidence; affected regressions explained.
- [ ] Scoped tests/race checks and diff checks pass; only changed Go files formatted.
- [ ] Required offline reports/manual review performed where this slice changes output.
- [ ] Remaining gaps, exact artifact locations and next checkpoint recorded.
- [ ] No out-of-scope change, fabricated protection or stale-response replay.
- [ ] User receives the completed slice and any remaining unchecked acceptance items.

## 10. Open questions

No deliverable-specific blocker remains beyond user review of the proposed design/scope.
Node.js GraphQL is first new server support; Apollo and grpc-js are concrete bounded
implementation choices, adjustable by user steering. The user's future GraphQL test app
is deferred evaluation, not a dependency. Resolve any materially new choice before its
dependent code; do not block unrelated earlier deliverables on later evaluation data.

## 11. Run log protocol

Use `.zazz/ephemeral/common-entry-dispatch-07-grpc-run-log.md`, append-only.
Record standards verification, user decisions, start revision, phase
completions, commands/results, deviations, report paths/hashes, review/rework and recovery.
At session start read AGENTS, this contract, relevant design sections and the log; verify
Git state and choose the next incomplete checkpoint. If the ignored log is absent,
reconstruct from Git and tracked acceptance evidence; never infer success from missing notes.

## 12. Agent implementation prompt

```text
Implement only the user-authorized deliverable in the existing common-entry-and-dispatch
worktree. Read AGENTS.md, .zazz/specifications/common-entry-dispatch-07-grpc.md,
its linked design/dependencies, and the deliverable-specific ignored run log. Load spec-driven and
matched analysis/engineering standards. Work as one agent. Start with the first failing
acceptance test, keep HTTP/browser regressions green, and verify actual pipeline behavior.
Record context/coverage and exact evidence; use recovery checkpoints to fit the available
usage window. Do not begin another deliverable automatically when only this one is assigned.
Preserve sample apps, raw historical answers and user edits. No paid calls, target execution,
new worktrees, merge or integration push. Report acceptance results and next checkpoint.
```

## 13. Implementation and review change log

2026-11-03 — Initial contract from the user's design/spec request, Node.js GraphQL
priority, explicit browser-event/reporting requirement, retained HTTP sample compatibility,
and later application evaluation. Product implementation has not started.

2026-11-03 — User requested one ignored run log per deliverable; updated log paths
and resume prompts. This changes execution tracking only, with no product/test change.
