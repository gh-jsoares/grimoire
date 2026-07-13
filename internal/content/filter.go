package content

import (
	"strings"

	"github.com/gh-jsoares/grimoire/internal/document"
)

// FilterDocument returns a copy of doc containing only items that match the query.
func FilterDocument(doc *document.Document, query string) *document.Document {
	if query == "" {
		return doc
	}

	q := strings.ToLower(query)

	filtered := &document.Document{
		Path:     doc.Path,
		Filename: doc.Filename,
		Format:   doc.Format,
		Title:    doc.Title,
		Icon:     doc.Icon,
		Order:    doc.Order,
		Aliases:  doc.Aliases,
		Hidden:   doc.Hidden,
	}

	for _, sec := range doc.Sections {
		var matchedItems []document.Item

		items := sec.Items
		if len(items) == 0 {
			for _, col := range sec.Columns {
				items = append(items, col.Items...)
			}
		}

		for _, item := range items {
			if item.Type == "keybind-list" {
				// Filter individual entries within keybind-list
				var matchedEntries []document.KeybindEntry
				for _, e := range item.Entries {
					if strings.Contains(strings.ToLower(e.Description), q) ||
						strings.Contains(strings.ToLower(strings.Join(e.Keys, " ")), q) {
						matchedEntries = append(matchedEntries, e)
					}
				}
				if len(matchedEntries) > 0 {
					filtered := item
					filtered.Entries = matchedEntries
					matchedItems = append(matchedItems, filtered)
				}
			} else if ItemMatches(item, q) {
				matchedItems = append(matchedItems, item)
			}
		}

		if len(matchedItems) > 0 {
			newSec := document.Section{
				ID:     sec.ID,
				Title:  sec.Title,
				Icon:   sec.Icon,
				Order:  sec.Order,
				Layout: sec.Layout,
				Items:  matchedItems,
			}
			filtered.Sections = append(filtered.Sections, newSec)
		}
	}

	return filtered
}

// ItemMatches reports whether an item's content matches the lowercase query string.
func ItemMatches(item document.Item, q string) bool {
	switch item.Type {
	case "keybind-list":
		for _, e := range item.Entries {
			if strings.Contains(strings.ToLower(e.Description), q) {
				return true
			}
			if strings.Contains(strings.ToLower(strings.Join(e.Keys, " ")), q) {
				return true
			}
		}
	case "command":
		if strings.Contains(strings.ToLower(item.Command), q) {
			return true
		}
		if strings.Contains(strings.ToLower(item.Description), q) {
			return true
		}
	case "table":
		for _, row := range item.Rows {
			for _, val := range row.Values {
				if strings.Contains(strings.ToLower(val), q) {
					return true
				}
			}
		}
	case "text":
		text := item.Text
		if text == "" {
			text = item.Description
		}
		if strings.Contains(strings.ToLower(text), q) {
			return true
		}
	case "callout":
		text := item.Text
		if text == "" {
			text = item.Description
		}
		if strings.Contains(strings.ToLower(text), q) {
			return true
		}
	case "link":
		if strings.Contains(strings.ToLower(item.Label), q) {
			return true
		}
		if strings.Contains(strings.ToLower(item.URL), q) {
			return true
		}
	}

	return strings.Contains(strings.ToLower(item.Title), q)
}
