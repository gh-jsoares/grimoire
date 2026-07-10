package content

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gh-jsoares/grimoire/internal/document"
	"github.com/gh-jsoares/grimoire/internal/theme"
)

func RenderDocument(doc *document.Document, width, height int, t theme.Theme, activeCommand int, highlightQuery ...string) string {
	if len(doc.Sections) == 0 {
		noResults := lipgloss.NewStyle().
			Foreground(t.Muted).
			Width(width).
			Align(lipgloss.Center).
			Render("No results")
		return "\n\n" + noResults
	}

	query := ""
	if len(highlightQuery) > 0 {
		query = highlightQuery[0]
	}

	cols := columnCount(width)
	gap := 3
	dividerWidth := 1
	colWidth := (width - (cols-1)*(gap*2+dividerWidth)) / cols

	cmdCounter := 0

	var sectionBlocks []string
	for _, sec := range doc.Sections {
		block := renderSectionBlock(sec, colWidth, t, activeCommand, &cmdCounter, query)
		sectionBlocks = append(sectionBlocks, block)
	}

	if cols == 1 {
		result := ""
		for i, block := range sectionBlocks {
			if i > 0 {
				result += "\n\n"
			}
			result += block
		}
		return result
	}

	divider := lipgloss.NewStyle().Foreground(t.Border).Render("│")
	spacer := strings.Repeat(" ", gap)

	var rows []string
	for i := 0; i < len(sectionBlocks); i += cols {
		var rowParts []string
		for c := 0; c < cols; c++ {
			if c > 0 {
				rowParts = append(rowParts, spacer+divider+spacer)
			}
			var content string
			if i+c < len(sectionBlocks) {
				content = sectionBlocks[i+c]
			}
			styled := lipgloss.NewStyle().Width(colWidth).Render(content)
			rowParts = append(rowParts, styled)
		}
		row := lipgloss.JoinHorizontal(lipgloss.Top, rowParts...)
		rows = append(rows, row)
	}

	return strings.Join(rows, "\n\n")
}

func columnCount(width int) int {
	if width >= 140 {
		return 3
	}
	if width >= 70 {
		return 2
	}
	return 1
}

func renderSectionBlock(sec document.Section, width int, t theme.Theme, activeCommand int, cmdCounter *int, query string) string {
	var lines []string

	// Check if any item in section matches (for dimming the heading)
	sectionHasMatch := query == ""
	if query != "" {
		q := strings.ToLower(query)
		items := sec.Items
		if len(items) == 0 {
			for _, col := range sec.Columns {
				items = append(items, col.Items...)
			}
		}
		for _, item := range items {
			if ItemMatches(item, q) {
				sectionHasMatch = true
				break
			}
		}
	}

	if sectionHasMatch {
		titleStyle := lipgloss.NewStyle().Bold(true).Foreground(t.Accent)
		lines = append(lines, titleStyle.Render(sec.Title))
		underline := lipgloss.NewStyle().Foreground(t.Accent).Render(strings.Repeat("─", width))
		lines = append(lines, underline)
	} else {
		dimStyle := lipgloss.NewStyle().Foreground(t.Border)
		lines = append(lines, dimStyle.Render(sec.Title))
		lines = append(lines, dimStyle.Render(strings.Repeat("─", width)))
	}
	lines = append(lines, "")

	items := sec.Items
	if len(items) == 0 {
		for _, col := range sec.Columns {
			items = append(items, col.Items...)
		}
	}

	for idx, item := range items {
		var block string
		if item.Type == "command" {
			isActive := *cmdCounter == activeCommand
			block = renderCommandBlock(item, width, t, isActive)
			*cmdCounter++
			if query != "" && !ItemMatches(item, strings.ToLower(query)) {
				block = dimBlock(block, t)
			}
		} else if item.Type == "keybind-list" && query != "" {
			block = renderKeybindBlockHighlighted(item, width, t, query)
		} else {
			block = renderItemBlock(item, width, t)
			if query != "" && !ItemMatches(item, strings.ToLower(query)) {
				block = dimBlock(block, t)
			}
		}
		if block != "" {
			if idx > 0 {
				lines = append(lines, "")
			}
			lines = append(lines, block)
		}
	}

	return strings.Join(lines, "\n")
}

