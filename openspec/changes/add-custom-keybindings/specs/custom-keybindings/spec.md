## ADDED Requirements

### Requirement: Neovim-oriented default keymap

The system SHALL provide a modifier-light Neovim-oriented default keymap in which `q` backs out, `Q` quits, `H`/`L` select the previous/next change, `h`/`l` move horizontally, and `Tab`/`Shift+Tab` select artifact tabs. `PgUp`/`PgDown` SHALL perform page scrolling with `Ctrl+U`/`Ctrl+D` as aliases.

#### Scenario: Navigate through index
- **WHEN** the user presses `q` in a change viewer and then `l` on an index item
- **THEN** Dossier enters the index and opens the selected item

#### Scenario: Page-scroll aliases
- **WHEN** the user presses `PgDown` or `Ctrl+D` in a scrollable viewer
- **THEN** Dossier scrolls down one page

### Requirement: Context-specific custom keybindings

The system SHALL load ordered key lists for every keyboard action from the XDG TOML configuration. Unspecified actions SHALL retain existing default bindings.

#### Scenario: Override navigation
- **WHEN** the viewer `down` action is configured as `["n"]`
- **THEN** `n` performs the action and the replaced default keys do not

#### Scenario: Preserve unspecified action
- **WHEN** configuration overrides only viewer `down`
- **THEN** all other viewer actions retain their defaults

### Requirement: Complete keyboard action coverage

Every keyboard action in viewer, index, filter editor, spec viewer, and config viewer modes SHALL be routed through the active keymap. Printable filter text entry SHALL continue to accept characters not consumed by configured editing actions.

#### Scenario: Custom filter acceptance
- **WHEN** the filter `accept` action is rebound
- **THEN** the configured key accepts the filter and ordinary printable characters still update filter text

### Requirement: Context-local collision validation

The system SHALL reject a key assigned to multiple actions in the same context while allowing the same key in different contexts.

#### Scenario: Collision in viewer
- **WHEN** `viewer.down` and `viewer.up` both contain `n`
- **THEN** startup fails and names both conflicting actions

#### Scenario: Reuse across contexts
- **WHEN** `n` is used by one viewer action and one index action
- **THEN** the keymap is valid

### Requirement: Dynamic key help

All keyboard hints SHALL be rendered from the active keymap. Disabled read-only mutation actions SHALL remain omitted.

#### Scenario: Custom task toggle hint
- **WHEN** task toggle is rebound from `Space` to `x`
- **THEN** the tasks help bar advertises `x: toggle`

### Requirement: Strict keymap validation

Unknown actions, empty key lists, empty key names, and unsupported multi-key sequences SHALL fail startup with an actionable configuration error.

#### Scenario: Unknown action
- **WHEN** TOML contains a misspelled action name
- **THEN** dossier reports the unknown key through strict TOML validation
