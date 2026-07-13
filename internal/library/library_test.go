package library

import (
	"testing"

	"github.com/gh-jsoares/grimoire/internal/document"
)

func intPtr(i int) *int { return &i }

func TestSort_ByOrder(t *testing.T) {
	docs := []document.Document{
		{Title: "C", Order: intPtr(3)},
		{Title: "A", Order: intPtr(1)},
		{Title: "B", Order: intPtr(2)},
	}
	Sort(docs)
	if docs[0].Title != "A" || docs[1].Title != "B" || docs[2].Title != "C" {
		t.Errorf("got %s %s %s, want A B C", docs[0].Title, docs[1].Title, docs[2].Title)
	}
}

func TestSort_OrderBeforeNoOrder(t *testing.T) {
	docs := []document.Document{
		{Title: "Z"},
		{Title: "A", Order: intPtr(1)},
	}
	Sort(docs)
	if docs[0].Title != "A" {
		t.Errorf("ordered doc should come first, got %q", docs[0].Title)
	}
}

func TestSort_AlphabeticalWhenNoOrder(t *testing.T) {
	docs := []document.Document{
		{Title: "Zsh"},
		{Title: "Bash"},
		{Title: "Fish"},
	}
	Sort(docs)
	if docs[0].Title != "Bash" || docs[1].Title != "Fish" || docs[2].Title != "Zsh" {
		t.Errorf("got %s %s %s, want Bash Fish Zsh", docs[0].Title, docs[1].Title, docs[2].Title)
	}
}

func TestSort_SameOrderAlphabetical(t *testing.T) {
	docs := []document.Document{
		{Title: "B", Order: intPtr(1)},
		{Title: "A", Order: intPtr(1)},
	}
	Sort(docs)
	if docs[0].Title != "A" {
		t.Errorf("same order should sort alphabetically, got %q first", docs[0].Title)
	}
}

func TestSort_Empty(t *testing.T) {
	docs := []document.Document{}
	Sort(docs)
}

func TestSort_NilOrder(t *testing.T) {
	docs := []document.Document{
		{Title: "No order 1"},
		{Title: "Has order", Order: intPtr(5)},
		{Title: "No order 2"},
	}
	Sort(docs)
	if docs[0].Title != "Has order" {
		t.Errorf("expected 'Has order' first, got %q", docs[0].Title)
	}
}
