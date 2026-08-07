## MODIFIED Requirements

### Requirement: Sección Specs en el índice
In `ModeIndex`, project-level specifications loaded from `openspec/specs/` SHALL appear beneath the **Canonical Specs** project section, between Active Work and History. Each canonical spec row SHALL show its name and requirement count and SHALL remain visually distinct from change-local delta spec outputs. If no canonical specs exist, the section SHALL show `No canonical specifications available`.

#### Scenario: Canonical specs present
- **WHEN** project-level specs exist
- **THEN** Canonical Specs lists each project spec with its requirement count between Active Work and History

#### Scenario: Delta and canonical spec share capability name
- **WHEN** an active change has a delta for `change-index` and a canonical `change-index` spec exists
- **THEN** the delta appears under the change's artifact hierarchy and the canonical spec appears under Canonical Specs

#### Scenario: No canonical specs
- **WHEN** `openspec/specs/` is absent or empty
- **THEN** Canonical Specs shows `No canonical specifications available`

### Requirement: Specs no seleccionables en el índice
Canonical spec rows and their expanded requirements SHALL remain navigable hierarchy items. Primary activation on a canonical spec SHALL toggle its requirements; dedicated inspect SHALL open `ModeViewingSpec`. Primary activation or dedicated inspect on a requirement SHALL open its focused requirement view.

#### Scenario: Primary activation toggles requirements
- **WHEN** the cursor is on a canonical spec and the user invokes primary activation
- **THEN** its requirements expand or collapse without leaving the index

#### Scenario: Dedicated inspect opens full canonical spec
- **WHEN** the cursor is on a canonical spec and the user invokes inspect
- **THEN** Dossier enters `ModeViewingSpec` for that canonical spec

#### Scenario: Requirement opens focused view
- **WHEN** the cursor is on an expanded requirement and the user invokes primary activation or inspect
- **THEN** Dossier opens that requirement in focused mode
