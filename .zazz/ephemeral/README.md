# Ephemeral scratch

This directory is temporary working space for the current worktree and branch. Only this README is
tracked; everything else here is gitignored and MUST NOT be committed.

Typical contents: run logs, local handoffs, the Go build cache, previews and one-off helper scripts.
Scan and test results do not go here: they go to `<DEV_ROOT>/zazzcode/zazzles/_test-findings/`
(set `SCAN_FINDINGS_DIR`), outside every worktree. Contents are usually not shared between worktrees, and any of them may be deleted at
any time without breaking the branch.

Nothing important belongs here:

- Curated results, decision traces and answer keys go to tracked `docs/samples/` or `.zazz/evaluation/`.
- Contracts, plans and durable guidance go to their owning tracked documents.
- Target application source stays in its own checkout under `<DEV_ROOT>` (for example
  `<DEV_ROOT>/nodejs-goof`). Do not copy target source here; a scan that needs a filtered snapshot
  creates it temporarily and deletes it afterward.
- Anything another worktree needs belongs in tracked documentation, not here.
