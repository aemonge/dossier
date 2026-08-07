## Purpose

Represents each OpenSpec change as a generic schema-selected work container with dynamic artifacts and output files, allowing Dossier to support standard and custom workflows without hardcoded artifact names.

## ADDED Requirements

### Requirement: Change retains selected schema
Each active or archived change SHALL expose the schema name from its `.openspec.yaml` metadata. Dossier SHALL display that schema as work-type metadata and SHALL NOT infer a schema from the change name.

#### Scenario: Custom bugfix schema
- **WHEN** a change's metadata contains `schema: bugfix`
- **THEN** Dossier represents the change as schema `bugfix` regardless of its directory name

#### Scenario: Missing schema metadata
- **WHEN** a change has no readable schema value
- **THEN** Dossier labels its schema as unknown and continues discovering available files

### Requirement: Active artifact model follows OpenSpec status
For an active change, Dossier SHALL use structured OpenSpec status to obtain the selected schema, schema-ordered artifacts, dependency identifiers, artifact status, configured output pattern, and existing output paths. Artifact identifiers SHALL not be limited to proposal, design, specs, or tasks.

#### Scenario: Bugfix artifacts
- **WHEN** OpenSpec status returns proposal, reproduction, diagnosis, specs, and tasks for a bugfix change
- **THEN** Dossier exposes all five artifacts in that order with their reported dependencies and statuses

#### Scenario: Branched artifact dependencies
- **WHEN** two artifacts both require proposal and a later artifact requires both of them
- **THEN** Dossier retains all `requires` relationships without duplicating either artifact

#### Scenario: Artifact not yet generated
- **WHEN** schema status contains an artifact with no existing output path
- **THEN** Dossier retains the artifact and marks it with its reported blocked or ready status

### Requirement: Artifact outputs support one or many files
An artifact SHALL contain zero or more existing output files resolved beneath its change root. A single-file artifact SHALL be directly inspectable. A multi-file artifact SHALL expose each output file as a named child, preserving OpenSpec's reported path order. Dossier SHALL reject resolved paths outside the change root.

#### Scenario: Single proposal output
- **WHEN** the proposal artifact resolves only to `proposal.md`
- **THEN** the proposal has one directly inspectable output file

#### Scenario: Multiple delta specifications
- **WHEN** the specs artifact resolves to four `specs/<capability>/spec.md` files
- **THEN** the artifact exposes four child files identified by capability/path

#### Scenario: Output escapes change root
- **WHEN** structured status reports an output path outside the selected change directory
- **THEN** Dossier rejects that output, reports the unsafe path, and does not read it

### Requirement: Ownership hierarchy does not misrepresent dependency DAG
All schema artifacts SHALL be sibling children owned by their change and ordered according to OpenSpec status. Dependencies SHALL be displayed or exposed as metadata and SHALL NOT be represented by nesting one artifact beneath another, because an artifact may have multiple prerequisites.

#### Scenario: Specs and design both depend on proposal
- **WHEN** the schema graph branches from proposal into specs and design
- **THEN** proposal, specs, and design appear as sibling artifact rows under the change while specs/design identify proposal as a prerequisite

### Requirement: Filesystem fallback preserves navigation
If OpenSpec status is unavailable because the executable, schema, or structured response is unavailable, Dossier SHALL retain the schema name from metadata when readable and SHALL discover existing Markdown artifact files beneath the change root. Fallback artifacts SHALL be marked as discovered rather than assigned invented dependency or readiness information. The failure SHALL not prevent project, canonical-spec, or archive navigation.

#### Scenario: OpenSpec executable missing
- **WHEN** Dossier cannot execute OpenSpec while loading an active change
- **THEN** existing Markdown files remain navigable under that change and no dependency/status values are fabricated

#### Scenario: Unknown custom schema
- **WHEN** a change references a schema that no longer resolves but its artifact files remain on disk
- **THEN** Dossier shows the schema name and discovered files with a visible degraded-state indication

### Requirement: Archived changes preserve available artifact hierarchy
Archived changes SHALL retain their schema metadata and expose the artifact files that remain in their archived directory. When authoritative status is unavailable for archives, Dossier SHALL use schema information when resolvable and filesystem discovery otherwise, without making archived files writable.

#### Scenario: Archived custom workflow
- **WHEN** an archived spike contains proposal, research, and decision artifacts
- **THEN** all three are available beneath the archived change and remain read-only

#### Scenario: Archived schema definition removed
- **WHEN** an archived change references a schema that is no longer installed
- **THEN** its existing artifact files remain discoverable and inspectable

### Requirement: Schema catalog is available to dependent workflows
Dossier SHALL expose the available OpenSpec schema catalog, including schema name, description, artifact identifiers, and source when returned by OpenSpec, so change-creation workflows can require an explicit schema choice when multiple schemas are available.

#### Scenario: Multiple issue-type schemas installed
- **WHEN** OpenSpec reports feature, bugfix, task, spike, and refactor schemas
- **THEN** Dossier exposes all five choices in OpenSpec's returned order with their descriptions
