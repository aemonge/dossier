# Plan: cache status truth so enriched identities survive disk reloads

> Estimate: 1h–1h30m · high
> Evidence: simulation test (deleted after diagnosis) proved: startup swap
> loses cursor identity once; every disk edit triggers adopt (rebuild #1,
> IDs flap to filename stems) + async re-enrichment (rebuild #2, IDs flap
> back); steady state issues zero status calls — the loop is event-correct
> but state-lossy.
> Parent idea: none
> Human wait: run dossier over a repo while files change (e.g. another agent
> working) and see the index update in place — no flashing, no cursor jumps.

## Acceptance

- Editing an artifact file while dossier is open produces exactly one index
  rebuild; the selected row keeps its schema identity (`artifact:idea`
  stays `artifact:idea`) and the cursor stays on the equivalent row.
- A redundant `openspec status` result (identical enriched state) causes no
  rebuild and no artifact viewport reload.
- First enrichment at startup still swaps discovered IDs for schema IDs
  (unavoidable — no cache yet) but the cursor follows the row by
  (change name, output path) fallback instead of clamping to the top.
- In change-viewer mode, an enrichment that does not change the current
  change's enriched state does not re-render markdown.
- `go test -race ./...`, `go vet ./...` pass; existing tests stay green.

## Scope

- `internal/ui/model.go` — `lastStatusByName map[string]openspec.ChangeStatus`
  cache field.
- `internal/ui/enrichment.go` — cache write on successful enrichment;
  synchronous re-apply of cached status inside `adoptDiscoveredChange` /
  `adoptDiscoveredProject`; skip rebuild when enriched state is unchanged.
- `internal/ui/index.go` — `rebuildIndexPreservingCursor` identity-miss
  fallback by (change, output path) then (change, artifact row).
- `internal/ui/update.go` — `statusEnrichmentMsg` viewport reload only when
  the current change's state changed.
- Regression tests replacing the scratch simulation.

## Exclusions

- Replacing polling with fsnotify/eventing (larger redesign, separate slice).
- Degraded-mode schema-declared IDs without CLI (needs richer schemas data).
- Git tab refresh behavior (already equality-gated).
