# Execution: cached-status adoption with equality-gated rebuilds

> Estimate: 1h–1h30m · high
> Human wait: user watches the index during live file edits.

## Skills

- golang-code-style: small map-based cache, no new interfaces, value
  semantics preserved.
- golang-testing: simulation-style regression tests on real temp FS with a
  stubbed status client, asserting rebuild counts and identity stability.

## Workflow

> The enriched change becomes the persistent truth: disk reloads re-apply
> the cached status synchronously, and async status results only rebuild
> when they actually differ — flapping becomes impossible by construction.

## Steps

1. Add `lastStatusByName` to Model; initialize alongside the other
   enrichment maps; populate in `applyStatusEnrichment` on success.
2. In `adoptDiscoveredChange` (and the changed-entry path of
   `adoptDiscoveredProject`), when a cached status exists for the name,
   immediately `EnrichChange(fresh, cached)` before storing — one rebuild,
   schema IDs retained; still clear `enrichedFingerprints` so authoritative
   status is re-fetched async.
3. Gate refreshes on enriched-state equality: compare the change state
   fingerprint before/after enrichment; skip `rebuildIndexPreservingCursor`,
   cache invalidation, and (ModeNormal) `loadViewport` when unchanged.
4. Add cursor fallback in `rebuildIndexPreservingCursor`: on identity miss,
   re-match by (change name, output path), then by (change name, artifact
   row) before clamping.
5. Regression tests: single-rebuild edit cycle with stable identity; no-op
   status result rebuilds nothing; startup swap preserves cursor via output
   path fallback.

## Verification

- `go test -race -cover ./...`, `go vet ./...`, `gofmt -l` clean.
- Manual: `make build && ./dossier` while touching `2-plan.md` externally.

## Safety and authority

- UI-internal refactor; no public API change; read-only flows untouched.

## History

- declared snapshot: working tree with fix-schema-artifact-duplication complete
