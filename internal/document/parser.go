package document

import (
	"fmt"
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

type rawDocument struct {
	Format   int      `toml:"format"`
	Title    string   `toml:"title"`
	Icon     string   `toml:"icon"`
	Order    *int     `toml:"order"`
	Aliases  []string `toml:"aliases"`
	Hidden   bool     `toml:"hidden"`
	Sections []rawSection `toml:"sections"`
}

type rawSection struct {
	ID      string      `toml:"id"`
	Title   string      `toml:"title"`
	Icon    string      `toml:"icon"`
	Order   *int        `toml:"order"`
	Layout  string      `toml:"layout"`
	Items   []rawItem   `toml:"items"`
	Columns []rawColumn `toml:"columns"`
}

type rawColumn struct {
	Span     int       `toml:"span"`
	MinWidth int       `toml:"min_width"`
	MaxWidth int       `toml:"max_width"`
	Items    []rawItem `toml:"items"`
}

type rawItem struct {
	Type        string `toml:"type"`
	Title       string `toml:"title"`
	Description string `toml:"description"`

	// keybind-list
	Entries []rawKeybindEntry `toml:"entries"`

	// command
	Command  string      `toml:"command"`
	Language string      `toml:"language"`
	Copy     interface{} `toml:"copy"`

	// table
	Columns      []string        `toml:"columns"`
	Rows         []rawTableRow   `toml:"rows"`
	ColumnConfig []rawColConfig  `toml:"column_config"`

	// callout
	Style string `toml:"style"`
	Text  string `toml:"text"`

	// cards
	CardColumns int       `toml:"columns_count"`
	Cards       []rawCard `toml:"cards"`

	// link
	Label string `toml:"label"`
	URL   string `toml:"url"`
}

type rawKeybindEntry struct {
	Keys        []string `toml:"keys"`
	Description string   `toml:"description"`
	Command     string   `toml:"command"`
	Tags        []string `toml:"tags"`
}

type rawTableRow struct {
	Values []string `toml:"values"`
	Copy   string   `toml:"copy"`
}

type rawColConfig struct {
	Name     string      `toml:"name"`
	Width    interface{} `toml:"width"`
	Priority int         `toml:"priority"`
	Align    string      `toml:"align"`
}

type rawCard struct {
	Title       string `toml:"title"`
	Value       string `toml:"value"`
	Description string `toml:"description"`
	Copy        string `toml:"copy"`
}

func Parse(path string) (*Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var raw rawDocument
	if err := toml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	doc := &Document{
		Path:     path,
		Filename: filepath.Base(path),
		Format:   raw.Format,
		Title:    raw.Title,
		Icon:     raw.Icon,
		Order:    raw.Order,
		Aliases:  raw.Aliases,
		Hidden:   raw.Hidden,
	}

	if doc.Title == "" {
		name := filepath.Base(path)
		ext := filepath.Ext(name)
		doc.Title = name[:len(name)-len(ext)]
	}

	for _, rs := range raw.Sections {
		section := Section{
			ID:     rs.ID,
			Title:  rs.Title,
			Icon:   rs.Icon,
			Order:  rs.Order,
			Layout: rs.Layout,
		}
		if section.Layout == "" {
			section.Layout = "stack"
		}

		for _, rc := range rs.Columns {
			col := Column{
				Span:     rc.Span,
				MinWidth: rc.MinWidth,
				MaxWidth: rc.MaxWidth,
			}
			if col.Span == 0 {
				col.Span = 1
			}
			for _, ri := range rc.Items {
				col.Items = append(col.Items, convertItem(ri))
			}
			section.Columns = append(section.Columns, col)
		}

		for _, ri := range rs.Items {
			section.Items = append(section.Items, convertItem(ri))
		}

		doc.Sections = append(doc.Sections, section)
	}

	return doc, nil
}

func convertItem(ri rawItem) Item {
	item := Item{
		Type:        ri.Type,
		Title:       ri.Title,
		Description: ri.Description,
		Command:     ri.Command,
		Language:    ri.Language,
		Copy:        ri.Copy,
		Style:       ri.Style,
		Text:        ri.Text,
		Label:       ri.Label,
		URL:         ri.URL,
		CardColumns: ri.CardColumns,
	}

	// Handle "columns" field collision: for cards it's columns_count, for table it's the columns list
	if ri.Type == "table" {
		item.TableColumns = ri.Columns
	}

	for _, e := range ri.Entries {
		item.Entries = append(item.Entries, KeybindEntry(e))
	}

	for _, r := range ri.Rows {
		item.Rows = append(item.Rows, TableRow(r))
	}

	for _, cc := range ri.ColumnConfig {
		item.ColumnConfig = append(item.ColumnConfig, ColumnConfig(cc))
	}

	for _, c := range ri.Cards {
		item.Cards = append(item.Cards, Card(c))
	}

	return item
}
