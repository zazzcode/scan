# Vulnerable Flask App (Erlik 2) sample

[Samples](../README.md) · [HTML](findings.html) · [JSON](findings.json) · [SARIF](findings.sarif)

Target: [retained upstream fork](https://github.com/ikpehlivan/Vulnerable-Flask-App), GPL-3.0,
pinned at the original author's last commit `a67cbe5`. The author's repository is deleted; the
handoff verified the retained tree against two other forks. Upstream source excerpts retain their
authorship and licensing. The checkout stayed clean and unchanged.

The report has 19 route/method entries and 24 records, all assessed with `gpt-6-luna`. Four static
vulnerabilities match source review: percent-built SQL and raw HTML flows at `/welcome2`, `/hello`
and the `/get_users` exception return. `/welcome` returns JSON and remains an XSS negative.

The [eight-case inventory](cases.json) passes, including four reviewed positive discovery cases.
[Source review](reference.md), [hashes](source-hashes.json), [metadata](metadata.json),
[Semgrep](semgrep/baseline.json) and [Bearer](bearer/baseline.json) are retained. Semgrep reports 17
results across SQL, HTML responses and several families outside this slice; Bearer reports two,
SQL and pickle. Raw tool result totals are not independent vulnerability counts.

Database-result logging is Review because fields/content are unresolved. Logging a literal internal
mailbox proves the emission but not private/person-identifying sensitivity; its initial model
classification was stronger than the independent evidence, so a generic static blocker keeps Review.
Arbitrary `data` logging is not mislabeled a PII leak. Unknown subprocess return content remains
Review; command injection, traversal, deserialization and the other planned families are out of scope.

```sh
python3 docs/samples/evaluate-python.py vulnerable-flask-app docs/samples/vulnerable-flask-app/findings.json
```

Upstream license: [source license text](SOURCE-LICENSE.txt).
