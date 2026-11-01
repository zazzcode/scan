# Optional decision-model assessment

[Documentation index](README.md) · [Getting started](getting-started.md) · [Interpret results](interpreting-results.md)

## Start offline

Source extraction, deterministic practice findings, saved-report export, and local tests need no API
key. Offline findings are useful remediation tasks where the rule establishes the practice. A model
can help assess custom validation/protection semantics from supplied evidence; it cannot establish
missing source facts. Assessment completion is separate from independent review and runtime confirmation.

## Credentials and local environment

Configure only the provider you deliberately choose:

| Scanner provider | Environment variable | Selection |
|---|---|---|
| OpenAI Decisions | `OPENAI_API_KEY` | `--judge openai --judge-model gpt-6-luna` |
| Jev | `TYPESAFE_API_KEY` | `--judge jev --judge-model <exact-provider-model>` |

Obtain a key through the selected provider's account/console. Keep it in a password manager or local
ignored `.env`, never chat, tracked files, logs, or command arguments. The tracked
[.env.sample](../.env.sample) has empty credential fields. Preserve an existing `.env` rather than
replacing it. See [environment loading](getting-started.md#local-environment-file): scan reads process
environment variables and does not automatically load the file.

A funded API account/key is separate from source-transmission permission and the scan's request cap.
The scanner has no currency-budget enforcement; choose an explicit small limit. There is no automatic
provider fallback, and Jev requires an explicit model name. Agents must also obey the user's active provider/model restrictions.

## Scope a live run

Confirm the target root, entry file, selected route, and allowed source evidence before running:

```bash
(
  set -a
  . ./.env
  set +a
  ./bin/scan --strings --json --entry server.ts \
    --route /rest/user/login --judge openai --judge-model gpt-6-luna \
    --judge-limit 1 /path/to/app
)
```

This sends bounded source evidence to the selected provider. The cap is per run; there are no retries.
Additional uses beyond the cap remain unassessed and cause exit 1. A requested assessment failure,
oversized payload, or exhausted cap is retained in the report; a setup failure before scanning can
produce no report. Optional `--checks` files specify business and protection requirements; they are
not expected model answers.

## Read outcomes separately

Model choices include sufficient, insufficient, uncertain, and not applicable for the named question.
The raw answers, execution state, and evidence completeness are retained. A supported potential
vulnerability alone remains Review in current presentation. Existing deterministic poor-practice
findings do not disappear because the model was not run or could not establish exploitability.

For authentication/permissions, verify the selected account and key locally without printing it.
For quota/rate failures, inspect that provider's account limits. For an exit-1 JSON report, read the
assessment reason rather than treating it as an empty or protected scan. Saved HTML/SARIF exports
need neither credentials nor another model call.
