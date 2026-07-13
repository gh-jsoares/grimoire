package content

import (
	"testing"

	"github.com/gh-jsoares/grimoire/internal/document"
)

func TestFilterDocument_EmptyQuery(t *testing.T) {
	doc := &document.Document{Title: "test"}
	result := FilterDocument(doc, "")
	if result != doc {
		t.Error("empty query should return original document")
	}
}

func TestFilterDocument_MatchesCommand(t *testing.T) {
	doc := &document.Document{
		Title: "test",
		Sections: []document.Section{
			{
				ID:    "s1",
				Title: "Section 1",
				Items: []document.Item{
					{Type: "command", Command: "git push", Description: "Push changes"},
					{Type: "command", Command: "docker run", Description: "Run container"},
				},
			},
		},
	}

	result := FilterDocument(doc, "git")
	if len(result.Sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(result.Sections))
	}
	if len(result.Sections[0].Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result.Sections[0].Items))
	}
	if result.Sections[0].Items[0].Command != "git push" {
		t.Errorf("expected 'git push', got %q", result.Sections[0].Items[0].Command)
	}
}

func TestFilterDocument_NoMatch(t *testing.T) {
	doc := &document.Document{
		Title: "test",
		Sections: []document.Section{
			{
				ID:    "s1",
				Title: "Section 1",
				Items: []document.Item{
					{Type: "command", Command: "git push"},
				},
			},
		},
	}

	result := FilterDocument(doc, "nonexistent")
	if len(result.Sections) != 0 {
		t.Errorf("expected 0 sections, got %d", len(result.Sections))
	}
}

func TestFilterDocument_KeybindListFiltersEntries(t *testing.T) {
	doc := &document.Document{
		Title: "test",
		Sections: []document.Section{
			{
				ID:    "s1",
				Title: "Keys",
				Items: []document.Item{
					{
						Type: "keybind-list",
						Entries: []document.KeybindEntry{
							{Keys: []string{"Ctrl+c"}, Description: "Copy"},
							{Keys: []string{"Ctrl+v"}, Description: "Paste"},
							{Keys: []string{"Ctrl+x"}, Description: "Cut"},
						},
					},
				},
			},
		},
	}

	result := FilterDocument(doc, "paste")
	if len(result.Sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(result.Sections))
	}
	item := result.Sections[0].Items[0]
	if len(item.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(item.Entries))
	}
	if item.Entries[0].Description != "Paste" {
		t.Errorf("expected 'Paste', got %q", item.Entries[0].Description)
	}
}

func TestFilterDocument_CaseInsensitive(t *testing.T) {
	doc := &document.Document{
		Title: "test",
		Sections: []document.Section{
			{
				ID:    "s1",
				Title: "Section",
				Items: []document.Item{
					{Type: "command", Command: "Git Push"},
				},
			},
		},
	}

	result := FilterDocument(doc, "GIT")
	if len(result.Sections) != 1 {
		t.Errorf("expected case-insensitive match, got %d sections", len(result.Sections))
	}
}

func TestItemMatches_Command(t *testing.T) {
	item := document.Item{Type: "command", Command: "git push", Description: "Push to remote"}
	if !ItemMatches(item, "push") {
		t.Error("should match command text")
	}
	if !ItemMatches(item, "remote") {
		t.Error("should match description")
	}
	if ItemMatches(item, "docker") {
		t.Error("should not match unrelated query")
	}
}

func TestItemMatches_Text(t *testing.T) {
	item := document.Item{Type: "text", Text: "Hello world"}
	if !ItemMatches(item, "hello") {
		t.Error("should match text content")
	}
	if ItemMatches(item, "goodbye") {
		t.Error("should not match unrelated query")
	}
}

func TestItemMatches_TextFallbackToDescription(t *testing.T) {
	item := document.Item{Type: "text", Description: "Fallback text"}
	if !ItemMatches(item, "fallback") {
		t.Error("should match description when text is empty")
	}
}

func TestItemMatches_Link(t *testing.T) {
	item := document.Item{Type: "link", Label: "Documentation", URL: "https://example.com"}
	if !ItemMatches(item, "doc") {
		t.Error("should match label")
	}
	if !ItemMatches(item, "example") {
		t.Error("should match URL")
	}
}

func TestItemMatches_Table(t *testing.T) {
	item := document.Item{
		Type: "table",
		Rows: []document.TableRow{
			{Values: []string{"--verbose", "Enable verbose output"}},
			{Values: []string{"--quiet", "Suppress output"}},
		},
	}
	if !ItemMatches(item, "verbose") {
		t.Error("should match table row values")
	}
	if ItemMatches(item, "json") {
		t.Error("should not match unrelated query")
	}
}

func TestItemMatches_Callout(t *testing.T) {
	item := document.Item{Type: "callout", Text: "Warning: be careful"}
	if !ItemMatches(item, "careful") {
		t.Error("should match callout text")
	}
}

func TestItemMatches_KeybindList(t *testing.T) {
	item := document.Item{
		Type: "keybind-list",
		Entries: []document.KeybindEntry{
			{Keys: []string{"Ctrl+c"}, Description: "Copy"},
		},
	}
	if !ItemMatches(item, "copy") {
		t.Error("should match keybind description")
	}
	if !ItemMatches(item, "ctrl") {
		t.Error("should match keybind keys")
	}
}

func TestItemMatches_FallbackToTitle(t *testing.T) {
	item := document.Item{Type: "separator", Title: "Advanced Section"}
	if !ItemMatches(item, "advanced") {
		t.Error("should match item title as fallback")
	}
}
