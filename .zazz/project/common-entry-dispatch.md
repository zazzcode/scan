# Common entry and dispatch model — follow-up scope

**Decision, 2026-11-02:** implement this as a separate foundation refactor after the Python/terminology
PR, using existing HTTP and browser consumers. The current branch records the architecture decision;
it does not implement a common entry/dispatch model. Command injection remains the next vulnerability
family; this foundation MUST precede dynamic-route or GraphQL entry-framework expansion.

## Governing design

Read [AGENTS.md](../../AGENTS.md), the [orientation](../docs/zazzles-agent-orientation.md), then
[architecture: entry/dispatch](../architecture.md#multiple-entry-and-dispatch-mechanisms),
[repository graph](../architecture/repository-graph.md), [evidence/claims](../architecture/evidence-and-claims.md)
and [trace judgment](../architecture/trace-judgment.md). These documents own the semantic boundaries;
this brief defines a bounded follow-up rather than a second canonical schema.

## Outcome

Analysis entry points MUST be independent of HTTP routes. Separate logical entry identity,
registration evidence, dispatch bindings, callable identity, invocation context and optional transport.
The same callable can serve multiple entries; an entry can have multiple supported dispatch targets.
Entry kind does not determine input trust. Multiple registrations and target alternatives MUST retain
conditions, provenance and unresolved coverage. Source/control/value semantics remain explicit.

Current pressure points are `analysis.Entry.Route`, separate route/browser identities, method/path
checks, route-handler completeness roles and route-centric report grouping. A struct-only addition
or an output projection of unchanged route-only internals does not meet the outcome.

## Bounded implementation

- Create a tracked implementation contract before coding, based on this brief and real consumers.
  Use a feature worktree from reviewed integration containing the terminology migration; if starting
  earlier is explicitly requested, name the reviewed parent commit and dependency in the draft PR.
- Establish common entry identity and registration/dispatch/callable records. IDs are snapshot-scoped
  unless an edit-stability policy is independently specified/tested. Model one-to-many/many-to-one
  relationships and partial/unknown resolution without inventing target callables.
- Migrate existing JS/TS HTTP, Python HTTP and browser analysis to that representation. Keep
  method/path and mount details as HTTP metadata. Preserve source/hash evidence and bounded flows.
- Carry entry/dispatch invocation context through candidate identity and claim completeness. A guard,
  value binding or protection fact in one registration context cannot protect another context.
- Make checks/selectors and judgment-required roles work without an HTTP route. Preserve old HTTP
  requirements and source policies; a newly generic name cannot weaken a completeness requirement.
- Group JSON/HTML/SARIF records by actual entries and locations. Retain HTTP/browser presentation
  where useful, and explicitly handle unknown kinds; no fake method/path or lost candidate record.
- Specify report versioning and saved schema-1/schema-2 conversion at the input boundary. Keep
  historical source/answers intact. Version changed model inputs and reject stale decision replay.
  Remove replaced internal aliases and duplicate canonical state in the same PR.

## Acceptance evidence

- Existing HTTP and browser scans use the common model and preserve supported candidates, source
  evidence, defect/protection semantics and coverage gaps in JS, TS and Python.
- One callable exposed twice has one callable identity, separate registration/dispatch contexts and
  separate applicable controls. A guarded entry cannot lend protection to an unguarded entry.
- An entry with two supported callable targets retains both; an unresolved additional target prevents
  an all-path protection conclusion when relevant. Discovery gaps remain visible even with zero uses.
- A non-HTTP entry fixture exercises the real shared candidate, completeness and JSON/HTML/SARIF
  paths without setting method/path, RouteID or ClientID. A hand-built model fixture is sufficient;
  this is not a new GraphQL/queue analyzer. Reachability alone does not fabricate a value dependency.
- Selection and source provenance remain explicit. Unknown input is not automatically attacker-controlled;
  schema/type validation does not automatically establish application policy or sink protection.
- Saved old reports export without changing files/answers. Current serialized records have one
  canonical entry association; compatibility conversion is tested and does not silently drop metadata.
- Identical input yields deterministic identities/order and bounded work. No whole-repository path
  enumeration or unbounded factory/target expansion is introduced.
- Run paired raw/protected/unknown controls, local HTTP-stub judgment/policy/replay tests, existing
  source-reviewed Python case checks, offline CLI scans/exports, Go race tests, vet, build and format.
  Record evidence and remaining gaps in the owning contract/evaluation documents.

## Exclusions

No dynamic-route discovery, GraphQL schema/resolver adapter, new vulnerability family, authorization
engine, queue/CLI/scheduler analyzer, full IR replacement, graph database, independent reviewer,
automatic cache, sample remediation or paid calibration is part of this foundation refactor.
Future adapters need their own source/binding/flow contracts and acceptance evidence.
