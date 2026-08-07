## Why

Dossier currently presents OpenSpec as three flat sections and hardcodes every change to `proposal`, `design`, `specs`, and `tasks`. OpenSpec changes are actually generic work containers whose selected schema defines an artifact dependency graph, so the current model cannot represent bugfix, task, spike, refactor, ADR, event-driven, or other custom workflows and makes canonical project specs appear unrelated to the project hierarchy.

## What Changes

- Replace the flat index with a project hierarchy containing **Active Work**, **Canonical Specs**, and **History**.
- Make active and archived change rows expandable and display their selected schema plus available schema-defined artifacts and output files.
- Keep canonical project specs separate from change-local delta specs, label both explicitly, and allow canonical requirements to expand beneath their spec.
- Resolve active change artifact order, dependencies, output paths, and status from OpenSpec structured status; use filesystem discovery as a degraded fallback when the CLI or schema is unavailable.
- Represent artifact ownership as a tree without pretending the schema dependency DAG is a tree: artifacts remain sibling children of their change, with dependency/status metadata.
- Replace the fixed four-artifact change model and viewer navigation with dynamic artifacts while preserving specialized rendering for tasks, delta specs, Markdown, Git state, and archived read-only behavior.
- Preserve filtering, polling, selection, mouse behavior, explicit-path launch, and collapse state using stable hierarchical identities.
- Define contextual navigation so `Enter`/`Space` toggle expandable rows and `i` inspects the selected change, artifact, output file, canonical spec, or requirement.

## Non-goals

- Installing or defining feature/bugfix/task/spike schema packs.
- Visualizing artifact dependencies as a node-edge graph.
- Integrating Taskflow runs or `parentRunId` lineage into the Dossier tree.
- Editing schema definitions from Dossier.
- Treating canonical project specs as children of a single change.
- Inferring issue type from names when `.openspec.yaml` has no schema.
- Changing OpenSpec archive merge semantics.

## Capabilities

### New Capabilities

- `schema-aware-work-model`: Represent changes, selected schemas, dynamic artifacts, artifact dependencies/status, and multi-file outputs without fixed artifact fields.

### Modified Capabilities

- `openspec-loader`: Load schema metadata and arbitrary artifact outputs, using authoritative OpenSpec status with a filesystem fallback.
- `change-index`: Render a project hierarchy with Active Work, Canonical Specs, History, and stable nested identities.
- `index-specs-section`: Relabel project specifications as Canonical Specs and retain expandable requirements beneath them.
- `collapsible-sections`: Expand/collapse changes, artifact groups, canonical specs, and archived work in addition to section headers.
- `tui-viewer`: Navigate and render dynamic schema artifacts instead of four fixed tabs.
- `archive-viewer`: Render the available dynamic artifacts of archived changes read-only.
- `path-arg`: Load a single active or archived change with its schema-aware artifact model.

## Dependencies

- Depends on the current OpenSpec CLI structured contracts: `openspec list --json`, `openspec schemas --json`, and `openspec status --change <name> --json`.
- `add-index-workflow-actions` SHALL depend on this change so create, inspect, edit, validate, archive/reactivate, and undo operate on stable hierarchy identities.

## Impact

- `internal/openspec/loader.go`: replace fixed artifact fields with schema-aware change/artifact/file structures and degraded discovery.
- New focused OpenSpec status/schema adapter with injectable command execution and JSON fixtures.
- `internal/ui/model.go`, `index.go`, `viewer.go`, `viewport.go`, `view.go`, `mouse.go`, and polling logic: hierarchical item kinds, dynamic artifacts, stable identities, and dynamic viewer navigation.
- Existing task, spec-subnav, editor, archive, Git, keymap, and help behavior must move from fixed tab constants to artifact capabilities or identities.
- Tests and documentation for standard and custom schemas, missing schemas, missing CLI, multi-file outputs, DAG ordering, and archive fallback.
- No new Go module dependency is required.

## Estimate

Epic-sized foundation: 15–25 engineering days, delivered in independently reviewable slices before index lifecycle mutations are implemented.
