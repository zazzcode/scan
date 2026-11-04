# Cross-entry regression and behavior reports — deliverable contract

Worktree / branch: `common-entry-and-dispatch`. Feature: common entry/dispatch.
Milestone deliverable: 8 of 8. Project and Deliverable Owner: Michael.
Drafted: 2026-11-03. Execution tracking: deliverable-specific ignored local run log only.
Coordination: one Contributor Agent working in the assigned worktree; no delegation required.
Proposed review shape: one milestone draft PR into `stage`, human review/merge.
The user requested design/specification now and will request implementation later.
This document does not itself authorize implementation, publishing, paid calls or merging.

## 0. Capability

Verify the completed shared model against existing HTTP/browser samples and new event/dynamic/GraphQL/gRPC fixtures; produce reviewable before/after evidence and JSON/HTML/SARIF reports without paid calls.

## 1. Required reading

- [AGENTS](../../AGENTS.md#authorization-and-delivery) and
  [orientation](../docs/zazzles-agent-orientation.md#documentation-map).
- [Milestone index](common-entry-dispatch.md#delivery-order) and [Design §9](../architecture/common-entry-dispatch.md#9-delivery-and-evidence), [Python results](../evaluation/python-support-results.md#verification-and-provenance), [contributor checks](../../docs/contributing.md#contributing).
- Dependency: All seven implementation deliverables. [Previous deliverable](common-entry-dispatch-07-grpc.md)
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
| `cmd/scan/entry_pipeline_test.go; internal/report/entries_test.go; internal/judge/*_test.go` | new/modified | Mixed-entry integration, replay and policy regression |
| `docs/samples/evaluate-python.py` | modified only if required | Canonical report input compatibility without changing expected cases |
| `.zazz/evaluation/common-entry-dispatch-results.md` | new | Durable source-reviewed outcomes, commands and limitations |
| `.zazz/architecture.md; .zazz/architecture/common-entry-dispatch.md; .zazz/architecture/reporting-and-sarif.md, html-reporting.md` | modified | Actual implemented package/report contracts |
| `docs/running-scans.md, interpreting-results.md, analysis/README.md; .zazz/project/common-entry-dispatch.md; .zazz/specifications/README.md` | modified | Capability/compatibility status and discoverable results |
| `.zazz/specifications/common-entry-dispatch-08-acceptance.md` | modified if needed | Contract decisions and evidence links |
| `.zazz/ephemeral/common-entry-dispatch-08-acceptance-run-log.md` | ignored local | Append-only execution evidence |

One milestone PR preserves a single canonical migration and removes transitional state
before human integration. Each deliverable remains a distinct review/evidence unit in
that PR. Do not create additional worktrees or branches merely to represent checkpoints.
This proposed review shape remains reviewable by the user before implementation.

Exclusions: new vulnerability families, authorization engine, queues/jobs/CLI entry adapters,
new target languages, full IR, graph database, automatic cache, independent LLM reviewer,
paid calibration, application execution, sample remediation and unrelated formatting.

## 4. Decisions

**Decision 1.** Compare contextual source/sink/claim outcomes against pre-change baselines because new snapshot identities and grouping are intentional changes.

**Decision 2.** Separate offline discovery, stub-assessed policy and saved historical model evidence because they prove different properties.

**Decision 3.** Defer the user-built GraphQL application to later evaluation because fixture acceptance should be reproducible now without claiming application-level accuracy.

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

- **AC1 — TestMixedEntryPipeline.** One mixed JS/TS/Python target contains actual HTTP/browser-event/GraphQL/gRPC adapter output and dynamic registration. Entries/candidates/contexts remain correctly associated through canonical JSON, HTML and SARIF; related transports do not double count resolver uses. Verified by `TestMixedEntryPipeline` in the packages named in §7; new test names are planned contracts, not existing results.
- **AC2 — TestEntryProtectionCoverage.** Across entry kinds, guarded/unguarded shared callables, multiple targets plus unresolved alternatives, unknown applicable guards and zero-use gaps preserve claim-specific outcomes and cannot produce an entry-wide safe label. Verified by `TestEntryProtectionCoverage` in the packages named in §7; new test names are planned contracts, not existing results.
- **AC3 — RegressionCommands.** Run all Go race/vet/build checks, focused local stubs for both provider transports, original 29 Python saved-report case checks, supported saved-schema exports and offline existing HTTP/browser plus new-adapter CLI reports. Explain every normalized behavioral difference; no unexplained HTTP regressions. Verified by the exact race/vet/build, Python case-check and CLI/export commands in the milestone index, with command output and artifact hashes recorded in the results.
- **AC4 — ReportReview.** Visually inspect HTTP-only, event-only, GraphQL-only, gRPC-only and mixed HTML. Confirm entry-specific labels, metadata, location, sorting, disclosures, zero-use coverage, source excerpts and status parity; record paths/hashes and source-reviewed expected/observed matrix. Verified by the manual HTML inspection and cross-format status matrix in the milestone index, with reviewed artifact paths recorded in the results.
- **Scope/checks.** Scoped tests pass and the slice diff contains only relevant approved paths;
  verify with `git diff --check`, `git status --short` and slice-diff inspection.

## 7. Test strategy

- `TestMixedEntryPipeline`: One mixed JS/TS/Python target contains actual HTTP/browser-event/GraphQL/gRPC adapter output and dynamic registration. Entries/candidates/contexts remain correctly associated through canonical JSON, HTML and SARIF; related transports do not double count resolver uses.
- `TestEntryProtectionCoverage`: Across entry kinds, guarded/unguarded shared callables, multiple targets plus unresolved alternatives, unknown applicable guards and zero-use gaps preserve claim-specific outcomes and cannot produce an entry-wide safe label.

Use compact table-driven controls and actual parsed fixtures where adapters are involved.
Hand-built records prove model/presenter boundaries only; they cannot substitute for actual
adapter-to-candidate-to-claim/report acceptance. Reuse adjacent existing tests for unchanged
behavior and inspect golden diffs rather than adding tests that mirror helper calls.

Scoped command (declared Go toolchain, no live network/model credentials):

```sh
GOTOOLCHAIN=local GOPROXY=off GOCACHE="$PWD/.zazz/ephemeral/go-cache" go test ./...
git diff --check
```

At the deliverable's completed boundary run race tests for affected packages. Full
race/vet/build and sample/report validation are mandatory for the final milestone; the
index defines the commands and each later code change invalidates affected prior evidence.
Local loopback stubs test policy/transport, not model accuracy. A sandbox bind denial is
an environment limitation to resolve, not a passing test or a reason to use a paid endpoint.

## 8. TDD entry point and execution sequence

First failing test: `TestMixedEntryPipeline` with the observable contract in §6.

1. Start with TestMixedEntryPipeline and a source-reviewed expected entry/use matrix authored independently from scanner output.
2. Build current scanner once. Generate offline JSON/HTML/SARIF for existing HTTP/browser targets and new source fixtures; add local-stub assessed fixture runs for protection policy without live keys.
3. Run full validation once after final code changes. Re-run only affected checks if fixes are needed, followed by the final required integration gate. Inspect every golden diff.
4. Verify original saved-report hashes unchanged, replay rejection, removed legacy internal aliases and no fabricated HTTP non-HTTP fields. Record improved discovery/context/report behavior and remaining gaps in the evaluation document.

Illustrative implementation skeleton; names/bodies may adapt to verified existing code:

```go
// cmd/scan/entry_pipeline_test.go
// Build source-only mixed fixture; run real scanner with local judge stubs;
// compare canonical associations/statuses against reviewed expectations and export all views.
```

Recovery checkpoints: Mixed fixture integration; then full suite; then visual report review and durable evidence.
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

Use `.zazz/ephemeral/common-entry-dispatch-08-acceptance-run-log.md`, append-only.
Record standards verification, user decisions, start revision, phase
completions, commands/results, deviations, report paths/hashes, review/rework and recovery.
At session start read AGENTS, this contract, relevant design sections and the log; verify
Git state and choose the next incomplete checkpoint. If the ignored log is absent,
reconstruct from Git and tracked acceptance evidence; never infer success from missing notes.

## 12. Agent implementation prompt

```text
Implement only the user-authorized deliverable in the existing common-entry-and-dispatch
worktree. Read AGENTS.md, .zazz/specifications/common-entry-dispatch-08-acceptance.md,
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
