package backend

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"cli/internal/dns/models"
	"cli/internal/features/status"
)

// ConfigReader reads DNS configuration from different backends
type ConfigReader struct {
	sysOps SystemOps
}

// NewConfigReader creates a new ConfigReader
func NewConfigReader(sysOps SystemOps) *ConfigReader {
	return &ConfigReader{sysOps: sysOps}
}

// ReadDNSConfig reads DNS configuration from the specified backend
func (r *ConfigReader) ReadDNSConfig(ctx context.Context, backend models.Backend) (*status.StatusInfo, error) {
	switch backend {
	case models.BackendNetworkManager:
		return r.readNetworkManager(ctx)
	case models.BackendSystemdResolved:
		return r.readSystemdResolved(ctx)
	case models.BackendResolvConf:
		return r.readResolvConf(ctx)
	default:
		return nil, fmt.Errorf("unsupported backend: %s", backend)
	}
}

// readNetworkManager reads DNS configuration from NetworkManager
func (r *ConfigReader) readNetworkManager(ctx context.Context) (*status.StatusInfo, error) {
	info := &status.StatusInfo{
		Backend:    models.BackendNetworkManager,
		Interfaces: []status.InterfaceStatus{},
		Managed:    true,
		Warnings:   []string{},
	}

	// Get list of active connections
	cmd := exec.CommandContext(ctx, "nmcli", "-t", "-f", "NAME,DEVICE", "connection", "show", "--active")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get active connections: %w", err)
	}

	// Parse active connections
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}

		device := parts[1]
		if device == "" {
			continue
		}

		// Get DNS for this connection
		ipv4, ipv6, err := r.getDNSForConnection(ctx, parts[0])
		if err != nil {
			info.Warnings = append(info.Warnings, fmt.Sprintf("failed to get DNS for %s: %v", device, err))
			continue
		}

		if len(ipv4) > 0 || len(ipv6) > 0 {
			info.Interfaces = append(info.Interfaces, status.InterfaceStatus{
				Name: device,
				IPv4: ipv4,
				IPv6: ipv6,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading nmcli output: %w", err)
	}

	return info, nil
}

// getDNSForConnection gets DNS servers for a specific connection
func (r *ConfigReader) getDNSForConnection(ctx context.Context, connName string) ([]string, []string, error) {
	cmd := exec.CommandContext(ctx, "nmcli", "-t", "-f", "IP4.DNS,IP6.DNS", "connection", "show", connName)
	output, err := cmd.Output()
	if err != nil {
		return nil, nil, err
	}

	var ipv4, ipv6 []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "IP4.DNS") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 && parts[1] != "" {
				ipv4 = append(ipv4, strings.TrimSpace(parts[1]))
			}
		} else if strings.HasPrefix(line, "IP6.DNS") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 && parts[1] != "" {
				ipv6 = append(ipv6, strings.TrimSpace(parts[1]))
			}
		}
	}

	return ipv4, ipv6, scanner.Err()
}

// readSystemdResolved reads DNS configuration from systemd-resolved
func (r *ConfigReader) readSystemdResolved(ctx context.Context) (*status.StatusInfo, error) {
	info := &status.StatusInfo{
		Backend:    models.BackendSystemdResolved,
		Interfaces: []status.InterfaceStatus{},
		Managed:    true,
		Warnings:   []string{},
	}

	// Try resolvectl first, fall back to systemd-resolve
	var cmd *exec.Cmd
	if r.sysOps.CommandExists("resolvectl") {
		cmd = exec.CommandContext(ctx, "resolvectl", "status")
	} else if r.sysOps.CommandExists("systemd-resolve") {
		cmd = exec.CommandContext(ctx, "systemd-resolve", "--status")
	} else {
		return nil, fmt.Errorf("neither resolvectl nor systemd-resolve found")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get systemd-resolved status: %w", err)
	}

	// Parse the output
	info.Interfaces = r.parseSystemdResolvedOutput(string(output))

	return info, nil
}

// parseSystemdResolvedOutput parses systemd-resolved status output
func (r *ConfigReader) parseSystemdResolvedOutput(output string) []status.InterfaceStatus {
	var interfaces []status.InterfaceStatus
	var currentInterface *status.InterfaceStatus

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Look for link/interface lines
		if strings.HasPrefix(line, "Link ") {
			// Save previous interface if exists
			if currentInterface != nil && (len(currentInterface.IPv4) > 0 || len(currentInterface.IPv6) > 0) {
				interfaces = append(interfaces, *currentInterface)
			}

			// Extract interface name
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				// Format: "Link 2 (eth0)"
				name := strings.Trim(parts[2], "()")
				currentInterface = &status.InterfaceStatus{
					Name: name,
					IPv4: []string{},
					IPv6: []string{},
				}
			}
		} else if currentInterface != nil && strings.HasPrefix(line, "Current DNS Server:") {
			// Extract DNS server address
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				addr := strings.TrimSpace(parts[1])
				if strings.Contains(addr, ":") {
					currentInterface.IPv6 = append(currentInterface.IPv6, addr)
				} else {
					currentInterface.IPv4 = append(currentInterface.IPv4, addr)
				}
			}
		} else if currentInterface != nil && strings.HasPrefix(line, "DNS Servers:") {
			// Extract DNS server address
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				addr := strings.TrimSpace(parts[1])
				if strings.Contains(addr, ":") {
					currentInterface.IPv6 = append(currentInterface.IPv6, addr)
				} else {
					currentInterface.IPv4 = append(currentInterface.IPv4, addr)
				}
			}
		} else if currentInterface != nil && strings.HasPrefix(line, "- ") {
			// Additional DNS server in a list
			addr := strings.TrimPrefix(line, "- ")
			addr = strings.TrimSpace(addr)
			if addr != "" {
				if strings.Contains(addr, ":") {
					currentInterface.IPv6 = append(currentInterface.IPv6, addr)
				} else {
					currentInterface.IPv4 = append(currentInterface.IPv4, addr)
				}
			}
		}
	}

	// Add last interface
	if currentInterface != nil && (len(currentInterface.IPv4) > 0 || len(currentInterface.IPv6) > 0) {
		interfaces = append(interfaces, *currentInterface)
	}

	return interfaces
}

// readResolvConf reads DNS configuration from /etc/resolv.conf
func (r *ConfigReader) readResolvConf(ctx context.Context) (*status.StatusInfo, error) {
	info := &status.StatusInfo{
		Backend:    models.BackendResolvConf,
		Interfaces: []status.InterfaceStatus{},
		Managed:    false,
		Warnings:   []string{},
	}

	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return nil, fmt.Errorf("failed to open /etc/resolv.conf: %w", err)
	}
	defer file.Close()

	var ipv4, ipv6 []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		// Look for nameserver entries
		if strings.HasPrefix(line, "nameserver") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				addr := fields[1]
				if strings.Contains(addr, ":") {
					ipv6 = append(ipv6, addr)
				} else {
					ipv4 = append(ipv4, addr)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading /etc/resolv.conf: %w", err)
	}

	// Add single "system" interface for resolv.conf
	if len(ipv4) > 0 || len(ipv6) > 0 {
		info.Interfaces = append(info.Interfaces, status.InterfaceStatus{
			Name: "system",
			IPv4: ipv4,
			IPv6: ipv6,
		})
	}

	return info, nil
}
