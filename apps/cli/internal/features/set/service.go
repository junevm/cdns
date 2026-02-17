package set

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/junevm/cdns/apps/cli/internal/config"
	"github.com/junevm/cdns/apps/cli/internal/dns/backend"
	"github.com/junevm/cdns/apps/cli/internal/dns/models"
	"github.com/junevm/cdns/apps/cli/internal/dns/presets"
	"github.com/junevm/cdns/apps/cli/internal/logger"
	"github.com/junevm/cdns/apps/cli/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	// ErrInsufficientPrivileges is returned when user doesn't have required privileges
	ErrInsufficientPrivileges = errors.New("insufficient privileges: root/administrator access required")

	// ErrUserCancelled is returned when user cancels the operation
	ErrUserCancelled = errors.New("operation cancelled by user")
)

// ExitCode represents command exit codes
type ExitCode int

const (
	ExitSuccess         ExitCode = 0
	ExitValidationError ExitCode = 2
	ExitPermissionError ExitCode = 3
	ExitPartialFailure  ExitCode = 4
)

// SetOptions holds options for the set operation
type SetOptions struct {
	Interfaces []string
	Scope      string // "active", "all", "explicit"
	DryRun     bool
	Yes        bool // Skip confirmation
	Verbose    bool // Show verbose logs
	PresetName string
}

// Detector interface for backend detection
type Detector interface {
	Detect() (models.Backend, error)
}

// Service handles set operations
type Service struct {
	config   *config.Config
	logger   *slog.Logger
	detector *backend.Detector
	writer   *backend.ConfigWriter
	reader   *backend.ConfigReader
	styles   *ui.Styles
}

// NewService creates a new set service
func NewService(cfg *config.Config, logger *slog.Logger, sysOps backend.SystemOps) *Service {
	return &Service{
		config:   cfg,
		logger:   logger,
		detector: backend.NewDetector(sysOps),
		writer:   backend.NewConfigWriter(sysOps),
		reader:   backend.NewConfigReader(sysOps),
		styles:   ui.NewStyles(),
	}
}

// IsInteractive checks if the stdout is a terminal
func (s *Service) IsInteractive() bool {
	return ui.IsTTY()
}

// RunInteractiveSet runs the interactive DNS configuration flow using Bubble Tea TUI
func (s *Service) RunInteractiveSet(ctx context.Context, opts SetOptions) error {
	if !s.IsInteractive() {
		return errors.New("interactive mode requires a terminal")
	}

	// check if nmcli is available
	if _, err := exec.LookPath("nmcli"); err != nil {
		fmt.Printf("%s %s\n\n", s.styles.Error.Render("Error:"), "NetworkManager (nmcli) is not installed.")
		fmt.Println("This tool requires nmcli to manage DNS settings.")
		fmt.Println("Please install it using your package manager:")
		fmt.Println("  Fedora/RHEL: sudo dnf install NetworkManager")
		fmt.Println("  Ubuntu/Debian: sudo apt install network-manager")
		fmt.Println("  Arch: sudo pacman -S networkmanager")
		return fmt.Errorf("nmcli not found")
	}

	// fetch available interfaces
	var ifaces []string
	// Use backend reader to get interfaces (we can improve this later to use proper structured data)
	// For now, let's use a simple detection similar to what was inline before, but better
	detected, err := s.detectInterfaces(ctx)
	if err != nil {
		s.logger.Warn("failed to detect interfaces", slog.Any("error", err))
		// continue with empty list or fallback
	} else {
		ifaces = detected
	}

	// Initialize TUI model
	m := newModel(s.config, ifaces)
	p := tea.NewProgram(m)

	// Run TUI
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run interactive mode: %w", err)
	}

	finalState := finalModel.(model)
	if finalState.quitting {
		return ErrUserCancelled
	}

	// Extract results
	if finalState.isCustom {
		// Parse custom DNS
		dnsList := strings.Split(finalState.customDNS, ",")
		for i, d := range dnsList {
			dnsList[i] = strings.TrimSpace(d)
		}
		// Set interface selection
		if finalState.selectedIface != "All Interfaces" {
			opts.Interfaces = []string{finalState.selectedIface}
		}
		// Interactive mode already confirms via TUI flow (selection implies intent)
		opts.Yes = true
		return s.SetCustom(ctx, dnsList, opts)
	} else {
		// Set preset
		if finalState.selectedIface != "All Interfaces" {
			opts.Interfaces = []string{finalState.selectedIface}
		}
		// Interactive mode already confirms via TUI flow (selection implies intent)
		opts.Yes = true
		return s.SetPreset(ctx, strings.ToLower(finalState.selectedPreset), opts)
	}
}

