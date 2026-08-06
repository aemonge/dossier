## Why

Dossier's keyboard controls are hardcoded across several mode handlers and duplicated in help text. Users need to adapt bindings to their keyboard layout and established terminal workflows without recompiling.

## What Changes

- Add context-specific action bindings to the XDG TOML configuration shared with custom themes.
- Provide a coherent Neovim-oriented default keymap while preserving multiple keys per action.
- Route every keyboard action through the active keymap, including filtering, tab selection, navigation, task actions, Git actions, spec focus navigation, and config viewing.
- Validate duplicate keys within the same active context and reject unknown actions/invalid key names.
- Generate help-bar key labels from the active bindings.

## Non-goals

- Configuring mouse gestures.
- Multi-key chords or timed key sequences.
- Runtime keymap editing or hot reload.
- Changing action semantics.

## Dependencies

Depends on the XDG TOML settings loader introduced by `add-custom-themes`.

## Impact

- New keymap types and matching/help helpers.
- Updates to all keyboard mode handlers and help rendering.
- Expanded tests and configuration documentation.

## Estimate

Medium story: 2–3 engineering days.
