## MODIFIED Requirements

### Requirement: Navegación entre changes
The active change viewer SHALL navigate the selected change's available schema-defined artifacts in schema order and SHALL retain existing configured navigation between active changes. Opening a change from the hierarchy SHALL select the first inspectable artifact unless a specific artifact/output was selected. Switching changes SHALL select the first inspectable artifact of the destination change rather than assuming proposal exists.

#### Scenario: Change without proposal
- **WHEN** a custom-schema change has research as its first inspectable artifact and no proposal artifact
- **THEN** opening that change displays research

#### Scenario: Open selected output from index
- **WHEN** the user inspects a specific artifact output in the index
- **THEN** the viewer opens that exact output within its owning change

#### Scenario: Switch to different schema
- **WHEN** change navigation moves from a feature change to a spike
- **THEN** the viewer rebuilds navigation for the spike's artifacts and selects its first inspectable output

### Requirement: Tabs de artifact
The change viewer SHALL generate artifact navigation from the selected change's schema-ordered artifact model instead of fixed proposal/design/specs/tasks constants. Artifacts with no output SHALL remain visible with ready/blocked state but SHALL not render missing content. Multi-file artifacts SHALL provide output subnavigation. Positional numeric selection SHALL address visible schema artifacts in order, and next/previous artifact actions and mouse clicks SHALL support any artifact count. The active-work Git/code view SHALL remain a synthetic final navigation item and SHALL not appear as a schema artifact or in archived mode.

#### Scenario: Bugfix navigation
- **WHEN** a bugfix change exposes proposal, reproduction, diagnosis, specs, and tasks
- **THEN** the viewer shows those artifact labels in schema order

#### Scenario: Spike navigation
- **WHEN** a spike exposes proposal, research, and decision
- **THEN** no design, specs, or tasks labels are invented

#### Scenario: Ready artifact without output
- **WHEN** an artifact is ready but has no generated output
- **THEN** its navigation label shows ready state and selecting it presents an informative empty/status view

#### Scenario: Multiple artifact outputs
- **WHEN** the selected specs artifact contains multiple delta files
- **THEN** output subnavigation allows each file to be selected without concatenating unrelated artifacts

#### Scenario: Code view remains synthetic
- **WHEN** an active change is viewed inside a Git worktree
- **THEN** code navigation follows schema artifacts but is not included in artifact completion counts

### Requirement: Render de markdown con glamour
The TUI SHALL render any Markdown artifact output discovered from a standard or custom schema through Glamour at the content width. Artifact ID SHALL not determine whether generic Markdown can be rendered. Specialized task interaction SHALL remain available for a writable active artifact identified as tasks; archived tasks and all other Markdown outputs SHALL be scroll-only unless another capability explicitly provides interaction.

#### Scenario: Render custom diagnosis
- **WHEN** a bugfix diagnosis output is selected
- **THEN** its Markdown is rendered through Glamour and is scrollable

#### Scenario: Render spike decision
- **WHEN** a spike decision output is selected
- **THEN** Dossier renders it without requiring a hardcoded decision tab

#### Scenario: Active tasks retain task behavior
- **WHEN** the selected active artifact is tasks
- **THEN** existing task cursor, progress, and toggle behavior remains available

#### Scenario: Archived tasks remain read-only
- **WHEN** the selected archived artifact is tasks
- **THEN** it is rendered as scroll-only Markdown

## ADDED Requirements

### Requirement: Viewer identifies schema and artifact status
The change viewer header or artifact navigation SHALL identify the selected change's schema when known and SHALL present artifact readiness/completion without treating workflow schema as a project specification.

#### Scenario: Feature schema label
- **WHEN** an active change selects schema feature
- **THEN** the viewer identifies feature as the change workflow schema

### Requirement: Dynamic artifact selection survives refresh
Viewer selection SHALL be stored by artifact ID and output path. If the selected output remains after refresh it SHALL stay selected; if removed, Dossier SHALL choose the nearest available output in schema order without using a stale numeric index.

#### Scenario: Earlier artifact generated
- **WHEN** a new output appears before the selected artifact in schema order
- **THEN** the original selected artifact remains visible
