## Context

Dossier already has a value-driven Bubble Tea model, a live-reloading index, strict XDG TOML settings, context-specific configurable keymaps, read-only mutation guards, and asynchronous editor/Glamour integration. It currently mutates tasks and Git state but has no OpenSpec CLI adapter or change lifecycle transaction model.

OpenSpec archive is semantically larger than a directory move: it validates the change, merges delta specifications into `openspec/specs/`, and then moves the change to a date-prefixed archive directory. OpenSpec 1.7 provides structured non-interactive create, validate, and archive commands, but no corresponding reactivation or unarchive command. Reactivation is therefore a Dossier filesystem operation, while undo archive must account for project-spec changes.

This change depends first on `add-schema-aware-index-hierarchy`, then on the active custom-theme, custom-keybinding, and read-only changes described in `proposal.md`. Actions operate on stable hierarchy identities and dynamic artifact outputs rather than fixed item indices or artifact fields.

## Goals / Non-Goals

**Goals:**
- Keep index interaction responsive while external OpenSpec operations run.
- Use OpenSpec as the authority for create, validate, and archive semantics.
- Make every lifecycle mutation confirmable, testable, collision-safe, and recoverable.
- Provide exact same-session undo without reverting unrelated user changes.
- Preserve logical index selection across filesystem-driven section changes.

**Non-Goals:**
- Provide crash-persistent transactions or recovery journals.
- Replace OpenSpec's validation, schema selection, spec merge, or archive naming logic.
- Turn the index into a general file manager or expose every OpenSpec subcommand.
- Make archived artifacts writable.

## Decisions

### D1: Startup view is strict UI configuration with explicit-path precedence

Add `UIConfig` to the existing settings model with `StartView string` mapped from `[ui].start_view`. Merge it with the default `index` value and validate the final value against `index` and `change`. Pass the resolved startup view into root-project model construction. `NewSinglePath` bypasses this choice and continues to open its selected target directly.

Alternative considered: a startup-only CLI flag. Rejected because the requested behavior is a durable preference and the existing XDG settings system already provides strict validation. A future CLI override can be added independently if needed.

### D2: Index actions use an explicit modal state machine

Extend the index state with one action state covering idle, new-name input, archive confirmation, reactivation confirmation, undo confirmation, and operation pending. Key dispatch handles this state before ordinary index navigation, just as filter input is handled today. Pending operations carry a stable logical item identity rather than a numeric cursor.

`a` is the default contextual lifecycle action, `i` is dedicated inspect, `n` creates, `e` edits, `v` validates, and `u` undoes. Preserve separate primary (`Enter`) and toggle (`Space`) keymap actions from the hierarchy foundation: primary toggles expandable project sections, changes, multi-output artifacts, and canonical specs while inspecting leaf outputs/requirements; toggle only affects expandable rows. This contextual dispatch avoids assigning `Enter` to two actions in the same keymap context, which strict collision validation correctly rejects. Selected-item mouse activation reuses the same primary dispatch.

Alternative considered: immediately mutate on `a`. Rejected because archive updates project specifications and reactivation changes lifecycle state.

### D3: A focused OpenSpec CLI adapter owns subprocess contracts

Extend the schema/status adapter from `add-schema-aware-index-hierarchy` for:
- `openspec new change <name> --schema <schema> --json`
- `openspec validate <name> --type change|spec --json --no-interactive`
- `openspec archive <name> --yes --json`

The adapter executes the binary directly without a shell, sets the project root as its working directory, captures stdout/stderr, parses structured success output, and returns typed errors. It is invoked from asynchronous `tea.Cmd` functions and sends typed result messages back through `Update`.

The adapter is not consulted during startup. A missing executable therefore affects only CLI-backed actions and never blocks navigation. Tests inject a fake runner rather than requiring OpenSpec to be installed.

Alternative considered: reproduce create/archive behavior directly in Go. Rejected because that would duplicate evolving OpenSpec schema, validation, merge, and naming rules.

### D4: Create selects a workflow schema, then prompts for name

Use the schema catalog supplied by `add-schema-aware-index-hierarchy`. If multiple schemas are available, `n` first opens a searchable schema chooser showing name, description, source, and artifact IDs; if one is available, it is preselected. Then prompt for the change name. Invoke `openspec new change <name> --schema <schema> --json` explicitly so feature/bugfix/task/spike/refactor semantics are never inferred from naming. Dossier does not duplicate kebab-case or schema validation beyond preventing empty submission and displays errors returned by OpenSpec. On success, reload and select the new schema-tagged change.

Optional OpenSpec fields such as goal, description, and store are deliberately omitted from the first workflow. They can later be exposed through an action menu without changing creation correctness.

### D5: Reactivation is a guarded filesystem move

Reactivation derives the active target from the archived change's clean name and moves the exact archived directory to `openspec/changes/<clean-name>`. It refuses an existing target, rejects paths outside the current project's change roots, and retains the original archive path in the undo record. It does not alter project specifications because those merges represent already-published history.

A small filesystem lifecycle service performs and tests these operations independently of the Bubble Tea handler. Move behavior must account for platforms where rename cannot complete by using a verified copy-and-remove fallback consistent with OpenSpec's archive move behavior.

