## 1. Dependency baseline and startup configuration

- [ ] 1.1 Verify that `add-schema-aware-index-hierarchy`, `add-custom-themes`, `add-custom-keybindings`, and `add-read-only-mode` are archived or that their resulting hierarchy/settings/keymap/read-only contracts are established before implementation (estimate: 1 hour)
- [ ] 1.2 Add failing settings and model-construction tests for default index startup, `ui.start_view = "change"`, invalid values, no-active fallback, and explicit-path precedence (estimate: 3 hours)
- [ ] 1.3 Implement strict `UIConfig.StartView` loading and pass the resolved startup choice through root-project model construction without changing single-path startup (estimate: 3 hours)

## 2. Index keymap and inspection semantics

- [ ] 2.1 Add failing keymap tests for configurable new, lifecycle, inspect, primary, toggle, edit, validate, and undo actions, including default bindings and same-context collision detection (estimate: 2 hours)
- [ ] 2.2 Extend `settings.KeyIndex`, default profiles, merge behavior, and configuration examples with new/lifecycle/inspect/edit/validate/undo while retaining hierarchy primary/toggle actions (estimate: 2 hours)
- [ ] 2.3 Add failing UI and mouse tests that `i` inspects hierarchy rows, `Enter`/`Space` toggle expandable rows, `Enter` inspects artifact-output/requirement leaves, `?` opens project information, and custom bindings replace defaults (estimate: 3 hours)
- [ ] 2.4 Extend hierarchy row-capability dispatch with the new actions and move project-information behavior/help to its configured action (estimate: 3 hours)

## 3. OpenSpec CLI adapter

- [ ] 3.1 Add adapter tests with fake command execution for exact create, change/spec validate, and archive arguments, working directory, JSON fixtures, stderr propagation, malformed output, and missing executable errors (estimate: 5 hours)
- [ ] 3.2 Implement the focused OpenSpec CLI adapter with direct argument execution, typed results, and injectable runner boundaries (estimate: 5 hours)

## 4. Modal index actions and stable refresh

- [ ] 4.1 Add failing model tests for new-name input, lifecycle/undo confirmation, pending-operation duplicate suppression, cancellation, and ordinary-key isolation (estimate: 4 hours)
- [ ] 4.2 Implement the value-driven index action state machine and typed asynchronous result messages (estimate: 5 hours)
- [ ] 4.3 Add failing tests for restoring typed hierarchy identity across create, archive, reactivation, filtering, artifact enrichment, and index rebuilds (estimate: 3 hours)
- [ ] 4.4 Implement a shared full-hierarchy reload that transfers active/history identities instead of preserving stale numeric indices (estimate: 4 hours)

## 5. Create, validate, and edit actions

- [ ] 5.1 Add failing UI tests for schema catalog choice, one-schema preselection, explicit `--schema`, valid creation, cancellation, schema/name rejection, pending progress, missing CLI, and selecting the created Active Work row (estimate: 5 hours)
- [ ] 5.2 Implement searchable schema selection plus name prompt and asynchronous create workflow through the CLI adapter (estimate: 5 hours)
- [ ] 5.3 Add failing tests for validating active changes, project specs, and containing specs from requirement rows, including success and actionable failure output (estimate: 3 hours)
- [ ] 5.4 Implement contextual `v` validation without leaving or mutating the index (estimate: 3 hours)
- [ ] 5.5 Add failing tests for dynamic active artifact-output and canonical-spec edit-path resolution, archived/missing/unsafe output rejection, editor return enrichment, and hierarchy selection preservation (estimate: 3 hours)
- [ ] 5.6 Reuse the dynamic viewer editor launch path for contextual `e` editing from the hierarchy (estimate: 3 hours)

## 6. Reactivation filesystem service

- [ ] 6.1 Add filesystem-service tests for clean-name target derivation, root containment, destination collisions, successful moves, copy/remove fallback, and restoration after move failure (estimate: 5 hours)
- [ ] 6.2 Implement guarded archived-to-active moves behind an injectable lifecycle filesystem boundary (estimate: 5 hours)
- [ ] 6.3 Add failing UI tests for reactivation confirmation, cancellation, collision reporting, schema/artifact refresh, and transferring selection from History to Active Work (estimate: 3 hours)
- [ ] 6.4 Connect contextual `a` on archived changes to asynchronous reactivation (estimate: 3 hours)

## 7. Transactional archive and recovery

- [ ] 7.1 Add transaction tests for affected-spec discovery, absent/existing spec preimages, change-directory temp backup, fingerprints, and cleanup (estimate: 6 hours)
- [ ] 7.2 Implement scoped lifecycle snapshots for the selected change and its affected project specifications (estimate: 6 hours)
- [ ] 7.3 Add failure-injection tests for archive success, no-mutation failure, partial spec mutation, partial directory move, successful automatic restoration, and restoration failure diagnostics (estimate: 7 hours)
- [ ] 7.4 Implement transactional OpenSpec archive execution, post-state capture, automatic failure restoration, and exact manual-recovery reporting (estimate: 7 hours)
- [ ] 7.5 Add UI tests for schema-aware archive summaries, optional task warnings, delta-spec outputs, cancellation, pending state, transferring selection to History, and missing CLI (estimate: 4 hours)
- [ ] 7.6 Connect contextual `a` on active changes to confirmed asynchronous archive (estimate: 4 hours)

## 8. Safe session undo

- [ ] 8.1 Add undo-record tests for archive and reactivation, record replacement, successful cleanup, and process-lifetime disposal (estimate: 4 hours)
- [ ] 8.2 Add conflict-preflight tests for changed files, changed directories, occupied restoration targets, no-partial-write refusal, and retained undo availability (estimate: 5 hours)
- [ ] 8.3 Implement archive undo that restores exact project-spec preimages and the original active change, plus reactivation undo that restores the exact archive path (estimate: 7 hours)
- [ ] 8.4 Add UI tests and dispatch for contextual `u` confirmation, success, conflict reporting, one-use behavior, and selection transfer (estimate: 4 hours)

## 9. Read-only behavior, contextual help, and status presentation

- [ ] 9.1 Add failing tests that read-only mode omits and blocks new, lifecycle, edit, and undo while retaining inspect and validate (estimate: 3 hours)
- [ ] 9.2 Implement defense-in-depth read-only guards in both key dispatch and mutation service entry points (estimate: 2 hours)
- [ ] 9.3 Add failing rendering tests for hierarchy row-capability hints, custom labels, available undo, schema/name prompts, pending progress, and persistent recovery errors (estimate: 4 hours)
- [ ] 9.4 Implement contextual index help and action status/error rendering from the active keymap and action state (estimate: 4 hours)

## 10. Documentation and verification

- [ ] 10.1 Update README startup behavior, `[ui].start_view`, OpenSpec CLI prerequisite, keyboard reference, archive/reactivation semantics, and session undo safeguards (estimate: 3 hours)
- [ ] 10.2 Update complete example configurations with startup and index action settings (estimate: 2 hours)
- [ ] 10.3 Run strict OpenSpec validation, `make fmt`, focused and full race tests, `make lint`, `go vet ./...`, and `make build`; record and resolve all failures (estimate: 4 hours)
- [ ] 10.4 Review the implementation against every scenario in the three delta specs and document any deferred follow-up as a separate OpenSpec change (estimate: 3 hours)
