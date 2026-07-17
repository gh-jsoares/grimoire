# Configuration

Grimoire supports an optional `grimoire.toml` config file in your library directory (the same directory where `.grim` files live).

## Location

The config file is loaded from:

1. `$GRIMOIRE_HOME/grimoire.toml`
2. `$XDG_CONFIG_HOME/grimoire/grimoire.toml`
3. `~/.config/grimoire/grimoire.toml`

If no config file exists, sensible defaults are used.

## Grid System

Grimoire uses a 12-column responsive grid to lay out sections. Configure it with:

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

### Fields

| Field | Default | Description |
|-------|---------|-------------|
| `grid.columns` | `12` | Total grid columns |
| `grid.breakpoints[].name` | — | Breakpoint identifier (matches `span_<name>` in sections) |
| `grid.breakpoints[].min_width` | — | Minimum terminal width (chars) for this breakpoint to activate |

### How breakpoints work

The active breakpoint is the one with the highest `min_width` that is less than or equal to the current terminal width. For example, at 120 characters wide, the `md` breakpoint (min_width=90) is active.

Sections declare their span per breakpoint:

```toml
[[sections]]
id = "example"
title = "Example"
span = 6         # default: half width at all breakpoints
span_lg = 4      # override: one-third at lg (>=140 chars)
span_sm = 12     # override: full width at sm (<90 chars)
```

### Layout flow

Sections flow left-to-right into rows. When the next section's span exceeds remaining grid slots in the current row, a new row starts. Width is allocated proportionally based on span values within each row.

### Common patterns

| Span | Width | Use case |
|------|-------|----------|
| `12` | 100% | Full-width sections (default if no span set) |
| `6` | 50% | Two sections side by side |
| `4` | 33% | Three sections per row |
| `3` | 25% | Four sections per row |

### Defaults

If no `grimoire.toml` exists:
- Grid uses 12 columns
- Three breakpoints: `lg` (>=140), `md` (>=90), `sm` (>=0)
- Sections without a `span` field default to full width (12)
