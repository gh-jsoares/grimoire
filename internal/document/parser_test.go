package document

import (
	"path/filepath"
	"testing"
)

func TestParseTmux(t *testing.T) {
	path, _ := filepath.Abs("../../testdata/valid/tmux.grim")
	doc, err := Parse(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if doc.Format != 1 {
		t.Errorf("format = %d, want 1", doc.Format)
	}
	if doc.Title != "tmux" {
		t.Errorf("title = %q, want %q", doc.Title, "tmux")
	}
	if doc.Order == nil || *doc.Order != 10 {
		t.Errorf("order = %v, want 10", doc.Order)
	}
	if len(doc.Aliases) != 2 {
		t.Errorf("aliases = %v, want 2 items", doc.Aliases)
	}
	if len(doc.Sections) != 6 {
		t.Fatalf("sections = %d, want 6", len(doc.Sections))
	}

	s := doc.Sections[0]
	if s.ID != "navigation" {
		t.Errorf("section[0].id = %q, want %q", s.ID, "navigation")
	}
	if s.Layout != "stack" {
		t.Errorf("section[0].layout = %q, want %q", s.Layout, "stack")
	}
	if len(s.Items) != 1 {
		t.Fatalf("section[0].items = %d, want 1", len(s.Items))
	}
	if s.Items[0].Type != "keybind-list" {
		t.Errorf("item type = %q, want keybind-list", s.Items[0].Type)
	}
	if len(s.Items[0].Entries) != 3 {
		t.Errorf("entries = %d, want 3", len(s.Items[0].Entries))
	}
}

func TestParseGit(t *testing.T) {
	path, _ := filepath.Abs("../../testdata/valid/git.grim")
	doc, err := Parse(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if doc.Title != "Git" {
		t.Errorf("title = %q, want %q", doc.Title, "Git")
	}
	if len(doc.Sections) != 2 {
		t.Fatalf("sections = %d, want 2", len(doc.Sections))
	}

	// Check commands in branching section
	s := doc.Sections[1]
	if s.ID != "branching" {
		t.Errorf("section[1].id = %q, want %q", s.ID, "branching")
	}
	cmdCount := 0
	for _, item := range s.Items {
		if item.Type == "command" {
			cmdCount++
		}
	}
	if cmdCount < 3 {
		t.Errorf("branching section commands = %d, want >= 3", cmdCount)
	}
}

func TestValidationDuplicateID(t *testing.T) {
	doc := &Document{
		Filename: "test.grim",
		Format:   1,
		Sections: []Section{
			{ID: "a", Title: "A", Layout: "stack"},
			{ID: "a", Title: "B", Layout: "stack"},
		},
	}
	errs := Validate(doc)
	found := false
	for _, e := range errs {
		if e.Message == "duplicate section id" {
			found = true
		}
	}
	if !found {
		t.Error("expected duplicate section id error")
	}
}

func TestValidationBadFormat(t *testing.T) {
	doc := &Document{
		Filename: "test.grim",
		Format:   0,
		Sections: []Section{{ID: "x", Title: "X", Layout: "stack"}},
	}
	errs := Validate(doc)
	if len(errs) == 0 {
		t.Fatal("expected validation error for bad format")
	}
}
