package set

import (
	"fmt"
	"sort"
	"strings"

	"github.com/junevm/cdns/internal/config"
	"github.com/junevm/cdns/internal/dns/presets"
	"github.com/junevm/cdns/internal/ui"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// step indicates the current step in the wizard
type step int

const (
	stepSelectPreset step = iota
	stepInputCustom
	stepSelectInterface
	stepConfirm
	stepDone
)

type model struct {
	step step

	// Components
	table table.Model
	input textinput.Model

	// Data
	interfaces []string

	// Selection
	isCustom       bool
	selectedPreset string
	customDNS      string
	selectedIface  string

	// State
	width, height int
	styles        *ui.Styles
	config        *config.Config
	quitting      bool
}

func newModel(cfg *config.Config, interfaces []string) model {
	input := textinput.New()
	input.Placeholder = "1.1.1.1, 8.8.8.8"
	input.CharLimit = 100
	input.Width = 40

	// Initialize table
	t := table.New(
		table.WithFocused(true),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("39"))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{
		step:       stepSelectPreset,
		table:      t,
		input:      input,
		interfaces: interfaces,
		styles:     ui.NewStyles(),
		config:     cfg,
	}
	m.initPresetTable()
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.String() == "q" {
			m.quitting = true
			return m, tea.Quit
		}
		if msg.String() == "c" && m.step == stepSelectPreset {
			m.isCustom = true
			m.step = stepInputCustom
			m.input.Focus()
			return m, textinput.Blink
		}
		if msg.Type == tea.KeyEsc {
			if m.step == stepInputCustom || m.step == stepSelectInterface {
				m.step = stepSelectPreset
				m.initPresetTable()
			} else if m.step == stepConfirm {
				m.step = stepSelectInterface
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Recalculate table dimensions
		if m.step == stepSelectPreset {
			m.initPresetTable()
		} else if m.step == stepSelectInterface {
			m.initInterfaceTable()
		}
	}

	switch m.step {
	case stepSelectPreset:
		m.table, cmd = m.table.Update(msg)
		if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
			selected := m.table.SelectedRow()[0]
			if selected == "CUSTOM" {
				m.isCustom = true
				m.step = stepInputCustom
				m.input.Focus()
				return m, textinput.Blink
			} else {
				m.isCustom = false
				m.selectedPreset = selected
				m.step = stepSelectInterface
				m.initInterfaceTable()
			}
		}

	case stepInputCustom:
		m.input, cmd = m.input.Update(msg)
		if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
			val := strings.TrimSpace(m.input.Value())
			if val != "" {
				m.customDNS = val
				m.step = stepSelectInterface
				m.initInterfaceTable()
			}
		}

	case stepSelectInterface:
		m.table, cmd = m.table.Update(msg)
		if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
			m.selectedIface = m.table.SelectedRow()[0]
			// Skip explicit confirmation screen, it breaks the flow
			m.step = stepDone
			return m, tea.Quit
		}

	case stepConfirm:
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch strings.ToLower(msg.String()) {
			case "y", "enter":
				m.step = stepDone
				return m, tea.Quit
			case "n":
				m.quitting = true
				return m, tea.Quit
			}
		}
	}

	return m, cmd
}

