## MODIFIED Requirements

### Requirement: Pantalla de bienvenida sin changes activos
Unless Dossier is launched with an explicit path, the TUI SHALL use the configured startup view. The optional `[ui].start_view` setting SHALL accept only `"index"` and `"change"` and SHALL default to `"index"`. `"index"` SHALL open `ModeIndex` whether or not active changes exist. `"change"` SHALL open `ModeNormal` when active changes exist and SHALL fall back to `ModeIndex` when none exist. An explicit path launch SHALL open the selected path directly and SHALL take precedence over the startup-view setting. If the TUI enters `ModeNormal` while there are no active changes, it SHALL show an informational message with the available actions.

#### Scenario: Default startup with active changes shows the index
- **WHEN** Dossier starts without an explicit path, active changes exist, and `ui.start_view` is unspecified
- **THEN** the TUI opens `ModeIndex`

#### Scenario: Configured change startup opens a change
- **WHEN** Dossier starts without an explicit path, active changes exist, and `ui.start_view` is `"change"`
- **THEN** the TUI opens `ModeNormal` on the first active change

#### Scenario: Change startup without active changes falls back to index
- **WHEN** Dossier starts without an explicit path, no active changes exist, and `ui.start_view` is `"change"`
- **THEN** the TUI opens `ModeIndex`

#### Scenario: Explicit path takes precedence
- **WHEN** Dossier starts with an explicit active or archived change path and `ui.start_view` is `"index"`
- **THEN** the TUI opens the selected path directly instead of opening the root index

#### Scenario: Invalid startup view fails startup
- **WHEN** configuration sets `ui.start_view` to a value other than `"index"` or `"change"`
- **THEN** Dossier exits before launching the TUI and reports the invalid value and accepted values

#### Scenario: Sin changes activos desde ModeNormal
- **WHEN** the mode is `ModeNormal` and there are no active changes
- **THEN** the TUI shows an informational empty state with the available index and quit actions
