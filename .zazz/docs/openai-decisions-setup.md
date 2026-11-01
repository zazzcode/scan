# Decision-model API keys and scanner setup

**Official API documentation checked 2026-10-31.** The scanner supports explicit OpenAI Decisions
or TypeSafe Jev selection. Offline analysis needs no key. Tests use local HTTP stubs; live [two-case results](../evaluation/juice-shop/results.md) are recorded, but held-out accuracy,
actual billing and provider latency benchmarks remain unmeasured.

## Key ownership for solo development

For one developer running scan locally, choose **Owned by: You**, name the key `scan-local-dev`, and
use the **Default project** if that is the only project offered. A separate `scan` project can help
organize usage later; it is not required to continue. Use a service account later for unattended
CI/deployment that should have its own credentials. This is a project recommendation, not a runtime
requirement. Both personal and service-account keys belong to an API project; see the
[project-key API reference](https://developers.openai.com/api/reference/resources/admin/subresources/organization/subresources/projects/subresources/api_keys/methods/retrieve).

## Create the OpenAI key

- Sign in to your personal account at the [API dashboard](https://platform.openai.com/).
- Open [API keys](https://platform.openai.com/api-keys), select **Default project** (or your chosen
  project), choose **Owned by: You** if offered, and create the secret key named above.
- Store it in a password manager when displayed. Do not paste it into chat, source, screenshots,
  commits, or shell command arguments. Use an ordinary model-call key, not an organization admin key.
- If selecting restricted permissions, allow the Decisions endpoint for this project. Check project
  access/limits in the dashboard; do not assume a signed-in Codex session establishes API access.

Sources: [OpenAI quickstart](https://developers.openai.com/api/docs/quickstart) and
[API authentication](https://developers.openai.com/api/reference/overview#authentication).

### Subscription credits

Your Go program uses an API project key even when you have a GPT subscription. The
[official credit guide](https://developers.openai.com/api/docs/guides/chatgpt-api-credits) says eligible
paid Pro and Business subscribers can claim API credits, including Decisions usage, by linking an API
organization they own. Check the benefit in ChatGPT billing or API Platform billing, activate it for
the intended organization if eligible, and inspect its balance. This session has not checked your
subscription tier, eligibility, or balance. A subscription alone does not establish funded API access.

## Set the OpenAI key locally

Run these commands in a local **zsh** terminal; the hidden prompt keeps the literal key out of command
history. Only processes launched from this shell inherit the variable. Go does not auto-load `.env`.

```zsh
read -r -s 'OPENAI_API_KEY?OpenAI API key: '
printf '\n'
export OPENAI_API_KEY
```

Verify presence without printing the secret:

```zsh
[[ -n "$OPENAI_API_KEY" ]] && printf 'API key is set\n'
```

Do not add the key to a command submitted to an agent. Set it in your own terminal or secret store.
An already-running Codex desktop process does not inherit a variable exported in another terminal.

## Run the scanner

From the active scan worktree, resolve `<target-root>` and its registration file first. Start offline:

```bash
go run ./cmd/scan --strings --json --entry server.ts <target-root>
```

For Juice Shop, the entry is `server.ts`. For `vulnerable-typescript-app`, use `src/server.ts`.
`--routes` remains the original syntax-only inventory. `--strings` adds selected handler resolution,
request-value dependency witnesses and candidate SQL/outbound HTTP uses. Neither mode executes the
application or requests the URLs it discovers.

When ready to transmit the bounded source evidence to OpenAI, explicitly opt in. This example permits
at most **one** request for the selected login route; remaining candidates are marked `not_run` and the command exits 1 to expose
incomplete assessment:

```bash
go run ./cmd/scan --strings --json --entry server.ts \
  --route /rest/user/login --judge openai --judge-model gpt-6-luna --judge-limit 1 <target-root>
```

The [Decisions guide](https://developers.openai.com/api/docs/guides/decisions) specifies
`POST https://api.openai.com/v1/decisions`. The adapter sends named choice questions and validates
answer names, types and allowed choices. See the
[create-decision reference](https://developers.openai.com/api/reference/resources/decisions/methods/create).

Each request has a 30-second HTTP deadline, a 64 KiB request-payload cap, a 1 MiB response cap and no
retries. The request-count limit is not a currency budget; inspect pricing/credits before live runs.
Provider errors/refusals remain failed assessments. Raw validation/protection answers are retained;
checked mode blocks only on claim-relevant missing facts, while unchecked legacy candidates remain conservative. Findings Review is not yet
implemented, so these are advisory assessments, not reviewed findings or security certification.

## Configure Jev instead

Get a key from the [TypeSafe console](https://console.typesafe.ai/) following the
[official quickstart](https://docs.typesafe.ai/introduction/quickstart). Set it in zsh:

```zsh
read -r -s 'TYPESAFE_API_KEY?TypeSafe API key: '
printf '\n'
export TYPESAFE_API_KEY
```

Select a model from TypeSafe's model catalog and pass its exact name explicitly; scan does not
silently choose a moving Jev alias. Replace `<jev-model>` with that provider model name:

```bash
go run ./cmd/scan --strings --json --entry server.ts \
  --judge jev --judge-model <jev-model> --judge-limit 1 <target-root>
```

The adapter uses TypeSafe's [System One API](https://docs.typesafe.ai/api/system-one) and its
[choice-question schema](https://docs.typesafe.ai/primitives/choice). Its wire contract is separate
from OpenAI Decisions. No provider fallback occurs and one provider's key is never sent to another.
Question intent is shared; probability calibration and measured accuracy are not interchangeable.

## Live connection verified

On 2026-10-31, one actual OpenAI Decisions request for Juice Shop's selected login route completed
with `gpt-6-luna`: 2,247 input tokens and zero output tokens. The key was loaded privately from the
ignored `.env` into the scanner's child environment; no credential was printed. Luna returned
`insufficient` for validation and protection. Static gaps kept policy protection `uncertain`.
This is an access/integration smoke test, not a security-accuracy benchmark.

The [published Decisions rate](https://developers.openai.com/api/docs/guides/decisions#pricing-and-availability)
is $0.10 per million uncached input tokens, with no output-token charge. The base estimate for this
request is $0.0002247 before applicable cache-write/regional adjustments; actual billing was not
queried. The user's $10 starting balance supports modest bounded experiments but is not a spend cap
implemented by the scanner. Live work on this branch is restricted to `gpt-6-luna`; the Jev setup
above documents an alternative adapter, not authorization to call another provider.

## Troubleshooting and cleanup

- Missing key: set the variable for the explicitly selected provider in the same terminal.
- `400`/`404`: verify the model name, endpoint schema and project availability.
- `401`/`403`: check the key's project and permissions; do not print the key when debugging.
- `429`: inspect quota/credits and rate limits in the selected provider's dashboard.
- Exit 1 with a JSON report: inspect assessment status/reason for errors or the request-count cap.
- Exit 0 offline means analysis completed within the stated capability, not that the target is safe.

```zsh
unset OPENAI_API_KEY TYPESAFE_API_KEY
```

Source and secrets stay outside tracked files. For a persistent key, use your password manager or
secret store. Follow the [judge architecture](../architecture.md#provider-neutral-judge-boundary) when
adding a future provider; it needs its own key, wire validation, capability tests and evaluation.
