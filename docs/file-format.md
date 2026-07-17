# File Format

Grimoire reads `.grim` files — TOML with a specific schema.

## Structure

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
```

## Document Fields

| Field | Required | Description |
|-------|----------|-------------|
| `format` | Yes | Format version (currently `1`) |
| `title` | Yes | Display title |
| `icon` | No | Nerd Font icon |
| `order` | No | Tab sort order (lower = first) |
| `aliases` | No | Alternative names for resolution |

## Section Fields

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Unique identifier |
| `title` | Yes | Display title |
| `span` | No | Grid columns to occupy (1-12, default: 12 = full width) |
| `span_lg` | No | Span override at `lg` breakpoint (≥140 chars) |
| `span_md` | No | Span override at `md` breakpoint (≥90 chars) |
| `span_sm` | No | Span override at `sm` breakpoint (<90 chars) |
| `layout` | No | `"stack"` (default), `"columns"`, or `"grid"` |

### Responsive Grid

Grimoire uses a 12-column grid system. Sections flow left-to-right into rows; when the next section's span exceeds remaining slots, a new row starts.

```toml
[[sections]]
id = "basics"
title = "Basics"
span = 4          # one-third width on large terminals
span_sm = 12      # full width on small terminals
layout = "stack"
```

Common span values:
- `12` — full width (default)
- `6` — half width
- `4` — one-third
- `3` — one-quarter

Breakpoints and grid columns are configurable via `grimoire.toml` in your library directory:

```toml
[grid]
columns = 12

[[grid.breakpoints]]
name = "lg"
min_width = 140

[[grid.breakpoints]]
name = "md"
min_width = 90

[[grid.breakpoints]]
name = "sm"
min_width = 0
```

## Item Types

### `command`

Copyable command — navigable with `n`/`N`, yankable with `y`.

```toml
[[sections.items]]
type = "command"
command = "my-tool --help"
description = "Show help"
```

### `keybind-list`

Key/shortcut reference with aligned columns.

```toml
[[sections.items]]
type = "keybind-list"

[[sections.items.entries]]
keys = ["Ctrl+c"]
description = "Cancel"

[[sections.items.entries]]
keys = ["Ctrl+d"]
description = "Exit"
```

### `table`

Tabular data with columns.

```toml
[[sections.items]]
type = "table"
columns = ["Flag", "Description"]
rows = [
  ["--verbose", "Enable verbose output"],
  ["--quiet", "Suppress output"],
]
```

### `callout`

Styled note block. Styles: `note`, `info`, `tip`, `warning`, `danger`.

```toml
[[sections.items]]
type = "callout"
style = "warning"
text = "Be careful with this."
```

### `text`

Plain paragraph.

```toml
[[sections.items]]
type = "text"
text = "Some explanatory text."
```

### `separator`

Horizontal rule with optional label.

```toml
[[sections.items]]
type = "separator"
label = "Advanced"
```

### `link`

Labeled URL — accessible via `u` picker.

```toml
[[sections.items]]
type = "link"
label = "Documentation"
url = "https://example.com/docs"
description = "Official docs"
```
