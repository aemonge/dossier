# Plan: commit the validated work path-scoped

> Estimate: 5–10m · high
> Evidence: both source changes archived with all tasks [x] and VALID
> human validations; working tree holds the code, tests, changelog, docs.
> Parent idea: stable-index-refresh
> Human wait: push + `just publish`.

## Acceptance

- One commit contains exactly: internal/openspec + internal/ui changes,
  index_refresh_test.go, CHANGELOG en/es, AGENTS.md, and the two archived
  change dirs.
- The Human's in-progress scaffolding dirs are not staged.
- Clean tree for the touched paths after commit; HEAD subject uses the
  repo's conventional style.

## Scope

- git stage (explicit paths) + git commit with explicit -m flags.

## Exclusions

- Version bump / tag / publish — the Human runs `just publish`.
