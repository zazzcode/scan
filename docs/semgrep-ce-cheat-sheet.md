# Semgrep CE cheat sheet (local comparison runs)

Semgrep Community Edition (CE) is the second local comparison scanner, alongside
[Bearer](samples/nodejs-goof/bearer-runbook.md). Use it to compare Zazzles results on
the same pinned sample apps. Official docs:
[Semgrep CE quickstart](https://docs.semgrep.dev/getting-started/quickstart-ce) ·
[CLI reference](https://semgrep.dev/docs/cli-reference) ·
[rule registry](https://semgrep.dev/explore).

Locally installed with Homebrew at `/opt/homebrew/bin/semgrep` (1.183.0 when this sheet was written).
Check before each comparison and record the version with the results:

```bash
semgrep --version
```

## Ground rules

- Sample apps are **read-only scan targets** (see the [sample reference policy](samples/reference-policy.md)).
  Do not create `.semgrepignore`, rule files or output inside a target checkout; pass exclusions on
  the command line and write output under `_test-findings`.
- Always pass `--metrics=off`. With the default `auto`, Semgrep sends usage metrics whenever rules come
  from the Semgrep server. `--config auto` needs metrics and is not used for comparisons.
- Registry packs (`p/...`) are downloaded from semgrep.dev at run time and change over time. For a
  reproducible comparison, snapshot the pack to a YAML file first and scan with that file.
- Registry rules are covered by the Semgrep Rules License; the CE engine is LGPL-2.1. Use them to run
  comparisons; do not copy registry rules into Zazzles' source or detection logic.
- Do not use `semgrep login`, `semgrep ci` or `--pro`; those are the AppSec Platform and Pro Engine,
  not CE.

## Layout

```text
<DEV_ROOT>/zazzcode/zazzles/_test-findings/<app>/semgrep-<YYYYMMDD>-<n>/
  rules/            snapshotted rule packs (YAML)
  results.json      --json-output
  results.sarif     --sarif-output
  results.txt       --text-output (human summary)
  run.md            version, command, target revision, rule hashes
```

## Snapshot rules, then scan

```bash
APP=nodejs-goof                       # or juice-shop, vulnerable-typescript-app
TARGET="$HOME/Dev/$APP"
OUT="$HOME/Dev/zazzcode/zazzles/_test-findings/$APP/semgrep-$(date +%Y%m%d)-1"
mkdir -p "$OUT/rules"
for pack in javascript typescript xss owasp-top-ten; do
  curl -fsSL "https://semgrep.dev/c/p/$pack" -o "$OUT/rules/$pack.yaml"
done
shasum -a 256 "$OUT"/rules/*.yaml > "$OUT/rules/SHA256SUMS"
```

```bash
semgrep scan --metrics=off \
  --config "$OUT/rules" \
  --exclude node_modules --exclude 'exploits' --exclude '*.min.js' \
  --json-output "$OUT/results.json" \
  --sarif-output "$OUT/results.sarif" \
  --text-output "$OUT/results.txt" \
  --dataflow-traces \
  "$TARGET"
```

Then record the target revision and confirm the target is untouched:

```bash
git -C "$TARGET" rev-parse HEAD
git -C "$TARGET" status --porcelain --ignored
```

The status output must be empty. Use the same exclusions as the Zazzles sample run (for Goof, see
the [Bearer runbook](samples/nodejs-goof/bearer-runbook.md) scope) so file sets match.

## Useful options

| Need | Option |
|---|---|
| One pack directly (not reproducible) | `--config p/xss` |
| Local rule file or directory | `--config path/to/rules.yaml` (repeatable) |
| Only some files | `--include 'routes/**' --include 'views/**'` |
| Skip a noisy rule | `--exclude-rule <rule-id>` |
| Severity filter | `--severity ERROR` (repeatable: `INFO`, `WARNING`, `ERROR`) |
| Show taint paths in output | `--dataflow-traces` |
| Exit 1 when findings exist | `--error` |
| Larger files / slow rules | `--max-target-bytes 0`, `--timeout 30` |
| Per-rule timing | `--time` |
| Validate or test custom rules | `semgrep scan --validate --config rules.yaml`; `semgrep --test rules/` |

## Reading results

```bash
jq '.results | length' "$OUT/results.json"
jq -r '.results[] | [.check_id, .path, .start.line, .extra.severity] | @tsv' "$OUT/results.json"
jq -r '.results[] | select(.check_id | test("xss"; "i")) | "\(.path):\(.start.line) \(.check_id)"' "$OUT/results.json"
jq '.errors' "$OUT/results.json"     # parse errors and timeouts reduce coverage
```

## Comparing with Zazzles

- Compare against each app's independent answer key (for example the
  [Goof XSS key](samples/nodejs-goof/xss-answer-key.md)), case by case. Raw finding counts are
  not accuracy.
- CE taint analysis is intraprocedural and single-file; cross-file and interprocedural tracking require
  Semgrep Pro. Expect stored, cross-tier and helper-mediated flows to be absent or only matched as patterns.
  Record that as a capability difference, not a Zazzles win on its own.
- Keep comparator outputs out of Zazzles' detection inputs and model bundles.
- Save curated comparison summaries in `.zazz/evaluation/`, next to the
  [Bearer results](../.zazz/evaluation/bearer-results.md); keep raw outputs in `_test-findings`.
