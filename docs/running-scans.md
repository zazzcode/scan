# Running scans and exporting reports

[Documentation index](README.md) · [Getting started](getting-started.md) · [Interpret results](interpreting-results.md)

All commands below run from the scanner checkout root. The target application is inspected as source;
it is not cloned, installed, or started by zazzles.

## Run

After compiling and [configuring required model assessment](decision-models.md), set a results root
outside the worktree and run a model-assessed scan:

```bash
export SCAN_FINDINGS_DIR="$HOME/Dev/zazzcode/zazzles/_test-findings"
./bin/scan --analyze --entry server.ts \
  --judge openai --judge-model gpt-6-luna --judge-limit 10 /path/to/juice-shop
./bin/scan --analyze --entry src/server.ts \
  --judge openai --judge-model gpt-6-luna --judge-limit 10 /path/to/vulnerable-typescript-app
./bin/scan --analyze --entry server/routes.py \
  --judge openai --judge-model gpt-6-luna --judge-limit 10 /path/to/ThreatByte
```

Decision-model assessment is core to security judgment. The explicit request cap bounds a run;
uses beyond it remain unassessed and the report is incomplete. Inventory, evidence-only inspection
and saved-report rendering can run without calls, but do not complete model assessment.

Each command prints a new run folder. The layout keeps repositories and successive runs separate:

```text
_test-findings/
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

`SCAN_FINDINGS_DIR` is the output root, and `--output-dir <root>` overrides it for one analysis scan.
Use `--output-dir ""` to disable bundle output for one run. Paths are resolved relative to the working
directory; use an absolute path for shared storage. The repository folder is named from the target
root's basename, so clones sharing a basename share that folder but receive separate run folders.
No existing run is overwritten. This local directory is a configurable convention, not a hard-coded
scanner default. You may keep the variable in your ignored `.env`; load it in the invoking shell,
because the scanner does not automatically read `.env`. It is a local filesystem path, not a URL.

Open the printed `findings.html` path in your browser (on macOS: `open /path/to/run/findings.html`).
Entry rows start with Findings, then Review, then Protected counts (highest first), followed by
entry label A–Z and method; blank rows appear last. Click
Entry to sort by entry label, or the grouped Findings / Review / Protected header to sort by
those counts together (highest first, then path A–Z; blank rows last). Click again to reverse.
Expanded rows remain open when reordered.
The **Declaration location** column points to the route declaration. Each expanded finding has its own
**Finding location**. Row-end type pills identify review areas such as SQL Injection, NoSQL Injection,
and Input Validation, alongside their poor-practice/review/protected status. Select a pill to open its
glossary side panel with a definition, coverage, and remediation guidance. A type label alone does
not mean an injection is confirmed. The complete glossary is also available below the grid.

Expand a source row to see the HTTP method and path, findings, validation and operation-protection
explanations, code evidence, unresolved analysis, and remediation recommendations. HTML embeds its
styles, content, and a small sorting script, with native disclosure and glossary controls. No server
or external assets are required; without JavaScript the default findings-first view remains usable. It presents zazzles results only. JSON retains the richer evidence for agents; SARIF is
the standard interchange format for security-tool integrations. Generated reports contain source
excerpts, so keep them with the same access controls as that source.

Without an output root, stdout modes still work. For evidence-only inspection and offline rendering:

```bash
./bin/scan --analyze --json --output-dir "" --entry server.ts /path/to/app > findings.json
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

## What the graph contains

`scan --routes` parses supported JavaScript/TypeScript extensions and `.py` files. Python route
extraction covers a bounded Flask/Blueprint, RESTX and FastAPI subset; generated registration and
unresolved dispatch remain explicit coverage work. It skips `node_modules`, `vendor`, `dist`, `build`
and other dependency/output directories listed in `internal/walk`. `--list` also recognizes Go and
Java files; those languages do not yet have route analysis. See [Python limits](python-support.md).

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
go run ./cmd/scan --analyze --json --entry server.ts <target-root>
```

Entry paths are relative to the target repository. Use `server.ts` for Juice Shop and
`src/server.ts` for the vulnerable TypeScript sample. Omit `--entry` to consider all selected JS/TS
registration files. Explicit selection avoids conflating training snippets with production scope.

The separate schema3 evidence report contains typed `entries`, `registrations`,
`dispatch_bindings`, `callables`, `invocation_contexts`, `uses`, `files`, `gaps` and `limitations`.
Each use identifies its route, source fields, ordered local transfer evidence, candidate sink,
snapshot hashes and assessment state. Direct named relative ESM imports, imported handler factories,
local assignment, string expressions and simple object destructuring are supported. Branches,
unknown transformations, middleware, receiver ambiguity, unverified library bindings and external
response continuations remain explicit gaps. A SQL-shaped string logged to the console is not
reported as an SQL call; if it carries a hinted request field, it is a sensitive-data logging record. HTTP client registrations are excluded by the new receiver filter.

Use `--route /rest/user/login` to select an exact path pattern before model requests.
Offline assessment is `not_run`. For provider selection and key instructions, including personal-key
ownership and Jev configuration, see [decision-model setup](decision-models.md).
`--judge openai|jev` explicitly transmits bounded evidence; `--judge-limit` caps requests (default 10,
no retries). A requested judge failure, skipped oversized payload or exhausted request cap exits 1
with its report retained. Neither a successful judge call nor offline exit 0 certifies the target.
The graph design and full CFG/interprocedural analysis exceed this first implementation slice.

## Explicit route validation checks

Use a requirements file to declare the method, path, fields and separate validation/protection
constraints. This avoids guessing business requirements from a helper's name or a missing regex.
Two known target cases are tracked in `.zazz/evaluation/juice-shop/validation-checks.json`:

```bash
go run ./cmd/scan --analyze --json --entry server.ts \
  --checks .zazz/evaluation/juice-shop/validation-checks.json <target-root>
