package version

import (
	"fmt"
	"log/slog"
	"runtime"

	"github.com/junevm/cdns/internal/config"
	"github.com/junevm/cdns/internal/ui"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// Module provides the version feature as an Fx module
var Module = fx.Module("version",
	fx.Provide(NewService),
	fx.Provide(NewCommand),
	fx.Invoke(RegisterCommand),
)

// BuildInfo holds version information injected at build time
type BuildInfo struct {
	Version string
	Commit  string
	Date    string
	BuiltBy string
}

// Service handles the business logic for version feature
type Service struct {
	config    *config.Config
	logger    *slog.Logger
	buildInfo BuildInfo
	styles    *ui.Styles
}

// NewService creates a new version service
func NewService(cfg *config.Config, logger *slog.Logger, buildInfo BuildInfo) *Service {
	return &Service{
		config:    cfg,
		logger:    logger,
		buildInfo: buildInfo,
		styles:    ui.NewStyles(),
	}
}

// PrintVersion displays version information
func (s *Service) PrintVersion(verbose bool) {
	// Print large banner
	fmt.Print(ui.GetBanner())
	fmt.Println()

	s.logger.Debug("displaying version information", slog.Bool("verbose", verbose))

	// Basic version info
	fmt.Printf("%s version %s\n",
		s.styles.RenderBold("cdns"),
		s.styles.Success.Render(s.buildInfo.Version),
	)

	// Verbose mode shows additional details
	if verbose {
		fmt.Println()
		fmt.Printf("  %s: %s\n", s.styles.RenderBold("Commit"), s.buildInfo.Commit)
		fmt.Printf("  %s: %s\n", s.styles.RenderBold("Built"), s.buildInfo.Date)
		fmt.Printf("  %s: %s\n", s.styles.RenderBold("Built by"), s.buildInfo.BuiltBy)
		fmt.Printf("  %s: %s\n", s.styles.RenderBold("Go version"), runtime.Version())
		fmt.Printf("  %s: %s/%s\n", s.styles.RenderBold("Platform"), runtime.GOOS, runtime.GOARCH)
	}
}

// CommandParams holds dependencies for the version command
type CommandParams struct {
	fx.In

	Service *Service
	Logger  *slog.Logger
}

// CommandResult wraps the command for Fx
type CommandResult struct {
	fx.Out

	Cmd *cobra.Command `name:"version"`
}

// NewCommand creates the version cobra command
func NewCommand(params CommandParams) CommandResult {
	var verbose bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  `Display version information for this CLI application.`,
		Run: func(cmd *cobra.Command, args []string) {
			params.Service.PrintVersion(verbose)
		},
	}

	// Command-specific flags
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed version information")

	return CommandResult{Cmd: cmd}
}

// RegisterParams holds dependencies for command registration
type RegisterParams struct {
	fx.In

	RootCmd *cobra.Command
	Cmd     *cobra.Command `name:"version"`
}

// RegisterCommand registers the version command with the root command
func RegisterCommand(params RegisterParams) {
	params.RootCmd.AddCommand(params.Cmd)
}
