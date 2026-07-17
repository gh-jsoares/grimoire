package content

import (
	"testing"

	"github.com/gh-jsoares/grimoire/internal/document"
)

func TestStripAnsi_Plain(t *testing.T) {
	got := stripAnsi("hello world")
	if got != "hello world" {
		t.Errorf("got %q, want %q", got, "hello world")
	}
}

func TestStripAnsi_CSI(t *testing.T) {
	got := stripAnsi("\x1b[31mred\x1b[0m")
	if got != "red" {
		t.Errorf("got %q, want %q", got, "red")
	}
}

func TestStripAnsi_Multiple(t *testing.T) {
	got := stripAnsi("\x1b[1m\x1b[34mbold blue\x1b[0m plain")
	if got != "bold blue plain" {
		t.Errorf("got %q, want %q", got, "bold blue plain")
	}
}

func TestStripAnsi_Empty(t *testing.T) {
	got := stripAnsi("")
	if got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestComputeTableWidths_Basic(t *testing.T) {
	cols := []string{"Key", "Value"}
	rows := []document.TableRow{
		{Values: []string{"abc", "def"}},
		{Values: []string{"longer", "x"}},
	}
	widths := computeTableWidths(cols, rows, 80)
	if len(widths) != 2 {
		t.Fatalf("expected 2 widths, got %d", len(widths))
	}
	// "longer" = 6 chars + 4 padding = 10
	if widths[0] != 10 {
		t.Errorf("first col width = %d, want 10", widths[0])
	}
}

func TestComputeTableWidths_Overflow(t *testing.T) {
	cols := []string{"A", "B"}
	rows := []document.TableRow{
		{Values: []string{"a very long value that exceeds width", "another long value here too"}},
	}
	widths := computeTableWidths(cols, rows, 30)
	total := 0
	for _, w := range widths {
		total += w
		if w < 4 {
			t.Errorf("column width %d below minimum 4", w)
		}
	}
	if total > 30 {
		t.Errorf("total width %d exceeds available %d", total, 30)
	}
}

func TestComputeTableWidths_NoCols(t *testing.T) {
	widths := computeTableWidths(nil, nil, 80)
	if widths != nil {
		t.Errorf("expected nil, got %v", widths)
	}
}

func TestComputeTableWidths_InferFromRows(t *testing.T) {
	rows := []document.TableRow{
		{Values: []string{"a", "b", "c"}},
	}
	widths := computeTableWidths(nil, rows, 80)
	if len(widths) != 3 {
		t.Errorf("expected 3 widths inferred from rows, got %d", len(widths))
	}
}
