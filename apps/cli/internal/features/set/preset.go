package set

import (
	"fmt"
	"log/slog"

	"github.com/junevm/cdns/apps/cli/internal/config"
	"github.com/junevm/cdns/apps/cli/internal/ui"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// PresetCommandParams holds dependencies for preset command
type PresetCommandParams struct {
	fx.In

	Service *Service
	Logger  *slog.Logger
	Config  *config.Config
}

// PresetCommandResult wraps the command for Fx
type PresetCommandResult struct {
	fx.Out

	Cmd *cobra.Command `name:"set_preset"`
}

// NewPresetCommand creates the 'set preset' command
func NewPresetCommand(params PresetCommandParams) PresetCommandResult {
	var opts SetOptions

	cmd := &cobra.Command{
		Use:   "preset <name>",
		Short: "Apply a DNS preset",
		Long: `Apply a DNS preset to network interfaces.

Use 'cdns list' to see all available presets (Cloudflare, Google, Quad9, AdGuard, etc.).

Examples:
  # Apply Cloudflare DNS to active interfaces
  cdns set preset cloudflare

  # Apply Google DNS to specific interface
  cdns set preset google --interface eth0

  # Dry-run to preview changes
  cdns set preset quad9 --dry-run

  # Skip confirmation prompt
  cdns set preset opendns --yes`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			presetName := args[0]

			// Ensure privileges upfront (unless dry-run)
			if !opts.DryRun {
				if err := EnsurePrivileges(); err != nil {
					return err
				}
			}

			// Execute set preset
			err := params.Service.SetPreset(cmd.Context(), presetName, opts)
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

	return PresetCommandResult{Cmd: cmd}
}