```

Add `--judge openai --judge-model gpt-6-luna --judge-limit 2` for a capped live run. Checked output adds
relevant middleware/helpers/model context, local registration control edges, required evidence roles,
blocking vs coverage gaps, and a scoped finding status. `supported_potential` is an advisory finding
under the declared models, not runtime confirmation or completed independent review. See the
[results guide](interpreting-results.md). SCA remains planned. Automatic claim assembly for
novel repositories remains incomplete; see [discovery design](../.zazz/architecture/evidence-and-claims.md#discovery-on-novel-target-repositories).

## SARIF output

```bash
go run ./cmd/scan --analyze --sarif --entry server.ts <target-root> > findings.sarif
go run ./cmd/scan --export-sarif saved-evidence.json > findings.sarif
```

The second command converts saved scanner JSON offline without a judge request. SARIF retains
route identity, scoped status, source evidence, gaps and illustrative remediation templates. Examples
are unverified suggestions, not applied fixes. `--sarif` and `--json` are mutually exclusive.
The same source JSON can be exported as [standalone HTML](interpreting-results.md).

## XSS and remediation view

Current reports include bounded server response/template, stored-rendering and browser XSS analysis.
Server/API routes and Browser/client-side sources use separate headings; client rows identify
components/source files and output locations rather than HTTP methods/routes. Selecting an API entry
file limits server registrations; browser discovery still inspects supported scripts in the target.
See [the XSS analysis guide](analysis/xss.md) for protection gates and current limits.

Use **Hide protected results** to focus on Findings and Review. It hides protected records and
protected-only source entries; saved JSON/SARIF and total counts remain unchanged. A Protected result
is scoped to its exact output and evidence, never the entire route/application.

## Analysis naming and saved-report compatibility

`--analyze` selects bounded security analysis of route and browser candidate uses. The previous
`--strings` flag remains a compatibility alias. Both perform the same analysis; add an explicitly
configured `--judge` for the required model-assessed security workflow.

Fresh analysis JSON uses `schema_version: 3` and `analysis_guidelines` metadata in place of
`practice`. Route inventory JSON remains schema 1. HTML/SARIF export accepts saved analysis reports
from schemas1/2/3, including historical guideline metadata. Current candidate and finding
associations use `entry_id` and `invocation_context_id`; old route/client associations are
converted only when reading saved reports. Original schema provenance is retained in
`source_schema_version` for converted inputs. Saved sample reports retain their
original assessments and provenance; changed guideline text or model inputs cannot reuse those
answers as judgments of a new scan. Consumers reading fresh JSON should use `analysis_guidelines`.

## Common entry selection and tagged checks

`--entry` keeps its registration-file scope. `--route` selects HTTP path patterns only.
Use `--entry-kind http|browser|graphql|grpc` and `--entry-id <exact-snapshot-id>` with
`--analyze`; filters intersect and an empty selection exits 1. An invalid kind exits 2.
IDs identify this source snapshot; portable requirements should use typed identities.

```bash
./bin/scan --analyze --entry-kind browser --json --output-dir "" /path/to/app
./bin/scan --analyze --entry server.ts --entry-kind http --route /login --json /path/to/app
```

Checks remain a top-level array of 1–20 records, capped at 64 KiB. Each record uses
legacy `method`/`path` or a tagged `selector`, never both. Unknown fields and incompatible
selector fields are rejected. For example, a portable HTTP check is:

```json
[{"id":"login-sql","selector":{"kind":"http","method":"POST","path":"/login"},"kind":"sql_query","fields":["email"],"validation":"Apply only the application's declared email requirements","protection":"Keep email separate from SQL syntax"}]
```

Selector identities are HTTP `method`/`path`, GraphQL `type`/`field`, gRPC fully qualified
`service`/`method`, or browser `category`/`name` with optional `event`. Optional
`registration_file` narrows any selector to that relative registration file. Selection
never creates an adapter or an operation: GraphQL/gRPC discovery remains pending their
adapter deliverables. Check kinds remain `sql_query`, `nosql_query`, and `account_creation`.
Field names retain their source roles, such as `args.email` or `request.email`.

Every matched invocation context is evaluated separately. JSON `check_matches` records
entry/context IDs and use counts; `selected_scope` records CLI filters. Zero matches and
matched contexts without supported operations remain explicit gaps, never protection.
