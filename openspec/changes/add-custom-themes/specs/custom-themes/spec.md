## ADDED Requirements

### Requirement: XDG theme configuration

The system SHALL load optional TOML settings from `$XDG_CONFIG_HOME/dossier/config.toml`, using the platform user configuration directory when `XDG_CONFIG_HOME` is unset. `--config <path>` SHALL override the default path.

#### Scenario: Default file absent
- **WHEN** no default configuration file exists
- **THEN** dossier starts with its existing default appearance

#### Scenario: Explicit file absent
- **WHEN** `--config` names a missing file
- **THEN** dossier exits with an actionable error

### Requirement: Built-in Gruvbox Light Soft theme

The system SHALL provide `gruvbox-light-soft` as a built-in theme selectable from TOML or `--theme`, using a light Gruvbox palette without yellow accents and preserving the terminal background.

#### Scenario: Select built-in Gruvbox theme
- **WHEN** `--theme gruvbox-light-soft` is supplied
- **THEN** Dossier builds its UI, Glamour, and Chroma styles from the embedded Gruvbox Light Soft theme

### Requirement: Configurable Dossier palette

The system SHALL permit every Dossier-owned foreground and background color role to be overridden from TOML, including viewport and Git diff backgrounds.

#### Scenario: Light terminal palette
- **WHEN** configuration selects the light base and overrides UI colors
- **THEN** every Dossier UI component is built from that merged palette without dark hardcoded fallbacks

### Requirement: Configurable Glamour colors

The system SHALL permit foreground/background overrides for every supported Glamour element and code-block syntax category on top of a named Glamour base style.

#### Scenario: Markdown heading override
- **WHEN** configuration sets the Glamour `h1` foreground
- **THEN** rendered H1 headings use that color

#### Scenario: Markdown code token override
- **WHEN** configuration sets the Glamour code-block `keyword` foreground
- **THEN** keywords in rendered Markdown code blocks use that color

### Requirement: Configurable Chroma colors

The system SHALL permit token-level Chroma style entries on top of a named Chroma base style for Git diff syntax rendering.

#### Scenario: Diff keyword override
- **WHEN** configuration defines a `Keyword` Chroma entry
- **THEN** Git diff keyword tokens use that configured style

### Requirement: Strict validation

Malformed TOML, unknown keys, unknown theme/style names, invalid colors, unsupported Glamour elements, and invalid Chroma token entries SHALL cause startup to fail with an actionable error.

#### Scenario: Unknown key
- **WHEN** configuration contains a misspelled key
- **THEN** dossier reports the unknown key rather than ignoring it
