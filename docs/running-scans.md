# Running scans and exporting reports

[Documentation index](README.md) · [Getting started](getting-started.md) · [Interpret results](interpreting-results.md)

All commands below run from the scanner checkout root. The target application is inspected as source;
it is not cloned, installed, or started by zazzles.

## Run

After compiling, configure a results root outside the worktree and run a local scan:

```bash
export SCAN_FINDINGS_DIR="$HOME/Dev/zazzcode/zazzles/_test-findings"
./bin/scan --strings --entry server.ts /path/to/juice-shop
./bin/scan --strings --entry src/server.ts /path/to/vulnerable-typescript-app
```

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

`SCAN_FINDINGS_DIR` is the output root, and `--output-dir <root>` overrides it for one strings scan.
Use `--output-dir ""` to disable bundle output for one run. Paths are resolved relative to the working
directory; use an absolute path for shared storage. The repository folder is named from the target
root's basename, so clones sharing a basename share that folder but receive separate run folders.
No existing run is overwritten. This local directory is a configurable convention, not a hard-coded
scanner default. You may keep the variable in your ignored `.env`; load it in the invoking shell,
because the scanner does not automatically read `.env`. It is a local filesystem path, not a URL.

Open the printed `findings.html` path in your browser (on macOS: `open /path/to/run/findings.html`).
Route rows start with Findings, then Review, then Protected counts (highest first), followed by
route path A–Z and method; blank rows appear last. Click
Route signature to sort by path, or the grouped Findings / Review / Protected header to sort by
those counts together (highest first, then path A–Z; blank rows last). Click again to reverse.
Expanded rows remain open when reordered.
The **Route location** column points to the route declaration. Each expanded finding has its own
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
go run ./cmd/scan --strings --json --entry server.ts \
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
go run ./cmd/scan --strings --sarif --entry server.ts <target-root> > findings.sarif
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
