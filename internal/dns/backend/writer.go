package backend

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/junevm/cdns/internal/dns/models"
)

// ConfigWriter handles applying DNS configurations to the system
type ConfigWriter struct {
	sysOps SystemOps
}

// NewConfigWriter creates a new ConfigWriter
func NewConfigWriter(sysOps SystemOps) *ConfigWriter {
	return &ConfigWriter{sysOps: sysOps}
}

// Apply applies the DNS configuration using the specified backend
func (w *ConfigWriter) Apply(ctx context.Context, backend models.Backend, configs []models.DNSConfig) error {
	switch backend {
	case models.BackendNetworkManager:
		return w.applyNetworkManager(ctx, configs)
	case models.BackendSystemdResolved:
		return w.applySystemdResolved(ctx, configs)
	default:
		return fmt.Errorf("unsupported backend for writing: %s", backend)
	}
}

func (w *ConfigWriter) applyNetworkManager(ctx context.Context, configs []models.DNSConfig) error {
	for _, cfg := range configs {
		if cfg.Interface.Name == "" {
			continue
		}

		// Get active connection name
		connName, err := w.getNMConnection(ctx, cfg.Interface.Name)
		if err != nil {
			// Fallback: If we can't get connection, try device modify (transient)
			// This might happen if device is unmanaged or something weird.
			// Ideally we log this.
			return fmt.Errorf("failed to get active connection for %s: %w", cfg.Interface.Name, err)
		}

		// Set IPv4 DNS
		if len(cfg.DNS.IPv4) > 0 {
			dnsStr := strings.Join(cfg.DNS.IPv4, " ")
			// Modify connection for persistence
			// ipv4.ignore-auto-dns yes ensures DHCP doesn't override it
			cmd := exec.CommandContext(ctx, "nmcli", "connection", "modify", connName, "ipv4.dns", dnsStr, "ipv4.ignore-auto-dns", "yes")
			if output, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("failed to set IPv4 DNS for %s (conn: %s): %s: %w", cfg.Interface.Name, connName, strings.TrimSpace(string(output)), err)
			}
		}

		// Set IPv6 DNS
		if len(cfg.DNS.IPv6) > 0 {
			dnsStr := strings.Join(cfg.DNS.IPv6, " ")
			cmd := exec.CommandContext(ctx, "nmcli", "connection", "modify", connName, "ipv6.dns", dnsStr, "ipv6.ignore-auto-dns", "yes")
			if output, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("failed to set IPv6 DNS for %s (conn: %s): %s: %w", cfg.Interface.Name, connName, strings.TrimSpace(string(output)), err)
			}
		}

		// Reapply changes to the device (runtime)
		// This makes the changes effective immediately without interface bounce usually
		cmd := exec.CommandContext(ctx, "nmcli", "device", "reapply", cfg.Interface.Name)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to reapply configuration on device %s: %s: %w", cfg.Interface.Name, string(output), err)
		}
	}
	return nil
}

func (w *ConfigWriter) getNMConnection(ctx context.Context, iface string) (string, error) {
	// usage: nmcli -g GENERAL.CONNECTION device show <iface>
	// -g prints just the value
	cmd := exec.CommandContext(ctx, "nmcli", "-g", "GENERAL.CONNECTION", "device", "show", iface)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("nmcli error: %s: %w", strings.TrimSpace(string(out)), err)
	}
	res := strings.TrimSpace(string(out))
	if res == "" {
		return "", fmt.Errorf("no active connection found for interface %s", iface)
	}
	return res, nil
}

func (w *ConfigWriter) applySystemdResolved(ctx context.Context, configs []models.DNSConfig) error {
	for _, cfg := range configs {
		if cfg.Interface.Name == "" {
			continue
		}

		allDNS := append(cfg.DNS.IPv4, cfg.DNS.IPv6...)
		if len(allDNS) == 0 {
			continue
		}

		args := []string{"dns", cfg.Interface.Name}
		args = append(args, allDNS...)

		// note: resolvectl is transient.
		cmd := exec.CommandContext(ctx, "resolvectl", args...)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to set DNS for %s via resolvectl: %s: %w", cfg.Interface.Name, strings.TrimSpace(string(output)), err)
		}
	}
	return nil
}

// ResetToAutomatic resets the DNS configuration for the specified interfaces to automatic (DHCP)
func (w *ConfigWriter) ResetToAutomatic(ctx context.Context, backend models.Backend, interfaces []string) error {
	switch backend {
	case models.BackendNetworkManager:
		return w.resetNetworkManager(ctx, interfaces)
	case models.BackendSystemdResolved:
		return w.resetSystemdResolved(ctx, interfaces)
	default:
		return fmt.Errorf("unsupported backend for reset: %s", backend)
	}
}

func (w *ConfigWriter) resetNetworkManager(ctx context.Context, interfaces []string) error {
	for _, iface := range interfaces {
		connName, err := w.getNMConnection(ctx, iface)
		if err != nil {
			return fmt.Errorf("failed to get connection for %s: %w", iface, err)
		}

		// Reset IPv4
		cmd := exec.CommandContext(ctx, "nmcli", "connection", "modify", connName, "ipv4.dns", "", "ipv4.ignore-auto-dns", "no")
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to reset IPv4 DNS for %s: %s: %w", iface, strings.TrimSpace(string(output)), err)
		}

		// Reset IPv6
		cmd = exec.CommandContext(ctx, "nmcli", "connection", "modify", connName, "ipv6.dns", "", "ipv6.ignore-auto-dns", "no")
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to reset IPv6 DNS for %s: %s: %w", iface, strings.TrimSpace(string(output)), err)
		}

		// Reapply
		cmd = exec.CommandContext(ctx, "nmcli", "device", "reapply", iface)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to reapply configuration on %s: %s: %w", iface, strings.TrimSpace(string(output)), err)
		}
	}
	return nil
}

func (w *ConfigWriter) resetSystemdResolved(ctx context.Context, interfaces []string) error {
	for _, iface := range interfaces {
		// resolvectl revert <interface> resets interface-specific DNS settings
		cmd := exec.CommandContext(ctx, "resolvectl", "revert", iface)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to revert DNS for %s: %s: %w", iface, strings.TrimSpace(string(output)), err)
		}
	}
	return nil
}
