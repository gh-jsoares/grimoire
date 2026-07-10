package tabs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gh-jsoares/grimoire/internal/theme"
)

type Tab struct {
	Title string
	Icon  string
}

type Model struct {
	Tabs   []Tab
	Active int
	Theme  theme.Theme
	Width  int
}

func New(tabs []Tab, t theme.Theme) Model {
	return Model{Tabs: tabs, Theme: t}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m *Model) Next() {
	if m.Active < len(m.Tabs)-1 {
		m.Active++
	}
}

func (m *Model) Prev() {
	if m.Active > 0 {
		m.Active--
	}
}

func (m *Model) Set(idx int) {
	if idx >= 0 && idx < len(m.Tabs) {
		m.Active = idx
	}
}

func (m Model) View() string {
	if len(m.Tabs) <= 1 {
		return ""
	}

	var rendered []string
	for i, tab := range m.Tabs {
		label := tab.Title
		if tab.Icon != "" {
			label = tab.Icon + " " + label
		}

		if i == m.Active {
			rendered = append(rendered, m.Theme.ActiveTab.Render(label))
		} else {
			rendered = append(rendered, m.Theme.InactiveTab.Render(label))
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Bottom, rendered...)

	style := lipgloss.NewStyle().
		Width(m.Width).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(m.Theme.Border)

	return style.Render(row)
}
