package config

import (
	"os"
	"path/filepath"
)

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
