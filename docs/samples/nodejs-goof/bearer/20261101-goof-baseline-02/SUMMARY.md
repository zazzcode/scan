# Bearer baseline: Node.js Goof

**2026-11-01, run `20261101-goof-baseline-02`.** Target revision:
`add14ba59e98240d9e00a235dd7d42cd61ae9912`. See [metadata](metadata.json),
[source hashes/exclusions](nodejs-goof.scope.json), [execution record](bearer-nodejs-goof.execution.json),
[debug log](bearer-nodejs-goof.log), [version](bearer-version.txt), and
[scan help](bearer-scan-help.txt). The scanner source revision, external rules revision and binary
hash match the earlier pinned comparator experiment. That binary was built with Go 1.27.3
darwin/arm64 and readonly modules, default cgo; its command is retained by the earlier evaluation.

Valid execution, exit 1 with findings, 1.579 seconds. Twenty-two source/template/manifest files
were copied; 38 tracked files were excluded. Templates completed processing without file errors.
This is a static-source baseline. No SCA, target runtime or model calls were made.

[Raw findings](bearer-nodejs-goof.json): **43 records, 12 rules**.

| Severity | Count |
|---|---:|
| Critical | 15 |
| High | 2 |
| Medium | 3 |
| Low | 23 |

| Rule ID | Count |
|---|---:|
| `javascript_lang_logger_leak` | 21 |
| `javascript_lang_hardcoded_secret` | 9 |
| `javascript_express_nosql_injection` | 4 |
| `javascript_lang_insufficiently_random_values` | 2 |
| `javascript_express_hardcoded_secret` | 1 |
| `javascript_express_https_protocol_missing` | 1 |
| `javascript_lang_dangerous_insert_html` | 1 |
| `javascript_lang_format_string_using_user_input` | 1 |
| `javascript_express_default_cookie_config` | 1 |
| `javascript_express_helmet_missing` | 1 |
| `javascript_express_reduce_fingerprint` | 1 |

Five records are test fixtures; 38 are in application/support source. The sole XSS record at
`public/js/ga.js:3` is source-inspected as a false positive for the XSS claim: fixed script creation
with no attacker-selected HTML or destination. The other 42 records remain unadjudicated.
Neither the raw login EJS attribute nor stored Markdown paths received an XSS finding.
See the [comparison](../../../../../.zazz/evaluation/nodejs-goof-results.md) and
[independent answer key](../../xss-answer-key.md) for case matches and limitations.

Local paths are normalized in public outputs: `<BEARER_BINARY>`, `<RULES_ROOT>`, `<RUN_ROOT>`,
`<TARGET_SNAPSHOT>`, `<ZAZZLES_BINARY>` and `<WORKTREE>`. Raw relative finding paths, source excerpts,
hashes and fingerprints are preserved; only absolute-prefix substitutions were applied to Bearer
JSON/log/command files. Original outputs remain in ignored execution scratch. Attempt 01's failed
socket run is excluded. Upstream code excerpts are covered by the
[Apache-2.0 license](../../UPSTREAM-LICENSE.txt).
