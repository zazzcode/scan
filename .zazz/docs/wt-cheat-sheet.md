# Scan Worktree and Worktrunk Cheat Sheet

Fast reference for working in:

- `~/Dev/zazzcode/scan/.bare`
- `~/Dev/zazzcode/scan/stage` (integration worktree — **stage** is the integration branch)
- sibling feature and PR worktrees under `~/Dev/zazzcode/scan/`

## Layout

```text
scan/
├── .bare/
├── stage/
├── feature-taint-graph/
└── <other-worktrees>/
```

## The main idea

- `.bare` is the shared Git repo backend.
- `stage` is the integration worktree (GitHub `stage`, the repository's default branch). **`main` is
  the release branch** — stage merges to main at release time, and a version tag is cut on main. PRs
  target `stage`.
- each feature branch or PR gets its own sibling worktree directory.
- do not edit the checked-in `.gitignore` for machine-local files.
- use `.bare/info/exclude` for local-only ignore rules.

## Where to run Worktrunk

From the repo container:

```bash
cd ~/Dev/zazzcode/scan
wt -C .bare <command>
```

From inside a worktree:

```bash
wt -C ../.bare <command>
```

## Everyday commands

List worktrees:

```bash
wt -C .bare list
```

Switch to the integration worktree:

```bash
wt -C .bare switch stage
```

Switch back to the default/integration worktree:

```bash
wt -C .bare switch ^
```

Switch to an existing branch worktree:

```bash
wt -C .bare switch feature-taint-graph
```

Create a new branch and worktree from `stage`:

```bash
wt -C .bare switch --create my-new-branch --base stage
```

Remove a finished worktree:

```bash
wt -C .bare remove my-new-branch
```

Force-remove a worktree with local junk left in it:

```bash
wt -C .bare remove my-new-branch -D
```

## Keep `stage` current

Update the integration worktree from GitHub:

```bash
git -C ~/Dev/zazzcode/scan/stage pull origin stage
```

Then create new worktrees from updated `stage`:

```bash
wt -C ~/Dev/zazzcode/scan/.bare switch --create another-branch --base stage
```

Use this flow before starting a new branch so the new worktree is based on the latest `origin/stage`.

## Review a PR

Create or switch to a PR review worktree:

```bash
wt -C .bare switch pr:123
```

That tells Worktrunk to fetch the PR branch and open it as its own worktree. The local copy hook
populates ignored local files there too.

When review is done:

```bash
wt -C .bare remove <pr-branch-name>
```

If you are not sure what branch name Worktrunk used:

```bash
wt -C .bare list
```

## Create a new feature worktree

From the repo container:

```bash
cd ~/Dev/zazzcode/scan
wt -C .bare switch --create my-feature --base stage
```

Then work inside the new sibling directory:

```bash
cd ~/Dev/zazzcode/scan/my-feature
```

When ready to push:

```bash
git push -u origin my-feature
```

## Stacked branches

A stacked branch series is a chain where each branch's PR targets the prior branch instead of
`stage`. Each branch lives in its own sibling worktree. Example chain:

```text
stage → feature-graph-core → feature-taint-graph → feature-llm-broker
```

PR `-2` merges into `-1`, PR `-3` merges into `-2`, and so on. The topmost branch contains the
cumulative content of the whole stack and is where end-to-end testing happens.

Create the next branch in a stack from the current one (not from `stage`):

```bash
wt -C .bare switch --create feature-taint-graph --base feature-graph-core
```

### Verify the tip contains every parent's changes

After any parent in the stack is rebased and force-pushed, the topmost branch needs to be rebased onto
the new parent. To prove the tip is current with every parent — even after rebases rewrite SHAs — use
`git cherry`:

```bash
git cherry -v HEAD origin/feature-graph-core
git cherry -v HEAD origin/feature-taint-graph
```

Empty output means every parent commit is present (by ancestry or by patch-equivalence). Any line
starting with `+` is a real gap that needs investigation.

`git cherry` compares patch-ids (the hash of the diff), so it survives rebases. `git log
<parent>..HEAD` only checks ancestry and gives false positives after a parent rebase.

One-liner to check every parent in a stack:

```bash
for p in graph-core taint-graph; do
  out=$(git cherry HEAD "origin/feature-$p")
  [ -z "$out" ] && echo "$p: contained" || printf -- "%s: MISSING:\n%s\n" "$p" "$out"
done
```

### Rebase onto a rewritten parent

When a parent in the stack gets rebased and force-pushed, fetch with explicit refspecs so the
remote-tracking refs actually update (`git fetch origin` alone only updates `FETCH_HEAD`):

```bash
git fetch origin \
  '+refs/heads/feature-taint-graph:refs/remotes/origin/feature-taint-graph' \
  '+refs/heads/feature-llm-broker:refs/remotes/origin/feature-llm-broker'
```

Then rebase the current branch onto the new parent:

```bash
git rebase origin/feature-taint-graph
```

Patch-equivalent commits (changes already absorbed into the new parent) are skipped automatically.
The branch's own unique commits are replayed on top.

### Force-push with a pinned lease

After a rebase, push with `--force-with-lease` pinned to the verified remote SHA. Plain
`--force-with-lease` can fail with "stale info" if remote-tracking refs are not fresh, and falling
back to plain `--force` discards that safety check.

```bash
git ls-remote origin refs/heads/feature-llm-broker
# copy the SHA, then:
git push --force-with-lease=feature-llm-broker:<expected-remote-sha> \
  origin feature-llm-broker
```

### Inspect divergence

When `git cherry` reports a `+` line and you want to see exactly how a commit differs across two
branches:

```bash
git range-diff origin/feature-taint-graph...HEAD
```

### Stacking with `gh-stack` inside a single worktree

For deliverables that split into 2–3 dependent layers, keep the entire stack in one worktree and use
`gh-stack` to manage branches and PRs. Full reference: [using-gh-stack.md](using-gh-stack.md).

```bash
cd ~/Dev/zazzcode/scan
wt -C .bare switch --create feature-graph-core --base stage
cd ~/Dev/zazzcode/scan/feature-graph-core
gh stack init --base stage feature-graph-core feature-taint-graph
```

All branches share the same working directory. Switch between them with `gh stack` navigation
(`gh stack bottom`, `gh stack top`, `gh stack up` / `gh stack down`). Commit with standard `git add` /
`git commit`; stage deliberately so each branch contains only its layer's changes.

Push and create linked draft PRs:

```bash
gh stack submit --auto --draft
```

Routine sync (fetch, rebase, push, sync PR state):

```bash
gh stack sync
```

## `wt` vs `git worktree`

Use Worktrunk for this repo layout when you want the local setup to come across automatically:

```bash
wt -C .bare switch --create my-feature --base stage
```

Avoid plain `git worktree add` for normal day-to-day branch creation, because the local Worktrunk
hook does not run there.

## Local files in new worktrees

Worktrunk's `copy-ignored` hook copies machine-local files (`.env`, `.config/wt.toml`,
`.zazz/ephemeral/`) into new worktrees. These are gitignored, so they propagate only through this
hook. The scan project config lives at `stage/.config/wt.toml` (gitignored):

