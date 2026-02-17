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
	Name    string
	Slug    string
	Servers string
	Type    string
}

// ListPresets retrieves all available presets, sorted by name
func (s *Service) ListPresets() []PresetItem {
	var items []PresetItem

	// 1. Built-in presets
	builtins := presets.All()
	for slug, preset := range builtins {
		// Use correct casing for built-ins
		displayName := slug
		switch strings.ToLower(slug) {
		case "cloudflare":
			displayName = "Cloudflare"
		case "google":
			displayName = "Google"
		case "opendns":
			displayName = "OpenDNS"
		case "opendns-family":
			displayName = "OpenDNS Family"
		case "adguard":
			displayName = "AdGuard"
		case "adguard-family":
			displayName = "AdGuard Family"
		case "quad9":
			displayName = "Quad9"
		case "cleanbrowsing-family":
			displayName = "CleanBrowsing Family"
		case "cleanbrowsing-security":
			displayName = "CleanBrowsing Security"
		case "yandex-basic":
			displayName = "Yandex Basic"
		case "yandex-safe":
			displayName = "Yandex Safe"
		case "yandex-family":
			displayName = "Yandex Family"
		case "comodo":
			displayName = "Comodo"
		case "verisign":
			displayName = "Verisign"
		case "dnswatch":
			displayName = "DNS.WATCH"
		case "level3":
			displayName = "Level3"
		case "tencent":
			displayName = "Tencent"
		case "alibaba":
			displayName = "Alibaba"
		case "neustar":
			displayName = "Neustar"
		case "opennic":
			displayName = "OpenNIC"
		case "he":
			displayName = "Hurricane Electric"
		case "safeserve":
			displayName = "SafeServe"
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
			Name:    displayName,
			Slug:    slug,
			Servers: ips,
			Type:    "Built-in",
		})
	}

	// 2. Custom presets from config
	if s.config != nil && s.config.DNS.CustomPresets != nil {
		for name, ips := range s.config.DNS.CustomPresets {
			items = append(items, PresetItem{
				Name:    name, // Preserve user case
				Slug:    strings.ToLower(strings.ReplaceAll(name, " ", "-")),
				Servers: strings.Join(ips, ", "),
				Type:    "Custom",
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

	// Dynamic column assignment
	items := s.ListPresets()

	// Column Widths
	nameWidth := 20
	slugWidth := 25
	sourceWidth := 10

	// Borders (5) + Padding (2*4=8) = 13
	remainingWidth := termWidth - (nameWidth + slugWidth + sourceWidth + 15)
	if remainingWidth < 15 {
		remainingWidth = 15
	}
	serverWidth := remainingWidth

	headers := []string{"PRESET", "COMMAND ID", "SOURCE", "DNS SERVERS"}

	var rows [][]string
	for _, item := range items {
		rows = append(rows, []string{
			item.Name,
			item.Slug,
			item.Type,
			item.Servers,
		})
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
		case 1: // ID
			style = style.Width(slugWidth).Foreground(lipgloss.Color("86")) // Cyan for IDs
		case 2: // Source
			style = style.Width(sourceWidth).Faint(true)
		case 3: // Servers
			style = style.Width(serverWidth)
		}

		if row == 0 { // Header
			return style.
				Bold(true).
				Foreground(lipgloss.Color("39")) // Cyan to match banner
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
