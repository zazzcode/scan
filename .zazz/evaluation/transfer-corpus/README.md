# Helper, mutation and outbound controls

Six additional development cases for the [pinned Bearer comparison](../bearer-results.md).
Provisional source-inspected properties live in [manifest.json](manifest.json), outside scanned case
roots; independent human adjudication is pending. They are not held-out repository accuracy cases.
No target application or database is executed.

| Case | Property isolated |
|---|---|
| case08 | Imported custom allow-list enforced before quoted SQL interpolation |
| case09 | Same-file helper containing interpolated SQL |
| case10 | Imported helper containing interpolated SQL |
| case11 | Validated local replaced by unvalidated alternate request input |
| case12 | Imported helper using a bound SQL parameter |
| case13 | Request URL passed directly to Axios |

Requirements in checks.json express questions, not labels. Case13 has no checks file because the
current checked contract accepts SQL and account creation only; outbound use remains an unchecked
candidate. The runner never loads the expected manifest as tool input. Fresh scan reports are offline
`not_run`; saved original-case judgments do not apply to these new bundles.

Run [the comparison runner](../run-bearer-comparison.py) with `--corpus transfer` to select only
these cases. Raw reports remain in ignored execution scratch; durable observations belong in the
comparison results. Expected SQL protection is contextual to the actual alphabet and quoting in
case08, and to binding semantics in case12; neither establishes general route safety.
