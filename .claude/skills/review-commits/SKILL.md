---
name: review-commits
description: Walk a large uncommitted change through review one commit at a time — survey everything, group files into dependency-ordered commits, stage a group, brief it, wait for approval, commit, repeat. Use when the user says "help me review and commit all this", "stage it in groups and brief me", "walk me through the changes before committing", or runs /review-commits.
---

# review-commits

Turn a big working tree into a sequence of reviewed, coherent commits. The user wants to
*understand* the change, not just land it — every group gets a brief they can approve or push back
on before anything is committed.

## Workflow

### 1. Survey everything before grouping anything

```bash
git branch --show-current
git status --porcelain --untracked-files=all
```

Expand untracked directories to files — a `?? backend/internal/foo/` line hides however many files
are inside it. Count them; the user will ask "how many changes" and the answer must be per file.

Then **read every diff and every new file**. Modified: `git diff <path>`. New: read the file. Do
not brief from names — the brief must say what the code does, and the findings below only surface
from reading. If the tree is large, read in batches, but finish before proposing groups.

If the branch is the default branch, stop and ask for a feature branch first.

### 2. Propose the grouping

Order groups so nothing references something that doesn't exist yet. The usual shape:

1. Foundations with no dependencies — bucket/key schemes, error codes, constants
2. Schema — migrations, `.sql` queries, **and the generated sqlc output in the same group**
3. Standalone internal packages, one per group (`internal/imaging`, `internal/openrouter`, ...)
4. The feature core — handlers, service structs, config
5. Endpoints, events, subscribers, tests
6. Docs — `CLAUDE.md`, `.agents/*.md`

Present it as a table: group number, one-line name, file count. Then say two things plainly:

- Groups over ~12 files are hard to review; offer to split. If a split would need hunk-level
  staging of one file across two commits, say so — that risks commits that don't compile.
- The ordering makes dependencies right, but each commit is **not proven to build in isolation**
  unless verified in a temp worktree. Offer that at the end; don't claim it.

Start with group 1 without waiting — the user can redirect.

### 3. For each group: stage → confirm → brief → wait → commit

**Stage** exactly the group's files with `git add <paths>`, then show
`git status --short --untracked-files=no` (plus a grep for the group's untracked files) so the user
sees what is and isn't in.

**Brief** the group. Structure every brief the same way:

1. One sentence on what the group is and where it sits in the feature.
2. Per file or per cluster — what it does, in the code's own terms. Call out the **load-bearing
   decisions**: the one design choice the rest depends on (a key layout, a fail-closed rule, a
   deliberate workaround for a framework quirk). Quote line numbers as `file.go:NN` links.
3. **"Things worth your attention"** — a numbered list, most important first. This is the point of
   the exercise. See the checklist below for what belongs here.

Keep the brief scannable. The user reads it beside the diff in their IDE.

**Wait.** Do not commit until the user says so. If they ask a question, answer it — briefly if
they ask for brevity — and keep waiting. If they ask for a change, see step 4.

**Commit** on approval. Compose the message per the `commit-msg` skill's format rules
(`type(scope): subject`, prose body on *why*, no trailers the repo doesn't use). Use a heredoc.
Report the hash with `git log --oneline -1`, then immediately stage the next group and brief it in
the same turn.

### 4. When the user asks for a change mid-review

Make it. Then:

- Run the tests for the affected packages (`encore test ./path/...`), not the whole suite.
- Re-stage the group. **Check whether any file's net diff against HEAD became empty** — a helper
  added and removed within the branch drops out of the group. Say so and update the file count.
- Re-show the relevant part of the staged diff (`git diff --cached <path>`) so they see the final
  shape before approving.

If the change deletes something the *previous* commit referenced (a package, a function), grep for
every remaining reference before declaring it done; the build is broken until they're all gone.

### 5. Carry findings forward

Keep a running list of findings the user hasn't decided on. Re-mention each one in the group where
it can be fixed, and fold small fixes (a wrong comment, a stale doc line, an off-by-two in a
docstring) into that group's commit rather than a separate "fix comments" commit. Big product
questions stay open — restate them in the final summary rather than resolving them silently.

### 6. Finish

- `encore test ./...` on the final tree.
- `git status --short` must be clean.
- Final summary: the commit list, then **deploy blockers** (secrets set only for Local, migrations
  edited in place, orphaned secrets), then **open findings** the user chose not to act on.

## What belongs in "things worth your attention"

Read for these specifically. They only show up if you look.

- **Comment says one thing, code does another.** "Held for review" when the original is deleted;
  "runs on every path" when three early returns skip it; "ten slots" when the constant says 8.
- **Numbers in docs that don't match measured reality** — cost tables, latencies, limits.
- **Unused declarations** — an error code with zero references, a config value never read.
- **Deploy blockers** — `encore secret list` shows a new secret set for Local only; a migration
  edited in place that has already run somewhere; a Go SDK pinned at a beta version.
- **Breaking API changes** hiding in a field type — a bool that became required, a default that
  flipped from true to "must be stated".
- **Fail-open paths** in code that claims to fail closed, and the asymmetric case: two read-error
  handlers that fail safe in *opposite* directions.
- **Tests that spend money** — anything calling a paid API, even behind an env-var skip.
- **Framework workarounds that look like naive code** — probing keys one at a time because prefix
  listing is broken locally. Say it's deliberate so nobody "fixes" it.
- **Scope the feature quietly doesn't cover** — one plate covered per photo when forecourts have
  many; a "locked layer" rule that protects branded objects mounted on the car.

## Notes

- Ask before every commit. Approval for one group is not approval for the next.
- Never `git add -A` or `git commit -a`. Stage named paths only.
- Do not push.
- When the user pushes back on a design decision that was in an earlier approved plan, don't
  re-argue the plan — restate the concrete trade-off in one sentence and act on their call.
- If the user says "I deleted X", verify it's gone (`ls`) and restage before continuing; don't
  assume.
