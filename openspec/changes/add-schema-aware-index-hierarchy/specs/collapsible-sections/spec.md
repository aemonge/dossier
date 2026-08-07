## MODIFIED Requirements

### Requirement: Sections can be collapsed and expanded
The system SHALL support expansion and collapse for top-level project sections, active and archived change rows, multi-file artifact rows, and canonical spec rows. The configured primary action (`Enter` by default) and toggle action (`Space` by default) SHALL both toggle expandable rows. Toggle on a leaf row SHALL do nothing; primary activation on a leaf SHALL inspect it.

#### Scenario: Toggle top-level section
- **WHEN** the cursor is on Active Work and the user presses `Enter` or `Space` with the default keymap
- **THEN** all active change descendants collapse or expand

#### Scenario: Toggle change artifacts
- **WHEN** the cursor is on an active or archived change and the user presses `Enter` or `Space`
- **THEN** that change's artifact rows collapse or expand

#### Scenario: Toggle multi-file artifact
- **WHEN** the cursor is on an artifact with multiple output files and the user presses `Enter` or `Space`
- **THEN** its output-file rows collapse or expand

#### Scenario: Toggle canonical requirements
- **WHEN** the cursor is on a canonical spec and the user presses `Enter` or `Space`
- **THEN** its requirement rows collapse or expand

#### Scenario: Space on leaf does nothing
- **WHEN** the cursor is on an artifact output file or requirement and the user presses `Space`
- **THEN** no navigation or mutation occurs

### Requirement: Collapse state persists across rebuilds
The system SHALL preserve expansion state by stable hierarchy identity across polling, filtering, mode switches, sibling insertion/removal, and artifact-status refresh. State for a row that no longer exists SHALL be discarded safely.

#### Scenario: Expanded change survives poll
- **WHEN** a change is expanded and an artifact status refresh rebuilds its branch
- **THEN** the change remains expanded

#### Scenario: Sibling insertion does not transfer state
- **WHEN** a new change is inserted before an expanded change
- **THEN** expansion remains attached to the original change rather than its old numeric position

### Requirement: Filtering respects collapse
Filtering SHALL temporarily reveal the ancestors and matching descendants required to present results without overwriting the user's stored expansion state. Clearing the filter SHALL restore the pre-filter hierarchy expansion.

#### Scenario: Match inside collapsed change
- **WHEN** a filter matches an artifact output inside a collapsed change
- **THEN** the matching path and its ancestors are visible while filtering

#### Scenario: Clear filter restores collapse
- **WHEN** the filter is cleared after revealing a match inside a collapsed change
- **THEN** that change returns to its stored collapsed state

### Requirement: Help bar shows section toggle action
The help bar SHALL advertise both configured primary and toggle keys on expandable hierarchy rows and SHALL omit toggle on leaf rows.

#### Scenario: Expandable row help
- **WHEN** a change, section, multi-file artifact, or canonical spec row is selected
- **THEN** the help bar advertises `Enter/Space: toggle` under the default keymap

## REMOVED Requirements

### Requirement: Enter is a no-op on section headers
**Reason**: Enter becomes the contextual primary action for the schema-aware hierarchy and toggles every expandable row.

**Migration**: Users can rebind the primary action away from Enter; Space remains the dedicated toggle action.
