## MODIFIED Requirements

### Requirement: View artifacts of an archived change
In `ViewingArchive` mode, the TUI SHALL display every available artifact output discovered for the archived change, using schema order when the selected schema still resolves and safe filesystem order otherwise. Navigation and rendering SHALL use the same dynamic artifact model as active changes, but all archived outputs SHALL remain read-only. Dossier SHALL not invent fixed proposal/design/specs/tasks tabs for artifacts absent from that archive.

#### Scenario: View archived spike
- **WHEN** an archived spike contains proposal, research, and decision outputs
- **THEN** archive navigation exposes those three artifacts in schema order

#### Scenario: View archive with missing schema
- **WHEN** the archived change's schema no longer resolves
- **THEN** Dossier shows its safely discovered Markdown files with a degraded-state indication

#### Scenario: Archived tasks scroll read-only
- **WHEN** an archived tasks output is selected and the user navigates vertically
- **THEN** the content scrolls without task cursor or mutation behavior

### Requirement: Helpbar adaptado en modo archivo
In `ViewingArchive` mode, the help bar SHALL render dynamic artifact/output navigation and scrolling actions from the active keymap, omit editing and task mutation, and include the configured return-to-index action. It SHALL not advertise a fixed artifact count.

#### Scenario: Dynamic archive help
- **WHEN** an archived custom-schema artifact is selected
- **THEN** help advertises available next/previous artifact or output navigation, scroll, index, and quit actions without fixed `1-4` wording
