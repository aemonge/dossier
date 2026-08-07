## MODIFIED Requirements

### Requirement: Vista índice de pantalla completa
The TUI SHALL implement `ModeIndex` as a project hierarchy with three top-level sections in this order: **Active Work**, **Canonical Specs**, and **History**. Active Work SHALL contain active change rows and their expandable schema-defined artifacts/output files. Canonical Specs SHALL contain project-level specifications and expandable requirements. History SHALL contain archived change rows and their expandable available artifacts/output files. Canonical specs SHALL never be rendered as children of an individual change. Empty sections SHALL remain visible with an appropriate empty-state message. The full index background SHALL use the configured view background throughout the terminal viewport.

#### Scenario: Project hierarchy with all categories
- **WHEN** active changes, canonical specs, and archived changes exist
- **THEN** the index shows Active Work, Canonical Specs, and History in that order

#### Scenario: Canonical specs remain project-level
- **WHEN** an active change contains a delta spec for capability `authentication` and the project contains canonical spec `authentication`
- **THEN** the delta output appears beneath that change's specs artifact while the canonical spec appears separately beneath Canonical Specs

#### Scenario: Empty active work
- **WHEN** no active change directories exist
- **THEN** Active Work shows an empty-state message while Canonical Specs and History remain navigable

#### Scenario: Configured background fills hierarchy
- **WHEN** a view background is configured
- **THEN** the hierarchy and all remaining viewport whitespace use that background

### Requirement: Formato de cambios activos en el índice
Each active change SHALL be displayed as an expandable work row with its name, selected schema, planning-artifact lifecycle summary, and implementation task progress when available. The schema SHALL render as a colored badge separate from the change name, and right-side metadata SHALL remain aligned without prematurely truncating the name. Expanding the change SHALL reveal schema-ordered artifact rows. Artifact rows SHALL translate reported ready/blocked/done state into unambiguous planning-document labels, render those status badges as high-contrast colored chips, and right-align them. Artifact identifiers and prerequisite metadata SHALL use stable shared columns across visible rows; multi-file outputs SHALL expand to file rows.

#### Scenario: Feature work row
- **WHEN** an active change `add-export` uses schema `feature` with three of four artifacts done
- **THEN** its row identifies `add-export`, schema `feature`, and planning-artifact progress `3/4` separately from implementation task progress

#### Scenario: Authored planning artifact is not implementation completion
- **WHEN** OpenSpec reports the tasks artifact as `done` while its Markdown checklist has incomplete tasks
- **THEN** the artifact row identifies the planning document as authored and the change row separately reports incomplete implementation task progress

#### Scenario: Bugfix children follow schema order
- **WHEN** a bugfix schema reports proposal, reproduction, diagnosis, specs, and tasks
- **THEN** expanding the change shows those five artifact rows in that order

#### Scenario: Multi-file delta artifact
- **WHEN** a specs artifact has multiple capability outputs
- **THEN** expanding the artifact shows each capability spec beneath it

### Requirement: Formato de cambios archivados en el índice
Each archived change SHALL be displayed beneath History as an expandable read-only work row with clean name, archive date, and schema when known. Expanding it SHALL reveal its available artifact/file hierarchy. The name/date columns SHALL remain aligned within archived work rows.

#### Scenario: Archived schema displayed
- **WHEN** archive directory `2026-05-02-fix-cache` contains metadata selecting `bugfix`
- **THEN** its History row shows clean name `fix-cache`, date `02/05/2026`, and schema `bugfix`

#### Scenario: Archived artifacts expand read-only
- **WHEN** the user expands an archived spike
- **THEN** its proposal, research, and decision files appear as read-only descendants

### Requirement: Navegación en el índice
The cursor SHALL move with configured up/down actions through every visible section, change, artifact, output file, canonical spec, requirement, and archived descendant in rendered hierarchy order. Hidden descendants of collapsed ancestors SHALL not be cursor targets. The cursor SHALL not move beyond the first or last visible row.

#### Scenario: Traverse active hierarchy
- **WHEN** an active change and one of its artifacts are expanded
- **THEN** repeated downward navigation visits the change, artifact, output files, and following work rows in rendered order

#### Scenario: Collapsed descendants skipped
- **WHEN** a change is collapsed
- **THEN** navigation moves from that change row directly to the next visible row

