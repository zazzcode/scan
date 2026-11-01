# Bearer capability inventory — version 1

**2026-10-31.** Goal: equal-or-better detection with lower false positives, then equivalent
capabilities across CLI/Pro. This inventory records gaps, not achieved parity. Version identifiers,
configuration and source-inspected outcomes are in [comparison results](bearer-results.md).
Bearer source: `e3753de6347988bed15bdba80b9ff731bf028c80`; rules:
`82c335f781a61a6d8b052394f84fd1d93edec467`. Product statements below describe that pinned README,
not a verified current commercial offering. Pro was unavailable for local evaluation.

The [Node.js Goof target inventory](nodejs-goof-capabilities.md) adds source-inspected template,
stored Markdown, database and other cases on 2026-11-01. Its [initial baseline](nodejs-goof-results.md)
is complete (43 findings / 12 rules; limited XSS detection); use the
[runbook](../../docs/samples/nodejs-goof/bearer-runbook.md) to produce and curate results under
`docs/samples/nodejs-goof/bearer/<run-id>/`. These cases do not change the observed results below.

| Capability | Pinned comparator evidence | Scan state | Phase / acceptance still required |
|---|---|---|---|
| JS/TS Express source analysis | Exercised both samples and controls | Bounded implemented | Phases 3–7; framework/version and mount/alias negatives |
| Other CLI languages | README lists Go, Java, PHP, Python, Ruby; not exercised | Language detection only outside JS/TS security analysis | Phases 11–13 and future language contracts |
| Pro languages | README additionally lists C#, Kotlin, Elixir, VB.Net, Rust, Swift; unverified | Unsupported security analysis | Future adapter contracts and independently labeled cases |
| SQL injection | Exercised vulnerable, bound and custom-validator cases | Local candidates and optional scoped judgments | Phases 3/5–8; held-out precision/recall, false-safe and abstention measurement |
| Custom validator precision | Two protected allow-list controls flagged | One saved protected judgment; imported guard blocked offline | Phases 5/7/8; lower false positives without lost recall on independently labeled variants |
| Helper / cross-file analysis | Helper SQL detected; dynamic parameter rule means caller linkage unproven | Helper context selection exists; helper-contained sinks missed | Phase 6; actual argument/parameter/return witnesses and safe/unknown pairs |
| General control/interprocedural analysis | Pro README advertises Java/Python/C#/Go; unverified | Bounded local control edges, no complete CFG/heap solver | Phases 5–6; branches, aliases, mutation, recursion/async acceptance cases |
| Outbound / SSRF | Direct Axios control detected; axios.default sample missed | Candidate use, no destination-specific checked contract | Phases 3/5–8; URL/body/header roles, custom destination policies, response continuation |
| Command/file/HTML/auth/crypto risks | Multiple rules exercised in sample; not all adjudicated | Mostly outside sink catalog | Phases 3–8/13; per-family reviewed source/sink models and controls |
| Sensitive data / privacy | Pinned README and rule catalog; logger/secret rules exercised | No general classification or privacy-impact analysis | Future Phase 3/4/9 contracts; data-type labels and privacy-negative controls |
| Explanations / remediation | Rule metadata includes locations and guidance | Route/source/control evidence, raw facets, illustrative templates | Phase 8/9; independent Findings Review and verified remediation evaluation |
| JSON / SARIF and consumers | JSON exercised; help/source expose SARIF and other formats | JSON and bounded SARIF implemented/schema checked | Phase 9; actual IDE/GitHub ingestion, grouping, edit-resistant identity |
| Rules / configuration | Local external rules, filters/config/ignore flags inspected | Checks and bounded models; no comparable general rule authoring contract | Phases 3–4/9; versioned models/rules, config precedence and scope accounting |
| Baselines / incremental / CI | Help exposes diff, fingerprints, ignore, severity exits; cache source inspected, workflows untested | Offline export exists; baseline/incremental/CI gate incomplete | Phases 9–10; changed findings, incomplete runs, cache equivalence and consumer tests |
| SCA | Not evaluated; no capability claim | Planned ecosystem adapters | Follow-on SCA contracts; advisory existence independent of reachability |

The [implementation plan](../project/implementation-plan.md) owns phase boundaries. Extend this
inventory when a capability is exercised or changes; preserve version/edition and evaluation scope.
Measure false-positive denominators, recall and abstentions together. A smaller finding count or
agreement between tools does not establish product superiority.
