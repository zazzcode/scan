# Python adapter implementation and evaluation

[Contract](../specifications/python-support.md) · [User guide](../../docs/python-support.md) ·
[Sample reports](../../docs/samples/README.md)

## Delivered behavior

The Python adapter uses Tree-sitter Python v0.23.6 and supplies the existing graph, versioned decision
protocols and JSON/HTML/SARIF pipeline. It discovers Flask/Blueprint, RESTX Resource and FastAPI routes,
tracks a bounded value-flow subset, verifies sqlite3/logging/HTTP bindings and models Flask/Jinja and
explicit FastAPI HTML. Inline template scripts reuse the browser adapter, including verified jQuery
input-to-HTML flow. No application code, dependency installation, database query or outbound request
from the targets was executed.

Input validation remains separate from SQL/XSS protection. Source-backed forms can establish policy
for the matching handler; client constraints do not establish server enforcement. For ThreatByte
signup, the form declares email syntax, while the backend shows presence/password-match checks.
The live assessment reports insufficient validation on the insert, independently of its sufficient
SQL parameter binding. Missing business requirements remain Review. Current per-use validation
facets are not independent confirmed per-field defect counts.

## Pinned corpus and final results

| Target | Revision | Route/method entries | Browser entries | Records / assessed | Static SQL/XSS findings | Other outcomes |
| --- | --- | --- | --- | --- | --- | --- |
| ThreatByte | `ede28f2` | 23 | 1 | 51 / 51 | 4 | 32 protected, 13 needs-evidence, 1 construction-practice, 1 destination-review |
| python-insecure-app | `b188493` | 1 | 0 | 2 / 2 | 1 | 1 validation-review |
| vulnerable-flask-app | `a67cbe5` | 19 | 0 | 24 / 24 | 4 | 20 needs-evidence |

The nine static findings match source-reviewed positive discovery cases: two SQL injection cases,
ThreatByte reflected/stored/DOM XSS, FastAPI HTML/template-source XSS, and three Erlik raw-response
paths. All [22 scoped case checks](../../docs/samples/evaluate-python.py) pass. This is case discovery
and classification evidence, not exhaustive recall, zero false positives on every operation, runtime
exploitation, or independent human findings review. The same agent reviewed source and implemented
the adapter; expected labels were kept separate from detector code and model evidence.

## Accuracy iterations

- Initial source review corrected the handoff: Erlik `/welcome` returns JSON and is an XSS negative;
  `/welcome2` returns raw HTML text. The initial implementation found seven of the final nine
  reviewed SQL/XSS cases.
- The first sample loop exposed a missing exception-body grammar shape. Prefix environments and
  exception preconditions now preserve input assigned before a failing operation, recovering
  Erlik's exception-return case.
- Source review found ThreatByte's inline jQuery comment append, absent from Python-only and saved
  comparator results. Snapshot-preserving script extraction and the existing browser flow adapter
  recover this ninth case. Native DOM append(string), fake bindings and escaped siblings are controls.
- Added method feasibility so POST-only bodies are not assigned to GET entries. Deduplicated value
  evidence and bounded helper/path expansion. Relevant unknown pre-use calls block claims without
  blocking a query on an unrelated field.
- Paired tests cover escaped versus raw/unknown branches, script/URL contexts, Markup/safe versus
  actual escaping/forceescape, imports and keyword helpers, stored-field origins, literal versus
  parameterized SQL, logger bindings, password/email content versus derived facts and outbound slots.
- Source-backed form policy supplies concrete validation questions without treating HTML required/type
  attributes as backend enforcement. Standalone fields with unknown policy stay Review.

## Live calibration and disagreement

The user authorized `gpt-6-luna` only and ultimately removed the request-count cap. There were **78
live requests**, all HTTP 200, totaling **306,362 reported input/total tokens** and zero reported
output tokens for this Decisions endpoint. Every current record has a completed assessment. The
extra request is the changed literal-email logging claim. Exact request/response traces and the
append-only attempt ledger are local; reports retain answers, probabilities/confidence and usage.
Confidence never establishes a verdict. Unchanged requests were replayed only after exact payload,
model and question-version matching; changed evidence received a fresh assessment.

Two model overclaims were caught by structural policy:

- ThreatByte's dynamic update columns are assembled from explicit key allow-lists while values are
  bound. The model answered vulnerable, but unsupported comprehension/key reasoning remains a
  blocker. A construction concern is retained; no confirmed injection claim is emitted. Precise
  allow-list/comprehension modeling remains a follow-up.
- Erlik logs a literal internal role mailbox. The model called it personal data; email syntax alone
  does not establish privacy or a person-identifying value. A generic sensitivity blocker, tested
  against overconfident responses, keeps that candidate Review. Database-result and whole-session
  logging also remain Review. Request password/email/header content controls are tested separately.

Some escaped stored fields received model answers inconsistent with static facts or lacked complete
origins. They remain needs-evidence. Unknown subprocess return content is not converted into XSS
confirmation. Lack of model agreement does not erase an established construction-practice concern.

## Comparator observations

| Target | Semgrep CE 1.183.0, snapshot 20261101 | Bearer e3753de, rules 82c335f | Source-backed differences |
| --- | --- | --- | --- |
| ThreatByte | 7 results | 6 results | Semgrep reports two locations for the dashboard SQL case; neither saved baseline reports the three reviewed XSS cases or outbound URL. Other findings include MD5, traversal, debug and IaC. |
| python-insecure-app | 0 | 0 | Source review and Zazzles identify the HTML/template-source flow. README SSTI/SCA/secret claims are separate families. |
| vulnerable-flask-app | 17 results | 2 results | Both report SQL; Semgrep also covers raw responses and several out-of-scope families. Neither saved baseline reports render_template_string at line 51 or the logging cases. |

Normalized baselines and upstream attribution are under each sample directory. A comparator's result
count spans its own rules and can contain multiple locations for one issue; these totals are not
precision/recall scores. No additional rule-pack downloads were made.

## Verification and provenance

Native Go 1.27.3 darwin/arm64, with `GOTOOLCHAIN=local`, `GOPROXY=off` and a local ignored GOCACHE:

- `go test ./internal/... ./cmd/...` and `go test -race ./...` pass, including existing JS/TS tests.
- `go vet ./...` and `go build ./...` pass; changed Go files are gofmt-clean.
- Python route golden, mixed-language CLI/bundle test, deterministic evidence, paired precision
  controls and a local HTTP-stub Python evidence→decision→report test pass.
- All 22 case checks pass on curated reports. Source excerpt byte ranges and hashes were checked;
  canonical JSON/SARIF parse, every record has an assessment, and public roots contain no local paths.
- ThreatByte and Erlik remained clean at their pins. Finder metadata appeared in the original FastAPI
  checkout during final verification; it was preserved. Final analysis used a fresh clean upstream
  clone at the same pin with identical source hashes and two exact decision replays, no extra calls.

Raw runs remain under `_test-findings/<app>/run-<id>/`; the local results report names exact runs and
records source checks. The source clones were not remediated or dependency-modified. Work is committed
locally on `python-support-base`; there is no integration push, PR publication or merge.

## Remaining boundaries

No selected app uses a MongoDB driver, so Python NoSQL is unimplemented/out of corpus. Dynamic
registration, custom decorators, arbitrary helper/control/heap semantics, loops/comprehensions,
full template inheritance/macros, exact stored-row identity and deployed schema verification remain
incomplete. Outbound destinations are review candidates, not a new SSRF verdict engine. New
vulnerability families, SCA and other languages remain separate work. General novel-repository
coverage and independent accuracy measurements still require more labeled controls and review.
