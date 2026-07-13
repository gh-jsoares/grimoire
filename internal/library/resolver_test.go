package library

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gh-jsoares/grimoire/internal/document"
)

func testLibrary(t *testing.T) *Library {
	t.Helper()
	abs, _ := filepath.Abs("../../testdata/valid/tmux.grim")
	return &Library{
		Documents: []document.Document{
			{Path: abs, Filename: "tmux.grim", Title: "tmux", Aliases: []string{"tm"}},
			{Path: "/fake/git.grim", Filename: "git.grim", Title: "Git", Aliases: []string{"g"}},
		},
	}
}

func TestResolve_ExactFilename(t *testing.T) {
	lib := testLibrary(t)
	doc, err := Resolve("tmux", lib)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Title != "tmux" {
		t.Errorf("got title %q, want %q", doc.Title, "tmux")
	}
}

func TestResolve_ExactTitle(t *testing.T) {
	lib := testLibrary(t)
	doc, err := Resolve("Git", lib)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Filename != "git.grim" {
		t.Errorf("got filename %q, want %q", doc.Filename, "git.grim")
	}
}

func TestResolve_Alias(t *testing.T) {
	lib := testLibrary(t)
	doc, err := Resolve("tm", lib)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Title != "tmux" {
		t.Errorf("got title %q, want %q", doc.Title, "tmux")
	}
}

func TestResolve_CaseInsensitiveFilename(t *testing.T) {
	lib := testLibrary(t)
	doc, err := Resolve("TMUX", lib)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Title != "tmux" {
		t.Errorf("got title %q, want %q", doc.Title, "tmux")
	}
}

func TestResolve_CaseInsensitiveTitle(t *testing.T) {
	lib := testLibrary(t)
	doc, err := Resolve("git", lib)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Title != "Git" {
		t.Errorf("got title %q, want %q", doc.Title, "Git")
	}
}

func TestResolve_CaseInsensitiveAlias(t *testing.T) {
	lib := testLibrary(t)
	doc, err := Resolve("TM", lib)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Title != "tmux" {
		t.Errorf("got title %q, want %q", doc.Title, "tmux")
	}
}

func TestResolve_DirectPath(t *testing.T) {
	lib := testLibrary(t)
	abs, _ := filepath.Abs("../../testdata/valid/tmux.grim")
	doc, err := Resolve(abs, lib)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Title != "tmux" {
		t.Errorf("got title %q, want %q", doc.Title, "tmux")
	}
}

func TestResolve_GrimSuffix(t *testing.T) {
	lib := testLibrary(t)
	_, err := Resolve("nonexistent.grim", lib)
	if err == nil {
		t.Fatal("expected error for nonexistent .grim path")
	}
}

func TestResolve_NotFound(t *testing.T) {
	lib := testLibrary(t)
	_, err := Resolve("nonexistent", lib)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestResolve_DotPrefix(t *testing.T) {
	lib := testLibrary(t)
	// Dot prefix triggers direct path resolution
	_, err := Resolve("./nonexistent", lib)
	if err == nil {
		t.Fatal("expected error for nonexistent dot-prefixed path")
	}
}

func TestResolve_PathSeparator(t *testing.T) {
	lib := testLibrary(t)
	_, err := Resolve("some"+string(os.PathSeparator)+"path", lib)
	if err == nil {
		t.Fatal("expected error for path with separator")
	}
}
