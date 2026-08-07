## MODIFIED Requirements

### Requirement: User can open project config view
The TUI SHALL provide a configurable information action, bound to `?` by default, that opens a full-screen read-only view of `openspec/config.yaml` from both `ModeNormal` and `ModeIndex`. The default `i` key in `ModeIndex` SHALL remain available for inspecting the selected index item rather than opening project configuration.

#### Scenario: Open from index with default binding
- **WHEN** the user is in `ModeIndex` and presses `?` with the default keymap
- **THEN** the TUI transitions to `ModeViewingConfig` and renders the config content in the viewport

#### Scenario: Inspect remains available from index
- **WHEN** the user is in `ModeIndex` and presses `i` with the default keymap
- **THEN** Dossier inspects the selected index item and does not open project configuration

#### Scenario: Open from normal mode with default binding
- **WHEN** the user is in `ModeNormal` and presses `?` with the default keymap
- **THEN** the TUI transitions to `ModeViewingConfig` and renders the config content in the viewport

#### Scenario: Configured information binding
- **WHEN** the information action is rebound and the user presses its configured key in `ModeNormal` or `ModeIndex`
- **THEN** the TUI transitions to `ModeViewingConfig`

### Requirement: Config view help bar shows navigation hints
The help bar in `ModeViewingConfig` SHALL render its navigation and back actions from the active keymap. The help bars in `ModeIndex` and `ModeNormal` SHALL advertise the configured information action, using `?: info` under the default keymap.

#### Scenario: Help bar content in config view
- **WHEN** the config view is open
- **THEN** the help bar displays the configured scroll and back actions

#### Scenario: Default help bar hint in index and normal mode
- **WHEN** the user is in `ModeIndex` or `ModeNormal` with the default keymap
- **THEN** the help bar includes `?: info`

#### Scenario: Reconfigured information hint
- **WHEN** the information action uses a custom binding
- **THEN** the index and normal-mode help bars advertise that binding instead of `?`
