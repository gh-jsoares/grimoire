package document

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

type Section struct {
	ID      string
	Title   string
	Icon    string
	Order   *int
	Layout  string // "stack" | "columns" | "grid"
	Items   []Item
	Columns []Column
}

type Column struct {
	Span     int
	MinWidth int
	MaxWidth int
	Items    []Item
}

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
