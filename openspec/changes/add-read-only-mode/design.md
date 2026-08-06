## Context

Normal mode exposes three user-triggered mutation paths: `doToggle` writes `tasks.md`, `e` launches an editor that may write any artifact, and `s` stages or unstages Git paths. Read-only mode is explicit rather than inferred so behavior is predictable even on writable filesystems.

## Decisions

### D1: Read-only is immutable launch configuration

`--read-only` is parsed by the command and passed into model construction. The model retains a boolean for the lifetime of the process; there is no runtime toggle.

### D2: Guard every mutation boundary

Key handlers block editor launch and Git stage/unstage, while `doToggle` also checks read-only state directly. The direct task guard provides defense in depth for callers other than the keyboard dispatcher.

### D3: Keep observation features available

Git status polling, diff rendering, task navigation, and artifact navigation continue unchanged. Read-only controls dossier-initiated writes, not reads or external filesystem changes.

### D4: Make mode visible and hints truthful

The help bar carries a `read-only` indicator. Mutation key hints are omitted in this mode rather than advertising actions that do nothing.

### D5: Use a dedicated FlagSet and custom usage text

Command parsing uses a testable `flag.FlagSet`. Both `-h` and `--help` aliases are registered, while custom usage output displays conventional `-h, --help` and `--long-option` spellings.

## Risks / Trade-offs

- An external process can still mutate files; this is application read-only mode, not an OS sandbox.
- `$EDITOR` is completely disabled in read-only mode because dossier cannot guarantee an editor will refrain from writing.
- Silent no-op mutation keys prioritize safety and match existing archive behavior; the persistent mode indicator explains why hints are absent.
