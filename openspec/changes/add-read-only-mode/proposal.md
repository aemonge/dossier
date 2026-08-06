## Why

Dossier currently mutates project state through task checkbox toggles, external editor launches, and Git stage/unstage actions. Users need an explicit mode for safely reviewing unfamiliar or protected projects without accidental writes. The command help also presents Go's single-dash long-option style instead of conventional `-h, --help` syntax.

## What Changes

- Add a `--read-only` command-line flag.
- Propagate read-only state into the UI model.
- In read-only mode, block task toggles, external editor launches, and Git stage/unstage actions while preserving navigation, polling, status, and diff viewing.
- Show a persistent read-only indicator and omit unavailable mutation hints from the help bar.
- Support both `-h` and `--help`, and render command usage with conventional long-option spelling.

## Non-goals

- Inferring read-only mode from filesystem permissions.
- Preventing other processes from changing files while dossier is running.
- Removing the tasks or Git tabs.
- Changing archive-mode behavior.

## Impact

- `cmd/dossier/main.go`: flag parsing, usage output, and model construction.
- `internal/ui/model.go`: read-only model state.
- `internal/ui/viewer.go` and `internal/ui/tasks.go`: mutation guards.
- `internal/ui/view.go`: mode indicator and read-only key hints.
- `cmd/dossier/main_test.go` and `internal/ui/view_test.go`: regression coverage.

## Estimate

Small story: 3–5 hours including tests, documentation, and verification.
