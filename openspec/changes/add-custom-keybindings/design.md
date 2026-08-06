## Context

Bindings are direct `msg.String()` switch cases in `viewer.go`, `index.go`, `spec.go`, and `config.go`. Help text independently embeds labels in `view.go`. Keys are intentionally overloaded across contexts, so a global key-to-action map would reject valid behavior.

## Decisions

### D1: Typed context-specific actions

Define typed bindings for viewer, index, filter editor, spec viewer, and config viewer contexts. Each action accepts an ordered list of Bubble Tea key strings; the first key is the preferred help label.

### D2: Base plus sparse overrides

Named key styles provide a base, with `nvim` as the default. `--keystyle` overrides the configured base name while sparse `[keys.*]` action overrides still apply, mirroring `--theme` precedence. The Neovim profile uses `q` to back out, `Q` to quit, `h`/`l` to switch visible artifacts, `H`/`L` to move horizontally, and `Tab`/`Shift+Tab` to switch the less-visible change dimension. Page keys remain primary while `Ctrl+U`/`Ctrl+D` are aliases.

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
