## MODIFIED Requirements

### Requirement: Invocación con ruta explícita
The binary SHALL accept an optional path to an active or archived change directory. When provided, Dossier SHALL load that change's schema metadata and dynamic artifact outputs without scanning sibling changes, and SHALL open its schema-aware viewer directly. Archived paths SHALL remain read-only. Without a path, root-project startup behavior SHALL apply.

#### Scenario: Explicit custom active change
- **WHEN** the user launches Dossier with a bugfix change path
- **THEN** only that change loads and its proposal, reproduction, diagnosis, specs, and tasks artifacts are available according to its schema/status

#### Scenario: Explicit archived custom change
- **WHEN** the user launches Dossier with an archived spike path
- **THEN** its available proposal, research, and decision outputs open in read-only archive mode

#### Scenario: Explicit path with unavailable schema
- **WHEN** the selected change references an unavailable schema but contains readable Markdown files
- **THEN** Dossier opens those discovered files in degraded mode rather than failing solely because the schema is unavailable

### Requirement: Polling estable en modo path
When launched with an explicit path, polling SHALL refresh only the selected change's metadata, artifact status when available, and output files. It SHALL not scan sibling changes, canonical specs, or archive listings.

#### Scenario: Tick in schema-aware path mode
- **WHEN** an explicit-path viewer is open and an artifact output is generated
- **THEN** only the selected change hierarchy/viewer refreshes and the new output becomes available
