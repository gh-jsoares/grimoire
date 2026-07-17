package app

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gh-jsoares/grimoire/internal/config"
	"github.com/gh-jsoares/grimoire/internal/library"
)

func testLibrary(t *testing.T) *library.Library {
	t.Helper()
	lib, err := library.Load("../../testdata/valid")
	if err != nil {
		t.Fatalf("failed to load test library: %v", err)
	}
	if len(lib.Documents) == 0 {
		t.Fatal("no documents loaded")
	}
	return lib
}

func TestNewModel(t *testing.T) {
	lib := testLibrary(t)
	m := New(lib, Config{}, config.DefaultConfig().Grid)
	if len(m.Tabs.Tabs) != 2 {
		t.Errorf("expected 2 tabs, got %d", len(m.Tabs.Tabs))
	}
}

func TestViewAfterResize(t *testing.T) {
	lib := testLibrary(t)
	m := New(lib, Config{}, config.DefaultConfig().Grid)

	result, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
	model := result.(Model)
	view := model.View()

	if !strings.Contains(view, "Navigation") {
		t.Error("view missing 'Navigation' section heading")
	}
	if !strings.Contains(view, "Copy Mode") {
		t.Error("view missing 'Copy Mode' section heading")
	}
	if !strings.Contains(view, "tmux") {
		t.Error("view missing 'tmux' tab")
	}
}

func TestTabSwitch(t *testing.T) {
	lib := testLibrary(t)
	m := New(lib, Config{}, config.DefaultConfig().Grid)

	result, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
	result, _ = result.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'L'}})
	model := result.(Model)

	if model.ActiveDoc != 1 {
		t.Errorf("after L: doc = %d, want 1", model.ActiveDoc)
	}

	view := model.View()
	if !strings.Contains(view, "Git") {
		t.Error("view should show Git content after tab switch")
	}
}

func TestScrollClamped(t *testing.T) {
	lib := testLibrary(t)
	m := New(lib, Config{}, config.DefaultConfig().Grid)

	// At height=80 with content that fits, scroll should be clamped to 0
	result, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 80})
	result, _ = result.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	model := result.(Model)

	if model.scrollOffsets[0] != 0 {
		t.Errorf("scroll should be 0 when content fits, got %d", model.scrollOffsets[0])
	}
}

func TestScrollWhenNeeded(t *testing.T) {
	lib := testLibrary(t)
	m := New(lib, Config{}, config.DefaultConfig().Grid)

	// At height=5, content should overflow and scroll should work
	result, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 5})
	result, _ = result.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	model := result.(Model)

	if model.scrollOffsets[0] == 0 {
		t.Error("scroll should move when content overflows")
	}
}

func TestNarrowLayout(t *testing.T) {
	lib := testLibrary(t)
	m := New(lib, Config{}, config.DefaultConfig().Grid)

	result, _ := m.Update(tea.WindowSizeMsg{Width: 60, Height: 30})
	model := result.(Model)
	view := model.View()

	if view == "" {
		t.Error("view should render at narrow width")
	}
}
