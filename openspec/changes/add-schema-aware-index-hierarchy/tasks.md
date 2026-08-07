## 1. Dynamic work model and metadata

- [ ] 1.1 Add failing loader tests for schema metadata, malformed metadata diagnostics, arbitrary artifact IDs, multi-file outputs, dependencies, and artifact status (estimate: 5 hours)
- [ ] 1.2 Introduce schema-aware `Change`, `ChangeArtifact`, and `ArtifactOutput` structures with project-relative safe identities (estimate: 5 hours)
- [ ] 1.3 Add failing tests for task progress and existing standard-schema behavior derived through dynamic artifacts rather than fixed fields (estimate: 4 hours)
- [ ] 1.4 Migrate loader-level task/spec helpers to dynamic artifact capabilities while retaining temporary compatibility at call sites (estimate: 5 hours)

## 2. Safe filesystem discovery and archive fallback

- [ ] 2.1 Add failing discovery tests for arbitrary Markdown artifacts, nested outputs, deterministic order, missing schema, malformed metadata, hidden files, symlink escapes, file-count limits, and read failures (estimate: 6 hours)
- [ ] 2.2 Implement safe degraded artifact discovery rooted beneath an active or archived change (estimate: 7 hours)
- [ ] 2.3 Add archived custom-schema fixtures for feature, bugfix, spike, unknown schema, and date/clean-name handling (estimate: 4 hours)
- [ ] 2.4 Migrate archived change loading to schema metadata plus safe discovery without changing read-only guarantees (estimate: 5 hours)

## 3. OpenSpec schema/status adapter

- [ ] 3.1 Add injectable command-runner tests for `openspec list --json`, `openspec schemas --json`, and `openspec status --change <name> --json`, including malformed JSON, stderr, missing executable, unsafe paths, and unknown status values (estimate: 7 hours)
- [ ] 3.2 Implement typed schema catalog, change status, artifact path/status, and diagnostic parsing in a focused adapter (estimate: 7 hours)
- [ ] 3.3 Add failing merge tests that enrich discovered changes with authoritative schema order/status while preserving readable fallback outputs on partial failures (estimate: 5 hours)
- [ ] 3.4 Implement active-change status enrichment with root-containment validation and stale-result fingerprint guards (estimate: 6 hours)

## 4. Bounded asynchronous enrichment and caching

- [ ] 4.1 Add tests for bounded concurrency, first-frame filesystem skeleton, unchanged fingerprint cache hits, selective invalidation, change deletion during status, and CLI recovery (estimate: 6 hours)
- [ ] 4.2 Implement asynchronous status commands/messages with bounded concurrency and per-change cache keys (estimate: 7 hours)
- [ ] 4.3 Update tick polling so cheap metadata/output fingerprints trigger only affected status refreshes rather than all changes every 500 ms (estimate: 5 hours)
- [ ] 4.4 Add timing-independent UI tests that late enrichment preserves stable selection and expansion (estimate: 4 hours)

## 5. Hierarchical index model and rendering

- [ ] 5.1 Add failing model tests for typed identities covering sections, active/archive changes, artifacts, outputs, canonical specs, and requirements (estimate: 5 hours)
- [ ] 5.2 Replace flat index item kinds/indices with stable hierarchy nodes and parent/child ownership (estimate: 7 hours)
- [ ] 5.3 Add rendering tests for Active Work, schema badges, artifact progress/status/requires, delta output children, Canonical Specs, History, degraded diagnostics, and empty states (estimate: 7 hours)
- [ ] 5.4 Implement hierarchy rendering with DAG dependencies as metadata and canonical/delta spec visual distinction (estimate: 7 hours)
- [ ] 5.5 Add tests and implementation for logical cursor restoration when siblings/outputs are inserted, removed, reordered, or enriched (estimate: 5 hours)

## 6. Expansion, filtering, mouse, and contextual navigation

- [ ] 6.1 Add failing tests for Enter/Space expansion of sections, changes, multi-output artifacts, canonical specs, and archives, plus leaf no-op/inspect behavior (estimate: 5 hours)
- [ ] 6.2 Implement identity-keyed expansion state and shared primary/toggle/inspect dispatch (estimate: 6 hours)
- [ ] 6.3 Add filtering tests for schema names, artifact IDs, output paths, canonical specs, requirements, ancestor revelation, and restoration of pre-filter collapse (estimate: 5 hours)
- [ ] 6.4 Implement hierarchy-aware filtering that retains ownership ancestry without mutating stored expansion (estimate: 6 hours)
- [ ] 6.5 Replace line-mirroring mouse lookup with hierarchy render metadata or update both paths together; test two-phase click for every row capability (estimate: 6 hours)
- [ ] 6.6 Implement row-capability help for expandable, inspectable, degraded, read-only, and leaf states (estimate: 4 hours)

