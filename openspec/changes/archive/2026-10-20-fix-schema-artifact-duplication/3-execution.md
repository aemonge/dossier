# Execution: claim-then-filter merge in EnrichChange

> Estimate: 30–45m · high
> Human wait: user validates the live dossier index on a custom-schema change.

## Skills

- golang-code-style: table-free focused tests, value semantics, no new deps.
- golang-testing: `t.TempDir()` real-FS regression tests mirroring
  `TestEnrichChangeFromStatus` conventions.

## Workflow

> Single surgical edit in the merge loop keeps risk at the loader boundary —
> the UI consumes `Change.Artifacts` unchanged.

## Steps

1. In `EnrichChange` (internal/openspec/cli.go), while building status
   artifacts, record every accepted output `RelativePath` in a `claimed` set
   and keep an index of status artifacts by ID.
2. Replace the blind fallback append: for each fallback artifact whose ID is
   seen, append its unclaimed outputs to the matching status artifact; for
   unseen IDs, append the artifact with only unclaimed outputs; drop
   artifacts left with none.
3. Add regression tests: numbered-output schema (no phantoms, no duplicated
   paths), stray file survival, and same-ID donation.

## Verification

- `go test -race -cover ./...`, `go vet ./...`, `golangci-lint run ./...`.
- New loader test asserts the enriched artifact list equals the schema's IDs
  with each file present exactly once.

## Safety and authority

- No permission beyond editing internal/openspec and its tests; no public API
  shape change (struct fields untouched; only merge behaviour).

## History

- declared snapshot: 031b727 fix: report version from build info for go install builds
