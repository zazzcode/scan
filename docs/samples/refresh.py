#!/usr/bin/env python3
"""Regenerate the checked-in sample views offline from their canonical JSON."""
import argparse
import json
from pathlib import Path
import subprocess


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--scan', default='./bin/scan', help='compiled scan executable')
    names = ('juice-shop', 'vulnerable-typescript-app', 'nodejs-goof',
             'ThreatByte', 'python-insecure-app', 'vulnerable-flask-app')
    parser.add_argument('--sample', action='append', choices=names, help='refresh only selected samples')
    args = parser.parse_args()
    executable = str(Path(args.scan).resolve())
    for name in args.sample or names:
        folder = Path(__file__).resolve().parent / name
        source = folder / 'findings.json'
        report = json.loads(source.read_text())
        report['root'] = name  # Public sample label; never publish a developer's home path.
        source.write_text(json.dumps(report, indent=2, ensure_ascii=False) + '\n')
        for mode, suffix in (('--export-html', 'html'), ('--export-sarif', 'sarif')):
            rendered = subprocess.run([executable, mode, str(source)], check=True, capture_output=True).stdout
            if suffix == 'html':
                # Source whitespace remains exact in canonical JSON; normalize the published view.
                rendered = b'\n'.join(line.rstrip() for line in rendered.splitlines()).rstrip(b'\n') + b'\n'
            destination = folder / ('findings.' + suffix)
            temporary = destination.with_suffix(destination.suffix + '.tmp')
            temporary.write_bytes(rendered)
            temporary.replace(destination)
        print(f'Refreshed {name}: HTML and SARIF from saved JSON (no model calls)')


if __name__ == '__main__':
    main()
