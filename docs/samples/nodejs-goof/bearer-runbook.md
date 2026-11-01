# Bearer baseline runbook: Node.js Goof

[Target and artifact policy](README.md) ·
[Case inventory](../../../.zazz/evaluation/nodejs-goof-capabilities.md) ·
[Earlier comparator protocol](../../../.zazz/evaluation/bearer-comparison.md)

Run from the `xss-response` worktree. This procedure statically scans a temporary, filtered snapshot of the
pinned checkout's application source, templates, tests and package manifests. The checkout at
`<DEV_ROOT>/nodejs-goof` remains the only source copy; the repository stores analysis, not Goof source. It requires no npm install, MongoDB,
target server, exploit requests or model key. It does not measure dependency vulnerabilities (SCA),
container vulnerabilities, runtime monitoring or Bearer Pro.

## 1. Resolve the local tool and rules

Set paths to your development directory and an existing Bearer binary/rules checkout. The previous
experiment used Bearer source `e3753de6347988bed15bdba80b9ff731bf028c80` and bearer-rules
`82c335f781a61a6d8b052394f84fd1d93edec467`. Reuse these revisions for continuity;
a different revision starts a separately labeled comparator baseline.

```bash
export GOOF_DEV_ROOT="$HOME/Dev"
export GOOF_TARGET="$GOOF_DEV_ROOT/nodejs-goof"
export GOOF_BEARER_SOURCE="$GOOF_DEV_ROOT/bearer"
export GOOF_BEARER="$PWD/.zazz/ephemeral/bearer"
export GOOF_RULES="$PWD/.zazz/ephemeral/bearer-rules"
export GOOF_RUN_ID="$(date -u +%Y%m%dT%H%M%SZ)"
export GOOF_OUTPUT="$GOOF_DEV_ROOT/zazzcode/zazzles/_test-findings/nodejs-goof/bearer-$GOOF_RUN_ID"
git status --short --branch
git -C "$GOOF_TARGET" status --short
git -C "$GOOF_TARGET" rev-parse HEAD
git -C "$GOOF_BEARER_SOURCE" rev-parse HEAD
git -C "$GOOF_RULES" rev-parse HEAD
test -x "$GOOF_BEARER"
test -d "$GOOF_RULES/rules"
```

