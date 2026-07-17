package app

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gh-jsoares/grimoire/internal/clipboard"
	"github.com/gh-jsoares/grimoire/internal/config"
	"github.com/gh-jsoares/grimoire/internal/content"
	"github.com/gh-jsoares/grimoire/internal/document"
	"github.com/gh-jsoares/grimoire/internal/library"
	"github.com/gh-jsoares/grimoire/internal/tabs"
	"github.com/gh-jsoares/grimoire/internal/theme"
)

type Config struct {
	SingleDoc   bool
	NoIcons     bool
	Plain       bool
	InitTab     string
	InitSection string
}

type LinkEntry struct {
	Label string
	URL   string
}

type Model struct {
	Library *library.Library
	Config  Config
	GridCfg config.GridConfig
	Theme   theme.Theme

	ActiveDoc int
	Tabs      tabs.Model

	scrollOffsets map[int]int
	totalLines    int

	// Command navigation
	commands   []string // flat list of commands in current doc
	commandIdx int      // -1 = no selection

	// Link overlay
	links      []LinkEntry
	linkIdx    int
	showLinks  bool

	// Search
	searching    bool
	searchQuery  string
	searchFilter bool // true = filter mode, false = highlight mode

	// Status flash
	flashMsg string

	Width  int
	Height int
	ready  bool
}

func New(lib *library.Library, cfg Config, gridCfg config.GridConfig) Model {
	t := theme.Default()
	if cfg.Plain {
		t = theme.PlainTheme()
	}

	m := Model{
		Library:       lib,
		Config:        cfg,
		GridCfg:       gridCfg,
		Theme:         t,
		scrollOffsets: make(map[int]int),
		commandIdx:    -1,
		searchFilter:  false,
	}

	var tabList []tabs.Tab
	for _, doc := range lib.Documents {
		if doc.Hidden {
			continue
		}
		icon := doc.Icon
		if cfg.NoIcons {
			icon = ""
		}
		tabList = append(tabList, tabs.Tab{Title: doc.Title, Icon: icon})
	}
	m.Tabs = tabs.New(tabList, t)

	// Apply initial tab selection
	if cfg.InitTab != "" {
		for i, doc := range lib.Documents {
			if strings.EqualFold(doc.Title, cfg.InitTab) {
				m.Tabs.Set(i)
				m.ActiveDoc = i
				break
			}
			for _, alias := range doc.Aliases {
				if strings.EqualFold(alias, cfg.InitTab) {
					m.Tabs.Set(i)
					m.ActiveDoc = i
					break
				}
			}
		}
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

type clearFlashMsg struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Tabs.Width = m.Width
		m.ready = true
		m.computeTotalLines()
		m.buildCommands()
		m.buildLinks()
		return m, nil

	case tea.KeyMsg:
		return m, m.handleKey(msg)

	case clearFlashMsg:
		m.flashMsg = ""
		return m, nil
	}

	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) tea.Cmd {
	key := msg.String()

	// Link overlay mode
	if m.showLinks {
		return m.handleLinkKey(key)
	}

	// Search mode
	if m.searching {
		return m.handleSearchKey(msg)
	}

	switch key {
	case "q", "ctrl+c":
		return tea.Quit
	case "H", "left":
		m.Tabs.Prev()
		m.ActiveDoc = m.Tabs.Active
		m.computeTotalLines()
		m.buildCommands()
		m.buildLinks()
		m.commandIdx = -1
		return nil
	case "L", "right":
		m.Tabs.Next()
		m.ActiveDoc = m.Tabs.Active
		m.computeTotalLines()
		m.buildCommands()
		m.buildLinks()
		m.commandIdx = -1
		return nil
	case "1", "2", "3", "4", "5", "6", "7", "8", "9":
		idx := int(key[0]-'0') - 1
		if idx < len(m.Library.Documents) {
			m.Tabs.Set(idx)
			m.ActiveDoc = idx
			m.computeTotalLines()
			m.buildCommands()
			m.buildLinks()
			m.commandIdx = -1
		}
		return nil
	case "j", "down":
		m.scroll(2)
		return nil
	case "k", "up":
		m.scroll(-2)
		return nil
	case "ctrl+d":
		m.scroll(m.contentHeight() / 2)
		return nil
	case "ctrl+u":
		m.scroll(-m.contentHeight() / 2)
		return nil
	case "g", "home":
		m.scrollOffsets[m.ActiveDoc] = 0
		return nil
	case "G", "end":
		m.scrollToEnd()
		return nil
	case "n":
		m.nextCommand()
		return nil
	case "N":
		m.prevCommand()
		return nil
	case "y":
		return m.yankCommand()
	case "u":
		if len(m.links) > 0 {
			m.showLinks = true
			m.linkIdx = 0
		}
		return nil
	case "/":
		m.searching = true
		m.commandIdx = -1
		m.scrollOffsets[m.ActiveDoc] = 0
		return nil
	case "esc":
		if m.searchQuery != "" {
			m.searchQuery = ""
			m.scrollOffsets[m.ActiveDoc] = 0
			m.computeTotalLines()
			m.buildCommands()
		}
		return nil
	}

	return nil
}

