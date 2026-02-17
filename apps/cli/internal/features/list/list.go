package list

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"cli/internal/config"
	"cli/internal/dns/presets"
	"cli/internal/ui"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// Module provides the list feature as an Fx module
var Module = fx.Module("list",
	fx.Provide(NewService),
	fx.Provide(NewCommand),
	fx.Invoke(RegisterCommand),
)

// Service handles listing DNS servers
type Service struct {
	config *config.Config
	logger *slog.Logger
	styles *ui.Styles
}

// NewService creates a new list service
func NewService(cfg *config.Config, logger *slog.Logger) *Service {
	return &Service{
		config: cfg,
		logger: logger,
		styles: ui.NewStyles(),
	}
}

// PresetItem represents a listed DNS preset
type PresetItem struct {
	Name        string
	Servers     string
	Type        string
	Description string
}

// ListPresets retrieves all available presets, sorted by name
func (s *Service) ListPresets() []PresetItem {
	var items []PresetItem

	// 1. Built-in presets
	builtins := presets.All()
	for name, preset := range builtins {
		// Use correct casing for built-ins
		displayName := name
		switch strings.ToLower(name) {
		case "cloudflare":
			displayName = "Cloudflare"
		case "google":
			displayName = "Google"
		case "opendns":
			displayName = "OpenDNS"
		case "adguard":
			displayName = "AdGuard"
		case "quad9":
			displayName = "Quad9"
		default:
			// Fallback: capitalize
			if len(displayName) > 1 {
				displayName = strings.ToUpper(displayName[:1]) + displayName[1:]
			}
		}

		ips := strings.Join(preset.IPv4, ", ")
		if len(ips) == 0 && len(preset.IPv6) > 0 {
			ips = strings.Join(preset.IPv6, ", ")
		}

		items = append(items, PresetItem{
			Name:        displayName,
			Servers:     ips,
			Type:        "Built-in",
			Description: preset.Description,
		})
	}

	// 2. Custom presets from config
	if s.config != nil && s.config.DNS.CustomPresets != nil {
		for name, ips := range s.config.DNS.CustomPresets {
			items = append(items, PresetItem{
				Name:        name, // Preserve user case
				Servers:     strings.Join(ips, ", "),
				Type:        "Custom",
				Description: "User-defined",
			})
		}
	}

	// Sort by name case-insensitively
	sort.Slice(items, func(i, j int) bool {
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})

	return items
}

// PrintTable displays the presets in a formatted table
func (s *Service) PrintTable() error {
	termWidth, _, _ := ui.GetTerminalSize()
	// Fallback to 80 if cannot detect width or it's too small
	if termWidth <= 0 || termWidth < 40 {
		termWidth = 80
	}

	// Calculate column widths based on available space
	// Minimum required for fixed columns: Name (15) + Type (8) + Borders (10) + Padding (6) = 39
	// We'll adjust based on terminal width

	var nameWidth, typeWidth int
	var showDesc bool

	if termWidth < 60 {
		// Small screen: Compact layout
		nameWidth = 15
		typeWidth = 8
		showDesc = false
	} else {
		// Standard/Large screen
		nameWidth = 20
		typeWidth = 10
		showDesc = true
	}

	// Calculate remaining width for Servers and potentially Description
	fixedWidth := nameWidth + typeWidth + 10 // Approximation for borders/padding
	remainingWidth := termWidth - fixedWidth
	if remainingWidth < 10 {
		remainingWidth = 10
	}

	var serverColWidth, descColWidth int
	if showDesc {
		serverColWidth = int(float64(remainingWidth) * 0.6)
		descColWidth = remainingWidth - serverColWidth
	} else {
		serverColWidth = remainingWidth
		descColWidth = 0
	}

	items := s.ListPresets()

	var headers []string
	if showDesc {
		headers = []string{"PRESET", "SOURCE", "DNS SERVERS", "NOTES"}
	} else {
		headers = []string{"PRESET", "SOURCE", "DNS SERVERS"}
	}

	var rows [][]string
	for _, item := range items {
		row := []string{
			item.Name,
			item.Type,
			item.Servers,
		}
		if showDesc {
			row = append(row, item.Description)
		}
		rows = append(rows, row)
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("63"))).
		Headers(headers...).
		Rows(rows...)

	t.StyleFunc(func(row, col int) lipgloss.Style {
		style := lipgloss.NewStyle().Padding(0, 1)

		switch col {
		case 0: // Preset
			style = style.Width(nameWidth)
		case 1: // Source
			style = style.Width(typeWidth)
		case 2: // DNS Servers
			style = style.Width(serverColWidth)
		case 3: // Notes
			if showDesc {
				style = style.Width(descColWidth)
			}
		}

		if row == 0 { // Header
			return style.
				Bold(true)
		}

		return style
	})

	fmt.Println(t.Render())

	// Add footer info
	configPath := s.config.LoadedFrom
	if configPath == "" {
		configPath = "defaults (no config file found)"
	}
	fmt.Printf("\n%s\n", s.styles.RenderInfo(fmt.Sprintf("Config: %s", configPath)))
	fmt.Printf("%s\n", s.styles.RenderDim("Use 'cdns set <name>' to apply a preset."))

	return nil
}

// CommandParams holds dependencies for the list command
type CommandParams struct {
	fx.In
	Service *Service
}

// CommandResult wraps the command for Fx
type CommandResult struct {
	fx.Out
	Cmd *cobra.Command `name:"list"`
}

// NewCommand creates the list cobra command
func NewCommand(params CommandParams) CommandResult {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List DNS presets",
		Long:  `List all available DNS presets.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := params.Service.PrintTable(); err != nil {
				fmt.Println("Error listing presets:", err)
			}
		},
	}

	return CommandResult{Cmd: cmd}
}

// RegisterParams holds dependencies for command registration
type RegisterParams struct {
	fx.In
	RootCmd *cobra.Command
	Cmd     *cobra.Command `name:"list"`
}

// RegisterCommand registers the list command with the root command
func RegisterCommand(params RegisterParams) {
	params.RootCmd.AddCommand(params.Cmd)
}
