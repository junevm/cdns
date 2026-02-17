package cli

import (
	"github.com/junevm/cdns/internal/ui"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle      = lipgloss.NewStyle().MarginLeft(2).Bold(true).Foreground(lipgloss.Color("39"))
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type item struct {
	title, desc string
	cmd         string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type menuModel struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m menuModel) Init() tea.Cmd {
	return nil
}

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			m.quitting = true
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.cmd
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		// Adjust for banner height
		m.list.SetSize(msg.Width-h, msg.Height-v-7)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m menuModel) View() string {
	if m.quitting {
		return ""
	}
	return ui.GetBanner() + "\n" + m.list.View()
}

// RunMainMenu shows the interactive main menu and returns the selected command name
func RunMainMenu() (string, error) {
	items := []list.Item{
		item{title: "Configure DNS", desc: "Select a preset or enter custom IPs (Interactive)", cmd: "set"},
		item{title: "List Servers", desc: "View all available DNS presets", cmd: "list"},
		item{title: "Check Status", desc: "View current DNS settings and active interfaces", cmd: "status"},
		item{title: "Quick Reset", desc: "Restore previous DNS configuration", cmd: "reset"},
		item{title: "Version Info", desc: "Display application version and build details", cmd: "version"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "CDNS - Linux DNS Manager"
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	m := menuModel{list: l}

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	choice := finalModel.(menuModel).choice
	return choice, nil
}
