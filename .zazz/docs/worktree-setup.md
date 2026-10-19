# Zazz Worktree Setup

This document defines the required worktree structure for repos using the Zazz methodology, adapted for scan:
the integration branch is **`stage`**, and **`main`** is the release branch (stage merges to main at release time;
a version tag is cut on main).

The methodology is opinionated here on purpose. A consistent worktree model reduces ambiguity, keeps agent work isolated, and gives teams a clean recovery path when a deliverable or document branch goes in the wrong direction.

Worktrees are required by the methodology. This document describes the methodology's opinionated worktree operating model.

Background references:

- [Git worktree documentation](https://git-scm.com/docs/git-worktree)
- [Worktrunk CLI](https://worktrunk.dev/worktrunk/)

`git worktree` is the underlying Git feature. Worktrunk is an encouraged convenience CLI that makes worktree workflows easier, especially when builders and AI agents are working in parallel.

Supporting docs in this repo:

- [wt-cheat-sheet.md](wt-cheat-sheet.md)

## Required Model

Use a bare-repo container with sibling worktrees.

Scan's container layout (resolve the development root using the
[orientation](scan-agent-orientation.md#portable-workspace-and-dev_root)):

```text
scan/                         # container; not a development checkout
├── .bare/                    # shared Git directory
├── stage/                    # integration and current definition worktree
├── feature-<spec>/           # one dedicated worktree per spec, created as needed
└── main/                     # optional release worktree; need not exist now
```

Required conventions:

- the container directory is not itself the active development checkout
- `.bare/` is the shared Git directory
- each active worktree is a sibling directory under the container
- one active deliverable or document effort maps to one worktree in the normal case
- each spec has its own branch and worktree; generic stacked-lane examples do not override this policy
- do not use `/` in branch names
- use flat branch names so the branch name can map directly to a sibling worktree directory
- merges happen through PRs, not by locally merging feature branches into the integration worktree

Every spec is authored, refined, and implemented in its own feature worktree, with its own draft PR
into `stage`. The current definition-only `stage/` exception does not apply to specs. Implementation
requires human spec greenlight. Multiple agents assigned to the same spec coordinate in its worktree;
alternative implementations belong to separately scoped specs/worktrees. The generic stacked-lane
guide does not authorize sharing a worktree across scan specs.

## Why This Pattern

This pattern gives the methodology three important properties:

- isolation: each active effort has its own checkout and branch
- recoverability: a failed line of work can be abandoned cleanly
- consistency: builders and AI agents can use the same branch/worktree mental model

It also pairs well with the methodology's document model:

- durable docs stay in Git
- durable specs stay tracked; transient execution artifacts stay local under the declared document model
- worktrees keep those transient artifacts scoped to one effort at a time

## Integration Worktree

Keep one integration worktree for the repo's declared integration branch — in scan, `stage/`.

Its role is:

- represent the integration target branch
- stay clean and reviewable
- receive changes through PR merge

Do not use the integration worktree as the normal place for day-to-day feature implementation.

Scan's current concept/architecture-definition documentation work is explicitly assigned to `stage/`
by the Project Owner. Assigned methodology refinement may also happen there. This narrow exception is
declared in [AGENTS.md](../../AGENTS.md) and the [orientation](scan-agent-orientation.md); it does not
permit feature implementation, direct pushes, local merges, or releases. Later approved deliverables
use sibling worktrees.

## New Setup Paths

### Option A: Convert an existing clone

If you already have a normal clone and want to convert it into the bare + sibling-worktree layout:

1. Make sure your current checkout is clean or intentionally stashed/committed.
2. Rename `.git` to `.bare`.
3. Mark `.bare` as a bare repository.
4. Create the integration worktree from the repo's declared integration branch.

Example:

```bash
cd <repo-container>
git status

mv .git .bare
git --git-dir=.bare config --bool core.bare true

git --git-dir=.bare worktree add stage -b stage origin/stage
# or, if that branch does not exist yet:
git --git-dir=.bare worktree add stage -b stage origin/main
```

### Option B: Fresh setup

If you are starting from scratch:

```bash
mkdir -p <repo-container>
cd <repo-container>
git clone --bare <repo-url> .bare
git --git-dir=.bare worktree add stage -b stage origin/stage
# or, if that branch does not exist yet:
git --git-dir=.bare worktree add stage -b stage origin/main
```

## Verification Checklist

After setup, verify:

- `.bare/` exists and is the shared Git directory
- the integration worktree exists
- `git --git-dir=.bare rev-parse --is-bare-repository` returns `true`
- `git --git-dir=.bare worktree list` shows the expected worktrees
- the integration worktree is on the expected base branch

Useful commands:

```bash
git --git-dir=.bare rev-parse --is-bare-repository
git --git-dir=.bare worktree list
cd stage && git branch -vv
```

## Reference application checkout

Juice Shop is a separate clone, not a branch or worktree of scan. Its portable lookup convention is
recorded in the [orientation](scan-agent-orientation.md#portable-workspace-and-dev_root). Identify a
fixture by upstream repository URL plus an agreed commit, independently of its local location.
Future fixture configuration must allow a per-machine root override; no flag/environment contract
is implemented now. Do not modify, install, or run the reference app merely for source inspection.
The [acceptance-target plan](../project/implementation-plan.md#acceptance-targets) owns revision and
challenge mapping, inclusion/exclusion, and fixture evaluation scope.

## Branch Naming Rules

Because the methodology expects sibling worktree directories, branch names should also be worktree-safe directory names.

Recommended guidance:

- do not use `/` in branch names
- use flat branch names with hyphens
- prefer names that describe the feature, deliverable, or document effort
- keep the branch name and worktree directory the same whenever practical

If the repo adopts the recommended sibling-worktree model:

- `feature/rbac` is wrong for Zazz worktrees
- `docs/reorg-mw1` is wrong for Zazz worktrees
- `feature-rbac` and `docs-reorg-mw1` are correct

Preferred examples:

- `feature-rbac`
- `deliverable-zazz-142-role-management-ui`
- `proposal-role-management-options`
- `docs-reorg-mw1`

Avoid:

- `feature/rbac`
- `docs/reorg-mw1`
- `deliverable/ZAZZ-142-role-management-ui`

Git allows `/` in branch names, but that creates awkward nested directory paths if the same string is used for a sibling worktree directory. Zazz avoids that mismatch by standardizing on flat branch names when using sibling worktrees.

## Creating Feature, Deliverable, or Proposal Worktrees

From the repo container:

```bash
git --git-dir=.bare fetch origin
git --git-dir=.bare worktree add feature-rbac -b feature-rbac origin/stage
git --git-dir=.bare worktree add proposal-role-management-options -b proposal-role-management-options origin/stage
```

From inside the integration worktree:

```bash
git pull origin stage
git worktree add ../feature-rbac -b feature-rbac
git worktree add ../proposal-role-management-options -b proposal-role-management-options
```

## Deliverables Directory Policy

Scan uses Git-tracked deliverable specs under `.zazz/specifications/`. Author the spec and its
implementation in that spec's feature worktree and review both through its PR into `stage`.
Do not ignore the specification directory or store contracts only in local scratch.

Run logs, QA evidence, recovery notes, and handoffs belong in gitignored `.zazz/ephemeral/`, as declared
in [AGENTS.md](../../AGENTS.md#document-model). Machine-local files may use local exclude rules;
never add secrets or transient execution records to the PR. Load
[execution discipline](agent-execution-discipline.md#execution-records) for record/handoff procedures.

## Daily Operations

### List worktrees

```bash
git --git-dir=.bare worktree list
git --git-dir=.bare worktree list --verbose
```

### Keep a branch current

```bash
cd <repo-container>/feature-rbac
git fetch origin
git rebase origin/stage
```

### Commit and push

```bash
cd <repo-container>/feature-rbac
git add .
git commit -m "type: description"
git push origin feature-rbac
```

### Remove a worktree

```bash
cd <repo-container>
git --git-dir=.bare worktree remove feature-rbac
git --git-dir=.bare worktree remove --force feature-rbac
git --git-dir=.bare worktree prune
```

## Worktrunk Usage

Worktrunk is recommended when the environment has it installed.

Common patterns from the repo container:

```bash
wt -C .bare list
wt -C .bare switch --create feature-rbac
wt -C .bare switch proposal-role-management-options
wt -C .bare switch pr:123
wt -C .bare remove feature-rbac
wt -C .bare step commit
wt -C .bare step push origin/feature-rbac
wt -C .bare step rebase dev
```

Common pattern from inside a worktree:

```bash
wt -C ../.bare list
wt -C ../.bare switch pr:123
wt -C ../.bare remove feature-rbac
```

If the team uses Worktrunk, prefer its wrappers for:

- creating worktrees
- switching between active efforts
- checking out PRs
- routine commit/push/rebase flows

Use plain `git worktree` when:

- Worktrunk is unavailable
- the repo explicitly prefers native Git commands
- debugging a lower-level worktree issue

## Review and PR Flow

When pushing a branch for review:

- target the repo's declared integration branch
- avoid merging feature branches locally into the integration worktree
- prefer PR review as the integration path

If the repo has a fixed integration branch such as `stage`, keep that policy explicit in `AGENTS.md` or repo docs.

## Recovery Model

One of the reasons the methodology requires worktrees is recovery.

If a session of work:

- goes down the wrong path
- fails owner review
- reveals that the proposal, feature requirements document, or deliverable specification is wrong

then the worktree can be abandoned and the team can return to the governing documents, revise the contract, and start a new worktree for the corrected approach.

Useful commands:

```bash
git --git-dir=.bare worktree remove my-failed-branch
git --git-dir=.bare worktree remove --force my-failed-branch
git --git-dir=.bare worktree prune
```

## Operational Rules

1. Do implementation work only inside feature, deliverable, proposal, or docs worktrees.
2. Keep the integration worktree clean.
3. Open PRs from non-integration worktrees into the integration branch.
4. Do not let agents merge PRs.
5. Require a human reviewer to approve and merge.

## Relationship to the Methodology

This document is the operational companion to [zazz-methodology.md](zazz-methodology.md).

The methodology defines:

- why worktrees are required
- where human gates remain
- how worktrees support isolation and recovery

This document defines:

- the required bare-repo + sibling-worktree layout
- setup and conversion paths
- worktree-local exclusion guidance
- daily operations with `git worktree`
- the encouraged Worktrunk command model
