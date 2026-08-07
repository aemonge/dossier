## MODIFIED Requirements

### Requirement: Seleccionar un change con Enter
The index SHALL provide a configurable inspect action bound to `i` by default and preserve the hierarchy's contextual primary action bound to `Enter`. Inspect SHALL open the selected inspectable hierarchy identity without mutating or changing expansion. Primary activation SHALL toggle project sections, active/archived changes, multi-output artifacts, and canonical specs, while inspecting artifact-output and requirement leaves. A left-click on an already-selected row SHALL perform the same contextual primary action.

#### Scenario: Inspect active change
- **WHEN** the cursor is on an active work row and the user presses `i`
- **THEN** the schema-aware change viewer opens on its first inspectable artifact

#### Scenario: Primary action toggles active change
- **WHEN** the cursor is on an active work row and the user presses `Enter`
- **THEN** its schema-defined artifacts expand or collapse without leaving the index

#### Scenario: Inspect active artifact output
- **WHEN** the cursor is on an active artifact-output leaf and the user invokes inspect or primary activation
- **THEN** the viewer opens that exact output

#### Scenario: Inspect archived change
- **WHEN** the cursor is on archived work and the user invokes inspect
- **THEN** the dynamic read-only archive viewer opens

#### Scenario: Primary action toggles canonical spec
- **WHEN** the cursor is on a canonical spec and the user presses `Enter`
- **THEN** its requirement rows expand or collapse without leaving the index

#### Scenario: Click selected row uses primary action
- **WHEN** the user left-clicks an already-selected hierarchy row
- **THEN** Dossier toggles it when expandable or inspects it when it is a leaf

### Requirement: Helpbar del índice
The index help bar SHALL render navigation, filtering, sorting, expansion, project information, and contextual item actions from the active keymap and selected hierarchy-row capabilities. Leaf rows SHALL advertise `i/Enter: inspect`; expandable rows SHALL advertise `Enter/Space: toggle` and `i: inspect` when inspectable. Applicable actions from `n: new`, `a: archive` or `a: make active`, `e: edit`, `v: validate`, and `u: undo` SHALL be shown contextually. Actions that do not apply, are unavailable in read-only mode, or have no undo record SHALL be omitted.

#### Scenario: Active work help
- **WHEN** an active work row is selected with the default keymap and read-only mode is inactive
- **THEN** the help bar advertises toggle, inspect, validate, archive, and new-change actions

#### Scenario: Archived change help
- **WHEN** an archived change is selected with the default keymap and read-only mode is inactive
- **THEN** the help bar advertises inspect and make-active actions but does not advertise edit or validate

#### Scenario: Active output help
- **WHEN** an editable active artifact output is selected with the default keymap and read-only mode is inactive
- **THEN** the help bar advertises inspect and edit but no change lifecycle action

#### Scenario: Canonical specification help
- **WHEN** a canonical spec is selected with the default keymap and read-only mode is inactive
- **THEN** the help bar advertises `Enter/Space: toggle`, `i: inspect`, edit, validate, and new-change actions but no lifecycle action

#### Scenario: Section header help
- **WHEN** a section header is selected with the default keymap
- **THEN** the help bar advertises `Enter/Space: toggle` and does not advertise inspect

#### Scenario: Sort hint reflects current order
- **WHEN** the specification sort mode changes between name and suffix
- **THEN** the help bar advertises the action that switches to the other sort mode

#### Scenario: Custom action labels
- **WHEN** an index action is rebound in the active keymap
- **THEN** the index help bar advertises the configured binding instead of its default

## ADDED Requirements

### Requirement: Context-specific index actions
The index SHALL expose separately configurable actions for new change, lifecycle action, inspect, edit, validate, and undo while retaining the hierarchy's contextual primary and toggle actions. Their default bindings SHALL be `n`, `a`, `i`, `e`, `v`, and `u`, with `Enter` as primary and `Space` as toggle. Dispatch SHALL use selected row capabilities rather than hardcoded standard artifact kinds, preserving same-context key collision validation. Rebinding an action SHALL replace its default bindings according to the existing keymap merge and collision rules.

#### Scenario: Default mnemonic actions
- **WHEN** Dossier starts with the default keymap
- **THEN** the index uses `n`, `a`, `i`, `Enter`, `Space`, `e`, `v`, and `u` for new, lifecycle, inspect, primary, toggle, edit, validate, and undo respectively

#### Scenario: Rebound lifecycle action
- **WHEN** the index lifecycle action is configured to a key other than `a`
- **THEN** only the configured key starts archive or reactivation from the index

### Requirement: Create a new change from the index
The index SHALL allow a writable session to create active work by selecting an available OpenSpec workflow schema and entering a change name. If multiple schemas are available, Dossier SHALL present their names, descriptions, sources, and artifact IDs; with one schema it MAY preselect it. Creation SHALL invoke the installed OpenSpec CLI non-interactively with the explicit schema. On success, Dossier SHALL remain in the index, refresh Active Work, clear prompts, and select the new schema-tagged change. Canceling or receiving an invalid schema/name SHALL not create files.

