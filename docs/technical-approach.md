# Technical approach and design strengths

[Documentation](README.md) · [Philosophy](philosophy.md) · [Implemented architecture](../.zazz/architecture.md)

## Source-backed analysis with constrained semantic judgment

Zazzles establishes supported code relationships through static analysis, then uses a decision model
to evaluate whether the demonstrated controls satisfy a specific security requirement. Explicit
evidence gates constrain the resulting judgment. The intended output helps a developer or coding
agent decide what to repair, what protection to preserve and what uncertainty to investigate.

The Go scanner reads application source, resolves supported language/framework bindings and follows
values into security-sensitive operations. Each candidate carries source spans and hashes, relevant
helpers, branch conditions, operation operand roles, library contracts and unresolved facts. The
semantic judge receives this evidence together with versioned language **Analysis Guidelines**.
Go then checks the answers against claim-specific evidence requirements and consistency gates.

The model is part of security judgment, rather than only a prose explanation of a preexisting alert.
It does not generate the source relationships or fill missing code with assumptions. Source text is
evidence data; scanner-supplied guidelines define the interpretation instructions. JSON preserves
both evidence and raw answers, while HTML and SARIF present the policy-gated result.

A custom validation function illustrates the intended benefit. When its implementation, actual value
transfers and relevant conditions are supplied within supported paths, the judge can assess whether
the control is adequate for the operation without requiring a named sanitizer rule for that particular
function. This is bounded interpretation of supplied evidence, not an ability to discover arbitrary
missing code or prove every custom validator correct.

**Semantic judgment reduces the need to enumerate every acceptable implementation pattern; it does
not eliminate the need for analysis models.** Input discovery, handler resolution, supported
transfers and library argument roles still need explicit engineering. Discovery misses never reach
the judge; unresolved claim-relevant facts remain blockers. Reduced false findings and improved
coverage remain hypotheses to measure alongside false protection, abstention, latency and cost.

## The distinctive design choices

The useful distinction is the combination and enforcement of these boundaries. It is not a claim that
syntax parsing, dataflow, taint analysis, model-assisted security review or SARIF originated here.

| Design choice | Practical strength | Current boundary |
|---|---|---|
| Separate claims for validation, construction and use protection | A fix can address the exact defect while preserving an existing control | Reports still group several concerns within one operation record |
| Source-backed evidence roles and claim blockers | A confident model answer cannot turn absent bindings, helpers or policies into protection | Incorrect extraction or an incorrect answer can still produce a bad result |
| Analysis Guidelines plus actual helper/guard evidence | Semantic questions can consider an implementation rather than only a validator's name | Custom behavior is judged only within resolved, bounded paths |
| Findings, Review and Protected outcomes | Agents receive both repair work and scoped controls to preserve, plus explicit investigation gaps | Protected applies to the named claim, not the route or whole application |
| Entry/dispatch invocation context and bounded stored/cross-tier links | Shared handlers keep separate input/control facts for each exposure; review can connect supported producer and rendering contexts | Framework discovery, dynamic dispatch, storage identity and cross-tier coverage remain bounded |
| Versioned evidence, guideline documents and exact-request replay | A historical decision remains auditable and cannot silently move onto changed input | Replay is explicit; an automatic persistent decision cache is not implemented |
| Shared language-neutral analysis records and reporting | JS/TS and Python use the same claim gates and output semantics | Each language still needs its own correct framework and flow models |

For example, a signup handler can bind an email safely into SQL while enforcing only presence, even
though the matching form establishes an email-format requirement. Zazzles can preserve the SQL-binding
claim and separately report the unmet validation policy. Conversely, a missing business requirement
stays Review: the model cannot invent a restriction merely because the field has a familiar name.
A construction defect can also merit repair without establishing an injection exploit.

## Relationship to established SAST

The common-entry design follows established static-analysis principles: represent program structure,
resolve possible call targets, bind inputs to particular invocations and track values separately from
control reachability. Other tools need not use Zazzles's record names or identical boundaries to
implement related analysis. This comparison concerns architectural ideas, not equivalent coverage.

### Research references and inspiration

We studied the following SAST approaches for design ideas and evaluation. Conventional SAST includes
structured program models, framework semantics and dataflow analysis, not only textual matching.
These references inform Zazzles's engineering; they do not establish equivalent coverage or quality.

