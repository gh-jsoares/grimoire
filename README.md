# Grimoire

A terminal cheatsheet viewer built with Go and the [Charm](https://charm.sh) ecosystem. Renders `.grim` TOML files as styled reference cards with responsive multi-column layout.

## Features

- Responsive 1/2/3 column layout based on terminal width
- Tab switching between multiple cheatsheets
- Command navigation and clipboard yanking
- Link overlay picker
- Search with filter and highlight modes
- Nerd Font icon support
- Vim-style keybinds

## Install

### From source

```sh
go install github.com/gh-jsoares/grimoire/cmd/grimoire@latest
```

### With Make

```sh
git clone https://github.com/gh-jsoares/grimoire.git
cd grimoire
make install
```

### Homebrew (coming soon)

```sh
brew install gh-jsoares/tap/grimoire
```

## Usage

```sh
# Open default library (~/.config/grimoire/)
grimoire

# Open a specific file
grimoire tmux.grim

# Open by name or alias
grimoire git
grimoire tm

# Open a directory as library
grimoire ~/cheatsheets/

# Flags
grimoire --tab git        # Start on a specific tab
grimoire --no-icons       # Disable Nerd Font icons
grimoire --plain          # Minimal styling
grimoire --library ~/dir  # Custom library path
```

### Library resolution

Grimoire looks for `.grim` files in:

1. `$GRIMOIRE_HOME`
2. `$XDG_CONFIG_HOME/grimoire`
3. `~/.config/grimoire`

### Keybinds

| Key | Action |
|-----|--------|
| `H` / `L` | Switch tabs |
| `1-9` | Jump to tab |
| `j` / `k` | Scroll |
| `g` / `G` | Top / bottom |
| `Ctrl+d` / `Ctrl+u` | Half-page scroll |
| `n` / `N` | Next / previous command |
| `y` | Yank (copy) active command |
| `u` | Open link picker |
| `/` | Search |
| `Tab` (in search) | Toggle filter/highlight mode |
| `Esc` | Clear search |
| `q` | Quit |

### tmux popup

Add to your `tmux.conf`:

```sh
bind-key ? display-popup -E -w 90% -h 85% "grimoire"
```

## File format

Grimoire reads `.grim` files — TOML with a specific schema:

```toml
format = 1

title = "My Tool"
icon = ""
order = 10
aliases = ["mt"]

[[sections]]
id = "basics"
title = "Basics"
layout = "stack"

[[sections.items]]
type = "command"
command = "my-tool --help"
description = "Show help"

[[sections.items]]
type = "keybind-list"

[[sections.items.entries]]
keys = ["Ctrl+c"]
description = "Cancel"

[[sections.items]]
type = "separator"
label = "Advanced"

[[sections.items]]
type = "text"
text = "Some explanatory text."

[[sections.items]]
type = "link"
label = "Documentation"
url = "https://example.com/docs"
description = "Official docs"

[[sections.items]]
type = "callout"
style = "warning"
text = "Be careful with this."
```

### Item types

| Type | Purpose |
|------|---------|
| `keybind-list` | Key/shortcut reference with aligned columns |
| `command` | Copyable command (navigable with `n`/`N`, yankable with `y`) |
| `table` | Tabular data with columns |
| `callout` | Styled note (styles: note, info, tip, warning, danger) |
| `text` | Plain paragraph |
| `separator` | Horizontal rule with optional label |
| `link` | Labeled URL (accessible via `u` picker) |

## License

MIT
