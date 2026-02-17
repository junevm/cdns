package reset

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/junevm/cdns/internal/config"
	"github.com/junevm/cdns/internal/dns/backend"
	"github.com/junevm/cdns/internal/dns/models"
	"github.com/junevm/cdns/internal/features/status"
	"github.com/junevm/cdns/internal/ui"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// Module provides the reset feature as an Fx module
var Module = fx.Module("reset",
	fx.Provide(NewService),
	fx.Provide(NewCommand),
	fx.Invoke(RegisterCommand),
)

// Detector defines the interface for detecting the backend
type Detector interface {
	Detect() (models.Backend, error)
}

// Reader defines the interface for reading current config
type Reader interface {
	ReadDNSConfig(ctx context.Context, backend models.Backend) (*status.StatusInfo, error)
}

// DNSWriter matches the interface needed to apply DNS settings
type DNSWriter interface {
	ResetToAutomatic(ctx context.Context, backend models.Backend, interfaces []string) error
}

// Service handles the business logic for reset feature
type Service struct {
	config   *config.Config
	logger   *slog.Logger
	styles   *ui.Styles
	detector Detector
	reader   Reader
	writer   DNSWriter
}

// NewService creates a new reset service
func NewService(cfg *config.Config, logger *slog.Logger, sysOps backend.SystemOps) *Service {
	return &Service{
		config:   cfg,
		logger:   logger,
		styles:   ui.NewStyles(),
		detector: backend.NewDetector(sysOps),
		reader:   backend.NewConfigReader(sysOps),
		writer:   backend.NewConfigWriter(sysOps),
	}
}

// Reset restores the system default DNS configuration (Automatic/DHCP)
func (s *Service) Reset(ctx context.Context) error {
	s.logger.Debug("resetting DNS configuration to system default")

	// Detect backend
	b, err := s.detector.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect DNS backend: %w", err)
	}

	// Read current status to get active interfaces
	statusInfo, err := s.reader.ReadDNSConfig(ctx, b)
	if err != nil {
		return fmt.Errorf("failed to read current configuration: %w", err)
	}

	var interfaces []string
	for _, iface := range statusInfo.Interfaces {
		interfaces = append(interfaces, iface.Name)
	}

	if len(interfaces) == 0 {
		fmt.Printf("\n%s\n", s.styles.RenderWarning("No active interfaces found to reset."))
		return nil
	}

	// Apply configuration
	if err := s.writer.ResetToAutomatic(ctx, b, interfaces); err != nil {
		return fmt.Errorf("failed to reset configuration: %w", err)
	}

	s.logger.Debug("successfully reset DNS configuration",
		slog.String("backend", string(b)))

	fmt.Printf("\n%s\n", s.styles.RenderSuccess("DNS configuration reset to system default (Automatic/DHCP)"))

	return nil
}

// CommandResult wraps the reset command
type CommandResult struct {
	fx.Out

	Cmd *cobra.Command `name:"reset"`
}

// NewCommand creates the reset cobra command
func NewCommand(s *Service) CommandResult {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Restore previous DNS configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Reset(cmd.Context())
		},
	}
	return CommandResult{Cmd: cmd}
}

// RegisterCommandParams holds dependencies for command registration
type RegisterCommandParams struct {
	fx.In

	Root *cobra.Command
	Cmd  *cobra.Command `name:"reset"`
}

// RegisterCommand registers the command with root
func RegisterCommand(p RegisterCommandParams) {
	p.Root.AddCommand(p.Cmd)
}
