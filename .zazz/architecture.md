# Architecture

Tech stack and technical approach for `scan`.

**Status: draft, under active design iteration.**

## 1. Premise

**A capable LLM reasons about code better than a static analyzer, and handles the cases that break
static analysis.**

| SAST breaks on | Why | Who handles it here |
|---|---|---|
| Custom validators and sanitizers | Whether a regex or a hand-rolled check is *adequate* is semantic | LLM (Track B) |
| Framework indirection | Dynamic routing, DI, decorators, middleware ordering defeat syntactic flow models | LLM (Track B, C) |
| Dynamic dispatch, reflection, metaprogramming | Targets are not statically known | LLM (Track C) |
| Authorization and logic flaws | Not dataflow properties at all | LLM (Track C) |
| Business-logic-dependent exploitability | Requires understanding intent | LLM (Track C) |
| Stored and second-order flows | The source is not a request parameter | Partial — surfaced as unresolved |
| Whether a mitigation is *sufficient* | Inherently a judgment | LLM (Track B) |

We do not try to out-engineer SAST on its own terms. We do the deterministic work static analysis is
genuinely good at — parsing, resolution, path tracking, pattern matching — and hand the semantic
questions to the model.

## 2. Design principle

Split **deterministic structure** from **semantic judgment**.

1. Tree-sitter extracts structure — declarations, call sites, imports — without executing the target.
2. The scanner normalizes that into one language-agnostic graph and classifies sources, sinks, guards,
   and data types.
3. Deterministic code tracks the paths through that graph and produces findings for everything decidable
   by shape.
4. The model handles what requires judgment, reasoning over prepared context rather than raw files.

There is deliberately **no dataflow engine** — no taint propagation through parameters, returns, fields,
or aliases, and no hand-coded source-to-sink string tracking. That class of engineering is what the
premise replaces.

## 3. Code path tracking — the core capability

Everything else depends on this. If we cannot resolve what calls what across files, there is no graph
worth traversing, no path to hand a model, and no stable identity for scan-to-scan tracking. This is the
Graphify approach, and it is the load-bearing part of the system.

### Why it is the key

**It reduces what the model has to review.** The alternative — pointing a model at a repository — spends
its entire context window locating relevant code, and gives no guarantee it found the right code. A
tracked path hands the model exactly the code on that path, in order, with the guard and sink already
identified. The model spends its context on reasoning instead of search, which is both cheaper and more
accurate.

**It provides structure for analysis.** A path is an explicit claim: this entry point reaches this
operation, through these intermediate functions, with these guards present or absent, passing these
arguments. The model reasons *over that structure* rather than over loose files. It can cite the path, and
so can a human reviewing the finding. That structure is also what makes findings comparable across runs.

**It is the only way the three tracks share anything.** Track A needs call relationships to reason about
structural patterns. Track B is source→sink paths. Track C scopes are graph-derived neighbourhoods — a
route file, its handlers, the middleware chain around them. One graph, three uses.

### How it works

1. **Extract.** Parse each file with Tree-sitter and emit declarations, imports, and call sites.
2. **Resolve.** Turn each call-site reference into a node ID where the target is statically knowable.
   This is the hard part, and the most language-specific — see § 6.2.
3. **Connect.** Emit edges for calls, references, imports, and the type hierarchy — see § 7.
4. **Classify.** Mark sources, sinks, guards, and data types — see § 11.
5. **Annotate.** Record, at every call site, which argument expressions are passed and where they land —
   see § 10. This is what lets the model check flow rather than guess it.
6. **Traverse.** From each source, search forward for reachable sinks, bounded by depth.
7. **Prune.** Drop nodes that neither descend from a source nor reach a sink; keep only the minimal
   connecting subgraph; collapse paths differing by an uninteresting prefix.
8. **Pre-filter.** Drop paths where no call-site argument is source-derived or non-literal. A path of
   constants cannot carry untrusted data, and paying a model to discover that is waste.
9. **Prioritize.** When the candidate set must be capped, order by sink severity, data-class impact, and
   centrality. Graphify's *god nodes* — the most-connected concepts — are a useful signal: a path through
   a hub is more likely to be reachable in practice.

Traversal runs against the in-memory graph (§ 14). It is the hot loop of the whole tool, and a database
round trip in that loop is pure overhead.

### Where this can go wrong

- **Resolution correctness.** Barrels, path aliases, and monorepo layouts in TS/JS. If resolution is
  wrong, every downstream finding is noise. This is why it lands early in the plan.
- **Candidate explosion.** In a real service, nearly every route transitively reaches nearly every
  repository function. Steps 7–9 are doing all the work, and **the pre-filter is where real findings can
  be silently lost.** It is the least principled component in the design and deserves deliberate
  attention rather than tuning by feel. Dropped paths must be countable and reportable, never invisible.
