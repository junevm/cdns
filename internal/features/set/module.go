package set

import (
	"fmt"

	"github.com/junevm/cdns/internal/config"
	"github.com/junevm/cdns/internal/ui"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// Module provides the set feature as an Fx module
var Module = fx.Module("set",
	fx.Provide(NewService),
	fx.Provide(NewSetCommand),
	fx.Provide(NewPresetCommand),
	fx.Provide(NewCustomCommand),
	fx.Invoke(RegisterCommands),
)

// SetCommandResult wraps the parent set command
type SetCommandResult struct {
	fx.Out

	Cmd *cobra.Command `name:"set"`
}

// NewSetCommand creates the parent 'set' command
func NewSetCommand(params struct {
	fx.In
	Service *Service
	Config  *config.Config
}) SetCommandResult {
	var opts SetOptions

	cmd := &cobra.Command{
		Use:   "set [preset|ip...]",
		Short: "Set DNS servers",
		Long: `Set DNS servers using a preset name or custom IP addresses.

Examples:
  # Interactive mode
  cdns set

  # Set preset
  cdns set google

  # Set custom IPs
  cdns set 1.1.1.1 8.8.8.8

  # Set specific interface
  cdns set cloudflare --interface eth0`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Merge persistent flags from root
			if !opts.Verbose {
				val, _ := cmd.Flags().GetBool("verbose")
				opts.Verbose = val
			}

			// Ensure privileges upfront for better UX (unless dry-run)
			if !opts.DryRun {
				if err := EnsurePrivileges(); err != nil {
					return err
				}
			}

			err := params.Service.SmartSet(cmd.Context(), args, opts)
			if err != nil {
				// Get exit code
				exitCode := ExitCodeFromError(err)

				// Print user-friendly error unless it's just user cancel
				if err != ErrUserCancelled {
					fmt.Fprintln(cmd.ErrOrStderr(), ui.NewStyles().RenderError(err.Error()))
				}

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

	cmd.Flags().StringSliceVarP(&opts.Interfaces, "interface", "i", defaultInterfaces, "specify interface(s) (repeatable)")
	cmd.Flags().StringVar(&opts.Scope, "scope", "active", "interface scope: active, all, or explicit")
	cmd.Flags().Lookup("scope").DefValue = defaultScope
	opts.Scope = defaultScope // Ensure initialized with config value

	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "preview changes without applying")
	cmd.Flags().BoolVar(&opts.Yes, "yes", false, "skip confirmation prompts")

	return SetCommandResult{Cmd: cmd}
}

// RegisterCommandsParams holds dependencies for command registration
type RegisterCommandsParams struct {
	fx.In

	RootCmd   *cobra.Command
	SetCmd    *cobra.Command `name:"set"`
	PresetCmd *cobra.Command `name:"set_preset"`
	CustomCmd *cobra.Command `name:"set_custom"`
}

// RegisterCommands registers all set commands
func RegisterCommands(params RegisterCommandsParams) {
	// Add subcommands to set command
	params.SetCmd.AddCommand(params.PresetCmd)
	params.SetCmd.AddCommand(params.CustomCmd)

	// Add set command to root
	params.RootCmd.AddCommand(params.SetCmd)

	// Make set the default action when no subcommand is provided
	// This makes 'cdns' equivalent to 'cdns set' (Interactive TUI)
	if params.RootCmd.RunE == nil {
		params.RootCmd.RunE = params.SetCmd.RunE
	}
}
