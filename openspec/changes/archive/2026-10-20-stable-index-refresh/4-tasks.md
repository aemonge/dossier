# Tasks

Step states: `[ ]` Pending, `[-]` Running/interrupted/failed, `[x]` Complete.

## Task 1 — Index updates in place while files change underneath

> Execution: stable-index-refresh

- [x] Step 1 — Cache last status per change and re-apply it synchronously on
      disk-change adoption start: 2026-10-20T10:05Z · elapsed: 20m · check:
      go test ./internal/ui paths: internal/ui/model.go internal/ui/enrichment.go
- [x] Step 2 — Equality-gate rebuilds and viewport reloads on enrichment
      start: 2026-10-20T10:25Z · elapsed: 15m · check: go test ./internal/ui
      paths: internal/ui/enrichment.go internal/ui/update.go
- [x] Step 3 — Cursor fallback by (change, output path) on identity miss
      start: 2026-10-20T10:40Z · elapsed: 10m · check: go test ./internal/ui
      paths: internal/ui/index.go
- [x] Step 4 — Full suites + fmt/vet start: 2026-10-20T10:50Z · elapsed: 5m · check:
      go test -race ./... && go vet ./... paths: internal/ui
- [-] Step 5 — Human validation — dossier index stays put during live edits
      response: INVALID — first paint still swapped discovered→schema IDs, viewer
      mode flapped every 500ms tick, cursor fell to section header · outcome:
      FAILED · at: 2026-10-20T11:05Z
- [x] Step 6 — Fix viewer-mode flap: mergeReloadedChange compares discovery
      fingerprints and adopts as single writer (no enriched→discovered stomp)
      start: 2026-10-20T11:10Z · elapsed: 20m · check: go test ./internal/ui
      paths: internal/ui/model.go internal/ui/enrichment.go
- [x] Step 7 — Bounded synchronous bootstrap enrichment in newModel so first
      paint is schema-aware; degrade to async when CLI slow/absent
      start: 2026-10-20T11:20Z · elapsed: 15m · check: real-CLI smoke test
      paths: internal/ui/model.go
- [x] Step 8 — enterIndex cursor restore falls back to the captured anchor
      when the row identity no longer exists start: 2026-10-20T11:25Z ·
      elapsed: 5m · check: go test ./internal/ui paths: internal/ui/index.go
- [-] Step 9 — Human validation — dossier: no flash on open, none while edits
      happen underneath, cursor survives index↔viewer round-trips
      response: INVALID — flash gone but enriched rows stripped the numbered
      disk names (1-idea…) down to bare schema IDs · outcome: FAILED · at: 2026-10-20T11:55Z
- [x] Step 10 — Label artifact rows and tabs by the concrete output file
      stem (mono-output artifacts) while identities stay schema-based;
      include labels in index filtering start: 2026-10-20T12:00Z · elapsed: 20m ·
      check: go test ./internal/ui paths: internal/ui/index.go internal/ui/view.go
- [x] Step 11 — Human validation — index and tabs show the numbered disk
      names (1-idea, 2-plan, …) with statuses and requires intact, no flash
      response: VALID — “Yeah !!! That's it!” · outcome: VALID · at: 2026-10-20T12:30Z
- [x] Step 12 — Startup optimization: bootstrap uses one
      `openspec status --all --json` invocation (1 spawn instead of N);
      bounded per-change fallback for CLIs without --all. Measured 1.08s →
      0.56s on 8 active changes start: 2026-10-20T12:40Z · elapsed: 25m · check:
      go test ./internal/... + real bench paths: internal/openspec/cli.go internal/ui/model.go
- [x] Step 13 — Human validation — startup feels ~2× faster, no flash, no cache
      files on disk response: VALID · outcome: VALID · at: 2026-10-20T13:05Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/model.go · at 2026-10-20T08:12:23Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/model.go · at 2026-10-20T08:12:29Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/enrichment.go · at 2026-10-20T08:12:37Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/enrichment.go · at 2026-10-20T08:12:51Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/enrichment.go · at 2026-10-20T08:13:07Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/enrichment.go · at 2026-10-20T08:13:18Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/enrichment.go · at 2026-10-20T08:13:24Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/update.go · at 2026-10-20T08:13:28Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/index_actions.go · at 2026-10-20T08:13:36Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/index_actions.go · at 2026-10-20T08:14:08Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/index.go · at 2026-10-20T08:14:48Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T08:14:59Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/index.go · at 2026-10-20T08:15:22Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/index.go · at 2026-10-20T08:15:37Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T08:15:45Z

  gate-attributed · tool: write · paths: /home/aemonge/projects/dossier/internal/ui/index_refresh_test.go · at 2026-10-20T08:16:47Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/index_refresh_test.go · at 2026-10-20T08:16:53Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T08:16:58Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/index_refresh_test.go · at 2026-10-20T08:17:38Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T08:17:42Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T08:17:49Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T08:18:02Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T08:18:12Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/CHANGELOG.md · at 2026-10-20T08:18:46Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/CHANGELOG.es.md · at 2026-10-20T08:18:58Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/AGENTS.md · at 2026-10-20T08:19:08Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T08:19:17Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/model.go · at 2026-10-20T08:39:33Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/model.go · at 2026-10-20T08:39:45Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/model.go · at 2026-10-20T08:39:51Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/model.go · at 2026-10-20T08:40:01Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/model.go · at 2026-10-20T08:40:12Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/model.go · at 2026-10-20T08:40:22Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/index.go · at 2026-10-20T08:40:26Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/index.go · at 2026-10-20T08:40:31Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/model.go · at 2026-10-20T08:40:36Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T08:40:41Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T08:41:25Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T08:41:42Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T08:41:57Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/enrichment.go · at 2026-10-20T08:42:45Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T08:42:57Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T08:43:56Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T08:44:49Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/CHANGELOG.md · at 2026-10-20T08:45:29Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/CHANGELOG.es.md · at 2026-10-20T08:45:37Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/CHANGELOG.es.md · at 2026-10-20T08:45:51Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/AGENTS.md · at 2026-10-20T08:46:08Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T08:46:14Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/index.go · at 2026-10-20T09:06:35Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/index.go · at 2026-10-20T09:06:44Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/view.go · at 2026-10-20T09:06:55Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/index.go · at 2026-10-20T09:07:05Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T09:07:09Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T09:07:53Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T09:08:13Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/CHANGELOG.md · at 2026-10-20T09:08:59Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/CHANGELOG.md · at 2026-10-20T09:09:08Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/CHANGELOG.es.md · at 2026-10-20T09:09:19Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/AGENTS.md · at 2026-10-20T09:09:28Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T09:09:36Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T09:15:27Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T09:15:38Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T09:16:17Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/model.go · at 2026-10-20T09:17:34Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T09:17:37Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/internal/ui/index_refresh_test.go · at 2026-10-20T09:17:54Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T09:18:01Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T09:18:11Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T09:18:16Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T09:18:29Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/CHANGELOG.md · at 2026-10-20T09:18:48Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/CHANGELOG.es.md · at 2026-10-20T09:18:54Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/AGENTS.md · at 2026-10-20T09:19:01Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T09:19:06Z

  gate-attributed · tool: mem_save · paths: — · at 2026-10-20T11:16:44Z