func (m *model) initPresetTable() {
	width := m.width
	if width <= 0 {
		width = 80
	}

	presetW := 15
	serverW := int(float64(width-presetW-10) * 0.4) // 40% of remaining
	if serverW < 20 {
		serverW = 20
	}
	descW := width - presetW - serverW - 8 // spacing
	if descW < 10 {
		descW = 10
	}

	columns := []table.Column{
		{Title: "PRESET", Width: presetW},
		{Title: "DNS SERVERS", Width: serverW},
		{Title: "DESCRIPTION", Width: descW},
	}

	all := presets.All()
	var names []string
	for k := range all {
		names = append(names, k)
	}
	sort.Strings(names)

	var rows []table.Row

	// Inbuilt presets
	for _, name := range names {
		p := all[name]
		ips := strings.Join(p.IPv4, ", ")
		if len(ips) == 0 && len(p.IPv6) > 0 {
			ips = strings.Join(p.IPv6, ", ")
		}
		rows = append(rows, table.Row{CapitalizePresetName(name), ips, p.Description})
	}

	// Add custom presets from config
	if m.config != nil && m.config.DNS.CustomPresets != nil {
		var customNames []string
		for name := range m.config.DNS.CustomPresets {
			customNames = append(customNames, name)
		}
		sort.Strings(customNames)

		for _, name := range customNames {
			ips := strings.Join(m.config.DNS.CustomPresets[name], ", ")
			rows = append(rows, table.Row{CapitalizePresetName(name), ips, "User-defined preset"})
		}
	}

	// Add Custom option last
	rows = append(rows, table.Row{"CUSTOM", "---", "Enter custom DNS servers"})

	m.table.SetRows([]table.Row{})
	m.table.SetColumns(columns)
	m.table.SetRows(rows)
	// Select the first item (should be an inbuilt preset)
	m.table.SetCursor(0)
	m.table.SetHeight(8)
}

func (m *model) initInterfaceTable() {
	width := m.width
	if width <= 0 {
		width = 80
	}

	ifaceW := 15
	descW := width - ifaceW - 8
	if descW < 20 {
		descW = 20
	}

	columns := []table.Column{
		{Title: "INTERFACE", Width: ifaceW},
		{Title: "DESCRIPTION", Width: descW},
	}

	rows := []table.Row{
		{"All Interfaces", "Apply to all active network interfaces"},
	}
	for _, iface := range m.interfaces {
		rows = append(rows, table.Row{iface, "Network interface"})
	}

	m.table.SetRows([]table.Row{})
	m.table.SetColumns(columns)
	m.table.SetRows(rows)

	// Ensure height is sufficient to show header + rows
	// If height is too small, rows might be hidden by header
	height := len(rows) + 2 // Add padding for header/border
	if height < 5 {
		height = 5
	}
	m.table.SetHeight(height)

	// Default to selecting the active interface (index 1) if available
	// otherwise "All Interfaces" (index 0) remains selected
	if len(m.interfaces) > 0 {
		m.table.SetCursor(1)
	}
}

func (m model) View() string {
	if m.quitting || m.step == stepDone {
		return ""
	}

	var s strings.Builder
	s.WriteString(m.styles.Header.Render("DNS Configuration") + "\n\n")

	footer := "\n\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Update config.yaml to add custom presets • ↑/k up • ↓/j down • enter select")

	switch m.step {
	case stepSelectPreset:
		s.WriteString("Select a DNS provider:\n\n")
		s.WriteString(m.table.View())
		s.WriteString(footer)
		// Removed extra help as it's now in footer

	case stepSelectInterface:
		s.WriteString("Select Network Interface:\n\n")
		s.WriteString(m.table.View())
		s.WriteString(footer)

	case stepInputCustom:
		s.WriteString("Enter DNS Addresses (comma-separated):\n\n")
		s.WriteString(m.input.View())
		s.WriteString("\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("enter confirm • esc back • q quit"))

	case stepConfirm:
		s.WriteString(m.styles.Warning.Render("Confirm Changes") + "\n\n")
		s.WriteString(fmt.Sprintf("  Type:      %s\n", m.getTypeDisplay()))
		s.WriteString(fmt.Sprintf("  Target:    %s\n", m.selectedIface))
		s.WriteString("\nApply configuration? [y/N]")
		s.WriteString("\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("y/n confirm • esc back • q quit"))
	}

	return s.String()
}

func (m model) getTypeDisplay() string {
	if m.isCustom {
		return m.styles.RenderInfo(fmt.Sprintf("Custom (%s)", m.customDNS))
	}
	return m.styles.RenderInfo(fmt.Sprintf("Preset (%s)", m.selectedPreset))
}
