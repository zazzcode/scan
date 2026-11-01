# scan

`scan` is a Go static application security scanner under active development. It aims to help
developers and coding agents understand how request inputs reach security-sensitive operations,
what validation and protection already hold, and where code needs hardening. Reports should give
reviewers source evidence and explicit analysis gaps so they can judge each finding.

The approach combines Tree-sitter source extraction with code-flow analysis and optional decision-model
assessment of relevant evidence. Field validation and sink protection are separate questions: a valid
email address does not establish that a SQL query uses bound parameters. Unknown behavior stays
unknown rather than becoming a safety claim or a vulnerability by assumption.

## Goals and current scope

The goal is to report vulnerabilities, weak security practices, and sufficient protection by route,
with useful explanations and remediation guidance. Longer term, scan aims to match or exceed Bearer's
SAST and developer-workflow capabilities with fewer false positives, measured at comparable detection
coverage on independently labeled cases. That is a project objective; broad parity is not established.

Today, scan analyzes local JavaScript and TypeScript source. It inventories Express-style route
registrations, follows a bounded subset of handlers and request-value dependencies into candidate SQL
and outbound HTTP uses, and exports source evidence as JSON, standalone HTML, or SARIF. Optional OpenAI Decisions or Jev
judging assesses validation and protection from that evidence.

Full control-flow and interprocedural analysis, independent findings review, Python and Go analysis,
and dependency vulnerability scanning remain planned. Recognizing a file extension in the file list
does not establish analysis support. A clean scan does not certify that an application is secure.

## Setup

1. Install Go 1.27.2 or newer. Check with `go version`.
2. Install a C compiler and keep it on `PATH`. The Tree-sitter bindings use cgo. On macOS
   that is the Xcode command line tools (`xcode-select --install`).
