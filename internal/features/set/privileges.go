package set

import (
	"os"
	"os/exec"
	"runtime"
)

// HasPrivileges checks if the current process has sufficient privileges
// to modify system DNS settings
func HasPrivileges() bool {
	// On Unix-like systems, check if running as root
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		return os.Geteuid() == 0
	}

	// On Windows, we'd need to check for admin privileges
	// For now, assume yes (will be implemented when Windows support is added)
	if runtime.GOOS == "windows" {
		// TODO: Implement Windows admin check
		return true
	}

	// Unknown OS, assume no privileges
	return false
}

// EnsurePrivileges attempts to escalate privileges if necessary
func EnsurePrivileges() error {
	if HasPrivileges() {
		return nil
	}

	// Only support escalation on Linux/Darwin for now
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		return ErrInsufficientPrivileges
	}

	// Check if sudo is available
	if _, err := exec.LookPath("sudo"); err != nil {
		return ErrInsufficientPrivileges
	}

	// Get the path to the current executable
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	// Re-run the command with sudo
	// We reconstruct the command carefully:
	// 1. Start with the executable path
	// 2. Add existing arguments (flags, etc.)
	// 3. Ensure the 'set' command is present if missing (e.g. running from menu)

	newArgs := []string{executable}

	// Create a copy of args to inspect
	args := make([]string, len(os.Args))
	copy(args, os.Args)

	setCommandPresent := false
	// Start checking from index 1 to skip executable name
	for _, arg := range args[1:] {
		if arg == "set" {
			setCommandPresent = true
			break
		}
	}

	// Add existing arguments (skipping arg[0] which is executable path)
	if len(args) > 1 {
		newArgs = append(newArgs, args[1:]...)
	}

	// If "set" is missing, append it to force the subcommand execution
	if !setCommandPresent {
		newArgs = append(newArgs, "set")
	}

	// Execute sudo with the preserved environment and constructed arguments
	// Uses --preserve-env to keep environment variables like config location
	cmd := exec.Command("sudo", append([]string{"--preserve-env"}, newArgs...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	// If the command succeeded, the child process handled everything.
	// We must exit the parent process.
	os.Exit(0)
	return nil
}