```toml
pre-start = "/opt/homebrew/bin/wt step copy-ignored"
post-start = [
  "git -C {{ worktree_path }} config --worktree core.bare false",
  "git -C {{ worktree_path }} config --worktree core.worktree {{ worktree_path }}",
]

[step.copy-ignored]
exclude = ["graphify-out/"]
```

Why the post-start git lines matter: the shared bare repo sets `core.bare=true`, which leaks into
wt-created worktrees and makes ordinary git commands fail there with "fatal: this operation must be
run in a work tree" — including `wt remove` itself. The per-worktree overrides fix that. They also
require `extensions.worktreeConfig=true` in `.bare/config` (already set in this container).

`graphify-out/` is excluded from the copy because it is generated output — rebuild with
`graphify update .` in the new worktree. Never copy installed dependencies between worktrees; when
Go module or tooling caches matter, add a `post-start` rebuild step instead.

### Hook approvals

Hook commands must be approved once per project. Interactively, Worktrunk prompts on first `switch` —
approve and remember. Non-interactively (agent shells/CI), `wt hook approvals add` itself cannot
prompt; either run it in a real TTY or edit `~/.config/worktrunk/approvals.toml` directly:

```toml
[projects."github.com/zazzcode/scan"]
approved-commands = [
    "/opt/homebrew/bin/wt step copy-ignored",
    "git -C {{ worktree_path }} config --worktree core.bare false",
    "git -C {{ worktree_path }} config --worktree core.worktree {{ worktree_path }}",
]
```

## Local-only ignore rules

Repo-wide local excludes:

```bash
~/Dev/zazzcode/scan/.bare/info/exclude
```

Per-worktree local excludes:

```bash
~/Dev/zazzcode/scan/.bare/worktrees/<worktree-name>/info/exclude
```

Check why something is ignored:

```bash
git check-ignore -v <path>
```

## Good habits

- keep `stage` clean — no feature work in the integration worktree
- do feature work in sibling worktrees
- use flat branch names like `feature-taint-graph`, not `feature/taint-graph`
- use `.bare/info/exclude` for local-only notes and machine-specific config
- use `wt -C .bare list` often so you can see what worktrees already exist
- never merge into `stage` or `main` locally — PRs only; releases (stage → main + tag) are
  human-operated
- all work — even expedited fixes — flows through `stage`; `main` moves only at release promotion
