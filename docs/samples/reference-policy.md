# Static baselines and remediation copies

[Samples](README.md)

Keep the pinned Goof, Juice Shop and vulnerable TypeScript checkouts unchanged. Their source-reviewed
answer keys and file-hash manifest provide reusable expectations for future offline scans. Detector
results and model agreement are comparison observations, not the authority defining these labels.
Review the code once per changed case/revision; do not spend model tokens repeatedly classifying
unchanged baseline source merely to regenerate reports.

Before using a reference, verify its Git revision, clean state and recorded hashes. A changed file
invalidates the affected label until its source/flow/context is reviewed again. Expected labels
MUST remain outside scanned source and model evidence bundles and MUST NOT become target-name
exceptions in detector code.

## Remediate a copy

Create an independent clone/copy at a separate remediation path, then pin it to the baseline revision:

```bash
git clone --no-hardlinks "$BASELINE_TARGET" "$REMEDIATION_COPY"
git -C "$REMEDIATION_COPY" checkout --detach "$BASELINE_REVISION"
```

Set those variables to explicit local paths and the reference revision before running the commands.
This procedure has not been executed as part of the current baseline work. Fix only the copy; record
its baseline lineage, patch, new source hashes and scan results. Compare the same case IDs before
and after remediation, preserving unrelated behavior and protected controls. A vulnerable baseline
does not become protected because a different copy was fixed.

Re-exporting JSON/HTML/SARIF and comparing source-based expected labels are offline operations.
Saved model decisions may be replayed only when the complete request—including evidence, questions,
model and rules—matches after insignificant JSON formatting is removed. Changed source/evidence
cannot inherit old judgments. Fresh paid calibration requires explicit authorization, a small case
cap, the active provider/model restriction and retained request/response evidence.

## Protection evidence

A Protected record MUST identify the exact output, a verified renderer/sink binding, its browser
context, the encoder/sanitizer contract and its ordering on every modeled path. Unresolved render
options, trusted-wrapper return types, sanitizer configuration, path merging or unsupported source
bindings remain Review. Escaped syntax alone does not prove that an entire route/application is safe.
An intentionally vulnerable application can contain protected individual outputs; retain them as
negative controls rather than forcing a desired count of positive findings.

For logging, identify the sensitive data class and its content transfer to the log before claiming a
leak. Whole-container logging alone is an observation for Review. Record missing evidence explicitly.
