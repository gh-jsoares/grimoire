package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Grid.Columns != 12 {
		t.Errorf("expected 12 columns, got %d", cfg.Grid.Columns)
	}
	if len(cfg.Grid.Breakpoints) != 3 {
		t.Errorf("expected 3 breakpoints, got %d", len(cfg.Grid.Breakpoints))
	}
}

func TestActiveBreakpoint_Large(t *testing.T) {
	gc := DefaultConfig().Grid
	bp := gc.ActiveBreakpoint(160)
	if bp != "lg" {
		t.Errorf("expected lg at width 160, got %q", bp)
	}
}

func TestActiveBreakpoint_Medium(t *testing.T) {
	gc := DefaultConfig().Grid
	bp := gc.ActiveBreakpoint(100)
	if bp != "md" {
		t.Errorf("expected md at width 100, got %q", bp)
	}
}

func TestActiveBreakpoint_Small(t *testing.T) {
	gc := DefaultConfig().Grid
	bp := gc.ActiveBreakpoint(50)
	if bp != "sm" {
		t.Errorf("expected sm at width 50, got %q", bp)
	}
}

func TestActiveBreakpoint_ExactBoundary(t *testing.T) {
	gc := DefaultConfig().Grid
	bp := gc.ActiveBreakpoint(140)
	if bp != "lg" {
		t.Errorf("expected lg at width 140, got %q", bp)
	}
	bp = gc.ActiveBreakpoint(90)
	if bp != "md" {
		t.Errorf("expected md at width 90, got %q", bp)
	}
}

func TestActiveBreakpoint_Fallback(t *testing.T) {
	gc := GridConfig{
		Columns: 12,
		Breakpoints: []Breakpoint{
			{Name: "big", MinWidth: 200},
		},
	}
	bp := gc.ActiveBreakpoint(100)
	if bp != "big" {
		t.Errorf("expected fallback to 'big', got %q", bp)
	}
}

func TestActiveBreakpoint_Empty(t *testing.T) {
	gc := GridConfig{Columns: 12}
	bp := gc.ActiveBreakpoint(100)
	if bp != "" {
		t.Errorf("expected empty with no breakpoints, got %q", bp)
	}
}

func TestLoadConfig_Missing(t *testing.T) {
	cfg := LoadConfig("/nonexistent/path")
	if cfg.Grid.Columns != 12 {
		t.Errorf("expected defaults when file missing, got %d columns", cfg.Grid.Columns)
	}
}

func TestLoadConfig_Valid(t *testing.T) {
	dir := t.TempDir()
	data := []byte(`[grid]
columns = 6

[[grid.breakpoints]]
name = "wide"
min_width = 200

[[grid.breakpoints]]
name = "narrow"
min_width = 0
`)
	if err := os.WriteFile(filepath.Join(dir, "grimoire.toml"), data, 0644); err != nil {
		t.Fatal(err)
	}

	cfg := LoadConfig(dir)
	if cfg.Grid.Columns != 6 {
		t.Errorf("expected 6 columns, got %d", cfg.Grid.Columns)
	}
	if len(cfg.Grid.Breakpoints) != 2 {
		t.Errorf("expected 2 breakpoints, got %d", len(cfg.Grid.Breakpoints))
	}
	if cfg.Grid.Breakpoints[0].Name != "wide" {
		t.Errorf("expected first breakpoint 'wide', got %q", cfg.Grid.Breakpoints[0].Name)
	}
}

func TestLoadConfig_InvalidToml(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "grimoire.toml"), []byte("invalid{{{"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := LoadConfig(dir)
	if cfg.Grid.Columns != 12 {
		t.Errorf("expected defaults on invalid TOML, got %d columns", cfg.Grid.Columns)
	}
}
