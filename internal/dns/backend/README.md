# DNS Backend Detection

This package provides DNS backend detection functionality for Linux systems. It identifies which DNS management system is active on the current system and provides a consistent interface for working with different backends.

## Supported Backends

The detector supports the following backends in priority order:

1. **NetworkManager** - Modern network configuration manager
2. **systemd-resolved** - systemd's DNS resolver service
3. **resolv.conf** - Traditional unmanaged /etc/resolv.conf

## Usage

### Basic Detection

```go
package main

import (
    "fmt"
    "log"
    
    "cli/internal/dns/backend"
)

func main() {
    // Create detector with default system operations
    detector := backend.NewDetector(backend.NewDefaultSystemOps())
    
    // Detect active backend
    activeBackend, err := detector.Detect()
    if err != nil {
        log.Fatalf("Failed to detect backend: %v", err)
    }
    
    fmt.Printf("Active backend: %s\n", activeBackend)
}
```

### Detection with Reason

To get detailed information about why a particular backend was selected:

```go
backend, reason, err := detector.DetectWithReason()
if err != nil {
    log.Fatalf("Failed to detect backend: %v", err)
}

fmt.Printf("Backend: %s\nReason: %s\n", backend, reason)
```

### Testing with Mock System Operations

For testing, you can inject a mock implementation of `SystemOps`:

```go
type mockSystemOps struct {
    commandExists  func(string) bool
    serviceRunning func(string) (bool, error)
    // ... other methods
}

func (m *mockSystemOps) CommandExists(cmd string) bool {
    if m.commandExists != nil {
        return m.commandExists(cmd)
    }
    return false
}

// Create detector with mock
detector := backend.NewDetector(&mockSystemOps{
    commandExists: func(cmd string) bool {
        return cmd == "nmcli"
    },
    serviceRunning: func(service string) (bool, error) {
        return service == "NetworkManager", nil
    },
})
```

## Detection Logic

### NetworkManager

Detected if:
- `nmcli` command is available in PATH
- NetworkManager service is running (checked via systemctl)

### systemd-resolved

Detected if:
- `resolvectl` or `systemd-resolve` command is available in PATH
- systemd-resolved service is running (checked via systemctl)

### resolv.conf

Detected if:
- `/etc/resolv.conf` exists
- `/etc/resolv.conf` is a regular file (not a symlink)
- No higher-priority backend is active

## Requirements

- **Detection only**: Does not require root privileges
- **Modification**: Backend-specific operations may require appropriate permissions

## Design Principles

1. **Dependency Injection**: All system operations are abstracted behind the `SystemOps` interface for testability
2. **Deterministic**: Same system state always produces the same detection result
3. **Priority-based**: Higher-priority backends are preferred when multiple are available
4. **Non-destructive**: Detection never modifies system state

## Error Handling

The detector returns errors in the following cases:

- System service status check fails (not just inactive)
- File system operations fail unexpectedly
- No supported backend is found
- `/etc/resolv.conf` is a symlink (managed by another service)

## Integration

This package is designed to work with the DNS configuration system in `internal/dns/models`. The detected backend type matches the `Backend` enum defined in that package.
