package config

import (
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

// Breakpoint defines a named responsive breakpoint with a minimum terminal width.
type Breakpoint struct {
	Name     string `toml:"name"`
	MinWidth int    `toml:"min_width"`
}

// GridConfig holds the 12-column grid settings.
type GridConfig struct {
	Columns     int          `toml:"columns"`
	Breakpoints []Breakpoint `toml:"breakpoints"`
}

// Config is the top-level configuration loaded from grimoire.toml.
type Config struct {
	Grid GridConfig `toml:"grid"`
}

// DefaultConfig returns sensible defaults for the grid system.
func DefaultConfig() Config {
	return Config{
		Grid: GridConfig{
			Columns: 12,
			Breakpoints: []Breakpoint{
				{Name: "lg", MinWidth: 140},
				{Name: "md", MinWidth: 90},
				{Name: "sm", MinWidth: 0},
			},
		},
	}
}

// LoadConfig reads grimoire.toml from the given library path.
// If the file does not exist or cannot be parsed, defaults are returned.
func LoadConfig(libraryPath string) Config {
	cfg := DefaultConfig()
	if libraryPath == "" {
		return cfg
	}

	data, err := os.ReadFile(filepath.Join(libraryPath, "grimoire.toml"))
	if err != nil {
		return cfg
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig()
	}

	// Ensure sane defaults if config is partially specified
	if cfg.Grid.Columns <= 0 {
		cfg.Grid.Columns = 12
	}
	if len(cfg.Grid.Breakpoints) == 0 {
		cfg.Grid.Breakpoints = DefaultConfig().Grid.Breakpoints
	}

	return cfg
}

// ActiveBreakpoint returns the name of the breakpoint that applies for the given width.
// It picks the breakpoint with the highest min_width that is <= width.
// If no breakpoint matches (all min_widths exceed width), falls back to the smallest.
func (gc GridConfig) ActiveBreakpoint(width int) string {
	best := ""
	bestMin := -1
	for _, bp := range gc.Breakpoints {
		if width >= bp.MinWidth && bp.MinWidth > bestMin {
			best = bp.Name
			bestMin = bp.MinWidth
		}
	}
	if best != "" || len(gc.Breakpoints) == 0 {
		return best
	}

	// Fallback: use the breakpoint with the lowest min_width
	best = gc.Breakpoints[0].Name
	lowest := gc.Breakpoints[0].MinWidth
	for _, bp := range gc.Breakpoints[1:] {
		if bp.MinWidth < lowest {
			best = bp.Name
			lowest = bp.MinWidth
		}
	}
	return best
}
