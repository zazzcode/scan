# VulnPy source reference

The [living detection guide](../../../.zazz/evaluation/vulnpy-detection-guide.md) owns trace rationale,
expected classifications, protected controls to add and future-family work. [Source hashes](source-hashes.json)
bind this reviewed reference to the application revision. Labels are evaluation-only inputs.

The SQL answer key is supported by actual format-to-sqlite operations and tests, not filenames:
execute at `src/vulnpy/trigger/sqli.py:37`, executemany at 48 and executescript at 59. Input is captured
by nested callbacks and interpolated before the database call. Executemany binds a separate value;
it does not protect its already-formatted SQL text.

Reflected XSS is in generated Flask/FastAPI views returning raw HTML, not in the empty do_raw trigger.
A separate Jinja environment/template path exists in `apps/flask_app.py` and its template uses safe.
Six current-family source positives remain known scanner misses; four warning controls ensure the
generated routes are not silently presented as covered. The [case evaluator](../evaluate-python.py)
retains those failures. Future framework/dispatch work must not relabel them to make tests appear green.