- **Anonymous node identity.** Inline arrow handlers are the most common real handler shape and the
  easiest to name unstably. Unstable naming means unstable fingerprints means scan-to-scan churn.

## 4. Three detection tracks

Every finding belongs to exactly one track, and the track is recorded on the finding.

### Track A — pattern rules (deterministic end to end, no LLM)

A rule matches a known-bad code pattern. The rule *is* the verdict. No model call, no verdict churn,
byte-stable across runs.

This track carries most of the OWASP Top 10. Using Bearer's published rule catalog as the reference for
what is pattern-matchable:

| Category | Pattern findings that need no LLM |
|---|---|
| A01 | Path traversal from user input into `path.join`/`sendFile`; open redirect from a request value into `res.redirect`; sensitive data exposed in a response |
| A02 | Weak hashing (`md5`, `sha1`), weak ciphers (`DES`, ECB), `InsecureSkipVerify`, `rejectUnauthorized: false`, plaintext HTTP endpoints |
| A03 (SQL) | `eval`, `new Function`, `child_process.exec` with a template literal, `dangerouslySetInnerHTML`, raw SQL built by concatenation with a request-derived value, `yaml.load` without `SafeLoader`, `pickle.loads` |
| A03 (NoSQL) | A request-derived object passed straight into a Mongo filter (`Model.find(req.body)`, `findOne(req.query)`), `$where` present, `$regex` built from user input |
| A04 | Missing encryption of sensitive fields; persistent cookies carrying sensitive data |
| A05 | CORS `origin: '*'` with credentials, debug mode enabled, missing security headers, cleartext secrets in a cookie or JWT |
| A07 | Hard-coded credentials, weak password hashing, improper certificate validation |
| A08 | Deserialization of untrusted data |
| A09 | Sensitive data written to a log |
| A10 | SSRF where the outbound URL argument is literally a request-derived member access |

A06 (vulnerable and outdated components) is deliberately **out of scope** — it is dependency scanning, a
different tool. Bearer omits it for the same reason.

Track A findings are deterministic. They gate CI by default.

The SQL and NoSQL rows are not interchangeable, and we get to test that directly. The first two acceptance
targets are TypeScript apps — one over MongoDB, one over SQL. Same language, two sink classes, and a
measurable hypothesis: **SQL paths should mostly settle deterministically**, because the presence of bound
parameters is close to decisive, while **NoSQL paths should mostly escalate to Track B**, because Mongo
has no equivalent escape hatch. If the split does not come out that way in practice, § 4 needs revisiting.
See § 11.

### Track B — path rules (deterministic discovery, LLM judgment)

A rule finds a source→sink path via § 3. The path's *existence* is deterministic; whether the data is
sufficiently validated is the model's judgment. Identity comes from the path, never from the verdict.

Track B exists for the cases rules genuinely cannot decide: custom validators, regex quality, ORM and
filter parameterization, template escaping, framework allow-lists, and whether data actually flows given
the call-site argument bindings.

### Track C — scope review (insight where paths cannot reach)

Some findings are not paths at all. "This handler fetches by ID without checking ownership." "This chain
authenticates but never authorizes." "This error path returns the full user document." Those are
authorization, logic, and design flaws, and they are where whole-program analyzers are structurally blind,
because they are not dataflow properties.

Track C hands the model a **scope** — a route file, a controller, a middleware chain, an auth module —
derived from the same graph. No source→sink path is required and none is constructed. Scope selection is
the Track C analogue of path pruning, and is *Open* below.

This is the mode where semantic review most clearly exceeds static analysis, so it is first-class rather
than a fallback.

### Why the split matters

- **Determinism where it counts.** Track A findings are stable scan-to-scan and never appear as churn.
  Track B's *existence* is stable; only its verdict is not.
- **Cost.** Track A and the whole path-tracking pipeline cost zero tokens. Only Track B paths and Track C
  scopes reach a model.
- **Trust differentiation.** Consumers can gate hard on Track A and treat low-confidence Track B and C
  findings as advisory. The finding says how it knows.
- **Degraded mode is a product, not a fallback.** With no API key configured, Track A and the full
  deterministic pipeline still work.
- **A different error profile from SAST.** A dataflow analyzer over-reports inside the scope it models;
  Tracks B and C reach further semantically and under-report outside what the graph and the scopes cover.
  Complementary, not competing.

### What stays deterministic

