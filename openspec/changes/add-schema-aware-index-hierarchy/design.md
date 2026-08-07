## Context

OpenSpec separates current truth (`openspec/specs/`) from proposed work (`openspec/changes/<id>/`) and lets each change select a workflow schema in `.openspec.yaml`. A schema defines an artifact dependency DAG with arbitrary IDs, generated output patterns, and an apply gate. Current OpenSpec structured status returns `schemaName`, schema-ordered artifacts, `requires`, status, output patterns, and resolved existing paths for active changes. `openspec schemas --json` exposes the available schema catalog.

Dossier currently models `Change` with fixed `Proposal`, `Design`, `Tasks`, `Specs`, and `SpecFiles` fields. The loader parses schema metadata into a temporary struct but discards it; the viewer and keymap then assume four fixed artifact tabs. The index separately lists active changes, project specs, and archives without showing the ownership relationship between a change and its delta artifacts.

Artifact dependencies form a DAG, not a strict tree. The UI needs a hierarchy for ownership/navigation while preserving dependencies as metadata rather than falsely nesting or duplicating artifacts.

## Goals / Non-Goals

**Goals:**
- Represent standard and custom OpenSpec workflows without artifact-name assumptions.
- Make project truth, active work, change-local deltas, and history visually unambiguous.
- Keep hierarchy selection/expansion stable across asynchronous enrichment and polling.
- Preserve useful specialized behavior for tasks, multi-file specs, Git, archive read-only mode, and explicit paths.
- Remain navigable in a degraded filesystem-only mode when OpenSpec status is unavailable.

**Non-Goals:**
- Render the dependency DAG as a graph visualization.
- Define, install, or modify issue-type schemas.
- Integrate Taskflow flows/runs; that requires explicit change/run linkage and is a separate capability.
- Generalize arbitrary non-Markdown artifact rendering in this change.

## Decisions

### D1: Replace fixed artifact fields with a schema-aware domain model

Introduce structures equivalent to:

```text
Change
  Name, Path, Created, DisplayDate, Schema
  Artifacts []ChangeArtifact
  Diagnostic

ChangeArtifact
  ID, OutputPattern, Status, Requires []string
  Outputs []ArtifactOutput
  DiscoverySource

ArtifactOutput
  RelativePath, DisplayName, Content, Present
```

Canonical `ProjectSpec` remains a separate project-level model because it represents merged truth, not one change's output. Task progress is derived by capability from an artifact/output identified as tasks rather than stored as a top-level `Change.Tasks` field.

Alternative considered: add dynamic artifacts alongside the four fixed fields. Rejected because it creates two sources of truth and ensures new code continues accidentally depending on fixed fields.

### D2: OpenSpec structured status is authoritative for active changes

A focused injectable CLI adapter owns these read contracts:

- `openspec list --json`
- `openspec schemas --json`
- `openspec status --change <name> --json`

For each active change, status provides schema order, dependencies, status, output patterns, and resolved existing output paths. Dossier validates every resolved path is beneath the reported change root before reading it. The adapter parses only documented fields and retains stderr/JSON diagnostics.

Alternative considered: parse every project/global/package schema directly. Rejected because it would duplicate OpenSpec schema resolution, shadowing, validation, glob semantics, and future format changes.

### D3: Load quickly, enrich asynchronously, and cache status

Initial filesystem discovery creates the project/change skeleton immediately from directories and `.openspec.yaml`. Active changes are enriched with structured status through bounded asynchronous commands (maximum four concurrent by default). The UI may first render discovered files, then rebuild the affected branch when status arrives.

Cache status by project root, change name, metadata fingerprint, relevant output-tree fingerprint, and OpenSpec version/schema catalog generation. The 500 ms UI tick checks cheap filesystem fingerprints and only reruns status for affected changes. It never launches status for every unchanged change on every tick.

Alternative considered: synchronous status for all changes before the first frame. Rejected because projects with many changes would have slow, failure-prone startup.

### D4: Filesystem fallback is deliberately less authoritative

When status fails, recursively discover readable Markdown files beneath the change root while excluding metadata, hidden/tool state, and unsafe symlink escapes. Group conventional top-level files by stem and nested outputs by their nearest meaningful artifact directory; mark every fallback node as `discovered`. Preserve the metadata schema name when readable, but do not invent `requires`, ready/blocked state, or schema order.

Fallback guarantees access to planning content; it does not claim to reconstruct the schema DAG. Diagnostics remain visible until authoritative enrichment succeeds.

### D5: Archived work uses metadata plus safe discovery

OpenSpec status addresses active changes, not dated archive paths. Archived changes therefore load schema metadata and safely discover actual outputs. If the schema remains resolvable in the catalog, its artifact ID order may order matching outputs; otherwise stable relative-path order is used. Archived outputs remain read-only and never participate in apply readiness.

Alternative considered: temporarily reactivate an archive to ask OpenSpec for status. Rejected as an unacceptable mutation for read-only inspection.

### D6: The index is an ownership tree with stable typed identities

Use typed hierarchy nodes:

