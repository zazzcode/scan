# Technical approach and design strengths

[Documentation](README.md) · [Philosophy](philosophy.md) · [Implemented architecture](../.zazz/architecture.md)

## Source-backed analysis with constrained semantic judgment

Zazzles combines bounded static program analysis with a decision model that answers explicit security
questions about the extracted evidence. The intended output helps a developer or coding agent decide
what to repair, what protection to preserve and what uncertainty to investigate.

The Go scanner reads application source, resolves supported language/framework bindings and follows
values into security-sensitive operations. Each candidate carries source spans and hashes, relevant
helpers, branch conditions, operation operand roles, library contracts and unresolved facts. The
semantic judge receives this evidence together with versioned language **Analysis Guidelines**.
Go then checks the answers against claim-specific evidence requirements and consistency gates.

The model is part of security judgment, rather than only a prose explanation of a preexisting alert.
It does not generate the source relationships or fill missing code with assumptions. Source text is
evidence data; scanner-supplied guidelines define the interpretation instructions. JSON preserves
both evidence and raw answers, while HTML and SARIF present the policy-gated result.

## The distinctive design choices

The useful distinction is the combination and enforcement of these boundaries. It is not a claim that
syntax parsing, dataflow, taint analysis, model-assisted security review or SARIF originated here.

| Design choice | Practical strength | Current boundary |
|---|---|---|
| Separate claims for validation, construction and use protection | A fix can address the exact defect while preserving an existing control | Reports still group several concerns within one operation record |
| Source-backed evidence roles and claim blockers | A confident model answer cannot turn absent bindings, helpers or policies into protection | Incorrect extraction or an incorrect answer can still produce a bad result |
| Analysis Guidelines plus actual helper/guard evidence | Semantic questions can consider an implementation rather than only a validator's name | Custom behavior is judged only within resolved, bounded paths |
| Findings, Review and Protected outcomes | Agents receive both repair work and scoped controls to preserve, plus explicit investigation gaps | Protected applies to the named claim, not the route or whole application |
| Route/browser context and bounded stored/cross-tier links | Review can connect an input, operation and rendering context across supported boundaries | Dynamic registrations, storage identity and cross-tier coverage remain incomplete |
| Versioned evidence, guideline documents and exact-request replay | A historical decision remains auditable and cannot silently move onto changed input | Replay is explicit; an automatic persistent decision cache is not implemented |
| Shared language-neutral analysis records and reporting | JS/TS and Python use the same claim gates and output semantics | Each language still needs its own correct framework and flow models |

For example, a signup handler can bind an email safely into SQL while enforcing only presence, even
though the matching form establishes an email-format requirement. Zazzles can preserve the SQL-binding
claim and separately report the unmet validation policy. Conversely, a missing business requirement
stays Review: the model cannot invent a restriction merely because the field has a familiar name.
A construction defect can also merit repair without establishing an injection exploit.

## Relationship to established SAST

Conventional SAST includes much more than textual matching. CodeQL documents local/global dataflow,
taint tracking and source-to-sink path queries. These are established capabilities, not unique Zazzles
features. [CodeQL dataflow documentation](https://codeql.github.com/docs/writing-codeql-queries/about-data-flow-analysis/)

Semgrep documents source, sink, sanitizer and propagation models. Its Community Edition and
proprietary analysis have different interprocedural/interfile capabilities; a local CE comparison
must not stand in for every Semgrep edition. [Semgrep taint analysis](https://docs.semgrep.dev/writing-rules/data-flow/taint-mode/overview)

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

## Evidence for the approach

The [Python evaluation](../.zazz/evaluation/python-support-results.md) records 29 passing scoped case
checks across the original three pinned applications, including nine source-reviewed SQL/XSS positive
cases. It also records model overclaims that Go policy blocked and the validation calibration needed
to avoid stronger requirements being invented. These examples support the value of explicit policy
and evidence boundaries; they do not establish general accuracy.

The added VulnPy corpus exposes six known positive coverage misses and four warning controls. Those
results remain separate from the original 29-case success claim. The [XSS evaluation](../.zazz/evaluation/xss-implementation-results.md)
and sample answer guides record additional implemented behavior and gaps. Local Semgrep CE/Bearer
baselines are pinned, scope-specific observations, not full-product rankings. CodeQL documentation
here is an architecture reference, not a measured CodeQL benchmark.

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

## Future entry types

Dynamic routes and GraphQL require both framework adapters and a focused evolution of today's
HTTP-shaped entry/check records. The shared design separates logical entry points, registrations,
dispatch bindings, callables and invocation context; HTTP routes are one view. The existing HTTP
and browser paths should be the first consumers of that foundation before more entry types land.
Sink evidence, Analysis Guidelines, judgment gates and report
formats remain reusable. The [architecture extension assessment](../.zazz/architecture.md#future-dynamic-routes-and-graphql)
names entry identity, registration/dispatch, resolver input provenance, discovery completeness
and bounded traversal as the pressure points. Neither capability is implemented by the naming migration.
