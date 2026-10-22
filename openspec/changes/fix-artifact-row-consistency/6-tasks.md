# Tasks

Step states: `[ ]` Pending, `[-]` Running/interrupted/failed, `[x]` Complete.

## Task 1 — Stable artifact identity across surfaces

> Spec: fix-artifact-row-consistency

- [x] Step 1 — Failing label tests: ID primary, stem dimmed secondary, suppressed when stem equals ID; tabs and filter agree
      start: 2026-10-22T14:15Z · elapsed: 3m · check: `go test ./internal/ui/ -run 'TestRenderArtifactRows'` — red for the right reasons (stem-as-label at column 4)
      paths: internal/ui/label_consistency_test.go
      commit: 75a7086 — fix-artifact-row-consistency: Task 1 Step 1 — label tests [x]
- [x] Step 2 — Refactor to `artifactLabels(primary, secondary)`; update rows, width, filter, tabs; retire stem-pinned test
      start: 2026-10-22T14:25Z · elapsed: 8m · check: `go test ./internal/ui/ ./internal/openspec/` — green (only gate-wrapper git test fails, exec 126, environmental)
      paths: internal/ui/index.go, internal/ui/index_refresh_test.go
      commit: pending — receipt below

## Task 2 — Strays distinct and uncounted

> Spec: fix-artifact-row-consistency

- [ ] Step 1 — Failing stray tests: `[stray]` badge, alpha after the Spec's steps, excluded count, ` · N strays` suffix
      start: · elapsed: · check: `go test ./internal/ui/ -run 'TestRenderStray'`
      paths: internal/ui/label_consistency_test.go
      commit: · fix-artifact-row-consistency: Task 2 Step 1 — stray tests [x]
- [ ] Step 2 — Source-only planning count and stray rendering in `activePlanningText` / badge
      start: · elapsed: · check: `go test ./internal/ui/`
      paths: internal/ui/index.go, internal/openspec/cli.go
      commit: · fix-artifact-row-consistency: Task 2 Step 2 — stray rendering [x]

## Task 3 — Honest degraded mode

> Spec: fix-artifact-row-consistency

- [ ] Step 1 — Failing degraded tests: `[discovered]` badges, `files N · status unavailable`, middle-elided diagnostic
      start: · elapsed: · check: `go test ./internal/ui/ -run 'TestRenderDegraded'`
      paths: internal/ui/label_consistency_test.go
      commit: · fix-artifact-row-consistency: Task 3 Step 1 — degraded tests [x]
- [ ] Step 2 — Discovery stops hardcoding done; degraded header and diagnostic rendering
      start: · elapsed: · check: `go test ./internal/ui/ ./internal/openspec/`
      paths: internal/openspec/loader.go, internal/openspec/cli.go, internal/ui/index.go
      commit: · fix-artifact-row-consistency: Task 3 Step 2 — honest degraded rendering [x]

## Task 4 — Verification

> Spec: fix-artifact-row-consistency

- [ ] Step 1 — Full gate: gofmt, vet, build, lint, full suite with race
      start: · elapsed: · check: `go vet ./... && go test ./...`
      paths: .
      commit: · fix-artifact-row-consistency: Task 4 Step 1 — full verification [x]
- [ ] Human validation — run dossier against `~/galactica` and `~/projects/kernel-rust`, observe the refined Requirements' final claim
      response: · outcome: · at:
