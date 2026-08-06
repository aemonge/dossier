## Context

UI colors currently live in `internal/ui/themes.go`; Git diff backgrounds also have fallback literals in `internal/ui/git.go`. Glamour uses named standard styles and Chroma diff rendering uses named Chroma styles. Configuration is currently limited to OpenSpec project metadata and CLI flags.

## Decisions

### D1: XDG TOML configuration

Use `${XDG_CONFIG_HOME:-os.UserConfigDir()}/dossier/config.toml`; an explicit `--config` path overrides it. A missing default file is not an error, while a missing explicit file is.

### D2: Base plus sparse overrides

A theme selects a built-in base (`none`, `dark`, `light`, `dracula`, or `gruvbox-light-soft`) and then applies sparse user overrides. CLI `--theme` selects the base when explicitly supplied; TOML supplies it otherwise.

### D3: Embedded TOML built-ins

Dossier-owned palette values and renderer style names are stored in embedded TOML files. Go code validates and converts strings to Lip Gloss colors.

### D4: Complete color surfaces

The TOML schema exposes every Dossier UI palette role, viewport/diff backgrounds, every Glamour top-level element foreground/background, Glamour code-block token categories, and every Chroma token type accepted by Chroma v2. Existing named Glamour/Chroma styles remain useful bases.

### D5: Strict startup validation

Unknown TOML keys and invalid override targets fail startup with the file path and offending key. Configuration errors never silently fall back to an unintended dark palette.

## Risks / Trade-offs

- Glamour and Chroma token names are library API vocabulary and may evolve across dependency upgrades.
- Sparse overrides inherit future base-theme changes; users wanting a fixed palette should specify every relevant role.
- `--theme` remains compatible but becomes an explicit base selector rather than the only theme source.
