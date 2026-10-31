# Two-route validation evaluation

**2026-10-31; known target integration checks, not a blind accuracy benchmark.** Target: OWASP Juice
Shop `v20.3.0`, commit `f189f521266f42c6bd05884468292000cee852ed`. Requirements:
[validation-checks.json](validation-checks.json). No target server was started or attacked.

## Final focused run

```bash
go run ./cmd/scan --strings --json --entry server.ts \
  --checks .zazz/evaluation/juice-shop/validation-checks.json \
  --judge openai --judge-model gpt-6-luna --judge-limit 2 <juice-shop-root>
```

Credentials were loaded privately from the ignored `.env`; neither the key nor file contents were
printed. The adapter made one request per candidate with no automatic retry or fallback.

| Case | Source context records | Local control edges | Input tokens | Raw validation / protection | Policy finding |
|---|---:|---:|---:|---|---|
| POST login: email / SQL use | 42 | 0 | 8,548 | insufficient / insufficient | supported_potential |
| POST registration: required credentials / account creation | 47 | 11 | 10,593 | insufficient / insufficient | supported_potential |

Both responses identified `gpt-6-luna`, had zero output tokens and completed successfully. No declared
check blocker remained. Status means a supported advisory potential issue under the supplied scope
and declared library models; it is not runtime-confirmed exploitability, complete graph proof or
completed independent Findings Review. Zero login control edges in this artifact does not claim that
the full login route has no branches; this increment emits local edges for registration middleware.

The preceding broad-context attempt used 16,751 + 15,371 input tokens. Luna also answered insufficient,
but helper expansion limits correctly retained needs_evidence. This exposed over-expansion into
unrelated challenge bookkeeping. Input-relevant helper selection reduced that context before the
final two requests. The old responses were not reused as judgments on the changed bundles.

Four requests in this iteration used **51,263 input tokens**, with a base estimate of **$0.0051263**
at the published $0.10/1M Decisions input rate, before applicable cache-write/regional adjustments.
Actual billing/balance was not queried. The earlier connection smoke test is separate (2,247 tokens).
Raw local reports remain ignored execution evidence under `.zazz/ephemeral/`.

## Offline evidence and limits

Regression checks cover vulnerable/protected/unknown route/model fixtures, method-specific checks,
namespace and concise factories, custom helper/model context, imported binding evidence, shadowed
helpers, required missing facts, unknown applicable middleware, irrelevant coverage gaps, relevant
helper limits and deduplicated payloads. Checked-evidence and legacy goldens are retained in testdata.
Fixed HTTP answers test the adapters/policy; they do not measure model recognition of safe variants.

The second target, `vulnerable-typescript-app`, remains an independent source-only sanity check: its
SQL-looking login string is logged, not executed as SQL, and must not acquire a fake database sink.
Its outbound proxy URL is a useful later destination-validation case. No SCA audit was run.

Next accuracy work needs held-out custom validators and safe parameter-binding/persistence variants,
false-positive and false-safe measurements, stricter framework/version/mutation coverage, and an
independent findings review. Runtime confirmation additionally requires an isolated application test
that checks actual authentication or persisted account state, not merely a challenge flag.

## Reporting and original controls

The saved final report exports offline to SARIF with route identity, source locations, supported
SQL dependency paths, registration control evidence, raw facets, scoped status and illustrative
remediation templates. It passes the official SARIF 2.1.0 schema. Examples are not verified fixes.

The [original seven-case comparison](../corpus/README.md) matched six expected outcomes initially.
A parameter-binding uncertainty exposed omitted import provenance in the submitted bundle; one
changed-evidence recheck then recognized protection while retaining insufficient field validation.
Those fixtures are development controls, not a held-out repository accuracy benchmark. Their eight
requests used 15,424 input tokens; total recorded session usage including the connection check and
two-route iterations is 68,934 input tokens (base estimate $0.0068934 before adjustments). Actual
billing was not queried. Independent review and runtime verification remain unperformed.