func dimBlock(block string, t theme.Theme) string {
	dimStyle := lipgloss.NewStyle().Foreground(t.Border)
	var dimmed []string
	for _, line := range strings.Split(block, "\n") {
		stripped := stripAnsi(line)
		dimmed = append(dimmed, dimStyle.Render(stripped))
	}
	return strings.Join(dimmed, "\n")
}

func stripAnsi(s string) string {
	var out strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\x1b' {
			// Skip CSI sequences (\x1b[...m) and OSC sequences
			i++
			if i < len(s) && s[i] == '[' {
				i++
				for i < len(s) && s[i] != 'm' {
					i++
				}
				if i < len(s) {
					i++ // skip 'm'
				}
			} else if i < len(s) && s[i] == ']' {
				for i < len(s) && s[i] != '\\' {
					i++
				}
				if i < len(s) {
					i++
				}
			}
		} else {
			out.WriteByte(s[i])
			i++
		}
	}
	return out.String()
}

func renderItemBlock(item document.Item, width int, t theme.Theme) string {
	switch item.Type {
	case "keybind-list":
		return renderKeybindBlock(item, width, t)
	case "command":
		return renderCommandBlock(item, width, t, false)
	case "table":
		return renderTableBlock(item, width, t)
	case "callout":
		return renderCalloutBlock(item, width, t)
	case "text":
		return renderTextBlock(item, width, t)
	case "separator":
		return renderSeparator(item, width, t)
	case "link":
		return renderLinkBlock(item, width, t)
	default:
		return ""
	}
}

func renderKeybindBlock(item document.Item, width int, t theme.Theme) string {
	if len(item.Entries) == 0 {
		return ""
	}

	keyWidth := 0
	for _, e := range item.Entries {
		kw := len(strings.Join(e.Keys, " "))
		if kw > keyWidth {
			keyWidth = kw
		}
	}
	if keyWidth > width/2 {
		keyWidth = width / 2
	}
	keyWidth += 4

	keyStyle := lipgloss.NewStyle().
		Foreground(t.KeyColor).
		Width(keyWidth)

	descStyle := lipgloss.NewStyle().
		Foreground(t.Foreground)

	var lines []string
	for _, e := range item.Entries {
		keyStr := strings.Join(e.Keys, " ")
		key := keyStyle.Render(keyStr)
		desc := descStyle.Render(e.Description)
		lines = append(lines, key+desc)
	}

	return strings.Join(lines, "\n")
}

