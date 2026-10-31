#!/usr/bin/env python3
"""Run pinned local tools without credentials; preserve reports, timing and failures."""

import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import time


def revision(root):
    return subprocess.check_output(
        ["git", "-C", str(root), "rev-parse", "HEAD"], text=True
    ).strip()


def snapshot(root, destination, excluded):
    """Copy tracked JS/TS/JSON inputs; record exact scope and content hashes."""
    skipped = {"node_modules", "vendor", "dist", "build", "out", "target", ".next",
               ".nuxt", ".output", ".cache", "coverage", ".venv", "venv", "__pycache__"}
    paths = subprocess.check_output(["git", "-C", str(root), "ls-files", "-z"]).decode().split("\0")
    manifest = {"revision": revision(root), "excluded": excluded, "files": {}}
    for relative in sorted(filter(None, paths)):
        path = Path(relative)
        source = root / path
        if (relative in excluded or skipped.intersection(path.parts)
                or path.suffix not in (".js", ".jsx", ".ts", ".tsx", ".json")
                or source.is_symlink()):
            continue
        target = destination / path
        target.parent.mkdir(parents=True, exist_ok=True)
        shutil.copyfile(source, target)
        manifest["files"][relative] = hashlib.sha256(target.read_bytes()).hexdigest()
    destination.with_suffix(".scope.json").write_text(json.dumps(manifest, indent=2) + "\n")
    return destination


def run(name, command, output, environment):
    start = time.monotonic()
    timed_out = False
    with (output / (name + ".json")).open("w") as stdout, (
        output / (name + ".log")
    ).open("w") as stderr:
        try:
            result = subprocess.run(
                command, env=environment, stdout=stdout, stderr=stderr, timeout=600
            )
        except subprocess.TimeoutExpired:
            timed_out = True
            result = subprocess.CompletedProcess(command, -1)
    log = (output / (name + ".log")).read_text()
    # Bearer can return 0 after worker failure: exit status alone is insufficient.
    errors = re.findall(r"^.*(?:error processing|error spawning|no files to scan| ERR | FTL ).*$", log, re.M)
    if timed_out:
        errors.append("Execution exceeded 600 seconds")
    try:
        report = json.loads((output / (name + ".json")).read_text())
    except json.JSONDecodeError:
        report = None
        errors.append("Missing or invalid JSON report")
    valid = result.returncode in (0, 1) and not errors
    if name.startswith("scan-"):
        valid = valid and result.returncode == 0 and isinstance(report, dict) and report.get("schema_version") == 1
    else:
        valid = valid and isinstance(report, dict) and all(isinstance(v, list) for v in report.values())
    record = {
        "command": command,
        "exit_code": result.returncode,
        "seconds": round(time.monotonic() - start, 3),
        "valid_execution": valid,
        "errors": errors,
    }
    (output / (name + ".execution.json")).write_text(json.dumps(record, indent=2) + "\n")
    print(name, record["seconds"], "seconds", "valid" if valid else "INVALID", flush=True)
    if not valid:
        raise RuntimeError("Incomplete execution: " + name)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--bearer", required=True, type=Path)
    parser.add_argument("--rules", required=True, type=Path, help="bearer-rules checkout")
    parser.add_argument("--bearer-source", type=Path, help="checkout used to build Bearer")
    parser.add_argument("--scan", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--juice-shop", type=Path)
    parser.add_argument("--second-target", type=Path)
    parser.add_argument("--corpus", choices=("all", "original", "transfer", "none"), default="all")
    parser.add_argument("--snapshot-samples", action="store_true", help="use matching tracked JS/TS/JSON copies")
    parser.add_argument("--sample-exclude", action="append", default=[], help="exact relative path omitted from both copies")
    args = parser.parse_args()
    output = args.output.resolve()
    if output.exists() and any(output.iterdir()):
        raise RuntimeError("Choose a fresh output directory to preserve earlier reports")
    output.mkdir(parents=True, exist_ok=True)
    scratch = output / "tmp"
    scratch.mkdir(exist_ok=True)
    # Allow-list, rather than deleting a few known secret variable names.
    environment = {"PATH": os.defpath, "TMPDIR": str(scratch)}
    empty_config = output / "empty.yml"
    empty_config.write_text("{}\n")
    empty_ignore = output / "empty.ignore"
    empty_ignore.write_text("{}\n")
    metadata = {
        "scan_revision": revision(Path(__file__).resolve().parents[2]),
        "rules_revision": revision(args.rules),
        "bearer_revision": revision(args.bearer_source) if args.bearer_source else "not supplied",
        "bearer_sha256": hashlib.sha256(args.bearer.read_bytes()).hexdigest(),
        "environment_keys": sorted(environment),
        "rules_mode": "pinned external rules; built-in and remote defaults disabled",
    }
    (output / "metadata.json").write_text(json.dumps(metadata, indent=2) + "\n")
    evaluation = Path(__file__).resolve().parent
    targets = []
    if args.corpus in ("all", "original"):
        targets += [(p.name, p, "server.ts") for p in sorted((evaluation / "corpus").glob("case*"))]
    if args.corpus in ("all", "transfer"):
        targets += [(p.name, p, "server.ts") for p in sorted((evaluation / "transfer-corpus").glob("case*"))]
    if args.juice_shop:
        targets.append(("juice-shop", args.juice_shop.resolve(), "server.ts"))
    if args.second_target:
        targets.append(("second-target", args.second_target.resolve(), "src/server.ts"))
    for name, target, entry in targets:
        if args.snapshot_samples and name in ("juice-shop", "second-target"):
            destination = output / (name + "-source")
            if destination.exists():
                raise RuntimeError("Choose a fresh output directory for snapshots")
            target = snapshot(target, destination, args.sample_exclude)
        if name.startswith("case"):
            sources = {str(p.relative_to(target)): hashlib.sha256(p.read_bytes()).hexdigest()
                       for p in sorted(target.iterdir()) if p.is_file()}
            (output / (name + ".source-hashes.json")).write_text(json.dumps(sources, indent=2) + "\n")
        command = [str(args.bearer.resolve()), "scan", "--disable-version-check",
                   "--disable-default-rules", "--external-rule-dir", str(args.rules.resolve() / "rules"),
                   "--config-file", str(empty_config), "--ignore-file", str(empty_ignore),
                   "--disable-domain-resolution", "--skip-test=false", "--skip-git-ignore",
                   "--ignore-git",
                   "--force", "--parallel", "2", "--hide-progress-bar", "--no-color",
                   "--format", "json", "--debug", str(target)]
        run("bearer-" + name, command, output, environment)
        command = [str(args.scan.resolve()), "--strings", "--json", "--entry", entry]
        checks = target / "checks.json"
        if name == "juice-shop":
            checks = evaluation / "juice-shop/validation-checks.json"
        if checks.exists():
            run("scan-discovery-" + name, command + [str(target)], output, environment)
            command += ["--checks", str(checks)]
        run("scan-" + name, command + [str(target)], output, environment)


if __name__ == "__main__":
    main()