// SmartSet determines whether to use preset or custom DNS based on input
func (s *Service) SmartSet(ctx context.Context, args []string, opts SetOptions) error {
	if opts.Verbose {
		s.logger = logger.New(config.LoggerConfig{Level: "debug", Format: "text"})
		s.logger.Debug("verbose logging enabled")
	}

	if len(args) == 0 {
		return s.RunInteractiveSet(ctx, opts)
	}

	// Check if first arg is a preset (Built-in or Custom)
	presetName := strings.ToLower(args[0])
	_, isBuiltin := presets.Get(presetName)
	isCustom := false
	if s.config != nil && s.config.DNS.CustomPresets != nil {
		_, isCustom = s.config.DNS.CustomPresets[presetName]
	}

	if isBuiltin || isCustom {
		return s.SetPreset(ctx, presetName, opts)
	}

	// Check if args are valid IPs (Custom DNS)
	if err := ValidateDNSAddresses(args); err == nil {
		return s.SetCustom(ctx, args, opts)
	}

	return fmt.Errorf("invalid argument '%s': not a known preset or valid IP address", args[0])
}

// detectInterfaces returns a list of active network interfaces
func (s *Service) detectInterfaces(ctx context.Context) ([]string, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}

	// Attempt to detect connected interfaces via nmcli
	cmd := exec.CommandContext(ctx, "nmcli", "-t", "-f", "DEVICE,STATE", "device", "status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("nmcli error: %s: %w", strings.TrimSpace(string(output)), err)
	}

	var interfaces []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")
		if len(parts) >= 2 && parts[1] == "connected" {
			interfaces = append(interfaces, parts[0])
		}
	}
	return interfaces, nil
}

func CapitalizePresetName(name string) string {
	switch name {
	case "opendns":
		return "OpenDNS"
	case "adguard":
		return "AdGuard"
	default:
		if len(name) == 0 {
			return name
		}
		return strings.ToUpper(name[:1]) + name[1:]
	}
}

// SetPreset applies a DNS preset
func (s *Service) SetPreset(ctx context.Context, presetName string, opts SetOptions) error {
	presetName = strings.ToLower(presetName)

	// 1. Check custom presets from config first
	if s.config != nil && s.config.DNS.CustomPresets != nil {
		if ips, ok := s.config.DNS.CustomPresets[presetName]; ok {
			s.logger.Debug("using custom preset from config", slog.String("name", presetName))
			opts.PresetName = strings.ToUpper(presetName)
			return s.setDNS(ctx, ips, opts)
		}
	}

	// 2. Check built-in presets
	if preset, ok := presets.Get(presetName); ok {
		// Combine IPv4 and IPv6 addresses
		dnsAddresses := append(preset.IPv4, preset.IPv6...)
		opts.PresetName = CapitalizePresetName(presetName)
		return s.setDNS(ctx, dnsAddresses, opts)
	}

	return fmt.Errorf("validation failed: %w: %s", ErrInvalidPresetName, presetName)
}

// SetCustom applies custom DNS servers
func (s *Service) SetCustom(ctx context.Context, dnsAddresses []string, opts SetOptions) error {
	// Validate DNS addresses
	if err := ValidateDNSAddresses(dnsAddresses); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return s.setDNS(ctx, dnsAddresses, opts)
}

