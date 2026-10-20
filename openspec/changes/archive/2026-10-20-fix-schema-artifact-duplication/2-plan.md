# Plan: dedupe EnrichChange fallbacks by output path

> Estimate: 30–45m · high
> Evidence: reproduced in a scratch loader test — ramona-plan numbered outputs
> produce 4 duplicated files after `EnrichChange`; `TestEnrichChangeFromStatus`
> shows the intended merge (fallback appended only when its ID is unseen).
> Parent idea: none
> Human wait: run dossier on a repo with a custom-schema change and confirm each
> artifact/file appears once in the index and tab bar.

## Acceptance

- Enriching a change whose schema IDs differ from filename stems (e.g. artifact
  `idea` generating `1-idea.md`) yields exactly one artifact per schema ID and
  one output row per file; no phantom `1-idea`-style artifacts remain.
- Fallback artifacts whose files are only partially claimed keep their
  unclaimed files visible (e.g. stray `notes.md` stays as its own artifact).
- A fallback artifact sharing an ID with a status artifact donates its
  unclaimed outputs to that status artifact instead of vanishing or
  duplicating.
- `go test -race ./...`, `golangci-lint run ./...`, and `go vet ./...` pass;
  existing merge test (`TestEnrichChangeFromStatus`) stays green unmodified.
- Human verification: user runs dossier in a ramona-plan change repo and sees
  no duplicated rows.

## Scope

- `internal/openspec/cli.go` — `EnrichChange` fallback merge: claim output
  relative paths from status artifacts, filter fallback outputs, merge
  same-ID leftovers, skip empty fallbacks.
- `internal/openspec/cli_test.go` — regression tests for the three merge
  behaviours above.

## Exclusions

- Degraded (no CLI) display of schema-declared artifact IDs for discovered
  files — discovery has no schema filename data; separate future work.
- Renaming artifact IDs, changing fingerprints, or touching index rendering.
