// Package config handles runtime configuration and path resolution.
package config

import (
	"os"
	"path/filepath"
)

// ResolveLibraryPath returns the library directory from $GRIMOIRE_HOME, $XDG_CONFIG_HOME/grimoire, or ~/.config/grimoire.
func ResolveLibraryPath() string {
	if p := os.Getenv("GRIMOIRE_HOME"); p != "" {
		return p
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "grimoire")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "grimoire")
}
