package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/junevm/cdns/apps/cli/internal/cli"
	"github.com/junevm/cdns/apps/cli/internal/config"
	"github.com/junevm/cdns/apps/cli/internal/dns/backend"
	"github.com/junevm/cdns/apps/cli/internal/features/list"
	"github.com/junevm/cdns/apps/cli/internal/features/reset"
	"github.com/junevm/cdns/apps/cli/internal/features/set"
	"github.com/junevm/cdns/apps/cli/internal/features/status"
	"github.com/junevm/cdns/apps/cli/internal/features/version"
	"github.com/junevm/cdns/apps/cli/internal/logger"
	"github.com/junevm/cdns/apps/cli/internal/ui"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// Build-time variables injected by GoReleaser
var (
	buildVersion = "dev"
	buildCommit  = "none"
	buildDate    = "unknown"
	buildBy      = "unknown"
)

func main() {
	// Create Fx application with all modules
	app := fx.New(
		// Suppress Fx logs for CLI applications
		fx.NopLogger,

		// Provide core dependencies
		fx.Provide(
			NewRootCommand,
			NewConfig,
			NewLogger,
			NewBuildInfo,
			NewSystemOps,
			NewDetector,
			NewConfigReader,
		),

		// Register feature modules
		version.Module,
		status.Module,
		set.Module,
		reset.Module,
		list.Module,

		// Lifecycle hooks
		fx.Invoke(RegisterLifecycleHooks),
		fx.Invoke(RunCLI),
	)

	// Run with graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		styles := ui.NewStyles()
		fmt.Fprintln(os.Stderr, styles.RenderError(fmt.Sprintf("Failed to start application: %v", err)))
		os.Exit(1)
	}

	<-ctx.Done()

	stopCtx, stopCancel := context.WithTimeout(context.Background(), fx.DefaultTimeout)
	defer stopCancel()

	if err := app.Stop(stopCtx); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to stop application gracefully: %v\n", err)
		os.Exit(1)
	}
}

// NewRootCommand creates the root cobra command
func NewRootCommand(cfg *config.Config, log *slog.Logger) *cobra.Command {
	deps := cli.Dependencies{
		Config: cfg,
		Logger: log,
	}
	return cli.NewRootCmd(deps)
}

// NewConfig loads configuration using Koanf
func NewConfig() (*config.Config, error) {
	loader := config.NewLoader()

	// 1. Check environment variable for config file override
	configFile := os.Getenv("CDNS_CONFIG_FILE")

	// 2. If not set, use default path (~/.config/cdns/config.yaml)
	if configFile == "" {
		defaultPath, err := config.DefaultConfigPath()
		if err != nil {
			// Fallback to local config if home dir can't be found
			configFile = "config.yaml"
		} else {
			configFile = defaultPath
		}
	}

	// 3. Ensure the config file exists (create if missing)
	if err := config.EnsureConfigFile(configFile); err != nil {
		// Log error but try to continue with defaults if possible
		fmt.Fprintf(os.Stderr, "Warning: failed to ensure config file: %v\n", err)
	}

	// 4. Load configuration
	// If the config file exists and is readable, load it. Otherwise, load with empty path to use defaults.
	loadPath := configFile
	if _, err := os.Stat(configFile); err != nil {
		loadPath = ""
	}

	cfg, err := loader.Load(loadPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return cfg, nil
}

// NewLogger creates a structured logger
func NewLogger(cfg *config.Config) *slog.Logger {
	return logger.New(cfg.Logger)
}

// NewBuildInfo provides build information to the application
func NewBuildInfo() version.BuildInfo {
	return version.BuildInfo{
		Version: buildVersion,
		Commit:  buildCommit,
		Date:    buildDate,
		BuiltBy: buildBy,
	}
}

// NewSystemOps creates a new SystemOps instance
func NewSystemOps() backend.SystemOps {
	return backend.NewDefaultSystemOps()
}

// NewDetector creates a new backend detector
func NewDetector(sysOps backend.SystemOps) status.Detector {
	return backend.NewDetector(sysOps)
}

// NewConfigReader creates a new DNS config reader
func NewConfigReader(sysOps backend.SystemOps) status.Reader {
	return backend.NewConfigReader(sysOps)
}

// RunCLI executes the CLI application
func RunCLI(lc fx.Lifecycle, rootCmd *cobra.Command, log *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Debug("CLI application starting")

			// Execute CLI in a goroutine to not block Fx lifecycle
			go func() {
				if err := cli.ExecuteContext(rootCmd); err != nil {
					exitCode := 1
					// If error has exit code, just use it
					if strings.HasPrefix(err.Error(), "exit:") {
						fmt.Sscanf(err.Error(), "exit:%d", &exitCode)
					} else {
						// Render error in user-friendly way
						styles := ui.NewStyles()
						fmt.Fprintln(os.Stderr, styles.RenderError(err.Error()))
						exitCode = 1
					}
					os.Exit(exitCode)
				}
				// Signal shutdown after CLI completes
				p, _ := os.FindProcess(os.Getpid())
				_ = p.Signal(syscall.SIGTERM)
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Debug("CLI application stopped")
			return nil
		},
	})
}

// RegisterLifecycleHooks registers global lifecycle hooks
func RegisterLifecycleHooks(lc fx.Lifecycle, log *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Debug("application started",
				slog.String("pid", fmt.Sprintf("%d", os.Getpid())),
			)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Debug("application stopped")
			// Ensure all logs are flushed
			return nil
		},
	})
}