#### Scenario: Create bugfix change
- **WHEN** the user invokes new change, selects `bugfix`, enters `fix-stale-cache`, and accepts
- **THEN** Dossier invokes OpenSpec with schema `bugfix` and the refreshed Active Work selects `fix-stale-cache`

#### Scenario: One schema available
- **WHEN** only `spec-driven` is available
- **THEN** Dossier may preselect it but still creates with that explicit schema

#### Scenario: Cancel creation
- **WHEN** the new-change prompt is active and the user cancels it
- **THEN** the prompt closes and no change is created

#### Scenario: Invalid change name
- **WHEN** the user submits a name rejected by OpenSpec
- **THEN** the index remains open, reports the validation error, and creates no partial change

### Requirement: Edit selected writable artifacts
The index SHALL allow a writable session to launch the configured external editor for a safe selected active artifact output, the first inspectable output of selected active work, or a selected canonical project specification. A selected requirement SHALL edit its containing canonical spec. Archived work, sections, missing outputs, blocked/ready artifacts without output, and unsafe paths SHALL not offer or perform edit. Editor fallback, terminal suspension, and immediate reload behavior SHALL match artifact editing from the dynamic change viewer.

#### Scenario: Edit active work
- **WHEN** an active work row is selected and the user invokes edit
- **THEN** Dossier opens its first safe inspectable artifact output in the external editor

#### Scenario: Edit selected custom artifact
- **WHEN** an active diagnosis output is selected and the user invokes edit
- **THEN** Dossier opens that exact diagnosis file

#### Scenario: Edit project specification
- **WHEN** a project specification or one of its requirements is selected and the user invokes edit
- **THEN** Dossier opens that project's `spec.md` in the external editor

#### Scenario: Archived change remains read-only
- **WHEN** an archived change is selected and the user invokes edit
- **THEN** no editor is launched and the archived change remains unchanged

#### Scenario: Refresh after index edit
- **WHEN** the editor exits after changing an artifact opened from the index
- **THEN** Dossier immediately reloads the affected item and preserves its index selection when it still exists

### Requirement: Validate selected OpenSpec items
The index SHALL allow active changes and project specifications to be validated through the installed OpenSpec CLI without leaving the index. Dossier SHALL identify whether the selected item is a change or specification, use non-interactive structured output, and show a concise success result or actionable validation failures while preserving selection. Archived changes, requirements as independent items, and section headers SHALL not be validated independently; a requirement SHALL validate its containing project specification.

#### Scenario: Validate active change
- **WHEN** an active change is selected and the user invokes validate
- **THEN** Dossier validates it as a change and reports the result in the index

#### Scenario: Validate project specification from requirement
- **WHEN** a requirement is selected and the user invokes validate
- **THEN** Dossier validates the containing project specification and reports the result

#### Scenario: Validation failure
- **WHEN** OpenSpec reports validation errors
- **THEN** Dossier remains in the index and displays actionable failure details without changing project files

### Requirement: Archive active change from the index
For a selected active change, the contextual lifecycle action SHALL present a confirmation containing the change name, task completion summary when available, and affected project specifications. Confirming SHALL invoke the official non-interactive OpenSpec archive operation, including validation and project-spec updates. Success SHALL move the change into the date-prefixed archive, refresh the index, and select the corresponding archived item. Canceling SHALL leave all files unchanged.

#### Scenario: Confirm archive
- **WHEN** an active change is selected, the user invokes the lifecycle action, and confirms
- **THEN** OpenSpec validates and archives the change, applicable project specifications are updated, and the archived item is selected

#### Scenario: Archive incomplete change warning
- **WHEN** the selected active change has incomplete tasks and the user invokes the lifecycle action
- **THEN** the confirmation identifies the incomplete task count before archive can be confirmed

#### Scenario: Cancel archive
- **WHEN** archive confirmation is visible and the user cancels
- **THEN** the active change and project specifications remain unchanged

### Requirement: Reactivate archived change from the index
For a selected archived change, the contextual lifecycle action SHALL present a confirmation to make it active. Confirming SHALL move the archived directory back to `openspec/changes/<clean-name>` without reversing project-spec content merged by an earlier archive. Dossier SHALL refuse reactivation if the target active path already exists. Success SHALL refresh the index and select the active item.

#### Scenario: Confirm reactivation
- **WHEN** an archived change is selected, its clean active target is available, and the user confirms the lifecycle action
- **THEN** the change moves to Active Work and the corresponding active hierarchy row is selected

#### Scenario: Reactivation preserves project specifications
- **WHEN** an archived change is reactivated
- **THEN** project specifications previously updated by archive remain unchanged

#### Scenario: Active target collision
- **WHEN** the clean active target already exists and the user attempts reactivation
- **THEN** Dossier reports the collision and leaves both directories unchanged

### Requirement: Confirm lifecycle mutations
Archive and reactivation SHALL use a dedicated confirmation state. With the default keymap, `a` or `Enter` SHALL confirm and `Esc` SHALL cancel. While confirmation is active, unrelated index actions SHALL not execute and repeated confirmation input SHALL start at most one operation.

