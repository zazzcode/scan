# VulnPy Python sample

[Samples](../README.md) · [HTML](findings.html) · [JSON](findings.json) · [SARIF](findings.sarif)

Target: [Contrast-Security-OSS/vulnpy](https://github.com/Contrast-Security-OSS/vulnpy), MIT,
reviewed at `b81885e58283a3449b578d20c20557aff637e5ea`. The target remained clean and unchanged.
Its vulnerable function library is exposed through several frameworks and dynamically built routes.

The offline Zazzles scan records 11 static entries, six unassessed records and four explicit
registration/path coverage warnings. The core generated SQL/XSS routes remain coverage gaps.
This is not a completed semantic calibration or a clean result. No new source was sent to a model.

Six source-positive cases remain known misses. Four controls check explicit coverage warnings.
The evaluator intentionally exits nonzero while those positives are missing:

```sh
python3 docs/samples/evaluate-python.py vulnpy docs/samples/vulnpy/findings.json
```

Pinned local Semgrep CE found eight results and Bearer found 21; neither baseline identifies the
three SQL trigger variants. Broader findings cover future vulnerability families. Counts are not
precision/recall. The first sandbox-blocked Bearer attempt is not used as a baseline.

See the living [detection guide and checklist](../../../.zazz/evaluation/vulnpy-detection-guide.md),
[source reference](reference.md), [expected cases](cases.json), [hashes](source-hashes.json),
[metadata](metadata.json), [Semgrep baseline](semgrep/baseline.json),
[Bearer baseline](bearer/baseline.json), and [upstream license](SOURCE-LICENSE.txt).
