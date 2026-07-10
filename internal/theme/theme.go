package theme

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Background   lipgloss.Color
	Foreground   lipgloss.Color
	Muted        lipgloss.Color
	Border       lipgloss.Color
	ActiveBorder lipgloss.Color
	Accent       lipgloss.Color
	KeyColor     lipgloss.Color
	CommandColor lipgloss.Color

	ActiveTab   lipgloss.Style
	InactiveTab lipgloss.Style
	StatusBar   lipgloss.Style
	StatusKey   lipgloss.Style

	Plain bool
}

func Default() Theme {
	t := Theme{
		Background:   lipgloss.Color("235"),
		Foreground:   lipgloss.Color("252"),
		Muted:        lipgloss.Color("242"),
		Border:       lipgloss.Color("238"),
		ActiveBorder: lipgloss.Color("62"),
		Accent:       lipgloss.Color("180"), // warm gold for headings
		KeyColor:     lipgloss.Color("115"), // teal/cyan for keys
		CommandColor: lipgloss.Color("114"), // green for commands
	}

	t.ActiveTab = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("252")).
		Background(lipgloss.Color("62")).
		Padding(0, 1)

	t.InactiveTab = lipgloss.NewStyle().
		Foreground(lipgloss.Color("242")).
		Padding(0, 1)

	t.StatusBar = lipgloss.NewStyle().
		Foreground(lipgloss.Color("242")).
		Background(lipgloss.Color("236"))

	t.StatusKey = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("252"))

	return t
}

func PlainTheme() Theme {
	t := Theme{Plain: true}
	t.Foreground = lipgloss.Color("252")
	t.Muted = lipgloss.Color("242")
	t.Border = lipgloss.Color("238")
	t.Accent = lipgloss.Color("252")
	t.KeyColor = lipgloss.Color("252")
	t.CommandColor = lipgloss.Color("252")
	t.ActiveTab = lipgloss.NewStyle().Bold(true).Padding(0, 1)
	t.InactiveTab = lipgloss.NewStyle().Padding(0, 1)
	t.StatusBar = lipgloss.NewStyle()
	t.StatusKey = lipgloss.NewStyle().Bold(true)
	return t
}