func (m *Model) handleSearchKey(msg tea.KeyMsg) tea.Cmd {
	key := msg.String()

	switch key {
	case "esc":
		m.searching = false
		m.searchQuery = ""
		m.scrollOffsets[m.ActiveDoc] = 0
		m.computeTotalLines()
		m.buildCommands()
		return nil
	case "enter":
		m.searching = false
		m.buildCommands()
		return nil
	case "tab":
		m.searchFilter = !m.searchFilter
		m.scrollOffsets[m.ActiveDoc] = 0
		m.computeTotalLines()
		m.buildCommands()
		return nil
	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			m.scrollOffsets[m.ActiveDoc] = 0
			m.computeTotalLines()
			m.buildCommands()
		}
		return nil
	case "ctrl+c":
		return tea.Quit
	default:
		if len(key) == 1 && key[0] >= 32 {
			m.searchQuery += key
			m.scrollOffsets[m.ActiveDoc] = 0
			m.computeTotalLines()
			m.buildCommands()
		}
		return nil
	}
}

func (m *Model) handleLinkKey(key string) tea.Cmd {
	switch key {
	case "q", "esc":
		m.showLinks = false
		return nil
	case "j", "down":
		if m.linkIdx < len(m.links)-1 {
			m.linkIdx++
		}
		return nil
	case "k", "up":
		if m.linkIdx > 0 {
			m.linkIdx--
		}
		return nil
	case "y", "enter":
		if m.linkIdx >= 0 && m.linkIdx < len(m.links) {
			url := m.links[m.linkIdx].URL
			if err := clipboard.Copy(url); err != nil {
				m.flashMsg = "Copy failed: " + err.Error()
			} else {
				m.flashMsg = "Copied: " + url
			}
			m.showLinks = false
			return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
				return clearFlashMsg{}
			})
		}
		return nil
	}
	return nil
}

func (m *Model) scroll(delta int) {
	offset := m.scrollOffsets[m.ActiveDoc] + delta
	maxOffset := m.totalLines - m.contentHeight()
	if maxOffset < 0 {
		maxOffset = 0
	}
	if offset > maxOffset {
		offset = maxOffset
	}
	if offset < 0 {
		offset = 0
	}
	m.scrollOffsets[m.ActiveDoc] = offset
}

func (m *Model) scrollToEnd() {
	maxOffset := m.totalLines - m.contentHeight()
	if maxOffset < 0 {
		maxOffset = 0
	}
	m.scrollOffsets[m.ActiveDoc] = maxOffset
}

func (m Model) activeDocument() *document.Document {
	doc := &m.Library.Documents[m.ActiveDoc]
	if m.searchQuery != "" && m.searchFilter {
		return content.FilterDocument(doc, m.searchQuery)
	}
	return doc
}

func (m *Model) computeTotalLines() {
	if m.Width == 0 || m.ActiveDoc >= len(m.Library.Documents) {
		m.totalLines = 0
		return
	}
	doc := m.activeDocument()
	rendered := content.RenderDocument(doc, m.Width-2, m.Height, m.Theme, -1, m.GridCfg)
	m.totalLines = strings.Count(rendered, "\n") + 1
}

func (m *Model) buildCommands() {
	m.commands = nil
	if m.ActiveDoc >= len(m.Library.Documents) {
		return
	}
	doc := m.activeDocument()
	for _, sec := range doc.Sections {
		items := sec.Items
		if len(items) == 0 {
			for _, col := range sec.Columns {
				items = append(items, col.Items...)
			}
		}
		for _, item := range items {
			if item.Type == "command" && item.Command != "" {
				m.commands = append(m.commands, item.Command)
			}
		}
	}
	if m.commandIdx >= len(m.commands) {
		m.commandIdx = -1
	}
}

func (m *Model) buildLinks() {
	m.links = nil
	if m.ActiveDoc >= len(m.Library.Documents) {
		return
	}
	doc := &m.Library.Documents[m.ActiveDoc]
	for _, sec := range doc.Sections {
		items := sec.Items
		if len(items) == 0 {
			for _, col := range sec.Columns {
				items = append(items, col.Items...)
			}
		}
		for _, item := range items {
			if item.Type == "link" && item.URL != "" {
				label := item.Label
				if label == "" {
					label = item.URL
				}
				m.links = append(m.links, LinkEntry{Label: label, URL: item.URL})
			}
		}
	}
}