| Concern | Deterministic | Rationale |
|---|---|---|
| Language detection, parsing | Yes | Tree-sitter AST |
| Module/import graph | Yes | Resolution is a lookup, not a judgment |
| Code path tracking | Yes | Graph traversal |
| Source, sink, guard identification | Yes | Pattern matching on known APIs |
| Sink class | Yes | `db.raw` is a SQL sink, always |
| Call-site argument bindings | Yes | Syntactic extraction |
| Guard present/absent on a path | Yes | Structural fact |
| Data-type classification (PII, PHI, credentials) | Yes | Name and source patterns, never values |
| Pattern findings (Track A) | Yes | Rule is the verdict |
| Finding identity / fingerprint | Yes | Never depends on a verdict |
| Scan-to-scan diff | Yes | Set membership |
| **Validation sufficiency** | **No — LLM** | Semantic judgment |
| **Does the data actually flow** | **No — LLM** | Judgment given arg bindings |
| **Exploitability in context** | **No — LLM** | Judgment |
| **Authorization and logic gaps** | **No — LLM** | Not a dataflow property; Track C's reason to exist |

Rule of thumb: **if it can be decided by matching a known shape, it is deterministic code. If it requires
judging whether a human's mitigation is adequate, it is the model.**

## 5. What the graph must carry

The model can only judge what it is shown. The graph's job is to make that judgment **checkable rather
than guessable**.

1. **Call-site argument annotations.** Without these, the model is asked to infer inter-procedural flow
   from whole function bodies — the dataflow problem, handed to the model, and it will hallucinate. See
   § 10.
2. **Off-path guard definitions.** The validator usually lives in another module. If the bundle holds only
   path nodes, the model calls it unvalidated. See § 11.
3. **Sink class and mitigation idiom.** SQL parameterization, shell escaping, HTML encoding, path
   traversal, and SSRF allow-listing are different questions. See § 11.
4. **Path evidence.** Ordered `file:line` spans per node, so every verdict is auditable.
5. **Guard status on the path.** Whether a guard sits on the path is a structural fact and must never be
   something the model has to discover.

## 6. Language abstraction

Four things must be settled before the first language is implemented, because retrofitting them is
expensive. The test for each: **Python and Java must be able to fill it in without changing the core.**

### 6.1 Language adapter

Each language supplies exactly: a grammar, the file extensions it claims, source queries, sink queries,
guard queries, and a resolver. Nothing else language-specific may leak into the core. If `internal/trace`
or `internal/graph` ever needs to know the language, the boundary is wrong.

### 6.2 Resolver

Module resolution is the most language-specific piece and the most likely to be botched — and it is what
makes § 3 work or fail. A resolver that assumes "an import path resolves relative to the file" works for
JS and Python and quietly breaks for Go and Java.

```go
// Resolver turns a callee reference in a file into a graph node ID.
// Returning ok=false is expected and normal: dynamic targets and external
// dependencies are not resolvable, and that is recorded rather than dropped.
type Resolver interface {
	Resolve(fromFile string, ref CalleeRef) (nodeID string, ok bool)
}
```

Language-specific pressure to design against now:

- **JS/TS:** `tsconfig` path aliases, barrel re-exports (`index.ts`), monorepo workspaces,
  `require`/`import` of a computed value, and the `node_modules` boundary — never traverse dependencies.
- **Python:** `__init__.py` re-exports (the barrel analogue), relative imports, `importlib`.
- **Go:** packages and receiver methods. **Java:** classpath, packages, overloads.

### 6.3 Node identity

*This is where JS/TS will punish a naive model.* Handlers are frequently inline and nameless —
`router.get('/x', async (req, res) => {...})` — so "node = named function" misses exactly the nodes that
matter. The scheme needs deterministic synthetic naming for anonymous callables, stable when *other* lines
in the file change.

Two stress tests from the other languages, designed against now:

- **Java overloads:** two methods named `create` in one class. A `file::Class.create` ID collapses them.
  The scheme needs an optional signature/arity discriminator.
- **Python:** `self` methods, `@staticmethod`/`@classmethod`, and decorators that *wrap* a handler. The
  wrapped function is still the node; the decorator is a guard annotation.

### 6.4 Rule engine expressiveness

The rule engine must express **Python decorators and Java annotations** even though only JS rules are
written first. If decorators and annotations cannot be expressed as data, source detection becomes
per-language Go code — the corner to avoid.

*Open:* the rule format. Bearer's model is worth studying: declarative rule definitions for the common
case, with an escape hatch for complex detection logic, and a published rule catalog documenting the code
pattern per rule. Our lean is `.scm` structural queries plus a small declarative schema, with Go detectors
only where a rule genuinely cannot express the check.

## 7. Edge vocabulary

`CALLS` alone is insufficient. Middleware is registered *by reference* — `app.use(authMiddleware)` — and
that is how Go and Python register handlers too.

