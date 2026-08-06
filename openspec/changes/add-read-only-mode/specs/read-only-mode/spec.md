## ADDED Requirements

### Requirement: Explicit read-only launch mode

The command SHALL accept `--read-only` and launch the TUI in read-only mode without changing which artifacts or Git information can be viewed.

#### Scenario: Launch with read-only flag
- **WHEN** the user runs `dossier --read-only`
- **THEN** the model is configured as read-only

### Requirement: Block dossier-initiated mutations

In read-only mode, the system SHALL NOT toggle task checkboxes, launch an external editor, or stage/unstage Git files. Navigation, live reload, Git status, and Git diff viewing SHALL remain available.

#### Scenario: Task toggle is blocked
- **WHEN** read-only mode is active and the user presses `Space` on a task
- **THEN** `tasks.md` and the in-memory task state remain unchanged

#### Scenario: Editor launch is blocked
- **WHEN** read-only mode is active and the user presses `e` on an artifact tab
- **THEN** no editor process is launched

#### Scenario: Git mutation is blocked
- **WHEN** read-only mode is active and the user presses `s` on a changed file
- **THEN** the Git index and displayed status remain unchanged

### Requirement: Read-only state is visible

The help bar SHALL identify read-only mode and SHALL omit key hints for disabled mutations.

#### Scenario: Tasks help in read-only mode
- **WHEN** read-only mode is active on the tasks tab
- **THEN** the help bar shows a read-only indicator and does not advertise task toggling or editing

#### Scenario: Git help in read-only mode
- **WHEN** read-only mode is active on the Git file list
- **THEN** the help bar shows a read-only indicator and does not advertise stage/unstage

### Requirement: Conventional help aliases and display

The command SHALL accept both `-h` and `--help`, and its usage output SHALL display help as `-h, --help` while spelling long-only options with two leading hyphens.

#### Scenario: Short help
- **WHEN** the user runs `dossier -h`
- **THEN** usage information is printed without launching the TUI

#### Scenario: Long help
- **WHEN** the user runs `dossier --help`
- **THEN** the same usage information is printed without launching the TUI
