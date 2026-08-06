**English** | **[Español](README.es.md)**

# dossier

A keyboard-driven terminal UI for reading and navigating [OpenSpec](https://github.com/openspec) project artifacts — proposals, designs, specs, and tasks.

> Built with OpenSpec. This repository contains 12 project-level spec files and 20+ archived changes that document the complete development history of the tool.

<p align="center">
  <img src="docs/dossier.gif" alt="dossier demo" />
</p>

---

## Features

- Navigates all active changes and their artifacts from a single interface
- Renders Markdown with full syntax highlighting
- Toggles task checkboxes (`- [ ]` / `- [x]`) in-place, writing directly to `tasks.md`
- Live-reloads on disk changes (500 ms polling)
- Opens any artifact in `$EDITOR`
- Offers `--read-only` mode for safe review without task, editor, or Git index mutations
- Supports XDG TOML configuration for custom UI, Glamour, and Chroma colors
- Supports context-specific custom keybindings with dynamic help labels
- Accepts a path argument to view a single change directory without a full project

---

## Installation

**Requirements:** terminal with ANSI color support. Go 1.25 or later if building from source.

```bash
# Homebrew
brew tap fselich/tap
brew install dossier

# From source
git clone https://github.com/fselich/dossier
cd dossier
make build    # produces ./dossier
make install  # installs via go install

# Using go install
go install github.com/fselich/dossier/cmd/dossier@latest
```

---

## Usage

Run from the root of an OpenSpec project:

```bash
dossier
```

View a single change directory by path:

```bash
dossier /path/to/openspec/changes/my-change
```

Review without allowing dossier to toggle tasks, launch an editor, or stage/unstage Git files:

```bash
dossier --read-only
dossier --read-only /path/to/openspec/changes/my-change
```

Use `dossier -h` or `dossier --help` for all command-line options. Read-only mode remains visible in the help bar; navigation, live reload, Git status, and diff viewing stay available.

### Configuration

Dossier loads optional TOML configuration from the XDG path:

```text
${XDG_CONFIG_HOME:-~/.config}/dossier/config.toml
```

Use another file with `dossier --config <path>`. A missing default file is ignored; a missing explicit file or invalid setting produces a startup error. `--theme` overrides the configured base theme while preserving custom color overrides. Built-in themes are `none`, `dark`, `light`, `dracula`, and `gruvbox-light-soft`:

```bash
dossier --theme gruvbox-light-soft
```

Copy the complete light-terminal example, which documents every Dossier color and key action plus Glamour/Chroma token overrides:

```bash
mkdir -p "${XDG_CONFIG_HOME:-$HOME/.config}/dossier"
cp examples/config.example.toml "${XDG_CONFIG_HOME:-$HOME/.config}/dossier/config.toml"
```

A Pi-aligned Gruvbox Light Soft example is also available at `examples/config.gruvbox-light-soft.toml`.

Keybindings are context-specific, so a key may be reused in different modes but cannot be assigned to two actions in the same mode. Help labels automatically reflect the configured bindings.

### Keyboard reference

#### Normal mode (viewing a change)

| Key | Action |
|---|---|
| `j` / `down` | Scroll down (or move task cursor down) |
| `k` / `up` | Scroll up (or move task cursor up) |
| `PgDown` / `Ctrl+D` | Scroll one page down |
| `PgUp` / `Ctrl+U` | Scroll one page up |
| `h` / `l` | Scroll horizontally left / right |
| `H` / `L` | Previous / next change |
| `Tab` / `Shift+Tab` | Next / previous artifact tab |
| `1`–`5` | Select an artifact tab directly |
| `Space` | Toggle task under cursor (tasks tab only) |
| `e` | Open artifact in `$EDITOR` |
| `?` | Open configuration information |
| `q` / `Esc` | Enter index mode |
| `Q` | Quit |

#### Index mode (change and spec navigator)

| Key | Action |
|---|---|
| `j` / `down` | Move cursor down |
| `k` / `up` | Move cursor up |
| `l` / `Enter` | Open selected change, spec, or archived change |
| `Space` | Expand / collapse a project spec |
| `/` | Filter index items |
| `?` | Open configuration information |
| `q` / `Q` / `Esc` | Quit from the root index |

#### Archive mode (viewing an archived change)

| Key | Action |
|---|---|
| `j` / `k` | Scroll |
| `PgDown` / `Ctrl+D` | Scroll one page down |
| `PgUp` / `Ctrl+U` | Scroll one page up |
| `Tab` / `Shift+Tab` | Next / previous artifact tab |
| `1`–`4` | Select an artifact tab directly |
| `q` / `Esc` | Return to index |
| `Q` | Quit |

#### Spec viewer mode

| Key | Action |
|---|---|
| `j` / `k` | Scroll |
| `PgDown` / `Ctrl+D` | Scroll one page down |
| `PgUp` / `Ctrl+U` | Scroll one page up |
| `q` / `Esc` | Return to index |
| `Q` | Quit |

In requirement focus mode:

| Key | Action |
|---|---|
| `h` / `l` | Previous / next requirement |
| `j` / `k` | Scroll |
| `PgDown` / `Ctrl+D` | Scroll one page down |
| `PgUp` / `Ctrl+U` | Scroll one page up |
| `q` / `Esc` | Return to index |
| `Q` | Quit |

---

## Project structure

dossier expects an `openspec/` directory at the project root:

```
openspec/
├── changes/
│   ├── <change-name>/
│   │   ├── .openspec.yaml   # Required: identifies the directory as a change
│   │   ├── proposal.md
│   │   ├── design.md
│   │   ├── tasks.md         # GFM checkbox syntax: - [ ] / - [x]
│   │   └── specs/
│   │       └── <spec-name>/
│   │           └── spec.md
│   └── archive/
│       └── YYYY-MM-DD-<name>/
└── specs/
    └── <spec-name>/
        └── spec.md          # Requirements parsed from: ### Requirement: <name>
```
