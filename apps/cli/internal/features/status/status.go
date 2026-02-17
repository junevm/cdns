package status

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/junevm/cdns/apps/cli/internal/config"
	"github.com/junevm/cdns/apps/cli/internal/dns/models"
	"github.com/junevm/cdns/apps/cli/internal/ui"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// Module provides the status feature as an Fx module
var Module = fx.Module("status",
	fx.Provide(NewService),
	fx.Provide(NewCommand),
	fx.Invoke(RegisterCommand),
)

// Detector defines the interface for backend detection
type Detector interface {
	Detect() (models.Backend, error)
	DetectWithReason() (models.Backend, string, error)
}

// Reader defines the interface for reading DNS configuration
type Reader interface {
	ReadDNSConfig(ctx context.Context, backend models.Backend) (*StatusInfo, error)
}

// StatusInfo holds the current DNS status
type StatusInfo struct {
	Backend    models.Backend    `json:"backend"`
	Interfaces []InterfaceStatus `json:"interfaces"`
	Managed    bool              `json:"managed"`
	Warnings   []string          `json:"warnings"`
}

// InterfaceStatus holds DNS information for a network interface
type InterfaceStatus struct {
	Name string   `json:"name"`
	IPv4 []string `json:"ipv4"`
	IPv6 []string `json:"ipv6"`
}

// Service handles the business logic for status feature
type Service struct {
	config   *config.Config
	logger   *slog.Logger
	styles   *ui.Styles
	detector Detector
	reader   Reader
}

// NewService creates a new status service
func NewService(cfg *config.Config, logger *slog.Logger, detector Detector, reader Reader) *Service {
	return &Service{
		config:   cfg,
		logger:   logger,
		styles:   ui.NewStyles(),
		detector: detector,
		reader:   reader,
	}
}

// GetStatus retrieves the current DNS status
func (s *Service) GetStatus(ctx context.Context) (*StatusInfo, error) {
	s.logger.Debug("retrieving DNS status")

	// Detect backend
	backend, err := s.detector.Detect()
	if err != nil {
		return nil, fmt.Errorf("failed to detect DNS backend: %w", err)
	}

	s.logger.Debug("detected backend", slog.String("backend", string(backend)))

	// Read DNS configuration from the detected backend
	status, err := s.reader.ReadDNSConfig(ctx, backend)
	if err != nil {
		return nil, fmt.Errorf("failed to read DNS configuration: %w", err)
	}

	return status, nil
}

// FormatStatus formats the status information for output
func (s *Service) FormatStatus(status *StatusInfo, jsonFormat bool) (string, error) {
	if jsonFormat {
		return s.formatJSON(status)
	}
	return s.formatHuman(status), nil
}

// formatJSON formats status as JSON
func (s *Service) formatJSON(status *StatusInfo) (string, error) {
	data, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal status to JSON: %w", err)
	}
	return string(data), nil
}

// formatHuman formats status in human-readable format (Visual & Concise)
func (s *Service) formatHuman(status *StatusInfo) string {
	if len(status.Interfaces) == 0 {
		return s.styles.RenderDim("No active network interfaces found.")
	}

	termWidth, _, _ := ui.GetTerminalSize()
	if termWidth <= 0 {
		termWidth = 80
	}

	// Calculate available width for DNS servers column
	// Fixed widths approx: Interface (15) + Backend (15) + Status (12) + Borders/Padding (14) = ~56
	dnsColWidth := termWidth - 60
	if dnsColWidth < 20 {
		dnsColWidth = 20
	}

	var rows [][]string
	for _, iface := range status.Interfaces {
		allIPs := append(iface.IPv4, iface.IPv6...)
		dnsString := "None"
		if len(allIPs) > 0 {
			dnsString = strings.Join(allIPs, ", ")
		}

		statusDot := s.styles.Success.Render("●")
		statusText := "Active"
		if len(allIPs) == 0 {
			statusDot = s.styles.Error.Render("●")
			statusText = "Inactive"
		}

		rows = append(rows, []string{
			iface.Name,
			fmt.Sprintf("%s", string(status.Backend)),
			dnsString,
			fmt.Sprintf("%s %s", statusDot, statusText),
		})
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("63"))).
		Headers("INTERFACE", "BACKEND", "DNS SERVERS", "STATUS").
		Rows(rows...)

	t.StyleFunc(func(row, col int) lipgloss.Style {
		style := lipgloss.NewStyle().Padding(0, 1)

		if col == 2 { // DNS Servers
			style = style.Width(dnsColWidth)
		}

		if row == 0 { // Header
			return style.
				Bold(true).
				Foreground(lipgloss.Color("205")).
				Align(lipgloss.Center)
		}

		// Content styles
		switch {
		case col == 2: // DNS Servers
			return style.Foreground(lipgloss.Color("86"))
		case col == 0: // Interface
			return style.Bold(true)
		default:
			return style
		}
	})

	var output strings.Builder
	output.WriteString("\n" + s.styles.Header.Render("Current DNS Status") + "\n\n")
	output.WriteString(t.Render())
	output.WriteString("\n")

	// Managed Status
	if !status.Managed {
		output.WriteString("  " + s.styles.RenderWarning("(Unmanaged by this tool)"))
	}

	// Warnings
	if len(status.Warnings) > 0 {
		output.WriteString("\n")
		for _, warning := range status.Warnings {
			output.WriteString(fmt.Sprintf("  %s %s\n", s.styles.Error.Render("!"), s.styles.RenderDim(warning)))
		}
	}

	return output.String()
}

// CommandParams holds dependencies for the status command
type CommandParams struct {
	fx.In

	Service *Service
	Logger  *slog.Logger
	Config  *config.Config
}

// CommandResult wraps the command for Fx
type CommandResult struct {
	fx.Out

	Cmd *cobra.Command `name:"status"`
}

// NewCommand creates the status cobra command
func NewCommand(params CommandParams) CommandResult {
	var jsonFormat bool

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Display current DNS configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Get status
			status, err := params.Service.GetStatus(ctx)
			if err != nil {
				return err
			}

			// Format and display output
			output, err := params.Service.FormatStatus(status, jsonFormat)
			if err != nil {
				return err
			}

			fmt.Println(output)
			return nil
		},
	}

	// Command-specific flags
	cmd.Flags().BoolVar(&jsonFormat, "json", false, "Output in JSON format")

	return CommandResult{Cmd: cmd}
}

// RegisterParams holds dependencies for command registration
type RegisterParams struct {
	fx.In

	RootCmd *cobra.Command
	Cmd     *cobra.Command `name:"status"`
}

// RegisterCommand registers the status command with the root command
func RegisterCommand(params RegisterParams) {
	params.RootCmd.AddCommand(params.Cmd)
}