func (m *Model) nextCommand() {
	if len(m.commands) == 0 {
		return
	}
	m.commandIdx++
	if m.commandIdx >= len(m.commands) {
		m.commandIdx = 0
	}
}

func (m *Model) prevCommand() {
	if len(m.commands) == 0 {
		return
	}
	if m.commandIdx <= 0 {
		m.commandIdx = len(m.commands) - 1
	} else {
		m.commandIdx--
	}
}

func (m *Model) yankCommand() tea.Cmd {
	if m.commandIdx < 0 || m.commandIdx >= len(m.commands) {
		return nil
	}
	cmd := m.commands[m.commandIdx]
	if err := clipboard.Copy(cmd); err != nil {
		m.flashMsg = "Copy failed: " + err.Error()
	} else {
		m.flashMsg = "Copied: " + cmd
	}
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
		return clearFlashMsg{}
	})
}

func (m Model) contentHeight() int {
	h := m.Height - 5
	if h < 1 {
		h = 1
	}
	return h
}

func (m Model) View() string {
	if !m.ready || m.Width == 0 {
		return ""
	}

	if m.Width < 30 || m.Height < 10 {
		return lipgloss.NewStyle().
			Width(m.Width).
			Height(m.Height).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(m.Theme.Muted).
			Render("Terminal too small")
	}

	// Link overlay replaces the view entirely
	if m.showLinks {
		return m.renderLinkOverlay()
	}

	doc := m.activeDocument()

	highlightQuery := ""
	if m.searchQuery != "" && !m.searchFilter {
		highlightQuery = m.searchQuery
	}
	rendered := content.RenderDocument(doc, m.Width-2, m.Height, m.Theme, m.commandIdx, m.GridCfg, highlightQuery)
	lines := splitLines(rendered)
	m.totalLines = len(lines)

	offset := m.scrollOffsets[m.ActiveDoc]
	maxOffset := len(lines) - m.contentHeight()
	if maxOffset < 0 {
		maxOffset = 0
	}
	if offset > maxOffset {
		offset = maxOffset
	}
	if offset < 0 {
		offset = 0
	}

	viewHeight := m.contentHeight()
	end := offset + viewHeight
	if end > len(lines) {
		end = len(lines)
	}
	visible := lines[offset:end]

	for len(visible) < viewHeight {
		visible = append(visible, "")
	}

	origDoc := &m.Library.Documents[m.ActiveDoc]
	var parts []string
	parts = append(parts, m.titleView(origDoc))
	parts = append(parts, strings.Join(visible, "\n"))
	parts = append(parts, m.statusView())

	return strings.Join(parts, "\n")
}

func (m Model) renderLinkOverlay() string {
	overlayWidth := m.Width / 2
	if overlayWidth < 40 {
		overlayWidth = 40
	}
	if overlayWidth > m.Width-4 {
		overlayWidth = m.Width - 4
	}

	innerWidth := overlayWidth - 6

	var rows []string
	for i, link := range m.links {
		label := link.Label
		if len(label) > innerWidth-4 {
			label = label[:innerWidth-7] + "..."
		}
		url := link.URL
		if len(url) > innerWidth-2 {
			url = url[:innerWidth-5] + "..."
		}

		marker := "  "
		if i == m.linkIdx {
			marker = "▸ "
		}

		labelStyle := lipgloss.NewStyle().Foreground(m.Theme.Foreground)
		urlStyle := lipgloss.NewStyle().Foreground(m.Theme.Muted)
		if i == m.linkIdx {
			labelStyle = labelStyle.Foreground(m.Theme.Accent).Bold(true)
			urlStyle = urlStyle.Foreground(m.Theme.KeyColor)
		}

		row := marker + labelStyle.Render(label)
		rows = append(rows, row)
		rows = append(rows, "   "+urlStyle.Render(url))
		if i < len(m.links)-1 {
			rows = append(rows, "")
		}
	}

	boxContent := strings.Join(rows, "\n")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.Theme.ActiveBorder).
		Padding(1, 2).
		Width(overlayWidth).
		Render(boxContent)

	hint := lipgloss.NewStyle().Foreground(m.Theme.Muted).Render("j/k navigate  y/enter copy  esc close")
	box += "\n" + lipgloss.NewStyle().Width(overlayWidth).Align(lipgloss.Center).Render(hint)

	// Title
	title := lipgloss.NewStyle().
		Foreground(m.Theme.Accent).
		Bold(true).
		Render("  Links  ")

	// Center the whole thing vertically and horizontally
	boxHeight := strings.Count(box, "\n") + 1
	padTop := (m.Height - boxHeight) / 2
	if padTop < 0 {
		padTop = 0
	}

	var output strings.Builder
	// Fill top with empty lines
	for i := 0; i < padTop; i++ {
		output.WriteString(strings.Repeat(" ", m.Width) + "\n")
	}

	// Center the title above the box
	titleLine := lipgloss.NewStyle().Width(m.Width).Align(lipgloss.Center).Render(title)
	output.WriteString(titleLine + "\n")

	// Center each line of the box
	for _, line := range strings.Split(box, "\n") {
		lineWidth := lipgloss.Width(line)
		pad := (m.Width - lineWidth) / 2
		if pad < 0 {
			pad = 0
		}
		output.WriteString(strings.Repeat(" ", pad) + line + "\n")
	}

	// Fill remaining height
	currentLines := padTop + 1 + boxHeight
	for i := currentLines; i < m.Height; i++ {
		output.WriteString(strings.Repeat(" ", m.Width) + "\n")
	}

	return output.String()
}

