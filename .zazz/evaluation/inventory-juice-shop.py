#!/usr/bin/env python3
"""Index all pinned Juice Shop challenges and every Bearer finding for source review.

Requires Python 3 and Ruby's standard-library YAML parser; never runs target code.
Reference matches and snippet overlap are discovery aids, not correctness labels.
"""

import argparse
from collections import Counter
import hashlib
import json
from pathlib import Path
import re
import subprocess


def write_rows(path, rows):
    path.write_text("".join(json.dumps(row, sort_keys=True) + "\n" for row in rows))


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--target", required=True, type=Path)
    parser.add_argument("--bearer-report", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    args = parser.parse_args()
    root = args.target.resolve()
    output = args.output.resolve()
    output.mkdir(parents=True, exist_ok=True)
    revision = subprocess.check_output(["git", "-C", str(root), "rev-parse", "HEAD"], text=True).strip()
    if revision != "f189f521266f42c6bd05884468292000cee852ed":
        raise RuntimeError("Review decisions require the pinned Juice Shop revision")
    subprocess.run(["git", "-C", str(root), "diff", "--quiet", "HEAD", "--"], check=True)
    challenges = json.loads(subprocess.check_output(
        ["ruby", "-ryaml", "-rjson", "-e", "puts JSON.generate(YAML.load_file(ARGV[0]))",
         str(root / "data/static/challenges.yml")], text=True
    ))
    keys = {row["key"] for row in challenges}
    assert len(keys) == len(challenges), "Duplicate challenge key"
    key_pattern = re.compile(r"\b(?:" + "|".join(re.escape(key) for key in sorted(keys)) + r")\b")
    references = {key: [] for key in keys}
    marked_lines = {}
    tracked = subprocess.check_output(["git", "-C", str(root), "ls-files", "-z"]).decode().split("\0")
    source_hashes = {}
    for filename in sorted(filter(None, tracked)):
        path = root / filename
        if path.suffix not in (".ts", ".js", ".html", ".yml") or path.is_symlink():
            continue
        source = path.read_text()
        source_hashes[filename] = hashlib.sha256(path.read_bytes()).hexdigest()
        for number, line in enumerate(source.splitlines(), 1):
            for key in set(key_pattern.findall(line)):
                references[key].append({"file": filename, "line": number,
                                        "kind": "test" if filename.startswith("test/") or ".spec." in filename
                                        else "source_reference"})
            if "vuln-code-snippet vuln-line" in line:
                marked_lines[(filename, number)] = sorted(set(key_pattern.findall(line)))
    inventory = []
    challenge_lines = {}
    for number, line in enumerate((root / "data/static/challenges.yml").read_text().splitlines(), 1):
        match = re.match(r"  key: (\w+)", line)
        if match:
            challenge_lines[match[1]] = number
    for row in challenges:
        key = row["key"]
        inventory.append({"key": key, "name": row["name"], "category": row["category"],
                          "description": row["description"], "difficulty": row["difficulty"],
                          "documented_at": {"file": "data/static/challenges.yml", "line": challenge_lines[key]},
                          "disabled_environments": row.get("disabledEnv", []),
                          "references": references[key], "independent_label_review": "pending",
                          "sast_detectability": "not_inferred_from_challenge_existence"})
    decisions = json.loads(Path(__file__).with_name("juice-shop").joinpath("review-decisions.json").read_text())
    reviews = {(d["rule"], d["file"], d["line"]): d for d in decisions}
    seen_reviews = set()
    report = json.loads(args.bearer_report.read_text())
    ledger = []
    for severity, findings in sorted(report.items()):
        for finding in findings:
            filename, line = finding["filename"], finding["line_number"]
            assert (root / filename).is_file(), filename
            identity = (finding["id"], filename, line)
            scope = ("test" if filename.startswith("test/") or ".spec." in filename or ".test." in filename
                     else "training_snippet" if filename.startswith("data/static/codefixes/") else "application_or_support")
            review = reviews.get(identity)
            if review:
                seen_reviews.add(identity)
            ledger.append({"rule": finding["id"], "file": filename, "line": line,
                           "severity": severity, "fingerprint": finding["fingerprint"], "scope": scope,
                           "disposition": review["disposition"] if review else "pending_source_review",
                           "reason": review["reason"] if review else "Not adjudicated; finding presence is not confirmation.",
                           "review_evidence": review["evidence"] if review else [],
                           "documented_vulnerable_line_overlap": marked_lines.get((filename, line), []),
                           "source_sha256": hashlib.sha256((root / filename).read_bytes()).hexdigest()})
    assert seen_reviews == set(reviews), "Review decision did not match a finding"
    write_rows(output / "expected-challenges.jsonl", sorted(inventory, key=lambda x: x["key"]))
    write_rows(output / "bearer-review.jsonl", sorted(ledger, key=lambda x: (x["file"], x["line"], x["rule"])))
    summary = {"target_revision": revision, "challenge_count": len(inventory),
               "bearer_report_sha256": hashlib.sha256(args.bearer_report.read_bytes()).hexdigest(),
               "bearer_findings": len(ledger), "dispositions": dict(Counter(x["disposition"] for x in ledger)),
               "scope_counts": dict(Counter(x["scope"] for x in ledger)),
               "challenge_source_sha256": source_hashes["data/static/challenges.yml"],
               "review_state": "agent source inspection; independent human review and runtime tests pending"}
    (output / "inventory-summary.json").write_text(json.dumps(summary, indent=2) + "\n")
    print(json.dumps(summary, indent=2))


if __name__ == "__main__":
    main()
