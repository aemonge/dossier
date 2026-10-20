# Execution: single conventional commit, path-scoped staging

> Estimate: 5–10m · high
> Human wait: push + publish.

## Skills

- golang-code-style: none (release mechanics only).

## Workflow

> Path-scoped staging keeps the Human's concurrent scaffolding out; explicit
> -m flags satisfy the devbox-git commit policy.

## Steps

1. Stage the session's files by explicit path list.
2. Commit with subject `fix: render custom schemas once, stably, and fast`
   plus a body enumerating the four user-visible fixes.
3. Verify clean status for the staged paths and show the commit.

## Verification

- git status shows no residual modifications for the session paths.
- Commit present on HEAD with both archive dirs included.

## Safety and authority

- No code changes in this change; commit only.

## History

- declared snapshot: 031b727 fix: report version from build info for go install builds