| Tool/reference | Approach and ideas studied | Role in Zazzles research |
|---|---|---|
| [CodeQL](https://codeql.github.com/docs/writing-codeql-queries/about-data-flow-analysis/) | Local/global dataflow, source/sink models, taint steps and path queries | Architecture and semantic-model reference; not a measured CodeQL benchmark or a Zazzles runtime dependency |
| [Joern / Joern Scan](https://docs.joern.io/scan/) and [Code Property Graph](https://docs.joern.io/code-property-graph/) | Shared program representation across language frontends; graph relationships, traversals and predefined security queries | Architecture reference reviewed for common representations and query-driven scanning; no Joern comparison run or integration is claimed |
| [Bearer](https://docs.bearer.com/guides/dataflow/) | Go/Tree-sitter language adapters and flow/reporting patterns studied in its pinned source; reports exposing classifications, source locations and components | Read-only implementation study plus pinned local comparison scans; see the [research record](../.zazz/project/05-prior-art-and-feasibility.md#bearer-inspection-provenance) and [baseline runbook](samples/nodejs-goof/bearer-runbook.md) |
| [Semgrep](https://docs.semgrep.dev/writing-rules/data-flow/taint-mode/overview) | Pattern rules and taint models for sources, propagation, sanitization and sinks | Rule/model reference plus pinned Community Edition comparison scans; CE evidence does not establish the capabilities or performance of every edition |

Bearer and Semgrep CE baselines are scope-specific observations on selected applications, not
rankings of whole products. CodeQL and Joern are architectural research references here, not executed
benchmarks. Studying these approaches does not imply their engines, source or rule sets are embedded
in Zazzles. The [prior-art record](../.zazz/project/05-prior-art-and-feasibility.md) retains research
provenance and the boundaries of implementation reuse.

### Relationship to Joern Scan

There is a broad similarity: extract program structure, preserve relationships and source locations,
evaluate security properties, and report findings. In the
[documented Joern Scan workflow](https://docs.joern.io/scan/), Joern generates a Code Property Graph,
runs a selected bundle of predefined queries against it and reports matches. Its queries can express
structural and flow-related analysis; the approach is more than a text-pattern scanner.

Zazzles is not Joern Scan with a model added. Its current static engine uses Tree-sitter, typed Go
entry/dispatch records and bounded language-specific flow/evidence assembly, rather than Joern's CPG
backend or a full repository IR. Go records deterministic findings where the construction is
established. The decision model assesses claim-specific security questions about supplied code,
including controls and operation roles; Go then applies evidence-completeness and consistency gates.
It is part of judgment, not just an explanation appended to an alert.

The model cannot create a missing call/argument binding, assume a generated route exists, or turn an
unresolved trace into protection. Thus the distinction involves both the representation and the
judgment contract, not merely adding AI to an otherwise equivalent scanner. Neither architecture
implies equal language coverage, trace depth, accuracy or runtime cost. This comparison concerns the
documented Joern Scan pipeline, not every possible extension built on Joern.

Model assistance also has prior art. IRIS combines static analysis with LLM-assisted specification
inference and contextual path review. Semgrep describes AI-assisted triage and remediation on top of
its deterministic engine. Zazzles's hybrid structure therefore does not establish research novelty
or exclusive use of semantic reasoning. [IRIS authors' repository](https://github.com/iris-sast/iris),
[Semgrep's technical description](https://semgrep.dev/blog/2024/the-tech-behind-semgrep-assistant/)

Zazzles's intended advantage is a precise workflow around source-backed, separately named security
claims: explicit evidence gaps, semantic interpretation of supplied code, preservation of sufficient
protection and targeted repair guidance for agents. These are implemented design strengths; the size
of any advantage over another scanner must be measured on matched cases. Mature tools may have
broader framework coverage, deeper flow models and production workflows that Zazzles does not yet have.

## Analysis knowledge and proposed packs

Today, analysis knowledge ships in the binary: compiled language/framework and flow adapters,
verified API contracts, embedded versioned Analysis Guidelines, claim questions and Go evidence
policy. Application `--checks` files supply requirements and selection; they do not add a parser or
new vulnerability detector. Public analysis guides govern development but are not loaded as prompts
by editing their Markdown.

The proposed [analysis pack MVP](../.zazz/project/analysis-pack-mvp.md) will organize these existing
components under explicit IDs, versions, applicability and regression coverage. The static engine
will use references to supported models and evidence obligations; the judge will receive selected
interpretation guidance and applicable contracts alongside the candidate's trace. Executable adapters
remain compiled, and mandatory Go gates remain authoritative. Language, entry framework, library and
claim components should compose without duplicating a pack for every combination.

This is a separate future initiative/worktree after common-entry/dispatch merges and before command
injection. No pack loader, public rule language, downloadable registry or third-party plugin system
is implemented or authorized by this documentation. Changed effective judgment inputs must retain
exact-request replay protection and cannot inherit historical answers.

## Evidence for the approach

The [Python evaluation](../.zazz/evaluation/python-support-results.md) records 29 passing scoped case
checks across the original three pinned applications, including nine source-reviewed SQL/XSS positive
cases. It also records model overclaims that Go policy blocked and the validation calibration needed
to avoid stronger requirements being invented. These examples support the value of explicit policy
and evidence boundaries; they do not establish general accuracy.

The added VulnPy corpus exposes six known positive coverage misses and four warning controls. Those
results remain separate from the original 29-case success claim. The [XSS evaluation](../.zazz/evaluation/xss-implementation-results.md)
and sample answer guides record additional implemented behavior and gaps. Local Semgrep CE/Bearer
baselines are pinned, scope-specific observations, not full-product rankings. CodeQL and Joern documentation
here are architecture references, not measured CodeQL/Joern benchmarks.

## What remains to establish

Broader or novel-repository coverage, lower false-positive rates at equal recall, safe handling of
custom validators and lower end-to-end cost/latency are evaluation goals. A novel combination may
have product value without being a new static-analysis algorithm; no exhaustive originality study
or general superiority claim has been completed.

Comparisons need pinned source/rules/models, matched vulnerability and framework scope, independently
reviewed vulnerable/protected/unknown controls and held-out cases. Measure precision, recall, false
protection, abstention/coverage, evidence usefulness and remediation correctness alongside time,
resources and model cost. Fewer alerts alone does not establish improvement. The same agent's source
review, implementation and calibration are not independent findings review.

Independent LLM Findings Review, a full repository IR/graph, automatic incremental caches, severity
CI gates and dependency scanning remain planned. Current analysis is bounded static evidence, not
runtime confirmation or a security certificate. Paid model judgment also introduces source-transmission,
service-availability and cost considerations; offline evidence inspection remains available.

## Common entry and dispatch architecture

The shared model separates these responsibilities:

| Record | Meaning | Why it is separate |
|---|---|---|
| Entry | Logical exposed operation, such as an HTTP endpoint or schema field | The same implementation can serve different operations; transport is optional metadata |
| Registration | Source/configuration wiring and generation provenance | A generated endpoint needs evidence of how it was exposed |
| Dispatch binding | A known or unresolved callable alternative, stage or condition | Discovery of an operation does not establish its implementation or value flow |
| Callable | Lexical implementation identity and source snapshot | Shared code retains one identity without acquiring a global protection flag |
| Invocation context | Actual parameters, captured values, input roles and applicable controls | A guard or binding from one exposure must not leak into another |
| Gap | A located unresolved fact or reached limit | Unknown alternatives and zero-use entries remain inspectable without invented security findings |

For example, one handler exposed twice can have one callable record, two registrations and distinct
invocation contexts. If only one registration runs an enforcing guard, the other retains its own
unprotected or unresolved claim. Separately, a discovered registration with no resolved callable
remains a coverage record; it cannot be interpreted as a safe operation.

Adapters own framework evidence and input roles. Language flow engines own supported transfers,
helper/callback binding and control facts. Shared claim policy and presenters interpret the resulting
evidence without manufacturing a missing edge. This makes reports and selectors reusable across entry
types while leaving each adapter responsible for proving its own source relationships.

The common records, existing HTTP/browser consumer migration and typed reports/selectors are
implemented on the current development branch. Explicit browser-event callback dispatch is next (D4).
D5 adds bounded Express and Python generated registration, finite lookup and closure/callback flows,
with real VulnPy trace/finding acceptance and Juice Shop regressions. Node.js Apollo GraphQL (D6) and
grpc-js unary gRPC (D7) are planned adapters. Typed metadata/presentation support already present
for GraphQL/gRPC is not a claim that their source discovery or input propagation is implemented.
The [README use-case matrix](../README.md#entry-use-cases-and-delivery-status) summarizes supported
and planned scenarios; the [common-entry design](../.zazz/architecture/common-entry-dispatch.md)
and [delivery contracts](../.zazz/specifications/common-entry-dispatch.md) own the detailed boundaries.

The architectural benefit is a common evidence and invocation contract rather than a separate
security interpretation for every framework. It does not automatically yield deeper analysis,
fewer false positives, complete dynamic-language coverage or proof that all runtime paths were
examined. Source-reviewed real applications and independent raw/protected/unknown controls must
establish those capabilities incrementally. This is a deliberate engineering combination of known
techniques and constrained semantic judgment; research novelty and superiority remain unproven.
