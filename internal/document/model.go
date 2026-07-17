// Package document defines the data model for .grim cheatsheet files.
package document

// Document represents a parsed .grim cheatsheet file.
type Document struct {
	Path     string
	Filename string
	Format   int
	Title    string
	Icon     string
	Order    *int
	Aliases  []string
	Hidden   bool
	Sections []Section
}

// Section is a named group of items within a document.
type Section struct {
	ID      string
	Title   string
	Icon    string
	Order   *int
	Span    map[string]int // keyed by breakpoint name, e.g. {"lg": 6, "md": 8, "sm": 12}
	Layout  string         // "stack" | "columns" | "grid"
	Items   []Item
	Columns []Column
}

type Column struct {
	Span     int
	MinWidth int
	MaxWidth int
	Items    []Item
}

// Item is a single content element within a section (command, keybind-list, table, etc.).
type Item struct {
	Type string // "keybind-list" | "command" | "table" | "callout" | "cards" | "text" | "separator" | "link"

	// Common
	Title string

	// keybind-list
	Entries []KeybindEntry

	// command
	Command     string
	Language    string
	Description string
	Copy        interface{} // bool or string

	// table
	TableColumns []string
	Rows         []TableRow
	ColumnConfig []ColumnConfig

	// callout
	Style string
	Text  string

	// cards
	CardColumns int
	Cards       []Card

	// link
	Label string
	URL   string

	// separator
	SeparatorLabel string
}

type KeybindEntry struct {
	Keys        []string
	Description string
	Command     string
	Tags        []string
}

type TableRow struct {
	Values []string
	Copy   string
}

type ColumnConfig struct {
	Name     string
	Width    interface{} // int or "fill"
	Priority int
	Align    string
}

type Card struct {
	Title       string
	Value       string
	Description string
	Copy        string
}