### Requirement: Seleccionar un change con Enter
The index SHALL provide contextual primary activation for every row. On an expandable row, activation SHALL toggle expansion without leaving the index. On a leaf output file or requirement, activation SHALL inspect that content. A dedicated inspect action SHALL inspect the selected inspectable row without changing expansion state: change inspection SHALL open its dynamic artifact viewer, artifact inspection SHALL open its first available output or status view, and canonical/requirement inspection SHALL use their existing viewers.

#### Scenario: Activate change toggles artifacts
- **WHEN** the cursor is on a collapsed active change and the user invokes primary activation
- **THEN** the change expands to show its schema-defined artifacts without leaving the index

#### Scenario: Inspect change opens dynamic viewer
- **WHEN** the cursor is on an active change and the user invokes dedicated inspect
- **THEN** Dossier opens the change viewer using schema-defined artifact navigation

#### Scenario: Activate artifact output inspects file
- **WHEN** the cursor is on a leaf artifact output and the user invokes primary activation
- **THEN** Dossier opens that exact output in the artifact viewer

#### Scenario: Activate canonical spec toggles requirements
- **WHEN** the cursor is on a canonical spec and the user invokes primary activation
- **THEN** its requirements expand or collapse without leaving the index

### Requirement: Helpbar del índice
The index help bar SHALL describe actions applicable to the selected hierarchy row using the active keymap. Expandable rows SHALL advertise primary/toggle expansion, inspectable rows SHALL advertise inspect, and unavailable leaf/container actions SHALL be omitted. Existing filter, sort, project information, read-only, and quit hints SHALL remain contextual.

#### Scenario: Change-row help
- **WHEN** an active change row is selected
- **THEN** help advertises expansion and inspection as distinct actions

#### Scenario: Artifact-file help
- **WHEN** a leaf artifact output is selected
- **THEN** help advertises inspection and does not advertise expansion

### Requirement: Actualización en tiempo real del índice
While in `ModeIndex`, Dossier SHALL poll for structural changes to active work, schemas/metadata, artifact outputs/status, canonical specs, and archived work. A detected change SHALL rebuild affected hierarchy branches while preserving unrelated expansion state and restoring the cursor by stable hierarchical identity. Polling SHALL not execute one OpenSpec status process per change every 500 ms when relevant metadata/output fingerprints are unchanged.

#### Scenario: Custom artifact appears
- **WHEN** a research output is generated for an expanded spike
- **THEN** it appears beneath the research artifact without restarting Dossier

#### Scenario: Schema status changes
- **WHEN** an artifact changes from blocked to ready
- **THEN** the corresponding artifact row updates while its change remains expanded

#### Scenario: Stable cursor after rebuild
- **WHEN** another change gains an output while the cursor is on an unchanged artifact
- **THEN** the cursor remains on the same logical artifact identity

## ADDED Requirements

### Requirement: Hierarchy rows use stable identities
Every hierarchy row SHALL have a stable identity composed from project-relative ownership and semantic identifiers rather than current slice positions. Identity SHALL distinguish active/archived changes, artifact IDs, output paths, canonical specs, and requirements.

#### Scenario: Earlier sibling inserted
- **WHEN** a new artifact output sorts before the selected output
- **THEN** rebuild preserves selection on the originally selected output

### Requirement: Filtering includes descendants and preserves ancestry
Index filtering SHALL match change names, schema names, artifact IDs, artifact output paths, canonical spec names, and requirement names case-insensitively. A matching descendant SHALL render with every ancestor required to understand its ownership, even when the ancestor text does not match. Filtering SHALL not permanently alter expansion state.

#### Scenario: Filter matches diagnosis artifact
- **WHEN** the user filters for `diagnosis`
- **THEN** matching diagnosis rows appear beneath their owning bugfix changes and Active Work remains visible

#### Scenario: Filter matches delta capability
- **WHEN** a change-local delta output path matches the query
- **THEN** its change and owning specs artifact remain visible above it

### Requirement: Dependency DAG is shown as metadata, not hierarchy
Artifact dependency relationships SHALL be represented as secondary status/prerequisite information while all artifacts remain sibling children of their owning change. The index SHALL not duplicate an artifact to place it under multiple prerequisites.

#### Scenario: Task requires specs and design
- **WHEN** tasks requires both specs and design
- **THEN** one tasks row appears under the change and identifies both prerequisites