| Kind | Meaning |
|---|---|
| `CALLS` | Direct invocation |
| `REFERENCES` | Function passed as a value: middleware, handlers, callbacks |
| `IMPORTS` / `RE_EXPORTS` | Module graph |
| `INHERITS` / `IMPLEMENTS` / `MIXES_IN` | Type hierarchy. Borrowed from Graphify, which resolves these across ~40 languages — a handler may be inherited, and an interface implementation may be the real sink |
| `REGISTERS_ROUTE` | Framework-specific source registration |
| `PASSES_DATA_TO` | Reserved; not populated in v1 (that would be dataflow) |

Every edge carries `Origin`: `EXTRACTED` (read from the syntax tree) or `INFERRED` (resolved by a later
pass). This is Graphify's distinction, and it lets reporting state how confident a path is.

## 8. Node model

Language-agnostic, so traversal, pruning, and prompting are written once.

```go
package graph

// NodeRole describes a node's security role in a call path.
type NodeRole string

const (
	RoleSource   NodeRole = "SOURCE"   // request entry points, HTTP inputs, CLI args
	RoleSink     NodeRole = "SINK"     // raw SQL, shell, filesystem, outbound HTTP, eval
	RoleGuard    NodeRole = "GUARD"    // validation, sanitization, authorization
	RoleInternal NodeRole = "INTERNAL" // ordinary function/method logic
)

// Node is a logical callable unit: a function, method, route handler, module, or
// an anonymous callable with a deterministic synthetic name.
type Node struct {
	ID         string   `json:"id"`         // "path/file.ts::Class.method#2"
	Language   string   `json:"language"`   // "typescript", "python", ...
	Role       NodeRole `json:"role"`
	FilePath   string   `json:"file_path"`
	StartLine  int      `json:"start_line"`
	EndLine    int      `json:"end_line"`
	Params     []Param  `json:"params"`     // ordered; the target of ArgBinding.Position
	Content    string   `json:"content"`    // source snippet
	SinkClass  string   `json:"sink_class,omitempty"`  // "sql", "nosql", "shell", "html", "path", "ssrf"
	Mitigation string   `json:"mitigation,omitempty"` // idiom to look for, e.g. "bound parameters"
}
```

`Params` and `SinkClass`/`Mitigation` are what make § 10 and § 11 possible.

## 9. Pipeline

```text
[ Target repository ]
        │
        ▼
1. Concurrent file walker        (goroutine worker pool)
        │
        ▼
2. Tree-sitter parser            (per-language grammar; one CST per file)
        │
        ▼
3. Structural resolver           (declarations, imports, call sites, arg bindings)
        │
        ▼
4. Classifier                    (SOURCE / SINK / GUARD roles; sink class; data types)
        │
        ▼
5. Code path tracking            (discovery, pruning, pre-filter, prioritization)  ← the spine
        │
        ├───────────────┬────────────────────┬────────────────────┐
        ▼               ▼                    ▼                    ▼
6a. Track A       6b. Track B         6c. Track C
    pattern rules     path review         scope review
    deterministic     over tracked        over graph-derived
    findings          paths               scopes
    (no LLM)          (LLM judgment)      (LLM review)
        │               │                    │
        │               ▼                    ▼
        │    7. Context bundler (snippets, arg bindings, off-path guards,
        │                        sink class + mitigation, scope context)
        │               │                    │
        │               ▼                    ▼
        │    8. LLM judgment (Track B: three-state verdict; Track C: findings)
        └───────────────┴────────────────────┴────────────────────┘
        │
        ▼
9. Reporter + findings log       (JSON + Markdown; diff against previous findings)
```

Stages 1–3 are parallel across files. Stage 4 is a single pass. Stage 5 is single-threaded with parallel
per-source search. Stages 7–8 are I/O-bound and parallel across bundles.

## 10. Call-site argument annotation

The highest-leverage piece of the schema. At every call site, record **syntactically** which argument
expressions are passed and which parameter position each lands in.

This is extraction, not analysis: one AST walk per call site. No fixpoint, no aliasing, no propagation. It
converts "the model guesses whether the value could flow" into "the model checks whether the argument
arriving at the sink is the one the source supplied."

```go
// ArgBinding records one argument at a call site and where it lands in the callee.
type ArgBinding struct {
	Position  int    `json:"position"`    // 0-based parameter position in the callee
	ParamName string `json:"param_name,omitempty"`
	Expr      string `json:"expr"`        // source text of the argument expression
	ExprKind  string `json:"expr_kind"`   // "literal", "param", "member", "call", "unknown"
	// RootParam is set only when the expression IS literally one of the enclosing
	// function's own parameters. This is the single flow-ish signal recorded, and
	// it is local to the call site: we never propagate beyond it.
	RootParam string `json:"root_param,omitempty"`
}
```

