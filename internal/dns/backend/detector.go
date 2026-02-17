package backend

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/junevm/cdns/internal/dns/models"
)

// SystemOps defines operations for interacting with the system
// This interface allows for dependency injection and testing
type SystemOps interface {
	CommandExists(cmd string) bool
	ServiceRunning(service string) (bool, error)
	FileExists(path string) bool
	IsRegularFile(path string) (bool, error)
	IsSymlink(path string) (bool, error)
}

// DefaultSystemOps implements SystemOps using real system calls
type DefaultSystemOps struct{}

// NewDefaultSystemOps creates a new DefaultSystemOps instance
func NewDefaultSystemOps() *DefaultSystemOps {
	return &DefaultSystemOps{}
}

// CommandExists checks if a command is available in PATH
func (d *DefaultSystemOps) CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// ServiceRunning checks if a systemd service is running
func (d *DefaultSystemOps) ServiceRunning(service string) (bool, error) {
	cmd := exec.Command("systemctl", "is-active", service)
	err := cmd.Run()
	if err != nil {
		// is-active returns non-zero if service is not active
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}
		// Other errors (command not found, permission denied, etc.)
		return false, fmt.Errorf("failed to check service status: %w", err)
	}
	return true, nil
}

// FileExists checks if a file or directory exists
func (d *DefaultSystemOps) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsRegularFile checks if the path is a regular file
func (d *DefaultSystemOps) IsRegularFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.Mode().IsRegular(), nil
}

// IsSymlink checks if the path is a symbolic link
func (d *DefaultSystemOps) IsSymlink(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return info.Mode()&os.ModeSymlink != 0, nil
}

// Detector handles backend detection
type Detector struct {
	sysOps SystemOps
}

// NewDetector creates a new Detector with the given SystemOps
func NewDetector(sysOps SystemOps) *Detector {
	return &Detector{sysOps: sysOps}
}

// Detect identifies and returns the active DNS backend
// Priority: NetworkManager > systemd-resolved > resolv.conf
func (d *Detector) Detect() (models.Backend, error) {
	backend, _, err := d.DetectWithReason()
	return backend, err
}

// DetectWithReason identifies the active DNS backend and returns a reason
func (d *Detector) DetectWithReason() (models.Backend, string, error) {
	// 1. Check for NetworkManager
	nmRunning, _ := d.sysOps.ServiceRunning("NetworkManager")
	if nmRunning {
		if d.sysOps.CommandExists("nmcli") {
			return models.BackendNetworkManager,
				"nmcli command available and NetworkManager service is running",
				nil
		}
		// NetworkManager is running but nmcli is missing
		return "", "", errors.New("NetworkManager is running but 'nmcli' command is missing.\n\n" +
			"To continue, please install the NetworkManager CLI tool:\n" +
			"  - Debian/Ubuntu: sudo apt install network-manager\n" +
			"  - Fedora/RHEL: sudo dnf install NetworkManager\n" +
			"  - Arch Linux: sudo pacman -S networkmanager")
	}

	// 2. Check for systemd-resolved
	hasResolvectl := d.sysOps.CommandExists("resolvectl")
	hasSystemdResolve := d.sysOps.CommandExists("systemd-resolve")

	if hasResolvectl || hasSystemdResolve {
		running, err := d.sysOps.ServiceRunning("systemd-resolved")
		if err != nil {
			return "", "", fmt.Errorf("failed to check systemd-resolved status: %w", err)
		}
		if running {
			cmdName := "resolvectl"
			if !hasResolvectl {
				cmdName = "systemd-resolve"
			}
			return models.BackendSystemdResolved,
				fmt.Sprintf("%s command available and systemd-resolved service is running", cmdName),
				nil
		}
	}

	// 3. Check for unmanaged resolv.conf
	if d.sysOps.FileExists("/etc/resolv.conf") {
		// Check if it's a symlink (managed by systemd-resolved or others)
		isSymlink, err := d.sysOps.IsSymlink("/etc/resolv.conf")
		if err != nil {
			return "", "", fmt.Errorf("failed to check if resolv.conf is a symlink: %w", err)
		}
		if isSymlink {
			return "", "resolv.conf is a symlink (managed by a service)",
				errors.New("resolv.conf is managed by a service")
		}

		// Check if it's a regular file
		isRegular, err := d.sysOps.IsRegularFile("/etc/resolv.conf")
		if err != nil {
			return "", "", fmt.Errorf("failed to check if resolv.conf is a regular file: %w", err)
		}
		if isRegular {
			return models.BackendResolvConf,
				"/etc/resolv.conf exists and is a regular file (not managed by a service)",
				nil
		}
	}

	// No supported backend found
	return "", "no supported DNS backend found",
		errors.New("no supported DNS backend found")
}