Alternative considered: invoke an OpenSpec unarchive command. Rejected because no such command exists in the supported CLI.

### D6: Archive uses a scoped transaction snapshot

Before archive, build the affected project-spec target set from the selected change's delta spec directories. Capture:
- the complete selected active change directory in a session temp backup;
- bytes and existence state for each affected `openspec/specs/<name>/spec.md` target;
- fingerprints of relevant pre-operation paths.

After OpenSpec succeeds, locate the archive destination from structured output, capture post-operation fingerprints, and retain the backup as the latest undo record. If the command fails, compare disk state with the captured pre-state and restore changed affected paths. Report both the command and restoration errors if recovery is incomplete.

The snapshot is scoped to paths OpenSpec can affect for the selected change rather than the entire repository. This avoids reverting unrelated work while remaining independent of Git state.

Alternative considered: `git restore`. Rejected because it could discard unrelated staged or unstaged changes and does not work outside Git repositories.

### D7: Undo validates post-state before restoring pre-state

An undo record contains the lifecycle kind, exact source/destination paths, affected project-spec preimages, post-operation fingerprints, and temp backup location when required. Before any undo write, validate all post-state fingerprints and destination collisions. Any mismatch aborts before mutation and retains the record.

Undo archive restores affected spec preimages and the active change directory from the retained backup, then removes the unchanged archived result. Undo reactivation moves the unchanged active directory back to its original archive path. Successful undo clears and deletes the record; a later successful lifecycle action first discards the older record. Process exit cleans temp backups.

The preflight is all-or-nothing for expected conflicts. Unexpected I/O failure during restoration is still possible and is reported with exact affected paths, as required by the transactional failure contract.

### D8: Editing reuses existing editor execution and reload behavior

Resolve the selected hierarchy identity to a canonical file:
- active change: the first inspectable dynamic artifact output in schema order;
- active artifact: its only output or currently selected/first output;
- active artifact-output leaf: that exact safe resolved path;
- canonical project spec or requirement: its project-level `spec.md`;
- archived change/artifact/output or section: no editable path.

Reuse the existing `$EDITOR`/`vi`, `tea.ExecProcess`, and editor-return reload path, extending the return message with enough identity to refresh and reselect the index item.

### D9: Refresh and help are based on logical action availability

Use the stable typed hierarchy identities supplied by the foundation, never a stale slice index. After any action result, reload active work, archives, canonical specs, and affected artifact status; rebuild filtered nodes; and locate the logical result identity. Archive/reactivation transfers selection between active and history identities for the same clean change name.

Help rendering asks whether each action is applicable to the selected row capability and session state. It shows `Enter/Space: toggle` for expandable rows and `i/Enter: inspect` for inspectable leaves, while dynamic active artifact outputs receive edit where safe. It omits mutation actions in read-only mode, archive-only restrictions, unavailable undo, and actions hidden while a modal prompt is active. Project information remains configurable and defaults to `?`, leaving `i` for dedicated inspect.

### D10: Action results use existing in-TUI error/status presentation

Typed result messages carry concise success text and structured errors. The index shows pending, success, validation, missing-CLI, collision, and recovery-failure states without switching modes. Status messages clear on the existing timed pattern where appropriate; recovery failures remain visible until dismissed or replaced.

## Risks / Trade-offs

- [Risk: OpenSpec CLI JSON contracts change across versions] → Keep parsing in one adapter, accept only fields needed for identity/results, preserve stderr for diagnostics, and test representative fixtures.
- [Risk: Archive changes an unexpected project path] → Derive expected targets from delta specs, compare disk after failure, refresh from disk unconditionally, and report unexpected state rather than using broad repository rollback.
- [Risk: External edits race with undo] → Fingerprint every affected post-state path and refuse undo before mutation on any mismatch.
- [Risk: Session temp backup consumes disk] → Keep only one lifecycle record, delete replaced/successful records, and clean the current record on normal exit.
- [Risk: Process termination leaves a temp directory] → Use the operating-system temp area; stale directories contain copies only and never become the source of truth.
- [Risk: Large changes make snapshots slower] → Copy only the selected change and affected target specs, and perform preparation/external work asynchronously.
- [Risk: Active OpenSpec changes overlap keymap/read-only/settings contracts] → Implement only after the declared dependency changes are archived or otherwise establish their resulting baseline.

## Migration Plan

1. Implement and archive `add-schema-aware-index-hierarchy`, then archive or establish the resulting baseline of `add-custom-themes`, `add-custom-keybindings`, and `add-read-only-mode`.
2. Add the startup setting with default `index`; users wanting legacy startup can set `[ui] start_view = "change"` before upgrading.
3. Add index actions and dynamic help without removing navigation, filtering, sorting, or mouse behavior.
4. Add CLI-backed actions and lifecycle transactions behind typed adapters and tests.
5. Update examples and README with the new default, configuration, prerequisites, and key reference.

Rollback is code-only: removing the feature restores legacy direct-change startup and index navigation. Lifecycle operations already completed before rollback remain ordinary valid OpenSpec filesystem state; no data migration is required.