3. Clone this repository and enter its root, the directory that contains `go.mod`. Contributors
   use the feature-worktree setup described under [Contributing](#contributing).
4. The first `go run` or `go test` downloads the modules in `go.mod`. No separate install step.

The program scans a directory you already have on disk. It does not clone, install, or start
the target application.

### Compile

Build a local executable from the directory containing `go.mod`:

```bash
go build -o bin/scan ./cmd/scan
./bin/scan --list /path/to/app
```

### Local environment file

The tracked [.env.sample](.env.sample) lists supported local environment settings with empty key
fields. For a new setup, copy it to `.env`. If `.env` already exists, merge only the needed settings;
do not replace your existing keys. `.env` stays Git-ignored.

```bash
cp .env.sample .env  # New setups only; do not overwrite an existing .env.
```

Set the output root in `.env`, adjusting the path for your checkout layout:

```bash
SCAN_FINDINGS_DIR="$HOME/Dev/zazzcode/scan/_test-finding"
```

Load your own trusted `.env` in bash or zsh before running scan:

```bash
set -a
. ./.env
set +a
```

The scanner reads environment variables and does not auto-load `.env`. This sample is for local
shell use; CI should set environment variables and secrets through its workflow configuration.

### Optional decision-model setup

Offline scanning and tests need no API key. For live assessment, configure `OPENAI_API_KEY` for
OpenAI Decisions or `TYPESAFE_API_KEY` for Jev. See the [decision API guide](.zazz/docs/openai-decisions-setup.md)
for account/key creation, provider-specific details, and troubleshooting.

You can keep the selected key in a local `.env` at the worktree root; `.env` and `.env.*` are Git-ignored.
The scanner reads process environment variables and **does not load `.env` automatically**. If you
already have a key there, reuse it. In bash or zsh, load your own trusted shell-compatible `.env` and
run a capped assessment in a subshell:

```bash
(
  set -a
  . ./.env
  set +a
  go run ./cmd/scan --strings --json --entry server.ts \
    --route /rest/user/login --judge openai --judge-model gpt-6-luna \
    --judge-limit 1 /path/to/app
)
```

This example transmits bounded source evidence and permits at most one request; additional candidates
remain unassessed and may cause exit 1. The request cap is not a currency budget. Keep keys out of
source, commits, logs, chat, and command arguments. Jev requires an explicit provider model name;
there is no automatic provider fallback.

## Run

After compiling, configure a results root outside the worktree and run a local scan:

```bash
export SCAN_FINDINGS_DIR="$HOME/Dev/zazzcode/scan/_test-finding"
./bin/scan --strings --entry server.ts /path/to/juice-shop
./bin/scan --strings --entry src/server.ts /path/to/vulnerable-typescript-app
```

Each command prints a new run folder. The layout keeps repositories and successive runs separate:

```text
_test-finding/
  juice-shop/
    run-<id>/
      findings.json
      findings.html
      findings.sarif
  vulnerable-typescript-app/
    run-<id>/
      findings.json
      findings.html
      findings.sarif
```

`SCAN_FINDINGS_DIR` is the output root, and `--output-dir <root>` overrides it for one strings scan.
Use `--output-dir ""` to disable bundle output for one run. Paths are resolved relative to the working
directory; use an absolute path for shared storage. The repository folder is named from the target
root's basename, so clones sharing a basename share that folder but receive separate run folders.
No existing run is overwritten. This local directory is a configurable convention, not a hard-coded
scanner default. You may keep the variable in your ignored `.env`; load it in the invoking shell,
because the scanner does not automatically read `.env`. It is a local filesystem path, not a URL.

Open the printed `findings.html` path in your browser (on macOS: `open /path/to/run/findings.html`).
Expand a source row to see the HTTP method and path, findings, validation and operation-protection
explanations, code evidence, unresolved analysis, and remediation recommendations. HTML embeds its
styles and content and uses native disclosure controls; no server, JavaScript, or external assets
are required. It presents scan results only. JSON retains the richer evidence for agents; SARIF is
the standard interchange format for security-tool integrations. Generated reports contain source
excerpts, so keep them with the same access controls as that source.

Without an output root, existing stdout modes still work. You can also regenerate views from JSON:

```bash
./bin/scan --strings --json --output-dir "" --entry server.ts /path/to/app > findings.json
./bin/scan --export-html findings.json > findings.html
./bin/scan --export-sarif findings.json > findings.sarif
```

Exports work offline without the target checkout, credentials, or new model calls, and ignore
`SCAN_FINDINGS_DIR`. Regenerate HTML and SARIF after JSON changes; they are snapshots. If an output
root is configured, `--json` and `--sarif` retain machine-readable stdout and also write the bundle;
the folder location goes to stderr. Text mode prints only the bundle location and artifact names.

For development, `go run ./cmd/scan` accepts the same flags as `./bin/scan`:

```bash
go run ./cmd/scan --list [--json] <path>
go run ./cmd/scan --routes [--json] <path>
```

`<path>` is the root of the application to scan. `--list` and `--routes` are separate modes.
Passing both, or omitting `<path>`, exits `2`. A missing directory exits `1`. Success exits `0`.

`--list` prints the supported source files and language counts. `--routes` prints one
`METHOD path file:line` line per route, then one gap line per gap. Add `--json` for indented
JSON on stdout.

```bash
go run ./cmd/scan --routes --json /path/to/app > graph.json
```

## Interpret results

Read the explanation and evidence before treating a record as a vulnerability. The JSON preserves
raw `entries`, `uses`, `gaps`, `files`, and `limitations`, and adds `presentation_version` and human-readable
`findings`. Each finding references its `use_id` and `route_id` and includes a title, explanation,
separate `input_validation` and `use_protection` text, uncertainty, and remediation recommendations.
A record's `record_type` distinguishes a finding, an unresolved analysis item, and supported protection.

| Display status | Meaning and next action |
|---|---|
| Observed code pattern (`observed_pattern`) | A specific construction or missing modeled check was observed. Review the named evidence and scope. This is not a confirmed exploit. |
| Analysis needs review (`needs_review`) | Scan identified a possible sensitive use but cannot establish its security behavior. Resolve the unknown library, value source, or check before making a security claim. |
| Evidence incomplete (`needs_evidence`) | Claim-relevant facts are missing. A model answer cannot close those gaps. |
| Supported potential issue (`supported_potential`) | An assessment supports the named issue within the declared evidence scope. Independent review and runtime confirmation remain separate. |
| Protection supported for this claim (`protected_for_claim`) | Supplied evidence and assessment support the named protection. Other fields, uses, and weaknesses can still require work. |
| Not applicable to this claim (`not_applicable`) | The named assessment does not apply; this is not a route-wide pass. |

Assessment state is separate: `not_run` means no semantic assessment, `failed` means the requested
assessment failed, and `completed` means an answer was received. A completed call is not itself a pass.
Input validation asks whether the field meets its business requirements. Operation protection asks
whether its use is safe for that particular database or other API. SQL parameter binding can preserve
query structure while field validation is still weak; NoSQL documents need slot/type/operator review.
A string inserted into an ORM search value is not automatically executable SQL. Free-text punctuation
is not by itself a defect, and logging SQL-looking text is not a database operation.

No validator identified means none was found in the analyzed slice; it does not prove all paths lack
validation. Review entry-level and report-level gaps as well as findings. A route with no reported
findings is not certified safe. Shared helper observations can repeat for several routes, so counts
are route/use records rather than unique vulnerabilities. Remediation is guidance; code examples are
illustrative and require driver-specific adaptation and regression checks.

Exit `0` means the requested operation completed, not “no security findings.” Exit `1` indicates an
execution or requested-assessment failure, including a request cap or payload limit; a retained JSON
report may still be useful. Exit `2` indicates invalid CLI usage. Exit codes currently do not implement
a vulnerability-severity CI policy.

## CI/CD artifacts

Set `SCAN_FINDINGS_DIR` to a runner-local output root or pass `--output-dir`. Run offline unless model
assessment is deliberately configured with the CI platform's secret store. A requested model failure
still writes the available report bundle and exits 1; an earlier setup/parse failure may produce no
bundle. Upload artifacts even after a failed scan. This example assumes the scanner binary has been
built and the target is checked out at `target/` with its entry file at `server.ts`:

```yaml
- name: Scan target
  env:
    SCAN_FINDINGS_DIR: ${{ runner.temp }}/scan-findings
  run: ./bin/scan --strings --entry server.ts target

- name: Save scan reports
  if: always()
  uses: actions/upload-artifact@v7
  with:
    name: scan-findings-${{ github.run_id }}-${{ github.run_attempt }}
    path: ${{ runner.temp }}/scan-findings/
    if-no-files-found: warn
    retention-days: 7
```

Download the artifact to open the HTML locally and retain JSON for agent analysis. Uploading an
artifact does not automatically import SARIF into a code-scanning service; use that service's SARIF
ingestion step separately when enabled. Review artifact access/retention because reports contain
source excerpts. Detection does not currently fail CI on finding severity; exit codes describe
execution and requested-assessment completion. See GitHub's
[artifact guide](https://docs.github.com/en/actions/tutorials/store-and-share-data) and
[upload action](https://github.com/actions/upload-artifact) for supported runner/edition requirements.

## What the graph contains

`scan --routes` parses `.js`, `.jsx`, `.mjs`, `.cjs`, `.ts`, `.tsx`, `.mts`, and `.cts`.
It skips `node_modules`, `vendor`, `dist`, `build`, and the other dependency and output
directories listed in `internal/walk`. `--list` still counts Python, Go, and Java files.
Those files are not mined for routes.

A route is a member call named `get`, `post`, `put`, `patch`, `delete`, `all`, or `use`
whose path is a string literal starting with `/`, or an array of such strings.
`receiver.route('/path').get(handler)` is one route per verb. Arguments after the path are
recorded in order. A member call with one function argument is unwrapped once, so
`wrapper(handler())` keeps both names. An inline function in the same file contributes its
direct call statements. A name with no function body in that file stays on the route and
adds an `unresolved_handler` gap.

The JSON document is `schema_version` 1 with `root`, `routes`, and `gaps`. Each route has
`id`, `method`, `path_pattern`, `framework` (`express`), `registration`, `handler`, and
`steps`. Spans are the file, 1-based lines, and byte offsets. Gap kinds are `dynamic_path`,
`unresolved_handler`, and `parse_error`.

## What it does not follow

The graph is the registration written in the target's own source. It is not an execution
trace.

- Express itself is not parsed. `node_modules` is skipped, so `app.post` stops at the call
  in the application. Router dispatch, `next()`, and Express internals are not steps.
- Imported middleware and handlers are names and source spans only. `login` imported from
  another file is an `unresolved_handler` gap. The body of that file is not attached.
- User middleware defined in the same file is followed only for direct call statements in
  that function. Calls nested in `try`, `if`, `await`, or a `const` initializer are omitted.
- A wrapper is unwrapped one level. The implementation of the wrapper is not opened.
- `app.use(someMiddleware())` with no string path is not a route.
- A path that is not a string literal is a `dynamic_path` gap and creates no route.
- The graph does not decide whether a value was validated or sanitized.

## Route-to-string-use analysis

```bash
go run ./cmd/scan --strings --json --entry server.ts <target-root>
```

Entry paths are relative to the target repository. Use `server.ts` for Juice Shop and
`src/server.ts` for the vulnerable TypeScript sample. Omit `--entry` to consider all selected JS/TS
registration files. Explicit selection avoids conflating training snippets with production scope.

The separate schema-v1 evidence report contains `entries`, `uses`, `files`, `gaps` and `limitations`.
Each use identifies its route, source fields, ordered local transfer evidence, candidate sink,
snapshot hashes and assessment state. Direct named relative ESM imports, imported handler factories,
local assignment, string expressions and simple object destructuring are supported. Branches,
unknown transformations, middleware, receiver ambiguity, unverified library bindings and external
response continuations remain explicit gaps. A SQL-shaped string logged to the console is not
reported as an SQL call. HTTP client registrations are excluded by the new receiver filter.

Use `--route /rest/user/login` to select an exact path pattern before model requests.
Offline assessment is `not_run`. For provider selection and key instructions, including personal-key
ownership and Jev configuration, see [decision API setup](.zazz/docs/openai-decisions-setup.md).
`--judge openai|jev` explicitly transmits bounded evidence; `--judge-limit` caps requests (default 10,
no retries). A requested judge failure, skipped oversized payload or exhausted request cap exits 1
with its report retained. Neither a successful judge call nor offline exit 0 certifies the target.
The graph design and full CFG/interprocedural analysis exceed this first implementation slice.

## Explicit route validation checks

Use a requirements file to declare the method, path, fields and separate validation/protection
constraints. This avoids guessing business requirements from a helper's name or a missing regex.
Two known target cases are tracked in `.zazz/evaluation/juice-shop/validation-checks.json`:

```bash
go run ./cmd/scan --strings --json --entry server.ts \
  --checks .zazz/evaluation/juice-shop/validation-checks.json <target-root>
```

Add `--judge openai --judge-model gpt-6-luna --judge-limit 2` for a capped live run. Checked output adds
relevant middleware/helpers/model context, local registration control edges, required evidence roles,
blocking vs coverage gaps, and a scoped finding status. `supported_potential` is an advisory finding
under the declared models, not runtime confirmation or completed independent review. See the
[evidence design](.zazz/architecture/evidence-and-claims.md) and
[results](.zazz/evaluation/juice-shop/results.md). SCA remains planned. Automatic claim assembly for
novel repositories remains incomplete; see [discovery design](.zazz/architecture/evidence-and-claims.md#discovery-on-novel-target-repositories).

## SARIF output

```bash
go run ./cmd/scan --strings --sarif --entry server.ts <target-root> > findings.sarif
go run ./cmd/scan --export-sarif saved-evidence.json > findings.sarif
```

The second command converts saved scanner JSON offline without a judge request. SARIF retains
route identity, scoped status, source evidence, gaps and illustrative remediation templates. Examples
are unverified suggestions, not applied fixes. `--sarif` and `--json` are mutually exclusive.
The same source JSON can be exported as [standalone HTML](.zazz/architecture/html-reporting.md).
See [report design](.zazz/architecture/reporting-and-sarif.md) and
[comparison results](.zazz/evaluation/corpus/README.md).

## Contributing

Use the same Go and C-toolchain prerequisites as user setup. Development uses a shared bare Git
repository with sibling feature worktrees; follow the [worktree setup guide](.zazz/docs/worktree-setup.md)
for a new checkout or conversion. Work in the assigned feature worktree and verify the location,
branch, and existing changes with `pwd` and `git status --short --branch` before editing.
`stage` is the integration branch; `main` is reserved for human-operated releases.

Read [AGENTS.md](AGENTS.md) for agent rules and the
[orientation](.zazz/docs/scan-agent-orientation.md) for documentation routing. Agree on bounded scope
and acceptance criteria before implementation, keeping the contract in `.zazz/specifications/`.
Put transient run logs and handoffs in ignored `.zazz/ephemeral/`. Keep changes focused and preserve
unrelated work already in the checkout.

Run the relevant scoped checks while developing. The Go CI checks are:

```bash
gofmt -l .
go vet ./...
go build ./...
go test -race ./...
```

`gofmt -l .` should print nothing; format only changed Go files. Use `go test ./...` for a quick local
pass. Tests use local fixtures and HTTP stubs, so no decision-model key is required. For documentation
changes, use the repository's [doc-check workflow](.agents/skills/doc-check/SKILL.md).

Submit focused draft PRs into `stage`, with scoped titles and validation evidence following the
[PR process](.zazz/standards/pr-process.md). Human reviewers approve and merge; agents do not merge
or promote releases. The maintained checkout uses the `zazzcode` Git identity and authentication
profile; agents verify them before committing or pushing.

## Project documentation

This README provides the introduction and getting-started path. Detailed design, plans, contracts,
and evaluation evidence live under `.zazz/`:

- [Project overview](.zazz/project.md) — tenets, objectives, and product direction
- [Architecture](.zazz/architecture.md) — analysis boundaries and provider responsibilities
- [Repository graph](.zazz/architecture/repository-graph.md) — inbound routes, internal paths, and outbound services
- [Evidence and claims](.zazz/architecture/evidence-and-claims.md) — validation, protection, and uncertainty
- [Implementation plan](.zazz/project/implementation-plan.md) — proposed delivery sequence
- [Specifications](.zazz/specifications/README.md) — tracked implementation contracts
- [Evaluation results](.zazz/evaluation/bearer-results.md) — pinned comparison evidence and limitations
- [Documentation map](.zazz/docs/scan-agent-orientation.md#documentation-map) — topic owners and further reading
