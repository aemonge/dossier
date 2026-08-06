## 1. Keymap model and configuration

- [x] 1.1 Add characterization tests for every default context/action binding
- [x] 1.2 Add failing tests for sparse TOML overrides, unknown actions, invalid keys, and context-local collisions
- [x] 1.3 Implement typed default keymaps, merging, matching, and validation
- [x] 1.4 Replace legacy bindings with the approved Neovim-oriented defaults
- [x] 1.5 Add named key-style base selection through TOML and `--keystyle`

## 2. Handler integration

- [x] 2.1 Refactor viewer mode to use semantic bindings with behavior tests
- [x] 2.2 Refactor index and filter modes to use semantic bindings with behavior tests
- [x] 2.3 Refactor spec and config modes to use semantic bindings with behavior tests

## 3. Dynamic help and documentation

- [x] 3.1 Add failing tests that help labels reflect custom bindings and read-only omissions
- [x] 3.2 Generate all keyboard help labels from the active keymap
- [x] 3.3 Document every action and a complete example keymap
- [ ] 3.4 Run OpenSpec validation, formatting, lint, race/coverage tests, vet, and build

## Estimate

2–3 engineering days.
