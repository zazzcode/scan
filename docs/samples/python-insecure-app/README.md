# Python insecure app sample

[Samples](../README.md) · [HTML](findings.html) · [JSON](findings.json) · [SARIF](findings.sarif)

Target: [trottomv/python-insecure-app](https://github.com/trottomv/python-insecure-app), pinned at
`b188493`. Final verification used a fresh clean upstream clone at the same revision, matching
the initial source hashes. Upstream source excerpts retain their authorship.

The report has one FastAPI route and two records, both assessed with `gpt-6-luna`. Request parameter
`name` enters HTML/template source, which is returned under an explicit HTML response contract.
The resulting static XSS finding is source-supported. Business validation remains uncertain rather
than treating a string type as HTML protection. The configured public-IP URL is a negative for
request-controlled destinations; response text is a separate external source.

The [three-case inventory](cases.json) passes (one reviewed positive discovery case).
[Source review](reference.md), [hashes](source-hashes.json), [metadata](metadata.json),
[Semgrep](semgrep/baseline.json) and [Bearer](bearer/baseline.json) are retained. Both comparator
snapshots report zero. The README's SSTI, dependency and secret claims are distinct from this XSS
slice; SCA/secrets/SSTI verdicts are not implemented here. No runtime exploit was executed.

```sh
python3 docs/samples/evaluate-python.py python-insecure-app docs/samples/python-insecure-app/findings.json
```

Upstream license: [source license text](SOURCE-LICENSE.txt).
