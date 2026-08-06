---
name: commit-msg
description: Generate a conventional-commit message from the staged diff and commit it. Use when the user says "write a commit message", "generate a commit", "commit my changes", or runs /commit-msg.
---

# commit-msg

Write a conventional-commit message for the currently staged changes and create the commit.

## Workflow

### 1. Verify something is staged

```bash
git diff --staged --stat
```

If the output is empty, **stop**. Do not commit, do not stage anything yourself. Tell the user:

> Nothing is staged. Stage the changes you want to commit (`git add ...`) and run this again.

### 2. Read the staged diff

```bash
git diff --staged
```

Read the whole diff. If it is very large, read `git diff --staged --stat` plus the diffs of the most
substantial files — enough to describe accurately what changed and why. Never guess at changes you
haven't read.

### 3. Compose the message

Format:

```
type(scope): short subject

- bullet of what changed
- bullet of why
```

Rules:

- **type** — one of: `feat`, `fix`, `refactor`, `chore`, `docs`, `style`, `test`. Pick the one that
  matches the dominant intent of the diff.
- **scope** — the area touched, derived from the diff paths (e.g. `auth`, `inspection`, `db`,
  `notifications`, `middleware`). Omit the parens entirely if the change spans too much to name one
  scope: `type: short subject`.
- **subject** — imperative mood, lowercase, no trailing period, **under 60 characters**.
- **body** — bullets are optional but encouraged. Cover *what* changed and *why*. Skip bullets that
  only restate the subject. Wrap body lines at ~72 characters.
- **Never** include a `Co-Authored-By` trailer, a "Generated with Claude Code" line, or any other
  trailer the user didn't ask for.

### 4. Commit

Use a heredoc so the multi-line body is preserved:

```bash
git commit -F - <<'EOF'
type(scope): short subject

- bullet of what changed
- bullet of why
EOF
```

Then report the result: show the short subject and the resulting commit hash
(`git log -1 --oneline`). If the commit fails (hooks, empty commit, etc.), report the error output
verbatim rather than retrying with a different message.

## Notes

- Commit only what is already staged — never run `git add`, `git commit -a`, or `--amend` unless the
  user explicitly asks.
- Do not push.
- Describe what the diff actually does, not what the branch name or ticket implies.

## Examples

```
feat(inspection): add inspection center CRUD endpoints

- new handlers/inspection package with list/create/update logic
- sqlc queries + generated code for inspection_centers
- keeps top-level get_car package as a thin route index
```

```
fix(auth): reject refresh tokens with a revoked jti

- lookup now checks revoked_at before issuing a new access token
- previously a logged-out session could still refresh
```

```
chore(db): regenerate sqlc output
```
