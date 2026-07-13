# Grimoire

[![CI](https://github.com/gh-jsoares/grimoire/actions/workflows/ci.yaml/badge.svg)](https://github.com/gh-jsoares/grimoire/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gh-jsoares/grimoire)](https://goreportcard.com/report/github.com/gh-jsoares/grimoire)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/gh-jsoares/grimoire)](https://github.com/gh-jsoares/grimoire/releases/latest)

A terminal cheatsheet viewer built with Go and the [Charm](https://charm.sh) ecosystem. Renders `.grim` TOML files as styled reference cards with responsive multi-column layout.

> **Note:** This project was vibe coded with [Claude](https://claude.ai) (Anthropic). The entire codebase was generated through conversational AI pair programming.

## Features

- Responsive 1/2/3 column layout based on terminal width
- Tab switching between multiple cheatsheets
- Command navigation and clipboard yanking
- Link overlay picker
- Search with filter and highlight modes
- Nerd Font icon support
- Vim-style keybinds

## Install

```sh
go install github.com/gh-jsoares/grimoire/cmd/grimoire@latest
```

See [docs/install.md](docs/install.md) for more installation methods.

## Usage

```sh
grimoire              # Open default library (~/.config/grimoire/)
grimoire tmux.grim    # Open a specific file
grimoire git          # Open by name or alias
grimoire ~/sheets/    # Open a directory as library
```

### Flags

```sh
grimoire --tab git        # Start on a specific tab
grimoire --no-icons       # Disable Nerd Font icons
grimoire --plain          # Minimal styling
grimoire --library ~/dir  # Custom library path
```

### Library Resolution

Grimoire looks for `.grim` files in:

1. `$GRIMOIRE_HOME`
2. `$XDG_CONFIG_HOME/grimoire`
3. `~/.config/grimoire`

## Keybinds

| Key | Action |
|-----|--------|
| `H` / `L` | Switch tabs |
| `j` / `k` | Scroll |
| `n` / `N` | Next / previous command |
| `y` | Yank command |
| `/` | Search |
| `q` | Quit |

See [docs/keybinds.md](docs/keybinds.md) for the full reference.

## File Format

Grimoire reads `.grim` files — TOML with a specific schema. See [docs/file-format.md](docs/file-format.md) for the full specification.

## Documentation

- [Installation](docs/install.md)
- [Keybinds](docs/keybinds.md)
- [File Format](docs/file-format.md)
- [tmux Integration](docs/tmux.md)

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT
