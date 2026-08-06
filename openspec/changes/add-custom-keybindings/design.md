## Context

Bindings are direct `msg.String()` switch cases in `viewer.go`, `index.go`, `spec.go`, and `config.go`. Help text independently embeds labels in `view.go`. Keys are intentionally overloaded across contexts, so a global key-to-action map would reject valid behavior.

## Decisions

### D1: Typed context-specific actions

Define typed bindings for viewer, index, filter editor, spec viewer, and config viewer contexts. Each action accepts an ordered list of Bubble Tea key strings; the first key is the preferred help label.

### D2: Base plus sparse overrides

Default bindings use a coherent Neovim-oriented profile: `q` backs out, `Q` quits, `H`/`L` switch changes, `h`/`l` move horizontally, and `Tab`/`Shift+Tab` switch artifacts. Page keys remain primary while `Ctrl+U`/`Ctrl+D` are aliases. TOML replaces only actions it names, allowing small personal keymap files.

### D3: Context-local collision validation

The same key may serve different actions in different contexts. Within one context, one key may not map to two actions. Printable filter input remains text input after configured editing controls are checked.

### D4: Generated help labels

Help-bar text uses the active binding labels rather than hardcoded key names. Read-only mode continues to omit disabled mutation actions.

### D5: No chords

Configuration accepts Bubble Tea key strings such as `ctrl+c`, `shift+tab`, `pgdown`, and single Unicode characters. Multi-key sequences are rejected as unsupported.

## Risks / Trade-offs

- Existing mode handlers require broad but mechanical refactoring; characterization tests protect current defaults.
- Some actions have contextual no-op behavior even when their key is configured.
- Very wide custom labels may make help text overflow, matching current behavior for narrow terminals.
