package set

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/junevm/cdns/apps/cli/internal/dns/presets"
)

var (
	// ErrEmptyDNSAddress is returned when DNS address is empty
	ErrEmptyDNSAddress = errors.New("DNS address cannot be empty")

	// ErrInvalidDNSAddress is returned when DNS address is not a valid IP
	ErrInvalidDNSAddress = errors.New("invalid DNS address")

	// ErrNoDNSAddresses is returned when no DNS addresses are provided
	ErrNoDNSAddresses = errors.New("at least one DNS address required")

	// ErrInvalidPresetName is returned when preset name is not found
	ErrInvalidPresetName = errors.New("invalid preset name")

	// ErrEmptyPresetName is returned when preset name is empty
	ErrEmptyPresetName = errors.New("preset name cannot be empty")

	// ErrInvalidInterfaceName is returned when interface name is invalid
	ErrInvalidInterfaceName = errors.New("invalid interface name")

	// ErrEmptyInterfaceName is returned when interface name is empty
	ErrEmptyInterfaceName = errors.New("interface name cannot be empty")
)

// ValidateDNSAddress validates a single DNS address (IPv4 or IPv6)
func ValidateDNSAddress(address string) error {
	if address == "" {
		return ErrEmptyDNSAddress
	}

	// Parse as IP address
	ip := net.ParseIP(address)
	if ip == nil {
		return fmt.Errorf("%w: %s", ErrInvalidDNSAddress, address)
	}

	return nil
}

// ValidateDNSAddresses validates a list of DNS addresses
func ValidateDNSAddresses(addresses []string) error {
	if len(addresses) == 0 {
		return ErrNoDNSAddresses
	}

	for _, addr := range addresses {
		if err := ValidateDNSAddress(addr); err != nil {
			return fmt.Errorf("invalid DNS address %q: %w", addr, err)
		}
	}

	return nil
}

// ValidatePresetName validates a preset name exists
func ValidatePresetName(name string) error {
	if name == "" {
		return ErrEmptyPresetName
	}

	// Canonicalize name for lookup
	name = strings.ToLower(name)

	_, ok := presets.Get(name)
	if !ok {
		return fmt.Errorf("%w: %s (use 'cdns list' to see all available presets)", ErrInvalidPresetName, name)
	}

	return nil
}

// ValidateInterfaceName validates a network interface name
// Valid names: alphanumeric, hyphens, underscores, max 15 chars
func ValidateInterfaceName(name string) error {
	if name == "" {
		return ErrEmptyInterfaceName
	}

	// Interface names on Linux are typically max 15 characters
	if len(name) > 15 {
		return fmt.Errorf("%w: name too long (max 15 characters)", ErrInvalidInterfaceName)
	}

	// Valid interface name pattern: alphanumeric, hyphens, underscores
	matched, err := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
	if err != nil {
		return fmt.Errorf("failed to validate interface name: %w", err)
	}

	if !matched {
		return fmt.Errorf("%w: %s (must contain only letters, numbers, hyphens, and underscores)", ErrInvalidInterfaceName, name)
	}

	return nil
}

// SeparateIPv4AndIPv6 separates a list of IP addresses into IPv4 and IPv6
func SeparateIPv4AndIPv6(addresses []string) (ipv4 []string, ipv6 []string) {
	for _, addr := range addresses {
		addr = strings.TrimSpace(addr)
		if addr == "" {
			continue
		}

		ip := net.ParseIP(addr)
		if ip == nil {
			continue // Skip invalid addresses
		}

		if ip.To4() != nil {
			ipv4 = append(ipv4, addr)
		} else {
			ipv6 = append(ipv6, addr)
		}
	}

	return ipv4, ipv6
}
