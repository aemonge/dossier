## Why

Dossier's built-in UI palette defaults to dark colors, which can be unreadable on a light terminal. Users need every Dossier-owned color and the Markdown/diff syntax palettes to be configurable without recompiling.

## What Changes

- Load optional TOML configuration from the XDG path `$XDG_CONFIG_HOME/dossier/config.toml` (falling back to `os.UserConfigDir()`) with `--config` override support.
- Move built-in Dossier palettes out of Go literals into embedded TOML theme files.
- Allow a configured base theme plus overrides for every Dossier UI color.
- Allow element and code-token color overrides for Glamour Markdown rendering.
- Allow token-level Chroma style overrides for Git diff syntax highlighting.
- Reject malformed files, unknown keys, invalid theme names, colors, elements, and token types with actionable startup errors.
- Preserve current behavior when no configuration file exists.

## Non-goals

- Runtime theme editing or hot reload.
- Downloading themes from the network.
- Configuring non-color layout attributes such as padding, borders, or margins.
- Automatically writing a configuration file.

## Impact

- New shared settings package and TOML dependency.
- Theme construction and renderer integration change from static globals to validated runtime values.
- Built-in theme TOML resources, tests, CLI help, and documentation.

## Estimate

Medium story: 1–2 engineering days.
