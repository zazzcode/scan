# Original comparison corpus

**2026-10-31; small development evaluation, not a blind accuracy benchmark.** These seven original
source fixtures contrast SQL interpolation, parameter binding, custom allow-list guards, length-only
checks, conditional required-field checks, enforced nonempty validation and unknown middleware.
They were written for this implementation. They are not unseen real repositories or runtime tests.

## Label isolation and reproduction

Each neutral case directory contains target source and a requirements file. Expected labels and
rationale live in [manifest.json](manifest.json), outside those target roots. The launcher never sends
the manifest to the judge. Requirements describe properties, not expected answers. Source and
requirements are visible to the judge; no README challenge labels are added as evidence.

```bash
go run ./cmd/scan --strings --json --entry server.ts \
  --checks .zazz/evaluation/corpus/case02/checks.json \
  .zazz/evaluation/corpus/case02
```

Add `--judge openai --judge-model gpt-6-luna --judge-limit 1` for one paid assessment. Tests/stubbed
responses remain offline. Live results are stochastic; these recorded outcomes are observations,
not golden assertions about future model answers.

## First seven Luna requests

| Case | Expected protection | Observed protection | Observed scoped status |
|---|---|---|---|
| case01 | insufficient | insufficient | supported_potential |
| case02 | sufficient | uncertain | needs_evidence |
| case03 | sufficient | sufficient | protected_for_claim |
| case04 | insufficient | insufficient | supported_potential |
| case05 | insufficient | insufficient | supported_potential |
| case06 | sufficient | sufficient | protected_for_claim |
| case07 | uncertain | uncertain | needs_evidence |

Six of seven matched both the expected protection and policy status. Case02 was uncertain rather
than protected. Inspection found that the parser verified the Sequelize import internally but sent
only constructor expressions as sink-binding evidence. The corrected bundle now carries the actual
import declaration and receiver assignment, with a regression check that those facts remain visible.
One additional request on case02 returned **sufficient protection / protected_for_claim**, while field
validation remained insufficient. No answer was overridden and the uncertainty was not hidden.
The first seven observations and the changed-evidence recheck are distinct; this is not an unbiased
7/7 success rate or proof that every model error is an evidence omission.

The eight requests used **15,424 input tokens** and zero output tokens (base input estimate
**$0.0015424**, before applicable adjustments; actual billing not queried). No retries or
provider fallback were used. Together with this session's previous two-route iterations and connection
check, recorded usage totals **68,934 input tokens** (base estimate **$0.0068934**).

The SQL allow-list example establishes protection only for its exact query/character constraints;
parameter binding remains the preferred remediation. Unknown middleware remains needs_evidence.
No source finding is runtime-confirmed or independently reviewed.

## Output and remaining evaluation

All seven initial live reports, the corrected control, seven offline reports and the two-route Juice
Shop report exported successfully against the official SARIF 2.1.0 schema (16 reports). No consumer
integration was tested. Raw reports remain ignored local execution evidence under `.zazz/ephemeral`.

Next measure false findings, false safe claims and abstentions on unseen repositories, more custom
validators, aliases/mutations, framework versions and asynchronous paths. Separate business-policy
requirements from intrinsic sink protection. The second sample application's README is useful for
candidate selection but contains stale source layout and dependency claims: its SQL-like login logs
text without executing SQL, while its proxy actually calls Axios. Source evidence governs sink
classification; historical audit counts were not verified and no dependency audit was run.
