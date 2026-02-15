# DNS Package

This package provides DNS preset configurations and domain models for the cdns CLI tool.

## Structure

- `models/` - Core domain models for DNS configuration
- `presets/` - Pre-configured DNS server definitions

## Models Package

The models package defines the core domain types used throughout the DNS changer:

### Types

- **Backend**: Represents the DNS management system (NetworkManager, systemd-resolved, resolv.conf, netplan)
- **DNSServer**: Holds DNS server IPv4 and IPv6 addresses
- **NetworkInterface**: Represents a network interface with its backend
- **DNSConfig**: Complete DNS configuration for an interface

### Usage Example

```go
import "cli/internal/dns/models"

// Create a DNS server configuration
dns := models.DNSServer{
    IPv4: []string{"1.1.1.1", "1.0.0.1"},
    IPv6: []string{"2606:4700:4700::1111"},
}

// Create a network interface
iface := models.NetworkInterface{
    Name: "eth0",
    Backend: models.BackendNetworkManager,
}

// Combine into a complete configuration
config := models.DNSConfig{
    Interface: iface,
    DNS: dns,
}
```

## Presets Package

The presets package provides pre-configured DNS server definitions for popular DNS providers.

### Available Presets

- **Cloudflare**: Fast and privacy-focused DNS (1.1.1.1)
- **Google**: Reliable public DNS (8.8.8.8)
- **Quad9**: Security and privacy-focused DNS (9.9.9.9)
- **OpenDNS**: Family-friendly DNS with filtering (208.67.222.222)
- **AdGuard**: DNS with ad and tracker blocking (94.140.14.14)

### Usage Example

```go
import "cli/internal/dns/presets"

// Get a specific preset
cloudflare := presets.Cloudflare()
google := presets.Google()

// Get a preset by name
dns, ok := presets.Get("cloudflare")
if !ok {
    // preset not found
}

// Get all available presets
all := presets.All()
for name, dns := range all {
    fmt.Printf("%s: %v\n", name, dns.IPv4)
}
```

### Preset Details

#### Cloudflare
- IPv4: 1.1.1.1, 1.0.0.1
- IPv6: 2606:4700:4700::1111, 2606:4700:4700::1001

#### Google
- IPv4: 8.8.8.8, 8.8.4.4
- IPv6: 2001:4860:4860::8888, 2001:4860:4860::8844

#### Quad9
- IPv4: 9.9.9.9, 149.112.112.112
- IPv6: 2620:fe::fe, 2620:fe::9

#### OpenDNS
- IPv4: 208.67.222.222, 208.67.220.220
- IPv6: 2620:119:35::35, 2620:119:53::53

#### AdGuard
- IPv4: 94.140.14.14, 94.140.15.15
- IPv6: 2a10:50c0::ad1:ff, 2a10:50c0::ad2:ff

## Testing

All packages include comprehensive tests following Go best practices:

```bash
# Test models
go test ./internal/dns/models -v

# Test presets
go test ./internal/dns/presets -v

# Test with coverage
go test ./internal/dns/... -cover
```

## Design Principles

1. **Immutability**: Presets return copies, preventing accidental modification
2. **Simplicity**: Models are plain structs with no business logic
3. **Framework-agnostic**: No dependencies on CLI frameworks or external libraries
4. **Type-safety**: Backend types use constants to prevent invalid values
5. **Testability**: 100% test coverage with table-driven tests
