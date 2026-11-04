# Bounded dynamic HTTP registration and dispatch — deliverable contract

Worktree / branch: `common-entry-and-dispatch`. Feature: common entry/dispatch.
Milestone deliverable: 5 of 8. Project and Deliverable Owner: Michael.
Drafted: 2026-11-03. Execution tracking: deliverable-specific ignored local run log only.
Coordination: one Contributor Agent working in the assigned worktree; no delegation required.
Proposed review shape: one milestone draft PR into `stage`, human review/merge.
The user requested design/specification now and will request implementation later.
This document does not itself authorize implementation, publishing, paid calls or merging.

## 0. Capability

Resolve finite Express literal-table/for-of registrations and direct literal-argument local factories through shared dispatch/context records, retaining unresolved alternatives, closure bindings and generation provenance.

## 1. Required reading

- [AGENTS](../../AGENTS.md#authorization-and-delivery) and
  [orientation](../docs/zazzles-agent-orientation.md#documentation-map).
- [Milestone index](common-entry-dispatch.md#delivery-order) and [Design §§2–5](../architecture/common-entry-dispatch.md#2-semantic-records-and-cardinality) and [§8](../architecture/common-entry-dispatch.md#8-bounded-framework-proposals).
- Dependency: 03-reporting-selectors; execute after 04-browser-events to keep one ordered lane. [Previous deliverable](common-entry-dispatch-04-browser-events.md)
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
| `internal/parse/dynamic_registrations.go, dynamic_registrations_test.go` | new | Finite registration expansion |
| `internal/parse/analysis.go, route_analysis.go, router_mounts.go, helper_resolution.go` | modified | Verified registration/target/closure integration |
| `internal/parse/testdata/typescript/dynamic-dispatch/; internal/parse/testdata/javascript/dynamic-dispatch/` | new | Finite generation and guard controls |
| `internal/report/entries.go; cmd/scan/analysis_test.go` | modified | Generation provenance and real pipeline checks |
| `.zazz/specifications/common-entry-dispatch-05-dynamic-http.md` | modified if needed | Contract decisions and evidence links |
| `.zazz/ephemeral/common-entry-dispatch-05-dynamic-http-run-log.md` | ignored local | Append-only execution evidence |

One milestone PR preserves a single canonical migration and removes transitional state
before human integration. Each deliverable remains a distinct review/evidence unit in
that PR. Do not create additional worktrees or branches merely to represent checkpoints.
This proposed review shape remains reviewable by the user before implementation.

Exclusions: new vulnerability families, authorization engine, queues/jobs/CLI entry adapters,
new target languages, full IR, graph database, automatic cache, independent LLM reviewer,
paid calibration, application execution, sample remediation and unrelated formatting.

## 4. Decisions

**Decision 1.** Expand literal finite tables and direct local factories because this yields source-reviewed bindings without executing target code.

**Decision 2.** Record generated registration origin separately from target resolution and security status because a generated route can be resolved, protected, unsafe or unknown.

**Decision 3.** Report unexpanded remainder at deterministic limits because truncation must not become an all-path protection claim.

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

- **AC1 — TestDynamicHTTPRegistrations.** A literal table/for-of and direct local literal-argument factory expose the expected entries, original spans, row/call/closure provenance and finite target alternatives. Verified by `TestDynamicHTTPRegistrations` in the packages named in §7; new test names are planned contracts, not existing results.
- **AC2 — TestDynamicHTTPContextIsolation.** Shared function/factory body exposed with different arguments or controls retains separate contexts and correct unsafe/protected/unknown results through real candidate and report paths. Verified by `TestDynamicHTTPContextIsolation` in the packages named in §7; new test names are planned contracts, not existing results.
- **AC3 — TestDynamicHTTPBounds.** Runtime registration data, mutation, recursive factories, unknown wrappers, unbound/fake receivers and over-limit expansions preserve gaps; known alternatives remain visible, zero-use unknown registrations remain reportable, repeated scans are deterministic. Verified by `TestDynamicHTTPBounds` in the packages named in §7; new test names are planned contracts, not existing results.
- **Scope/checks.** Scoped tests pass and the slice diff contains only relevant approved paths;
  verify with `git diff --check`, `git status --short` and slice-diff inspection.

## 7. Test strategy

- `TestDynamicHTTPRegistrations`: A literal table/for-of and direct local literal-argument factory expose the expected entries, original spans, row/call/closure provenance and finite target alternatives.
- `TestDynamicHTTPContextIsolation`: Shared function/factory body exposed with different arguments or controls retains separate contexts and correct unsafe/protected/unknown results through real candidate and report paths.
- `TestDynamicHTTPBounds`: Runtime registration data, mutation, recursive factories, unknown wrappers, unbound/fake receivers and over-limit expansions preserve gaps; known alternatives remain visible, zero-use unknown registrations remain reportable, repeated scans are deterministic.

Use compact table-driven controls and actual parsed fixtures where adapters are involved.
Hand-built records prove model/presenter boundaries only; they cannot substitute for actual
adapter-to-candidate-to-claim/report acceptance. Reuse adjacent existing tests for unchanged
behavior and inspect golden diffs rather than adding tests that mirror helper calls.

Scoped command (declared Go toolchain, no live network/model credentials):

```sh
GOTOOLCHAIN=local GOPROXY=off GOCACHE="$PWD/.zazz/ephemeral/go-cache" go test ./internal/parse ./internal/report ./cmd/scan
git diff --check
```

At the deliverable's completed boundary run race tests for affected packages. Full
race/vet/build and sample/report validation are mandatory for the final milestone; the
index defines the commands and each later code change invalidates affected prior evidence.
Local loopback stubs test policy/transport, not model accuracy. A sandbox bind denial is
an environment limitation to resolve, not a passing test or a reason to use a paid endpoint.

## 8. TDD entry point and execution sequence

First failing test: `TestDynamicHTTPRegistrations` with the observable contract in §6.

1. Start with TestDynamicHTTPContextIsolation on two factory registrations using different closure values/controls.
2. Implement finite expansion in deterministic order, sharing existing module binding and HTTP metadata models. Retain existing direct routes and finale persistence semantics.
3. Carry actual arguments/closure projections into invocation seeds; unsupported capture/mutation cannot acquire a safe status.
4. Exercise per-site, per-entry and scan budgets and compare existing HTTP sample/control outputs.

Illustrative implementation skeleton; names/bodies may adapt to verified existing code:

```go
// internal/parse/dynamic_registrations.go
// expandFiniteRegistration(site, bindings, budget) -> registrations, dispatches, gaps.
// Each derived record cites site + literal row/factory call + supported closure binding.
```

Recovery checkpoints: Literal table expansion; then factory/closure context and limits.
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

Use `.zazz/ephemeral/common-entry-dispatch-05-dynamic-http-run-log.md`, append-only.
Record standards verification, user decisions, start revision, phase
completions, commands/results, deviations, report paths/hashes, review/rework and recovery.
At session start read AGENTS, this contract, relevant design sections and the log; verify
Git state and choose the next incomplete checkpoint. If the ignored log is absent,
reconstruct from Git and tracked acceptance evidence; never infer success from missing notes.

## 12. Agent implementation prompt

```text
Implement only the user-authorized deliverable in the existing common-entry-and-dispatch
worktree. Read AGENTS.md, .zazz/specifications/common-entry-dispatch-05-dynamic-http.md,
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