Edge with annotations:

```go
type Edge struct {
	SourceID string       `json:"source_id"`
	TargetID string       `json:"target_id"`
	Kind     string       `json:"kind"`
	Origin   string       `json:"origin"` // "EXTRACTED" or "INFERRED"
	Args     []ArgBinding `json:"args,omitempty"`
}
```

**Boundary to hold:** `ExprKind` and `RootParam` describe what the expression *is*. Whether an expression
is *untrusted* is a rule or model interpretation, not something this schema decides.

*Open:* extending this with a **per-function local binding table** — `const id = req.body.id` then
`find(id)` — so the annotation follows one level of local assignment. Still purely syntactic, and it
closes the most common flow gap in real Express code.

## 11. Classification: sources, sinks, guards, data types

All four are data-driven from per-language query sets. Adding a rule must not require changing Go code.

- **Sources** — route registrations and handler signatures: Express/Fastify for JS/TS, FastAPI and Flask
  for Python, `net/http` and common routers for Go, Spring and servlet entry points for Java. Plus CLI
  entry points.
- **Sinks** — dangerous operations, each carrying a **sink class** (`nosql`, `sql`, `shell`, `path`,
  `ssrf`, `html`, `deserialization`, `log`) and the **mitigation idiom** to look for, so the prompt asks
  the right question for that sink type. For example: `nosql` looks for filter validation or type
  coercion; `sql` looks for bound parameters; `shell` looks for the argument-array form and an allow-list;
  `ssrf` looks for a host allow-list; `path` looks for `path.resolve` containment.

  **The SQL/NoSQL asymmetry is the sharpest illustration of the track boundary.** For SQL, the presence of
  bound parameters is close to decisive: the deterministic tier can often settle it and the model only
  confirms. MongoDB has no equivalent escape hatch — passing a request object as a filter is safe *only
  if* something upstream coerced or validated it, and whether a mongoose schema, a `zod` parse, or a
  hand-rolled `typeof x === 'string'` check is sufficient is a judgment. NoSQL findings therefore land on
  Track B far more often than SQL ones. Two TypeScript fixtures, one per database, exist to measure that.
- **Guards** — validation, sanitization, and authorization functions. First-class nodes so they can be
  pulled into a bundle even when off-path. Structural detection stays deterministic; only *sufficiency* is
  the model's call.
- **Data types** — classify by name and source pattern only, never by value: personal data, PII, PHI,
  credentials, payment data. Borrowed from Bearer's privacy model. This is an **impact** signal, not a
  detection: a sink handling PII ranks above one handling a public field. Pure pattern matching, so it
  stays deterministic.

*Open:* a curated **framework-facts store** — per-library notes on what actually mitigates each sink class
(for example, that mongoose casts per schema and the native driver does not). This is a deterministic
asset that would let the model reason from facts rather than guess at library behaviour.

## 12. Context bundling

One bundle per unit of review, token-bounded. For Track B the unit is a path; for Track C it is a scope.

- the path summary: node IDs, roles, ordered `file:line` spans
- the source snippet for each node on the path, in order
- **the `ArgBinding` annotations for each edge on the path** — the flow evidence
- **guard nodes reachable from the path, even off-path**, with their definitions
- the sink class and the mitigation idiom to look for
- the data types the path handles, as impact context
- the OWASP categories the path's sink set could implicate
- for Track C, the scope itself: the file or chain, its neighbours, and the middleware order around it

Bundles are truncated by a per-node and per-path token budget. Never truncate the path itself — drop the
whole path and report it as unexamined rather than silently passing it.

## 13. LLM judgment

We track paths deterministically; the model judges the prepared bundle. The model does **not** drive
traversal in v1.

- **Track B verdict is three-state:** `validated`, `not_validated`, `uncertain`, with a confidence and
  cited evidence. Not binary — forcing binary converts uncertainty into either a false positive or a
  silent false negative.
- **Evidence is mandatory.** The verdict must cite the validation it accepted, or the missing check. A
  verdict without a citation is unusable and is treated as `uncertain`.
