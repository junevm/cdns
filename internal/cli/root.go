package cli

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/junevm/cdns/internal/config"
	"github.com/junevm/cdns/internal/ui"

	"github.com/spf13/cobra"
)

// Dependencies holds all dependencies needed by CLI commands
type Dependencies struct {
	Config *config.Config
	Logger *slog.Logger
}

// NewRootCmd creates the root command with dependency injection
// This follows the Command Factory pattern, avoiding global state
func NewRootCmd(deps Dependencies) *cobra.Command {
	var cfgFile string

	rootCmd := &cobra.Command{
		Use:           "cdns",
		Short:         "A trusted, Linux-first DNS management CLI tool",
		Long:          ui.GetBanner() + "\n\nA trusted, Linux-first DNS management CLI tool",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no subcommand is provided, show the main menu
			choice, err := RunMainMenu()
			if err != nil {
				return err
			}

			if choice == "" {
				return nil
			}

			// Find and execute the selected subcommand
			for _, sub := range cmd.Commands() {
				if sub.Name() == choice {
					// Propagate context from parent command
					if cmd.Context() != nil {
						sub.SetContext(cmd.Context())
					}

					// Execute the command logic
					if sub.RunE != nil {
						return sub.RunE(sub, nil)
					}
					if sub.Run != nil {
						sub.Run(sub, nil)
						return nil
					}
					return nil
				}
			}

			return nil
		},
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.config/cdns/config.yaml)")
	rootCmd.PersistentFlags().String("log-level", "warn", "log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "show verbose logs")

	return rootCmd
}

// AddCommand adds a subcommand to the root command with panic recovery
func AddCommand(root *cobra.Command, cmd *cobra.Command) {
	// Wrap RunE with panic recovery middleware
	if cmd.RunE != nil {
		originalRunE := cmd.RunE
		cmd.RunE = WithPanicRecovery(originalRunE)
	}

	root.AddCommand(cmd)
}

// WithPanicRecovery wraps a cobra RunE function with panic recovery
// This is the panic middleware mentioned in the requirements
func WithPanicRecovery(fn func(*cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) (err error) {
		defer func() {
			if r := recover(); r != nil {
				// Get logger from command context if available
				logger := slog.Default()
				if cmd.Context() != nil {
					if ctxLogger, ok := cmd.Context().Value("logger").(*slog.Logger); ok {
						logger = ctxLogger
					}
				}

				// Log the panic securely
				logger.Error("panic recovered",
					slog.Any("panic", r),
					slog.String("stack", string(debug.Stack())),
				)

				// Return as error
				err = fmt.Errorf("command panicked: %v", r)
			}
		}()

		return fn(cmd, args)
	}
}

// ExecuteContext executes the root command with context
func ExecuteContext(rootCmd *cobra.Command) error {
	return rootCmd.Execute()
}