func renderKeybindBlockHighlighted(item document.Item, width int, t theme.Theme, query string) string {
	if len(item.Entries) == 0 {
		return ""
	}

	q := strings.ToLower(query)

	keyWidth := 0
	for _, e := range item.Entries {
		kw := len(strings.Join(e.Keys, " "))
		if kw > keyWidth {
			keyWidth = kw
		}
	}
	if keyWidth > width/2 {
		keyWidth = width / 2
	}
	keyWidth += 4

	keyStyle := lipgloss.NewStyle().Foreground(t.KeyColor).Width(keyWidth)
	descStyle := lipgloss.NewStyle().Foreground(t.Foreground)
	dimStyle := lipgloss.NewStyle().Foreground(t.Border)

	var lines []string
	for _, e := range item.Entries {
		keyStr := strings.Join(e.Keys, " ")
		matches := strings.Contains(strings.ToLower(e.Description), q) ||
			strings.Contains(strings.ToLower(keyStr), q)

		if matches {
			key := keyStyle.Render(keyStr)
			desc := descStyle.Render(e.Description)
			lines = append(lines, key+desc)
		} else {
			line := dimStyle.Width(keyWidth).Render(keyStr) + dimStyle.Render(e.Description)
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}

func renderCommandBlock(item document.Item, width int, t theme.Theme, active bool) string {
	cmdStyle := lipgloss.NewStyle().Foreground(t.CommandColor)

	if active {
		cmdStyle = cmdStyle.Bold(true)
	}

	var lines []string
	marker := "  "
	if active {
		marker = "▸ "
	}
	lines = append(lines, marker+cmdStyle.Render(item.Command))
	if item.Description != "" {
		descStyle := lipgloss.NewStyle().Foreground(t.Muted)
		lines = append(lines, "  "+descStyle.Render(item.Description))
	}
	return strings.Join(lines, "\n")
}

func renderTableBlock(item document.Item, width int, t theme.Theme) string {
	if len(item.Rows) == 0 {
		return ""
	}

	colWidths := computeTableWidths(item.TableColumns, item.Rows, width)

	keyStyle := lipgloss.NewStyle().Foreground(t.KeyColor)
	descStyle := lipgloss.NewStyle().Foreground(t.Foreground)

	var lines []string
	for _, row := range item.Rows {
		var cells []string
		for i, val := range row.Values {
			if i >= len(colWidths) {
				break
			}
			if i == 0 {
				cells = append(cells, keyStyle.Width(colWidths[i]).Render(val))
			} else {
				cells = append(cells, descStyle.Width(colWidths[i]).Render(val))
			}
		}
		lines = append(lines, strings.Join(cells, ""))
	}

	return strings.Join(lines, "\n")
}

func computeTableWidths(cols []string, rows []document.TableRow, width int) []int {
	n := len(cols)
	if n == 0 && len(rows) > 0 {
		n = len(rows[0].Values)
	}
	if n == 0 {
		return nil
	}

	maxWidths := make([]int, n)
	for _, row := range rows {
		for i, val := range row.Values {
			if i < n && len(val) > maxWidths[i] {
				maxWidths[i] = len(val)
			}
		}
	}

	for i := range maxWidths {
		maxWidths[i] += 4
	}

	total := 0
	for _, w := range maxWidths {
		total += w
	}
	if total > width {
		for i := range maxWidths {
			maxWidths[i] = maxWidths[i] * width / total
			if maxWidths[i] < 4 {
				maxWidths[i] = 4
			}
		}
	}

	return maxWidths
}

func renderCalloutBlock(item document.Item, width int, t theme.Theme) string {
	text := item.Text
	if text == "" {
		text = item.Description
	}
	style := lipgloss.NewStyle().Foreground(t.Muted).Italic(true)
	return style.Render(text)
}

func renderTextBlock(item document.Item, width int, t theme.Theme) string {
	text := item.Text
	if text == "" {
		text = item.Description
	}
	if text == "" {
		return ""
	}
	style := lipgloss.NewStyle().Foreground(t.Foreground)
	return style.Render(text)
}

func renderSeparator(item document.Item, width int, t theme.Theme) string {
	label := item.SeparatorLabel
	if label == "" {
		label = item.Label
	}

	lineStyle := lipgloss.NewStyle().Foreground(t.Border)

	if label == "" {
		return lineStyle.Render(strings.Repeat("─", width))
	}

	labelStyled := lipgloss.NewStyle().Foreground(t.Muted).Render(" " + label + " ")
	labelWidth := lipgloss.Width(labelStyled)
	sideWidth := (width - labelWidth) / 2
	if sideWidth < 2 {
		sideWidth = 2
	}

	left := lineStyle.Render(strings.Repeat("─", sideWidth))
	right := lineStyle.Render(strings.Repeat("─", sideWidth))
	return left + labelStyled + right
}

func renderLinkBlock(item document.Item, width int, t theme.Theme) string {
	label := item.Label
	if label == "" {
		label = item.URL
	}

	labelStyle := lipgloss.NewStyle().
		Foreground(t.KeyColor).
		Underline(true)

	line := labelStyle.Render(label)

	if item.Description != "" {
		desc := lipgloss.NewStyle().Foreground(t.Muted).Render("  " + item.Description)
		line += desc
	}

	return line
}
