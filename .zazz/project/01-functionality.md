# 01 — Functionality

What `scan` does, expressed as capabilities rather than implementation. Design context:
[`../architecture.md`](../architecture.md).

## 1. Purpose

Given a target repository, produce a ranked, evidence-backed list of potential OWASP Top 10
vulnerabilities, where each finding is tied to a concrete source-to-sink call path — or to a known-bad
code pattern — and states how it was detected.

## 2. Detection tracks

Every finding belongs to exactly one track, and the track is recorded on the finding.

| Track | Detection | Determinism | Carries |
|---|---|---|---|
| **A — pattern rules** | Rule matches a known-bad shape; the rule is the verdict | Deterministic; byte-stable across runs | Most of the OWASP Top 10 |
| **B — path rules** | Deterministic discovery of a source→sink path, then model judgment of whether the data is sufficiently validated | Discovery deterministic; verdict is not | Custom validators, ORM/filter parameterization, template escaping, allow-lists |
| **C — scope review** | The model reviews a scope — a route file, a controller, a middleware chain, an auth module — rather than a path | Discovery deterministic; findings are not | Authorization, logic, and design flaws that are not dataflow properties |

Consequences: Track A findings gate CI by default. Track B findings carry a confidence and can be
treated as advisory. The tool is useful with no model configured.

## 3. Functional domains

**Ingestion and structure**

- **FD-1 — Repository ingestion.** Walk a target repository, filter by language and path, respect ignore
  rules, read content concurrently.
- **FD-2 — Multi-language parsing.** Parse each file with the language-appropriate Tree-sitter grammar.
- **FD-3 — Structural resolution.** Extract declarations, imports, call sites, and the type hierarchy
  into a normalized graph.
- **FD-4 — Call-site argument binding.** Record, at each call site, which argument expressions are passed
  and which parameter position each lands in. This is the flow evidence Track B depends on.

**Classification**

- **FD-5 — Source identification.** Route registrations, handler signatures, CLI entry points.
- **FD-6 — Sink identification.** Dangerous operations, each assigned a sink class: `nosql`, `sql`,
  `shell`, `path`, `ssrf`, `html`, `deserialization`, `log`.
- **FD-7 — Guard identification.** Validation, sanitization, and authorization functions, as first-class
  nodes so they can be included in a bundle even when off-path.
- **FD-8 — Data-type classification.** PII, PHI, credentials, payment data — by name and source pattern
  only, never by value. An impact signal, not a detection.

**Detection**

- **FD-9 — Pattern rules (Track A).** Known-bad shape detection across the OWASP Top 10.
- **FD-10 — Code path tracking (Track B).** Source-to-sink reachability across files, dead-end pruning,
  deduplication, pre-filtering, and prioritization. The spine: every other capability consumes this graph.
- **FD-11 — Scope selection (Track C).** Select the route files, controllers, middleware chains, and auth
  modules worth reviewing, without requiring a source→sink path.

**Judgment and reporting**

- **FD-12 — Context bundling.** Token-bounded bundle per unit of review — a path or a scope: ordered
  snippets, argument bindings, sink class and mitigation idiom, data types, and off-path guard
  definitions.
- **FD-13 — Model judgment.** Track B: a three-state verdict (`validated`, `not_validated`, `uncertain`)
  with mandatory cited evidence and a confidence. Track C: findings with a severity and an anchor symbol.
- **FD-14 — Reporting.** JSON and Markdown, with OWASP category, CWE ID, severity, data classes,
  detection track, and ordered `file:line` evidence.
- **FD-15 — Findings tracking.** Stable fingerprint; report what is new, persisting, and resolved.
- **FD-16 — Baseline gating.** Fail CI only on findings absent from a Git-tracked baseline; per-finding
  suppression and rule-level skips.
- **FD-17 — Incremental rescan.** Reuse unchanged parse results and previously recorded paths; re-judge
  only what a change affects.

## 4. Vulnerability coverage

Track A targets, by category:

- **A01** — path traversal from user input into filesystem APIs; open redirect; sensitive data in a
  response.
- **A02** — weak hashing and ciphers; disabled certificate verification; plaintext transport.
- **A03 (SQL)** — raw SQL built by concatenation with request-derived data; `eval`; `new Function`;
  shell execution with interpolated input; unsafe template or YAML deserialization.
- **A03 (NoSQL)** — a request-derived object passed straight into a Mongo filter; `$where`; `$regex`
  built from user input.
- **A04** — missing encryption of sensitive fields; persistent cookies carrying sensitive data.
- **A05** — permissive CORS with credentials; debug mode enabled; missing security headers; cleartext
  secrets in a cookie or JWT.
- **A07** — hard-coded credentials; weak password hashing; improper certificate validation.
- **A08** — deserialization of untrusted data.
- **A09** — sensitive data written to a log.
- **A10** — SSRF where the outbound URL is literally request-derived.

**A06 is out of scope** — vulnerable and outdated components is dependency scanning, a different tool.

Track B targets the same categories where the decision is *whether a mitigation is sufficient* rather
than *whether a dangerous shape is present* — for example, whether a Mongo filter's validation is
adequate, or whether a custom sanitizer is sufficient for the sink it feeds.

Coverage grows rule by rule. The rule set is **data** — queries plus declarative rules — not code paths
inside the scanner, so adding a rule must not require changing Go code.

## 5. Supported languages

| Order | Language | Why it is in this position |
|---|---|---|
| 1 | TypeScript / JavaScript (Node.js) | Hardest resolution case: barrels, path aliases, inline anonymous handlers, middleware by reference |
| 2 | Python | Different axis: decorators as sources, `self`, no static types |
| 3 | Go | Receiver methods, interface dispatch |
| 4 | Java | Annotations, and method overloads stress the node ID scheme |

Each language contributes exactly three things: a grammar, source/sink/guard query sets, and a resolver.
Adding a language must not require changes to traversal, pruning, or the LLM stage.

## 6. Interfaces

- **CLI** — `scan <path>` with flags for language selection, output format, baseline, and LLM
  configuration. Subcommands for inspection: `--list`, `--graph`, `--classify`, `--paths`, `--scopes`,
  `--explain <path-id>`.
- **Machine-readable output** — a JSON findings document, stable enough for CI consumption, with a
  versioned schema.
- **Human-readable output** — a Markdown report grouped by severity and OWASP category, each finding
  showing its call path and its detection track.

## 7. Open questions

- Should the CLI make the "not a guarantee" caveat part of the normal output rather than a `--help` note?
- How are false positives suppressed across runs — baseline entries, inline annotations, or both?
- Is a Bearer-style privacy report a future output, given `data_classes` already exists?
- Should rule documentation (the pattern each rule matches) be generated from the rules themselves?
