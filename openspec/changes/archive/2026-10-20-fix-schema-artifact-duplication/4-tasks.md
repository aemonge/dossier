# Tasks

Step states: `[ ]` Pending, `[-]` Running/interrupted/failed, `[x]` Complete.

## Task 1 — EnrichChange renders custom-schema artifacts once

> Execution: fix-schema-artifact-duplication

- [x] Step 1 — Claim output paths from status artifacts and filter fallback
      artifacts by unclaimed paths start: 2026-10-20T09:15Z · elapsed: 10m · check:
      go test ./internal/openspec paths: internal/openspec/cli.go
- [x] Step 2 — Regression tests: numbered outputs, stray-file survival,
      same-ID donation start: 2026-10-20T09:25Z · elapsed: 5m · check: go test -race ./...
      paths: internal/openspec/cli_test.go
- [x] Step 3 — Human validation — dossier index on a ramona-plan change shows
      each artifact and file once response: VALID · outcome: VALID · at: 2026-10-20T09:40Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T07:17:20Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T07:18:20Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T07:19:01Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T07:19:20Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T07:19:31Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T07:19:59Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T07:20:13Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T07:20:28Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T07:21:14Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T07:21:23Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T07:21:32Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/CHANGELOG.md · at 2026-10-20T07:21:53Z

  gate-attributed · tool: edit · paths: /home/aemonge/projects/dossier/CHANGELOG.es.md · at 2026-10-20T07:22:11Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T07:22:17Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T07:22:26Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T07:31:20Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T07:31:29Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T07:31:39Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T07:34:09Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T07:34:30Z

  gate-attributed · tool: bash · paths: — · at 2026-10-20T07:36:08Z