- **Track C returns findings, not verdicts.** A scope review reports issues ("no ownership check before
  the fetch") with a severity and an anchor symbol, rather than judging a path.
- **Schema-constrained output**, low temperature, one request per bundle.
- **Provider-agnostic** over an OpenAI-compatible endpoint, so the base URL can point at a hosted provider
  or a self-hosted model.
- **Bounded concurrency and retry with backoff.** Fully stubbed in tests; no network in unit tests.
- **Version recorded.** Model ID and prompt version are stored per scan, so a verdict change is
  explainable rather than mysterious.

*Open:* an opt-in deep-dive mode where the model drives traversal via a graph-query tool. More flexible,
far more token-hungry, and unpredictable in coverage — same repo, two runs, different paths examined.
Deferred deliberately.

## 14. Persistence

Four distinct problems. Conflating them is the trap.

1. **Execution** — enumerate bounded source→sink paths repeatedly over the whole graph. A hot loop.
2. **Rescans** — skip unchanged work and correctly invalidate what a change affects.
3. **Findings tracking** — what is new, persisting, or resolved since last scan.
4. **Interchange** — a durable, reviewable artifact.

### Decision

- **Execution: in-memory Go graph.** `map[string]*Node` for lookup plus adjacency slices (`[]int32`
  indices) for forward and reverse traversal. BFS/DFS over flat slices is orders of magnitude faster than
  a SQL or Cypher round trip, and it keeps the graph engine free of native dependencies.
- **Interchange: JSON.** `graph.json` per run, deterministically ordered — Graphify's model.
- **Findings log: previous findings, persisted.** Diffing against it yields new / persisting / resolved as
  set membership. A flat file is sufficient and more reviewable than a database.
- **Suppression: a Git-tracked baseline** (`.scan-baseline.json`) listing accepted and suppressed
  fingerprints with a reason and an owner, plus rule-level skips (Bearer's `--skip-rule` idea). CI fails
  only on findings absent from the baseline.
- **Rescan cache: keyed by content hash and ruleset version.** A rule change invalidates classification
  even though the parse tree is unchanged, so keys are split by artifact (`parse:<hash>`,
  `classify:<hash>:<ruleset>`). *Open:* content-addressed directory or SQLite (`modernc.org/sqlite`, pure
  Go). Performance only — correctness does not depend on it.

### Rescan correctness

A per-file cache is necessary but not sufficient: a change to one file can alter paths running through
*other* files, so findings cannot be invalidated per file.

1. Walk the tree; hash each file's content.
2. Load cached nodes and edges for unchanged files; re-parse only changed and new files.
3. Build the in-memory graph from cached plus freshly parsed subgraphs.
4. Invalidate any recorded path with a node in a changed file, or whose endpoint files changed.
5. Re-run path tracking only for sources whose reachable subgraph intersects a change.
6. Re-run the model only on new or invalidated bundles. Unchanged paths keep their previous verdict.

Step 6 is where rescans pay off: the model is the only expensive stage.

### Finding identity

Line numbers, file paths, and enumeration order are all unstable, so identity is built from durable
content:

```text
fingerprint = sha256(rule_id || "\x00" || join(qualifiedSymbolChain, ">"))
```

- `rule_id` — the rule that produced the finding (Track A or Track B).
- `qualifiedSymbolChain` — the qualified symbol of each node on the path, in order, for example
  `handlers.CreateUser>service.Users.Create>db.Raw`.

Deliberate consequences:

- A function rename produces a new fingerprint, so the finding reads as one new plus one resolved. Surface
  that as "likely moved" when rule, sink symbol, and path shape still match.
- Moving a file does **not** change the fingerprint. Desired.
- Two structurally identical paths collapse into one finding. Desired.
- **The verdict is not part of the fingerprint.** A verdict change must never create or resolve a finding.

Track C has no path, so identity falls back to
`rule_id || "\x00" || scopeID || "\x00" || anchorSymbol`, where the anchor is the symbol the finding
attaches to (the handler, the middleware, the function). *Open:* whether that is stable enough for gating,
or whether Track C findings should be advisory and excluded from new-finding gating.

### Finding record

```
rule_id          which rule fired
detection        "pattern" (Track A) | "path" (Track B) | "scope" (Track C)
owasp_category   A01..A10
cwe_id           CWE mapping, as Bearer reports
severity         CRITICAL | HIGH | MEDIUM | LOW
data_classes     PII | PHI | credentials | payment | none
confidence       1.0 for Track A; model confidence for Track B
verdict          validated | not_validated | uncertain (Track B only)
evidence         ordered file:line spans
rationale        rule description, or the model's cited reasoning
```

`detection` is what lets a consumer gate hard on Track A and treat low-confidence Track B and C findings
as advisory — and it keeps the output honest about how each finding was reached.

### Gating policy

CI gates on **new fingerprints**. A verdict change on unchanged code is surfaced as an informational
status change, not as a new finding. Stating this explicitly prevents the first confusing CI failure.

## 15. Determinism

The structural graph, the tracked path set, and every Track A finding must be reproducible for identical
input, independent of goroutine scheduling. This is a correctness requirement: finding identity, the
rescan cache, and the JSON report all depend on it.

- Never iterate a Go map when producing ordered output; sort keys explicitly.
- Never let worker completion order determine slice order.
- No wall-clock, random seed, or host-specific value in graph or findings output.
- When parallelism is added to a stage, assert parallel output equals sequential output in a test.

Track B and C verdicts are **not** deterministic — inherent to the design, and why identity never depends
on them.

## 16. Prior art

### Graphify — `https://github.com/Graphify-Labs/graphify`

Apache-2.0. Turns a codebase plus docs, SQL schemas, configs, and PDFs into a queryable knowledge graph.

What we take:

- **The core validation.** Code is parsed with tree-sitter AST and is explicitly *"deterministic, no
  LLM"*, with graph build measured at **zero LLM credits**. It is the same split we are making, at a much
  larger scale, and it is evidence the deterministic tier carries the load.
- **The path-tracking model itself** — § 3 is this approach, specialised for security.
- **`EXTRACTED` / `INFERRED` edge tagging**, so every edge says whether it was read from source or
  resolved. Adopted directly.
- **`INHERITS` / `MIXES_IN` edges** — a reminder that the type hierarchy matters and that we were missing
  it.
- **`graph.json` as the artifact** — the graph is the durable output, queryable without re-reading files.
- **God nodes and communities** — most-connected concepts and Leiden clustering, computed without an LLM.
  Used here for path prioritization.
- **Rationale as first-class nodes** — `# NOTE:` / `# WHY:` comments and ADR citations become graph nodes.
  A security analogue is attractive: an in-code annotation like `// scan: safe - parameterized` becoming a
  node the model can see and cite.
- **Local-first.** Nothing leaves the machine for the code pass.

Where we differ: Graphify maps a codebase for comprehension. We map it to find a specific kind of path and
to have a model judge those paths, and we carry security-specific classification (source, sink, guard,
sink class, data types) that Graphify has no reason to model.

### Bearer — `https://github.com/bearer/bearer`

Open-source SAST (Go) covering the OWASP Top 10 and CWE Top 25, with dataflow analysis and a privacy
scanner. It is the tool we are deliberately not building, which makes it the best available reference for
*where rules beat a model*.

What we take:

- **The pattern-rule catalog as the Track A target.** Bearer's per-category rule list is close to a
  specification for what must be deterministic in our tool.
- **CWE IDs alongside OWASP categories**, plus severity tiers and per-rule documentation.
- **Rule-level suppression** (`--skip-rule`) in addition to per-finding suppression.
- **Data-type classification for impact** — 120+ data types across PII, PHI, and personal data,
  identified by matching patterns in source code and never reading values.
- **The rule-format lesson** — declarative rules for the common case with an escape hatch for complex
  detection.
- **A complement, not a competitor.** They model flows across files; we reach further semantically, and
  the error profiles differ. Running both and diffing is the cheapest available check on our Track A rules.

## 17. Tech stack

| Concern | Choice | Rationale |
|---|---|---|
| Language | Go 1.27+ | Real multi-core parallelism over file I/O and AST traversal; single static binary |
| Parsing | Tree-sitter | One uniform CST across languages; rules authored as `.scm` queries |
| TS binding | Behind an interface | `tree-sitter/go-tree-sitter` (cgo) or a pure-Go runtime; decide by benchmark |
| Graph engine | In-memory Go graph | Path tracking is a hot loop; a database there is pure overhead |
| Graph artifact | JSON | Diffable, reviewable, Graphify-compatible |
| Findings log | Flat file, diffed | Set membership; no database needed |
| Suppression | Git-tracked baseline | Reviewable in a PR |
| Rescan cache | Content-addressed or SQLite (*Open*) | Performance only |
| LLM client | `net/http` + `encoding/json`, OpenAI-compatible | Tiny API surface; no SDK dependency |
| CLI | stdlib `flag` | One command with flags |
| Logging | `log/slog` | Structured progress for long scans |
| Testing | stdlib `testing`, golden files, `httptest` | No test framework dependency |
| Formatting/lint | `gofmt`, `go vet` | Stdlib; add `golangci-lint` when rules stabilize |
| CI | GitHub Actions | build, vet, `go test -race`, `gofmt -l` gate |

## 18. Package layout

```text
cmd/scan/          CLI entry point
internal/walk/     repository traversal and file job production
internal/parse/    Tree-sitter grammar registry and parsing worker pool
internal/graph/    node/edge types, resolver, arg bindings, graph construction
internal/lang/     per-language adapters (grammar, queries, resolver)
internal/rules/    rule engine, rule catalog, classification (embedded .scm + declarative rules)
internal/patterns/ Track A pattern rules
internal/trace/    code path tracking: discovery, pruning, dedup, pre-filter, prioritization
internal/scopes/   Track C scope selection and scope bundles
internal/context/  token-bounded bundle construction
internal/llm/      OpenAI-compatible broker client
internal/findings/ identity, finding record, previous-findings log, baseline diffing, grouping
internal/report/   JSON and Markdown emitters
```

`internal/` is deliberate: nothing here is a public API yet. See
[`standards/code-structure.md`](standards/code-structure.md).

## 19. Testing strategy

- **Golden-file tests** for parsing and resolution: a fixture repository in, a normalized graph out.
- **Path-tracking tests** — the highest-value tests in the project. A hand-built graph in, expected paths
  out, including cycles, unreachable sinks, and guarded sinks.
- **Table-driven tests** for rule matching, covering Track A and classification.
- **Arg-binding tests** — the highest-risk extraction. A call site in, expected bindings out.
- **Determinism tests** — two runs, byte-identical `graph.json` and Track A findings.
- **Stubbed broker** for the model stage, including the `uncertain` and malformed-response paths.
- **End-to-end** on a small fixture repository with a stubbed broker, asserting the JSON report and the
  findings diff.
- **Optional differential testing** against Bearer on a fixture app, to check Track A rules.

See [`standards/testing.md`](standards/testing.md).

## 20. Roadmap

0. **Calibration** — hand-built vulnerable/safe pairs, hand-written bundles, measure what context makes
   the model reason well. Produces the prompt, the response schema, and the corpus that later phases need.
   No Go code.
1. **Language abstraction and graph schema** — the four interfaces in § 6, node/edge types, arg bindings.
2. **JS/TS adapter** — grammar, resolver (aliases, barrels, `node_modules` boundary), anonymous-node
   naming.
3. **Pipeline skeleton** — walker, parser worker pool, resolver, classifier, `graph.json` export.
4. **Track A pattern rules** — the deterministic OWASP set from § 4. A useful tool with no LLM at all,
   which is why it comes before the model work.
5. **Code path tracking** — discovery, pruning, pre-filter, prioritization. Measured on a real repository.
6. **Context bundling** — arg bindings plus off-path guards.
7. **LLM judgment (Track B)** — three-state verdict, cited evidence, stub and degraded mode.
8. **Track C scope review** — scope selection, scope bundles, scope findings.
9. **Reporting and CLI** — JSON and Markdown, CWE mapping, severity, flags, CI exit codes.
10. **Findings tracking** — identity, previous-findings log, Git-tracked baseline, CI gating.
11. **Rescan cache** — parse cache and invalidation.
12. **Python adapter** — decorators, `self`, relative imports. Validates the § 6 interfaces.
13. **Go and Java adapters** — receiver methods, interface dispatch, annotations, overloads.
14. **Rule expansion** — deeper framework coverage and broader OWASP coverage.

Steps 0–4 produce a deterministic, useful scanner with no model involved. Step 12 is the real test of
whether the abstractions hold.

## 21. Open questions

- **Rule format.** `.scm` plus declarative schema, and where a Go detector is justified. Bearer's model is
  the reference.
- **Pre-filter threshold.** What reliably drops constant-only paths without dropping real findings, and how
  to make dropped paths countable rather than invisible.
- **Scope selection for Track C.** Which files, controllers, and middleware chains are worth reviewing,
  and how that scales with repository size.
- **Track C finding identity.** Whether scope + anchor is stable enough for CI gating, or Track C findings
  should be advisory.
- **Node identity for anonymous callables.** Exact synthetic naming scheme, and stability under edits to
  unrelated lines.
- **Signature discrimination.** How much of a signature to encode in a node ID, given Java overloads.
- **Local binding table.** Extending `ArgBinding` with per-function local assignments.
- **Framework facts.** A curated per-library store of what actually mitigates each sink class.
- **Finding grouping.** Collapsing the same rule + same sink across many routes into one finding with many
  entry points.
- **Trust hints on argument expressions.** Do rules tag expressions as request-rooted, or is that the
  model's inference?
- **Cache medium.** Content-addressed directory vs. SQLite.
- **Tree-sitter binding.** cgo vs. pure Go, decided by benchmark on a large repository.
- **Baseline granularity.** Per repository or per scan target; accepted findings, suppressed fingerprints,
  rule-level skips, or all three.
- **In-code annotations.** Whether `// scan: safe - <reason>` comments become guard evidence the model can
  cite.
- **Deep-dive mode.** Whether model-driven traversal earns its token cost as an opt-in.
- **Findings schema versioning.** Whether it is versioned from the first release.