```text
section(active-work|canonical-specs|history)
change(active|archive, relative change path)
artifact(change identity, artifact ID)
artifact-output(change identity, artifact ID, relative output path)
canonical-spec(spec name)
requirement(spec name, requirement anchor/name)
```

Artifacts are direct siblings under a change in schema/status order. `requires` and status render as secondary metadata. Multi-output artifacts gain output children; single-output artifact rows may be directly inspectable without requiring a redundant child row, while identity still includes the output path.

Expansion maps, filter ancestry, cursor restoration, mouse hit testing, and polling reconciliation all use these identities rather than slice indices.

### D7: Section names describe semantics, not storage shorthand

Rename the top-level UI categories:

- `Active Work`: generic active OpenSpec changes, with schema badge.
- `Canonical Specs`: project truth from `openspec/specs/`.
- `History`: archived changes and their available artifacts.

Change-local specs are labeled `delta specs` or by their artifact ID under the owning change. This preserves OpenSpec's “Specs = what is true; Changes = what should change” model.

### D8: Enter is contextual primary; Space is toggle; i is inspect

Primary activation (`Enter` by default) toggles expandable sections, changes, multi-output artifacts, and canonical specs; on a leaf output/requirement it inspects. Toggle (`Space`) only toggles expandable nodes. Dedicated inspect (`i`) opens an inspectable row without changing expansion.

Selected-item mouse activation reuses primary dispatch. Help is generated from row capabilities, not hardcoded item-kind strings.

### D9: Change viewer navigation becomes dynamic

Replace fixed `Tab` constants for proposal/design/specs/tasks with artifact selection by `{artifactID, outputPath}`. The navigation bar is generated from schema artifacts, preserving ready/blocked nodes and multi-file subnavigation. Positional numeric shortcuts target schema positions rather than named artifacts; configured next/previous navigation handles arbitrary counts. Git/code remains a synthetic active-only destination outside artifact completion.

Generic Markdown outputs use one rendering path. Existing specialized behavior is selected by capability:

- artifact ID/output recognized as tasks → task parsing/toggling/progress in writable active mode;
- multi-file specs output → file subnavigation and requirement focus where applicable;
- archived outputs → read-only Glamour rendering;
- unknown Markdown artifact → scrollable Glamour rendering.

Alternative considered: keep fixed viewer tabs and inspect custom artifacts only from the index. Rejected because explicit-path launches and change inspection would remain incomplete and inconsistent.

### D10: Schema catalog is exposed for the dependent create workflow

Load and cache `openspec schemas --json` as a catalog service. This change only exposes the catalog in the model/API; `add-index-workflow-actions` owns the create prompt and SHALL use it to select a schema when multiple choices exist. With one schema, creation may preselect it while still invoking `openspec new change <name> --schema <schema> --json` explicitly.

### D11: Compatibility is delivered by capability adapters, not fixed fields

During implementation, migrate consumers in dependency order: loader/model, index, viewport/render cache, task behavior, specs subnavigation, editor path resolution, polling, Git synthetic tab, archive, and explicit path. Temporary compatibility helpers may exist within one delivery slice but fixed public `Change.Proposal/Design/Tasks/Specs` fields are removed before completion.

## Risks / Trade-offs

- [Risk: N+1 status subprocesses] → Render filesystem skeleton first, bound concurrency, fingerprint/cache results, and refresh only affected branches.
- [Risk: CLI JSON changes] → Isolate parsing, use minimal documented fields, keep fixtures for supported versions, and degrade to filesystem navigation.
- [Risk: DAG presented as hierarchy] → Make ownership hierarchy explicit and show dependencies as metadata; never dependency-nest or duplicate artifacts.
- [Risk: Archived schema removed] → Discover actual files and clearly mark degraded ordering/status.
- [Risk: Dynamic viewer breaks specialized behavior] → Migrate each capability with regression tests before removing fixed fields.
- [Risk: Very wide artifact sets overflow the tab bar] → Use horizontal clipping/scrolling or compact labels while next/previous navigation remains available.
- [Risk: Filesystem discovery reads unintended files] → Restrict to Markdown beneath the change root, reject escaping symlinks/resolved paths, and cap file count/size with diagnostics.
- [Risk: Concurrent status result arrives after change removal] → Apply enrichment only when stable change identity and fingerprint still match.

## Migration Plan

1. Introduce dynamic domain structures, safe discovery, CLI fixtures, and adapters behind existing behavior.
2. Build hierarchical index nodes and stable identity/collapse/filter behavior while retaining fixed viewer compatibility temporarily.
3. Migrate generic viewer, render cache, tasks/specs/editor/Git behavior to dynamic artifact identity.
4. Migrate archives and explicit-path loading, then remove fixed artifact fields and tab assumptions.
5. Update `add-index-workflow-actions` to depend on this change and use hierarchy identities/schema catalog.
6. Run strict OpenSpec validation and full Go verification before archive.

Rollback before fixed-field removal is straightforward. After migration, rollback requires reverting the complete change because mixed fixed/dynamic models are not a supported persisted state; no on-disk OpenSpec data migration occurs.
