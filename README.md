# scan

`scan` reads a local JavaScript or TypeScript tree and writes a JSON graph of Express-style
route registrations. Each route records the method, the path literal, and the handler names
written at that call. It does not follow those names into Express, into `node_modules`, or
into other files of the application.

- [Project overview](.zazz/project.md) — purpose, current state, product direction
- [Route graph spec](.zazz/specifications/tree-sitter-route-identification-mvp.md) — contract for this slice
- [Agent orientation](.zazz/docs/scan-agent-orientation.md) — workspace conventions, branches, commands, documentation map
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
