# Dependency auditing (SCA) after code-path analysis

**Planned direction, 2026-10-31; no SCA adapter is implemented or executed in this phase.** Add
software composition analysis after the route/source-to-use foundation is working and evaluated.
This expands the project from selected OWASP-related source risks to known vulnerable dependencies.
It does not replace the source scanner or require all target languages to be implemented first.

## Adapter boundary

Go owns orchestration, normalization, identity, reporting and correlation. Ecosystem-specific audit
commands are external analysis tools, not additional zazzles-authored analysis runtimes. Prefer the
package manager/lockfile already used by the target. Proposed initial adapters:

- JavaScript/TypeScript: [npm audit](https://docs.npmjs.com/cli/v11/commands/npm-audit/) or
  [pnpm audit](https://pnpm.io/cli/audit), selected by the target's declared toolchain/lockfile.
- Python: [pip-audit](https://github.com/pypa/pip-audit) over supported locked/resolved inputs.
- Other ecosystems: separate adapter contracts after tool/output/coverage evaluation.

An adapter contract pins supported tool/output versions, input formats, network behavior, exit-code
meaning, timeouts, cancellation and parser tests. Preserve raw report hashes and normalized advisory
provenance. No audit execution, package installation, auto-fix, target startup or lockfile mutation is
automatically authorized by the design. Use isolated, bounded read-only audit modes where supported;
resolution/build metadata may execute ecosystem tooling and needs explicit handling in that contract.

## Findings and graph correlation

Record ecosystem/package identity, installed/resolved version, advisory IDs/aliases, affected range,
reported severity/source, fixes and dependency path. Distinguish manifests' declared ranges from
actual resolved versions, direct/transitive dependencies, development/production scope, ignored
advisories and audit failures. A tool failure or unsupported lockfile cannot become a clean result.

Report dependency exposure independently of route reachability. Later correlate imported package
symbols and source call paths to a vulnerable API when advisory/library facts support that mapping.
Use reachable, not proven reachable, and unknown states with their evidence/assumptions. Unknown or
apparently unused code MUST NOT silently dismiss a known vulnerable installed dependency.

The judge may assist with contextual interpretation/remediation prioritization, but authoritative
advisory data and resolved versions establish the dependency match. An LLM opinion does not override
an advisory or invent that a package version is fixed. Source-code and SCA findings share the
[report model and SARIF export](reporting-and-sarif.md), retaining their distinct evidence.

## Milestones and acceptance

- **SCA-1:** pick one ecosystem/tool, pin JSON schema and normalize representative vulnerable,
  clean, ignored, missing-lockfile, unsupported-version, registry-failure and canceled outputs.
- **SCA-2:** bounded audit runner and package/dependency finding report with exact provenance;
  no automatic dependency upgrades.
- **SCA-3:** optional route/call-graph correlation with API-level advisory mappings and unknown
  reachability; retain uncorrelated dependency findings.

This is a named follow-on phase after the current code-path milestone, rather than a requirement to
finish the long-term Python/Go/Java source-analysis roadmap before auditing JavaScript dependencies.
