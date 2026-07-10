package library

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gh-jsoares/grimoire/internal/document"
)

func Load(dir string) (*Library, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("library path: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("library path is not a directory: %s", dir)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading library: %w", err)
	}

	lib := &Library{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".grim") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		doc, err := document.Parse(path)
		if err != nil {
			lib.Errors = append(lib.Errors, fmt.Errorf("%s: %w", entry.Name(), err))
			continue
		}

		validationErrs := document.Validate(doc)
		if len(validationErrs) > 0 {
			for _, ve := range validationErrs {
				lib.Errors = append(lib.Errors, ve)
			}
			continue
		}

		lib.Documents = append(lib.Documents, *doc)
	}

	Sort(lib.Documents)
	return lib, nil
}

func LoadSingle(path string) (*Library, error) {
	doc, err := document.Parse(path)
	if err != nil {
		return nil, err
	}

	validationErrs := document.Validate(doc)
	if len(validationErrs) > 0 {
		var errs []error
		for _, ve := range validationErrs {
			errs = append(errs, ve)
		}
		return &Library{Errors: errs}, fmt.Errorf("validation failed for %s", path)
	}

	return &Library{Documents: []document.Document{*doc}}, nil
}
