package document

import "fmt"

type ValidationError struct {
	File    string
	Section string
	Item    int
	Message string
}

func (e ValidationError) Error() string {
	if e.Section != "" && e.Item >= 0 {
		return fmt.Sprintf("%s: section %q, item %d: %s", e.File, e.Section, e.Item+1, e.Message)
	}
	if e.Section != "" {
		return fmt.Sprintf("%s: section %q: %s", e.File, e.Section, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.File, e.Message)
}

func Validate(doc *Document) []ValidationError {
	var errs []ValidationError
	add := func(section string, item int, msg string) {
		errs = append(errs, ValidationError{
			File:    doc.Filename,
			Section: section,
			Item:    item,
			Message: msg,
		})
	}

	if doc.Format != 1 {
		add("", -1, fmt.Sprintf("unsupported format version: %d (expected 1)", doc.Format))
	}

	sectionIDs := make(map[string]bool)
	for _, s := range doc.Sections {
		if s.ID == "" {
			add(s.Title, -1, "section missing required field: id")
			continue
		}
		if sectionIDs[s.ID] {
			add(s.ID, -1, "duplicate section id")
		}
		sectionIDs[s.ID] = true

		if s.Title == "" {
			add(s.ID, -1, "section missing required field: title")
		}

		switch s.Layout {
		case "stack", "columns", "grid":
		default:
			add(s.ID, -1, fmt.Sprintf("invalid layout: %q (expected stack, columns, or grid)", s.Layout))
		}

		if s.Layout == "columns" {
			for _, col := range s.Columns {
				if col.Span < 1 {
					add(s.ID, -1, "column span must be positive")
				}
				for i, item := range col.Items {
					validateItem(item, s.ID, i, &errs, doc)
				}
			}
		}

		for i, item := range s.Items {
			validateItem(item, s.ID, i, &errs, doc)
		}
	}

	return errs
}

func validateItem(item Item, sectionID string, idx int, errs *[]ValidationError, doc *Document) {
	add := func(msg string) {
		*errs = append(*errs, ValidationError{
			File:    doc.Filename,
			Section: sectionID,
			Item:    idx,
			Message: msg,
		})
	}

	switch item.Type {
	case "keybind-list":
		for _, e := range item.Entries {
			if len(e.Keys) == 0 {
				add("keybind-list entry requires at least one key")
			}
		}
	case "command":
		if item.Command == "" {
			add("command item requires a command field")
		}
	case "table":
		colCount := len(item.TableColumns)
		for _, row := range item.Rows {
			if colCount > 0 && len(row.Values) != colCount {
				add(fmt.Sprintf("table row has %d values, expected %d columns", len(row.Values), colCount))
			}
		}
	case "callout":
		switch item.Style {
		case "note", "info", "tip", "warning", "danger", "":
		default:
			add(fmt.Sprintf("invalid callout style: %q", item.Style))
		}
	case "cards", "text", "separator", "link":
		// minimal validation for M1
	case "":
		add("item missing required field: type")
	default:
		add(fmt.Sprintf("unsupported item type: %q", item.Type))
	}
}