## 7. Dynamic artifact viewer

- [ ] 7.1 Add viewer tests for feature, bugfix, spike, task-only, ready/blocked missing outputs, arbitrary artifact counts, and selected output identity (estimate: 7 hours)
- [ ] 7.2 Replace fixed artifact tab state with `{artifactID, outputPath}` selection and schema-ordered navigation (estimate: 8 hours)
- [ ] 7.3 Implement dynamic artifact labels/status, positional selection, next/previous navigation, mouse selection, overflow handling, and synthetic code destination (estimate: 7 hours)
- [ ] 7.4 Migrate Glamour loading and render-cache keys to generic output paths; test refresh/removal/insertion without raw-content flashes (estimate: 7 hours)
- [ ] 7.5 Add tests and implementation for opening a change, artifact, or exact output from hierarchy inspect while preserving return selection (estimate: 5 hours)

## 8. Specialized task, spec, editor, and Git behavior

- [ ] 8.1 Add regression tests that active tasks retain cursor/toggle/progress while archived tasks and unknown Markdown remain scroll-only (estimate: 5 hours)
- [ ] 8.2 Migrate task state and polling from fixed `Change.Tasks` to dynamic tasks artifact/output identity (estimate: 6 hours)
- [ ] 8.3 Add regression tests for multi-file delta spec subnavigation, selected spec refresh, and canonical requirement focus (estimate: 5 hours)
- [ ] 8.4 Migrate spec subnavigation and focus helpers to dynamic multi-file outputs (estimate: 6 hours)
- [ ] 8.5 Add editor tests for dynamic active outputs, unavailable outputs, archived rejection, return refresh, and safe resolved paths (estimate: 4 hours)
- [ ] 8.6 Migrate editor launch/return cache invalidation to artifact-output identity (estimate: 4 hours)
- [ ] 8.7 Add Git regression tests and keep code as a synthetic active-only destination outside schema artifacts/progress (estimate: 4 hours)

## 9. Archive and explicit-path viewers

- [ ] 9.1 Add UI tests for archived feature/bugfix/spike navigation, unavailable schema fallback, dynamic help, and read-only task/editor behavior (estimate: 5 hours)
- [ ] 9.2 Migrate archive viewer navigation/rendering to dynamic discovered artifacts (estimate: 5 hours)
- [ ] 9.3 Add CLI/model tests for explicit active/archived custom-schema paths, degraded fallback, and selected-change-only polling (estimate: 5 hours)
- [ ] 9.4 Migrate single-path construction and polling to the schema-aware model without sibling scans (estimate: 5 hours)

## 10. Remove fixed assumptions and integrate dependent change

- [ ] 10.1 Search production/tests for fixed `Proposal`, `Design`, `Tasks`, `Specs`, `SpecFiles`, fixed tab constants, and `1-4` assumptions; classify every remaining use (estimate: 3 hours)
- [ ] 10.2 Remove temporary fixed-artifact compatibility fields/helpers and update remaining tests to dynamic fixtures (estimate: 7 hours)
- [ ] 10.3 Update `add-index-workflow-actions` to depend on hierarchy identities and use the schema catalog during new-change creation (estimate: 5 hours)
- [ ] 10.4 Verify no Taskflow integration or schema-pack installation entered scope; capture either as separate follow-up changes if desired (estimate: 2 hours)

## 11. Documentation and verification

- [ ] 11.1 Update README with the OpenSpec hierarchy, canonical-vs-delta specs, schema badges, expansion/inspect controls, dynamic artifacts, and degraded mode (estimate: 4 hours)
- [ ] 11.2 Add test fixtures/documentation for at least spec-driven, bugfix, spike, task-only, multi-file, missing-schema, and archived workflows (estimate: 4 hours)
- [ ] 11.3 Run strict OpenSpec validation, `make fmt`, focused and full race tests, `make lint`, `go vet ./...`, and `make build`; resolve and record all failures (estimate: 5 hours)
- [ ] 11.4 Review implementation against every hierarchy scenario and verify the old flat/fixed model is no longer observable (estimate: 4 hours)
