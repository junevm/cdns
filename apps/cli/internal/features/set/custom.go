package set

import (
	"fmt"
	"log/slog"

	"cli/internal/config"
	"cli/internal/ui"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// CustomCommandParams holds dependencies for custom command
type CustomCommandParams struct {
	fx.In

	Service *Service
	Logger  *slog.Logger
	Config  *config.Config
}

// CustomCommandResult wraps the command for Fx
type CustomCommandResult struct {
	fx.Out

	Cmd *cobra.Command `name:"set_custom"`
}

// NewCustomCommand creates the 'set custom' command
func NewCustomCommand(params CustomCommandParams) CustomCommandResult {
	var opts SetOptions

	cmd := &cobra.Command{
		Use:   "custom <dns1> [dns2...]",
		Short: "Apply custom DNS servers",
		Long: `Apply custom DNS servers to network interfaces.

Supports both IPv4 and IPv6 addresses. At least one DNS server is required.

Examples:
  # Set custom DNS (IPv4 only)
  cdns set custom 8.8.8.8 8.8.4.4

  # Set custom DNS (mixed IPv4 and IPv6)
  cdns set custom 1.1.1.1 2606:4700:4700::1111

  # Set custom DNS on specific interface
  cdns set custom 9.9.9.9 --interface wlan0

  # Dry-run to preview changes
  cdns set custom 208.67.222.222 208.67.220.220 --dry-run

  # Skip confirmation prompt
  cdns set custom 94.140.14.14 --yes`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dnsAddresses := args

			// Ensure privileges upfront (unless dry-run)
			if !opts.DryRun {
				if err := EnsurePrivileges(); err != nil {
					return err
				}
			}

			// Execute set custom
			err := params.Service.SetCustom(cmd.Context(), dnsAddresses, opts)
			if err != nil {
				// Get exit code
				exitCode := ExitCodeFromError(err)

				// Print user-friendly error
				styles := ui.NewStyles()
				fmt.Fprintln(cmd.ErrOrStderr(), styles.RenderError(err.Error()))

				// Exit with appropriate code
				cmd.SilenceUsage = true
				return fmt.Errorf("exit:%d", exitCode)
			}

			return nil
		},
	}

	// Flags
	// Use defaults from config if available
	defaultInterfaces := params.Config.DNS.DefaultInterfaces
	if defaultInterfaces == nil {
		defaultInterfaces = []string{}
	}
	defaultScope := params.Config.DNS.DefaultScope
	if defaultScope == "" {
		defaultScope = "active"
	}

	cmd.Flags().StringSliceVar(&opts.Interfaces, "interface", defaultInterfaces, "specify interface(s) (repeatable)")
	cmd.Flags().StringVar(&opts.Scope, "scope", defaultScope, "interface scope: active, all, or explicit")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "preview changes without applying")
	cmd.Flags().BoolVar(&opts.Yes, "yes", false, "skip confirmation prompts")

	return CustomCommandResult{Cmd: cmd}
}
