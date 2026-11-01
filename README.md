# scan

`scan` analyzes a local JavaScript or TypeScript target repository. `--routes` inventories syntactic
Express-style registrations. `--strings` resolves a bounded subset of handlers and request-value
dependencies into candidate SQL/outbound HTTP uses, retaining source evidence and uncertainty.
Optional GPT Decisions or Jev judging assesses validation and protection from that evidence.

- [Project overview](.zazz/project.md) — purpose, current state, product direction
- [Route graph spec](.zazz/specifications/tree-sitter-route-identification-mvp.md) — contract for this slice
- [Agent orientation](.zazz/docs/scan-agent-orientation.md) — workspace conventions, branches, commands, documentation map
- [OpenAI Decisions setup](.zazz/docs/openai-decisions-setup.md) — API key and Go access example
- [Architecture](.zazz/architecture.md) — technical proposal and analysis boundaries
- [Phased plan](.zazz/project/implementation-plan.md) — future spec slices, dependencies, acceptance targets

## Setup

1. Install Go 1.27.2 or newer. Check with `go version`.
2. Install a C compiler and keep it on `PATH`. The Tree-sitter bindings use cgo. On macOS
   that is the Xcode command line tools (`xcode-select --install`).
3. Clone this repository and enter its root, the directory that contains `go.mod`.
4. The first `go run` or `go test` downloads the modules in `go.mod`. No separate install step.

The program scans a directory you already have on disk. It does not clone, install, or start
the target application.

## Run

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

Check the scanner itself with:

```bash
go test ./...
go test -race ./...
go vet ./...
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
The phase's two known target cases are tracked in `.zazz/evaluation/juice-shop/validation-checks.json`:

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
See [report design](.zazz/architecture/reporting-and-sarif.md) and
[comparison results](.zazz/evaluation/corpus/README.md).