Adjust the binary/rules paths if they live in another worktree's ignored scratch. If absent, follow
the [local setup evidence](../../../.zazz/evaluation/bearer-comparison.md#local-setup-evidence)
and Bearer's local `docs/contributing/code.md`: build the pinned source with readonly modules and
obtain the separate pinned rules checkout. Record build toolchain and flags in the run summary.
Do not substitute an unrecorded global CLI or silently fetch current rules.

## 2. Snapshot and run

The following block uses the existing comparison runner's execution validator. It rejects dirty
tracked/untracked target changes and dirty tool/rules checkouts so the recorded revisions identify
the inputs. It saves binary help before running; stop and adapt the procedure if the installed
binary lacks a listed flag. The 600-second scan deadline matches the previous experiment.

Unlike the older JS/TS/JSON-only comparison snapshots, this scope includes `.ejs`, `.hbs`, `.dust`
and `.html`. Omitting templates would discard Goof's principal XSS sink evidence. This remains a
selected source baseline: repository documentation, exploit fixtures, workflows, generated bundles,
existing `sarif.json`, and dependency/build directories are excluded. Templates being present does
not prove Bearer parses them; inspect its logs and findings.

```bash
python3 - <<'PY'
import hashlib
import importlib.util
import json
import os
from pathlib import Path
import shutil
import subprocess

spec = importlib.util.spec_from_file_location(
    "comparison", ".zazz/evaluation/run-bearer-comparison.py")
comparison = importlib.util.module_from_spec(spec)
spec.loader.exec_module(comparison)
target, source, binary, rules, output = (
    Path(os.environ[name]).resolve() for name in (
        "GOOF_TARGET", "GOOF_BEARER_SOURCE", "GOOF_BEARER",
        "GOOF_RULES", "GOOF_OUTPUT"))
for root in (target, source, rules):
    if subprocess.check_output(["git", "-C", str(root), "status", "--porcelain"]):
        raise SystemExit(f"Expected clean checkout: {root}")
if not binary.is_file() or not (rules / "rules").is_dir():
    raise SystemExit("Missing Bearer binary or external rules directory")
output.mkdir(parents=True, exist_ok=False)
scratch = output / "tmp"
scratch.mkdir()
environment = {"PATH": os.defpath, "TMPDIR": str(scratch)}
for args, name in [(["version"], "bearer-version.txt"),
                   (["scan", "--help"], "bearer-scan-help.txt")]:
    with (output / name).open("w") as stream:
        subprocess.run([str(binary), *args], env=environment,
                       stdout=stream, stderr=subprocess.STDOUT, check=True)
snapshot = scratch / "nodejs-goof-source"  # temporary; removed after the scan
snapshot.mkdir()
suffixes = {".js", ".jsx", ".ts", ".tsx", ".json", ".ejs", ".hbs", ".dust", ".html"}
skip_dirs = {"node_modules", "vendor", "dist", "build", "out", "target",
             ".next", ".nuxt", ".output", ".cache", "coverage", ".venv", "venv"}
excluded, hashes = [], {}
paths = subprocess.check_output(["git", "-C", str(target), "ls-files", "-z"])
for name in sorted(filter(None, paths.decode().split("\0"))):
    relative = Path(name)
    original = target / relative
    if (relative.suffix not in suffixes or skip_dirs.intersection(relative.parts)
            or relative.parts[0] in {"exploits", ".github"}
            or name in {"sarif.json", "public/js/bundle.js"}
            or original.is_symlink()):
        excluded.append(name)
        continue
    destination = snapshot / relative
    destination.parent.mkdir(parents=True, exist_ok=True)
    shutil.copyfile(original, destination)
    hashes[name] = hashlib.sha256(destination.read_bytes()).hexdigest()
if not hashes or "views/admin.ejs" not in hashes:
    raise SystemExit("Incomplete source scope: expected application and templates")
def save(name, value):
    (output / name).write_text(json.dumps(value, indent=2) + "\n")
save("nodejs-goof.scope.json", {
    "revision": comparison.revision(target), "files": hashes, "excluded": excluded,
    "included_suffixes": sorted(suffixes), "tests_included": True})
save("metadata.json", {
    "target_revision": comparison.revision(target),
    "bearer_revision": comparison.revision(source),
    "bearer_sha256": hashlib.sha256(binary.read_bytes()).hexdigest(),
    "rules_revision": comparison.revision(rules),
    "zazzles_revision": comparison.revision(Path.cwd()),
    "environment_keys": sorted(environment),
    "rules_mode": "pinned external rules; default rules disabled"})
config, ignore = output / "empty.yml", output / "empty.ignore"
config.write_text("{}\n")
ignore.write_text("{}\n")
command = [str(binary), "scan", "--disable-version-check", "--disable-default-rules",
           "--external-rule-dir", str(rules / "rules"), "--config-file", str(config),
           "--ignore-file", str(ignore), "--disable-domain-resolution",
           "--skip-test=false", "--skip-git-ignore", "--ignore-git", "--force",
           "--parallel", "2", "--hide-progress-bar", "--no-color", "--format", "json",
           "--debug", str(snapshot)]
try:
    comparison.run("bearer-nodejs-goof", command, output, environment)
finally:
    shutil.rmtree(snapshot)
print(f"Review baseline in {output}")
PY
```

The snapshot exists only for the scan and is deleted afterward; `nodejs-goof.scope.json` records its
revision, file hashes and exclusions. Scratch run folders MUST NOT retain or publish target source.
The child process receives only `PATH` and `TMPDIR`, excluding model/provider credentials.
Explicit empty config/ignore files and filesystem discovery prevent local config or ignored scratch
from silently changing the scope. This configuration follows the pinned experiment; it does not
establish zero network activity. That experiment found rule metadata requests even with version
checks disabled. No hosted report is selected.

## 3. Validate and review

Exit 0 or 1 may be valid (1 can mean findings); a Python error or `valid_execution: false` is an
incomplete run. The validator rejects malformed JSON, worker errors and no-files diagnostics.
Also inspect the complete log for skipped/failed files, loaded JS rules and whether templates were
processed. A syntactically valid empty report alone does not prove the source was analyzed.
The earlier experiment needed local worker sockets permitted by the execution environment.
If sockets are blocked, preserve the failed attempt and retry in an authorized environment with a
new run ID. Do not publish the failed or partial attempt as the baseline.

In `SUMMARY.md`, record date/run ID, all revisions, binary hash/version, build toolchain, file count,
actual omissions, findings by rule/severity, exit/duration and validity. Match each
[inventory case](../../../.zazz/evaluation/nodejs-goof-capabilities.md#case-matrix) to its rule IDs
and source/sink locations, or mark it not detected/unresolved. Count unique claims separately from
multiple rules at one sink. Separate tests from application findings; review false positives and
missing evidence. Neither no finding nor agreement between tools proves safety or exploitability.

## 4. Curate the baseline into samples

After reviewing validity and publication contents, copy only the files listed in the
[artifact policy](README.md#saved-results) into a new
`docs/samples/nodejs-goof/bearer/<run-id>/` directory. Keep original commands/absolute paths in
the `_test-findings` run folder; replace local prefixes in public metadata/logs with documented tokens such as
`<TARGET_SNAPSHOT>`, `<BEARER_BINARY>`, `<RULES_ROOT>` and `<RUN_ROOT>`, recording substitutions
in the summary. Preserve relative finding paths, hashes and findings. If any report path requires
normalization, retain the original in the run folder and describe that transformation explicitly.

Link the completed run from the sample README, update its pending status and fill the inventory's
observed Bearer column from the saved evidence. Do not change the earlier two-target Bearer results
or imply that this added target establishes parity. No zazzles/model run is part of this procedure.