func (m Model) titleView(doc *document.Document) string {
	var rows []string

	if len(m.Library.Documents) > 1 {
		var tabParts []string
		for i, tab := range m.Tabs.Tabs {
			label := tab.Title
			if tab.Icon != "" {
				label = tab.Icon + " " + label
			}
			if i == m.Tabs.Active {
				tabParts = append(tabParts, m.Theme.ActiveTab.Render(label))
			} else {
				tabParts = append(tabParts, m.Theme.InactiveTab.Render(label))
			}
		}
		tabRow := lipgloss.JoinHorizontal(lipgloss.Bottom, tabParts...)
		rows = append(rows, lipgloss.NewStyle().Width(m.Width).Align(lipgloss.Center).Render(tabRow))
	}

	title := doc.Title
	if doc.Icon != "" && !m.Config.NoIcons {
		title = doc.Icon + "  " + title
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.Theme.Accent).
		MarginTop(1).
		Width(m.Width).
		Align(lipgloss.Center)

	rows = append(rows, titleStyle.Render(strings.ToUpper(title)))

	lineWidth := len(title) + 8
	if lineWidth > m.Width-4 {
		lineWidth = m.Width - 4
	}
	underline := lipgloss.NewStyle().
		Foreground(m.Theme.Border).
		Width(m.Width).
		Align(lipgloss.Center).
		Render(strings.Repeat("═", lineWidth))
	rows = append(rows, underline)

	return strings.Join(rows, "\n")
}

func (m Model) statusView() string {
	if m.searching {
		prompt := lipgloss.NewStyle().Foreground(m.Theme.Accent).Render("/")
		query := lipgloss.NewStyle().Foreground(m.Theme.Foreground).Render(m.searchQuery)
		cursor := lipgloss.NewStyle().Foreground(m.Theme.Accent).Render("█")
		mode := "filter"
		if !m.searchFilter {
			mode = "highlight"
		}
		modeTag := lipgloss.NewStyle().Foreground(m.Theme.Muted).Render("  [tab: " + mode + "]")
		return m.Theme.StatusBar.Width(m.Width).Padding(0, 1).Render(prompt + query + cursor + modeTag)
	}

	if m.searchQuery != "" {
		mode := "filter"
		if !m.searchFilter {
			mode = "highlight"
		}
		clear := lipgloss.NewStyle().Foreground(m.Theme.Muted).Render("  (" + mode + " · / edit · esc clear)")
		filter := lipgloss.NewStyle().Foreground(m.Theme.Accent).Render("search: " + m.searchQuery)
		return m.Theme.StatusBar.Width(m.Width).Padding(0, 1).Render(filter + clear)
	}

	if m.flashMsg != "" {
		return m.Theme.StatusBar.Width(m.Width).Padding(0, 1).
			Foreground(m.Theme.Accent).Render(m.flashMsg)
	}

	hint := "q quit  H/L tabs  j/k scroll  n/N cmd  y yank  u links  / search"
	if len(m.Library.Documents) <= 1 {
		hint = "q quit  j/k scroll  n/N cmd  y yank  u links  / search"
	}
	if m.commandIdx >= 0 && len(m.commands) > 0 {
		hint += "  [" + m.commands[m.commandIdx] + "]"
		if len(hint) > m.Width-2 {
			hint = hint[:m.Width-5] + "..."
		}
	}
	return m.Theme.StatusBar.Width(m.Width).Padding(0, 1).Render(hint)
}

func splitLines(s string) []string {
	if s == "" {
		return []string{""}
	}
	return strings.Split(s, "\n")
}