#### Scenario: Confirm with contextual action key
- **WHEN** a lifecycle confirmation is visible and the user presses `a` with the default keymap
- **THEN** Dossier starts exactly one confirmed lifecycle operation

#### Scenario: Cancel with escape
- **WHEN** a lifecycle confirmation is visible and the user presses `Esc`
- **THEN** Dossier returns to normal index interaction without mutation

### Requirement: Undo latest lifecycle action safely
A writable session SHALL retain an undo record for only the latest successful archive or reactivation. Undoing an archive SHALL restore the change to its original active path and restore every project specification changed by that archive to its exact pre-archive state. Undoing a reactivation SHALL return the change to its exact original archive path without modifying project specifications. Undo SHALL be session-local, SHALL require confirmation, and SHALL be unavailable after it succeeds or after Dossier exits.

Before undo, Dossier SHALL verify that every affected path still matches the state produced by the recorded lifecycle action and that restoration targets are free. If any precondition fails, Dossier SHALL refuse the entire undo without partial restoration and report the conflicting path. A later successful lifecycle action SHALL replace the previous undo record.

#### Scenario: Undo archive restores change and specs
- **WHEN** the latest lifecycle action archived a change, affected paths are unchanged, and the user confirms undo
- **THEN** the change returns to its original active path and all project specifications changed by archive return to their pre-archive contents

#### Scenario: Undo reactivation restores original archive path
- **WHEN** the latest lifecycle action reactivated a change, affected paths are unchanged, and the user confirms undo
- **THEN** the change returns to the same date-prefixed archive path it occupied before reactivation

#### Scenario: Refuse undo after external modification
- **WHEN** a path affected by the latest lifecycle action changed afterward
- **THEN** Dossier reports that path, performs no undo mutation, and retains the undo record

#### Scenario: Undo is session-local
- **WHEN** Dossier restarts after a lifecycle action
- **THEN** no undo action is available for the previous session

#### Scenario: New lifecycle action replaces undo record
- **WHEN** a second lifecycle action succeeds before the first is undone
- **THEN** undo applies only to the second action

### Requirement: Transactional lifecycle failure handling
Dossier SHALL capture the pre-operation state required to reverse an archive before invoking OpenSpec. If archive or reactivation fails after changing any affected path, Dossier SHALL attempt to restore the captured pre-operation state before returning control to the index. If restoration cannot complete, Dossier SHALL report the operation error, the restoration error, and the affected paths requiring manual recovery. The index SHALL refresh from disk after every success or failure.

#### Scenario: Archive command fails before mutation
- **WHEN** the OpenSpec archive operation fails without changing files
- **THEN** Dossier reports the error and the index continues to show the active change

#### Scenario: Archive fails after partial mutation
- **WHEN** the archive operation fails after changing an affected path
- **THEN** Dossier attempts to restore all captured paths and reports whether restoration succeeded

#### Scenario: Recovery also fails
- **WHEN** both a lifecycle operation and its automatic restoration fail
- **THEN** Dossier reports both errors and identifies paths requiring manual recovery

### Requirement: Index actions remain responsive and preserve navigation
External OpenSpec operations SHALL execute without blocking rendering, input handling unrelated to the pending item, or the polling cycle. While an item action is pending, Dossier SHALL show progress and prevent a duplicate operation for that item. After completion, Dossier SHALL rebuild the index from disk and preserve or transfer selection to the logical result item when possible.

#### Scenario: Pending archive shows progress
- **WHEN** an archive command is still running
- **THEN** the index remains rendered, shows archive progress, and ignores duplicate archive confirmation

#### Scenario: Selection follows lifecycle result
- **WHEN** archive or reactivation succeeds
- **THEN** the cursor selects the same logical change in its destination section

### Requirement: Missing OpenSpec CLI does not block navigation
If the `openspec` executable is unavailable, Dossier SHALL continue to start and provide all filesystem-based navigation and inspection. Invoking create, validate, or archive SHALL report that the OpenSpec CLI is required and SHALL not mutate files. Reactivation and undo of a reactivation MAY remain available because they do not require the CLI.

#### Scenario: Start without OpenSpec executable
- **WHEN** Dossier starts in a project and the OpenSpec CLI cannot be found
- **THEN** the index and inspection workflows remain available

#### Scenario: Invoke CLI-backed action without executable
- **WHEN** the user invokes create, validate, or archive without an available OpenSpec executable
- **THEN** Dossier reports the missing prerequisite and leaves project files unchanged

### Requirement: Read-only mode blocks index mutations
In read-only mode, Dossier SHALL block new change, archive, reactivation, edit, and undo actions from the index. Inspect and validate SHALL remain available because they do not mutate project artifacts. Mutation prompts and hints SHALL not be displayed.

#### Scenario: Read-only active change
- **WHEN** read-only mode is active and an active change is selected
- **THEN** inspect and validate are available while edit and archive are unavailable

#### Scenario: Read-only new change
- **WHEN** read-only mode is active and the user invokes new change
- **THEN** no prompt or external command is started and no files change

#### Scenario: Read-only archived change
- **WHEN** read-only mode is active and an archived change is selected
- **THEN** inspect remains available while make-active and undo are unavailable
