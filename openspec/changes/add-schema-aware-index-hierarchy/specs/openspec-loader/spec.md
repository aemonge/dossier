## MODIFIED Requirements

### Requirement: Leer metadatos del change
The loader SHALL read `openspec/changes/<name>/.openspec.yaml` to obtain at least the selected schema and creation date. Missing or malformed optional metadata SHALL not prevent loading the change; unavailable values SHALL be empty/unknown and the loader SHALL retain a diagnostic for presentation.

#### Scenario: Metadatos presentes
- **WHEN** `.openspec.yaml` contains `schema: bugfix` and `created: 2026-05-01`
- **THEN** the change exposes schema `bugfix` and creation date `2026-05-01`

#### Scenario: Malformed metadata
- **WHEN** `.openspec.yaml` is malformed but artifact files exist
- **THEN** the change and discovered files remain available with unknown schema and a metadata diagnostic

### Requirement: Cargar artifacts disponibles
The loader SHALL load arbitrary artifacts and their existing output files from the schema-aware work model rather than attempting only proposal, design, tasks, and specs. For active changes it SHALL prefer authoritative structured OpenSpec status; for unavailable status and archived changes it SHALL use safe filesystem discovery. Missing outputs SHALL not cause the change load to fail.

#### Scenario: Custom schema artifacts
- **WHEN** an active change uses a spike schema with proposal, research, and decision artifacts
- **THEN** the loader returns those schema-ordered artifacts without requiring design, tasks, or specs

#### Scenario: Artifact ausente
- **WHEN** schema status includes a ready artifact whose output does not yet exist
- **THEN** the artifact remains in the model with no output files and its reported status

#### Scenario: Fallback discovery
- **WHEN** structured status is unavailable and the change contains `proposal.md` and `research.md`
- **THEN** both files are exposed as discovered artifact outputs without fabricated dependencies

### Requirement: Releer artifacts de un change en disco
The loader SHALL expose a function that, given an already-loaded change, refreshes its schema metadata, dynamic artifact status, and existing output content. Newly generated or removed outputs and schema/status changes SHALL be reflected without restarting. A failure to refresh authoritative status SHALL fall back safely and return a visible diagnostic instead of discarding readable files.

#### Scenario: Custom artifact updated
- **WHEN** `diagnosis.md` is externally modified after a bugfix change is loaded
- **THEN** refresh returns the change with the updated diagnosis output content

#### Scenario: Artifact generated between reloads
- **WHEN** a previously ready artifact gains its configured output file
- **THEN** refresh marks the output available and makes it inspectable

#### Scenario: Schema changed between reloads
- **WHEN** `.openspec.yaml` selects a different valid schema
- **THEN** refresh rebuilds the artifact model from the newly selected schema status

#### Scenario: Status refresh failure
- **WHEN** OpenSpec status fails during refresh but existing files are readable
- **THEN** refresh preserves navigation through discovered files and reports degraded status
