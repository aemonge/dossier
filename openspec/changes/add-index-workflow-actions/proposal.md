## Why

Dossier opens directly into a change viewer even though the index is now the central place for finding active changes, archived changes, specifications, and requirements. Making the index the primary workflow surface also creates a natural place for common OpenSpec actions without requiring users to leave the TUI.

## What Changes

- Make the index the default startup view, with `[ui] start_view = "index" | "change"` configuration and explicit path launches continuing to open their target directly.
- Add mnemonic, configurable index actions: `n` creates a change, `a` archives or reactivates the selected change, `i` inspects, `e` edits, `v` validates, and `u` undoes the latest lifecycle action.
- Make `Enter` the contextual primary action: it toggles expandable project sections, changes, artifact groups, and canonical specs while inspecting leaf outputs/requirements; `Space` also toggles expandable rows.
- Use the official OpenSpec CLI for create, validate, and archive so schema, validation, and project-spec merge behavior remain authoritative.
- Add confirmation, collision detection, asynchronous progress, actionable errors, and selection-preserving refresh for lifecycle operations.
- Make archive undo restore both the active change and project specifications changed by that archive; make reactivation undo return the change to its original archive path.
- Keep undo session-local and refuse it when affected files changed after the recorded operation.
- Extend read-only mode to block create, archive/reactivate, edit, and undo while retaining inspect and validate.
- Move the default project-information binding to `?` so `i` consistently means inspect in the index.

## Non-goals

- Persisting undo records across Dossier restarts.
- Bulk archive, bulk reactivation, permanent deletion, or arbitrary filesystem operations.
- Reversing an archive through Git reset/restore.
- Editing archived changes; archived artifacts remain read-only.
- Exposing OpenSpec apply, agent instructions, stores, or every CLI option in the primary index workflow.
- Replacing custom keybindings or adding runtime keymap editing.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `change-index`: Add contextual primary/inspect, create, edit, validate, archive/reactivate, and safe session undo actions to the schema-aware hierarchy with dynamic help and read-only behavior.
- `tui-viewer`: Make index-first startup the default and add a strict startup-view setting while preserving explicit path launches.
- `project-config-view`: Reserve `i` for index inspection and use the configurable information action, defaulting to `?`, for project configuration.

## Dependencies

- Depends on `add-schema-aware-index-hierarchy` for dynamic change/artifact identities, schema catalog, hierarchy navigation, and canonical-vs-delta semantics.
- Depends on `add-custom-themes` for strict XDG TOML loading.
- Depends on `add-custom-keybindings` for context-specific action bindings and generated help labels.
- Depends on `add-read-only-mode` for the immutable launch guard applied to new mutations.

## Impact

- `cmd/dossier/main.go`: startup configuration flow and OpenSpec command availability/errors.
- `internal/settings/`: `[ui].start_view` and additional index key actions.
- `internal/openspec/` or a focused CLI adapter: non-interactive OpenSpec command execution and JSON result parsing.
- `internal/ui/model.go`, `internal/ui/index.go`, `internal/ui/update.go`, and `internal/ui/view.go`: prompt state, asynchronous action results, confirmations, undo records, contextual help, and hierarchy-identity selection restoration.
- Tests, example configuration, README keyboard reference, and affected long-lived specifications.
- Runtime integration with the installed `openspec` executable; navigation remains available when it is missing.

## Estimate

Large story: 15–20 engineering days, intended for dependency-ordered delivery slices, including tests, documentation, command integration, transactional undo safeguards, and verification.