// setDNS is the internal method that applies DNS settings
func (s *Service) setDNS(ctx context.Context, dnsAddresses []string, opts SetOptions) error {
	// Validate interface names if provided
	for _, iface := range opts.Interfaces {
		if err := ValidateInterfaceName(iface); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	}

	// Detect backend
	backendObj, err := s.detector.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect DNS backend: %w", err)
	}

	s.logger.Debug("detected backend", slog.String("backend", string(backendObj)))

	// Identify target interfaces
	var targetInterfaces []string
	if len(opts.Interfaces) > 0 {
		targetInterfaces = opts.Interfaces
	} else {
		// Auto-detect active interfaces
		detected, err := s.detectInterfaces(ctx)
		if err != nil {
			s.logger.Warn("failed to detect active interfaces", slog.Any("error", err))
		}
		if len(detected) > 0 {
			targetInterfaces = detected
		} else {
			// Fallback (only if absolutely no detection possible)
			targetInterfaces = []string{"eth0"}
		}
	}

	// Prepare config for all target interfaces
	appliedConfigs := make([]models.DNSConfig, 0, len(targetInterfaces))

	// Clean and separate addresses
	ipv4, ipv6 := SeparateIPv4AndIPv6(dnsAddresses)

	for _, iface := range targetInterfaces {
		appliedConfigs = append(appliedConfigs, models.DNSConfig{
			Interface: models.NetworkInterface{Name: iface, Backend: backendObj},
			DNS:       models.DNSServer{IPv4: ipv4, IPv6: ipv6},
		})
	}

	// Dry-run mode: show what would change and exit
	if opts.DryRun {
		return s.showDryRun(backendObj, dnsAddresses, appliedConfigs, opts)
	}

	// Confirmation prompt (only if interactive and not suppressed by --yes)
	if !opts.Yes && s.IsInteractive() {
		confirmed, err := s.confirmChange(dnsAddresses, targetInterfaces, opts)
		if err != nil {
			return err
		}
		if !confirmed {
			return ErrUserCancelled
		}
	}

	// Apply DNS changes via backend
	// KISS: No "Applying..." spinner mess unless logic is slow. nmcli is fast.
	if err := s.writer.Apply(ctx, backendObj, appliedConfigs); err != nil {
		return fmt.Errorf("failed to apply DNS: %w", err)
	}

	s.logger.Debug("DNS settings applied",
		slog.Any("dns", dnsAddresses),
		slog.Any("interfaces", targetInterfaces),
		slog.String("backend", string(backendObj)))

	// Minimal feedback
	if s.IsInteractive() {
		fmt.Printf("%s Applied DNS (%s) to %s.\n",
			s.styles.Success.Render("✔"),
			s.styles.RenderInfo(strings.Join(dnsAddresses, ", ")),
			s.styles.RenderBold(strings.Join(targetInterfaces, ", ")))
	} else {
		// Even more minimal for non-interactive (piped output)
		if opts.PresetName != "" {
			fmt.Printf("DNS set to %s\n", opts.PresetName)
		} else {
			fmt.Printf("DNS set to %s\n", strings.Join(dnsAddresses, ", "))
		}
	}

	return nil
}

// showDryRun displays what would change without applying
func (s *Service) showDryRun(backend models.Backend, dnsAddresses []string, configs []models.DNSConfig, opts SetOptions) error {
	fmt.Printf("%s\n\n", s.styles.RenderBold("Dry-run mode: No changes will be applied"))
	fmt.Printf("Backend: %s\n", s.styles.RenderInfo(string(backend)))
	fmt.Printf("DNS servers to set:\n")
	for _, dns := range dnsAddresses {
		fmt.Printf("  - %s\n", s.styles.RenderInfo(dns))
	}

	fmt.Printf("Target Interfaces:\n")
	for _, cfg := range configs {
		fmt.Printf("  - %s\n", s.styles.RenderBold(cfg.Interface.Name))
	}

	return nil
}

// confirmChange prompts user to confirm the change
func (s *Service) confirmChange(dnsAddresses []string, interfaces []string, opts SetOptions) (bool, error) {
	fmt.Printf("\n%s\n", s.styles.RenderWarning("This will change DNS settings:"))
	fmt.Printf("  DNS: %s\n", s.styles.RenderInfo(strings.Join(dnsAddresses, ", ")))
	fmt.Printf("  Interfaces: %s\n", s.styles.RenderBold(strings.Join(interfaces, ", ")))

	fmt.Printf("\nContinue? [%s/%s]: ", s.styles.RenderBold("y"), "N")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read confirmation: %w", err)
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes", nil
}

// ExitCodeFromError determines the exit code from an error
func ExitCodeFromError(err error) ExitCode {
	if err == nil {
		return ExitSuccess
	}

	switch {
	case errors.Is(err, ErrInsufficientPrivileges):
		return ExitPermissionError
	case errors.Is(err, ErrInvalidDNSAddress),
		errors.Is(err, ErrNoDNSAddresses),
		errors.Is(err, ErrInvalidPresetName),
		errors.Is(err, ErrEmptyPresetName),
		errors.Is(err, ErrInvalidInterfaceName),
		errors.Is(err, ErrEmptyInterfaceName):
		return ExitValidationError
	default:
		// TODO: Add partial failure detection
		return ExitPartialFailure
	}
}
