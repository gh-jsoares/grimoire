package library

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gh-jsoares/grimoire/internal/document"
)

// Resolve finds a document by path, filename, title, or alias (with case-insensitive fallback).
func Resolve(arg string, lib *Library) (*document.Document, error) {
	// Direct path
	if strings.HasSuffix(arg, ".grim") || strings.Contains(arg, string(os.PathSeparator)) || strings.HasPrefix(arg, ".") {
		abs, err := filepath.Abs(arg)
		if err != nil {
			return nil, err
		}
		if _, err := os.Stat(abs); err == nil {
			for i := range lib.Documents {
				if lib.Documents[i].Path == abs {
					return &lib.Documents[i], nil
				}
			}
		}
		return nil, fmt.Errorf("document not found: %s", arg)
	}

	lower := strings.ToLower(arg)

	// Exact filename without extension
	for i := range lib.Documents {
		name := lib.Documents[i].Filename
		ext := filepath.Ext(name)
		if name[:len(name)-len(ext)] == arg {
			return &lib.Documents[i], nil
		}
	}

	// Exact title
	for i := range lib.Documents {
		if lib.Documents[i].Title == arg {
			return &lib.Documents[i], nil
		}
	}

	// Alias
	for i := range lib.Documents {
		for _, alias := range lib.Documents[i].Aliases {
			if alias == arg {
				return &lib.Documents[i], nil
			}
		}
	}

	// Case-insensitive filename
	for i := range lib.Documents {
		name := lib.Documents[i].Filename
		ext := filepath.Ext(name)
		if strings.ToLower(name[:len(name)-len(ext)]) == lower {
			return &lib.Documents[i], nil
		}
	}

	// Case-insensitive title
	for i := range lib.Documents {
		if strings.ToLower(lib.Documents[i].Title) == lower {
			return &lib.Documents[i], nil
		}
	}

	// Case-insensitive alias
	for i := range lib.Documents {
		for _, alias := range lib.Documents[i].Aliases {
			if strings.ToLower(alias) == lower {
				return &lib.Documents[i], nil
			}
		}
	}

	return nil, fmt.Errorf("no document matches %q", arg)
}
