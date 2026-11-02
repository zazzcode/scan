#!/usr/bin/env python3
"""Compare saved Python evidence with the separate source-reviewed case inventory."""
import argparse
import json
from pathlib import Path


def matches(case, uses):
    return [u for u in uses if u['kind'] == case['kind']
            and u['sink']['span']['file'] == case['file']
            and u['sink']['span']['start_line'] == case['line']]


def evaluate(case, uses, gaps=()):
    if case["expect"] == "coverage_gap":
        return any(g["kind"] == case["kind"] and g.get("at", {}).get("file") == case["file"] and g.get("at", {}).get("start_line") == case["line"] for g in gaps)
    found = matches(case, uses)
    expected = case['expect']
    if expected == 'absent':
        return not found
    if not found:
        return False
    if expected == 'candidate':
        return True
    if expected == 'raw_xss':
        return any(u.get('xss', {}).get('protection') == 'raw' for u in found)
    if expected == 'construction':
        return any(any(g['kind'] == 'query_text_interpolated' for g in u['gaps']) for u in found)
    if expected == 'bound':
        return all(not any(g['kind'] == 'query_text_interpolated' for g in u['gaps']) for u in found)
    if expected == 'not_confirmed':
        return all(u['assessment'].get('finding_status') != 'confirmed_vulnerability' for u in found)
    if expected == 'encoded':
        return all(u.get('xss', {}).get('protection') == 'encoded' for u in found)
    if expected == 'validation_defect':
        return any(u.get('completeness', {}).get('claim_kind') == 'validation' and u['assessment'].get('finding_status') == 'poor_practice' for u in found)
    if expected == 'validation_enforced':
        return all(u.get('completeness', {}).get('claim_kind') == 'validation' and u['assessment'].get('finding_status') == 'protected_for_claim' for u in found)
    if expected == 'form_policy':
        return any(any(e['kind'] == 'form_field_requirement' and 'email' in e['text']
                       for e in u['context']) for u in found)
    raise ValueError(f'Unknown expectation: {expected}')


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('app', choices=('ThreatByte', 'python-insecure-app', 'vulnerable-flask-app', 'vulnpy'))
    parser.add_argument('report', type=Path)
    args = parser.parse_args()
    cases = json.loads((Path(__file__).parent / args.app / 'cases.json').read_text())
    report = json.loads(args.report.read_text())
    results = [{'id': case['id'], 'pass': evaluate(case, report['uses'], report.get('gaps', []))} for case in cases['cases']]
    positive = [case for case in cases['cases'] if case['expect'] in ('construction', 'raw_xss', 'validation_defect')]
    print(json.dumps({'app': args.app, 'cases': results,
                      'passed': sum(row['pass'] for row in results), 'total': len(results),
                      'discovered_positive_cases': sum(evaluate(case, report['uses'], report.get('gaps', [])) for case in positive),
                      'reviewed_positive_cases': len(positive)}, indent=2))
    raise SystemExit(0 if all(row['pass'] for row in results) else 1)


if __name__ == '__main__':
    main()
