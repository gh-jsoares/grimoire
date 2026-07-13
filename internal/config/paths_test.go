package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveLibraryPath_GrimoireHome(t *testing.T) {
	t.Setenv("GRIMOIRE_HOME", "/custom/path")
	t.Setenv("XDG_CONFIG_HOME", "/xdg")

	got := ResolveLibraryPath()
	if got != "/custom/path" {
		t.Errorf("got %q, want /custom/path", got)
	}
}

func TestResolveLibraryPath_XDG(t *testing.T) {
	t.Setenv("GRIMOIRE_HOME", "")
	t.Setenv("XDG_CONFIG_HOME", "/xdg")

	got := ResolveLibraryPath()
	if got != "/xdg/grimoire" {
		t.Errorf("got %q, want /xdg/grimoire", got)
	}
}

func TestResolveLibraryPath_Default(t *testing.T) {
	t.Setenv("GRIMOIRE_HOME", "")
	t.Setenv("XDG_CONFIG_HOME", "")

	got := ResolveLibraryPath()
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, ".config", "grimoire")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
