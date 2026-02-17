package backend_test

import (
	"errors"
	"os"
	"testing"

	"github.com/junevm/cdns/internal/dns/backend"
	"github.com/junevm/cdns/internal/dns/models"
)

// mockSystemOps implements backend.SystemOps for testing
type mockSystemOps struct {
	commandExists  func(string) bool
	serviceRunning func(string) (bool, error)
	fileExists     func(string) bool
	isRegularFile  func(string) (bool, error)
	isSymlink      func(string) (bool, error)
}

func (m *mockSystemOps) CommandExists(cmd string) bool {
	if m.commandExists != nil {
		return m.commandExists(cmd)
	}
	return false
}

func (m *mockSystemOps) ServiceRunning(service string) (bool, error) {
	if m.serviceRunning != nil {
		return m.serviceRunning(service)
	}
	return false, nil
}

func (m *mockSystemOps) FileExists(path string) bool {
	if m.fileExists != nil {
		return m.fileExists(path)
	}
	return false
}

func (m *mockSystemOps) IsRegularFile(path string) (bool, error) {
	if m.isRegularFile != nil {
		return m.isRegularFile(path)
	}
	return false, nil
}

func (m *mockSystemOps) IsSymlink(path string) (bool, error) {
	if m.isSymlink != nil {
		return m.isSymlink(path)
	}
	return false, nil
}

func TestDetector_Detect(t *testing.T) {
	tests := []struct {
		name    string
		sysOps  backend.SystemOps
		want    models.Backend
		wantErr bool
	}{
		{
			name: "NetworkManager present and active",
			sysOps: &mockSystemOps{
				commandExists: func(cmd string) bool {
					return cmd == "nmcli"
				},
				serviceRunning: func(service string) (bool, error) {
					if service == "NetworkManager" {
						return true, nil
					}
					return false, nil
				},
			},
			want:    models.BackendNetworkManager,
			wantErr: false,
		},
		{
			name: "NetworkManager command exists but service not running",
			sysOps: &mockSystemOps{
				commandExists: func(cmd string) bool {
					return cmd == "nmcli"
				},
				serviceRunning: func(service string) (bool, error) {
					return false, nil
				},
				fileExists: func(path string) bool {
					return path == "/etc/resolv.conf"
				},
				isRegularFile: func(path string) (bool, error) {
					return true, nil
				},
				isSymlink: func(path string) (bool, error) {
					return false, nil
				},
			},
			want:    models.BackendResolvConf,
			wantErr: false,
		},
		{
			name: "systemd-resolved present and active (resolvectl)",
			sysOps: &mockSystemOps{
				commandExists: func(cmd string) bool {
					return cmd == "resolvectl"
				},
				serviceRunning: func(service string) (bool, error) {
					if service == "systemd-resolved" {
						return true, nil
					}
					return false, nil
				},
			},
			want:    models.BackendSystemdResolved,
			wantErr: false,
		},
		{
			name: "systemd-resolved present and active (systemd-resolve fallback)",
			sysOps: &mockSystemOps{
				commandExists: func(cmd string) bool {
					return cmd == "systemd-resolve"
				},
				serviceRunning: func(service string) (bool, error) {
					if service == "systemd-resolved" {
						return true, nil
					}
					return false, nil
				},
			},
			want:    models.BackendSystemdResolved,
			wantErr: false,
		},
		{
			name: "unmanaged resolv.conf",
			sysOps: &mockSystemOps{
				commandExists: func(cmd string) bool {
					return false
				},
				fileExists: func(path string) bool {
					return path == "/etc/resolv.conf"
				},
				isRegularFile: func(path string) (bool, error) {
					if path == "/etc/resolv.conf" {
						return true, nil
					}
					return false, nil
				},
				isSymlink: func(path string) (bool, error) {
					return false, nil
				},
			},
			want:    models.BackendResolvConf,
			wantErr: false,
		},
		{
			name: "resolv.conf is symlink (managed)",
			sysOps: &mockSystemOps{
				commandExists: func(cmd string) bool {
					return false
				},
				fileExists: func(path string) bool {
					return path == "/etc/resolv.conf"
				},
				isSymlink: func(path string) (bool, error) {
					if path == "/etc/resolv.conf" {
						return true, nil
					}
					return false, nil
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "no supported backend",
			sysOps: &mockSystemOps{
				commandExists: func(cmd string) bool {
					return false
				},
				fileExists: func(path string) bool {
					return false
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "service check error",
			sysOps: &mockSystemOps{
				commandExists: func(cmd string) bool {
					return cmd == "nmcli"
				},
				serviceRunning: func(service string) (bool, error) {
					return false, errors.New("systemctl error")
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "NetworkManager priority over systemd-resolved",
			sysOps: &mockSystemOps{
				commandExists: func(cmd string) bool {
					return true // both available
				},
				serviceRunning: func(service string) (bool, error) {
					// both running
					return true, nil
				},
			},
			want:    models.BackendNetworkManager,
			wantErr: false,
		},
		{
			name: "systemd-resolved priority over resolv.conf",
			sysOps: &mockSystemOps{
				commandExists: func(cmd string) bool {
					return cmd == "resolvectl"
				},
				serviceRunning: func(service string) (bool, error) {
					if service == "systemd-resolved" {
						return true, nil
					}
					return false, nil
				},
				fileExists: func(path string) bool {
					return path == "/etc/resolv.conf"
				},
				isRegularFile: func(path string) (bool, error) {
					return true, nil
				},
				isSymlink: func(path string) (bool, error) {
					return false, nil
				},
			},
			want:    models.BackendSystemdResolved,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := backend.NewDetector(tt.sysOps)
			got, err := detector.Detect()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Detect() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Detect() unexpected error: %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("Detect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetector_DetectWithReason(t *testing.T) {
	tests := []struct {
		name       string
		sysOps     backend.SystemOps
		want       models.Backend
		wantReason string
		wantErr    bool
	}{
		{
			name: "NetworkManager with reason",
			sysOps: &mockSystemOps{
				commandExists: func(cmd string) bool {
					return cmd == "nmcli"
				},
				serviceRunning: func(service string) (bool, error) {
					if service == "NetworkManager" {
						return true, nil
					}
					return false, nil
				},
			},
			want:       models.BackendNetworkManager,
			wantReason: "nmcli command available and NetworkManager service is running",
			wantErr:    false,
		},
		{
			name: "systemd-resolved with reason",
			sysOps: &mockSystemOps{
				commandExists: func(cmd string) bool {
					return cmd == "resolvectl"
				},
				serviceRunning: func(service string) (bool, error) {
					if service == "systemd-resolved" {
						return true, nil
					}
					return false, nil
				},
			},
			want:       models.BackendSystemdResolved,
			wantReason: "resolvectl command available and systemd-resolved service is running",
			wantErr:    false,
		},
		{
			name: "resolv.conf with reason",
			sysOps: &mockSystemOps{
				commandExists: func(cmd string) bool {
					return false
				},
				fileExists: func(path string) bool {
					return path == "/etc/resolv.conf"
				},
				isRegularFile: func(path string) (bool, error) {
					return true, nil
				},
				isSymlink: func(path string) (bool, error) {
					return false, nil
				},
			},
			want:       models.BackendResolvConf,
			wantReason: "/etc/resolv.conf exists and is a regular file (not managed by a service)",
			wantErr:    false,
		},
		{
			name: "unsupported with reason",
			sysOps: &mockSystemOps{
				commandExists: func(cmd string) bool {
					return false
				},
				fileExists: func(path string) bool {
					return false
				},
			},
			want:       "",
			wantReason: "no supported DNS backend found",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := backend.NewDetector(tt.sysOps)
			got, reason, err := detector.DetectWithReason()

			if tt.wantErr {
				if err == nil {
					t.Errorf("DetectWithReason() expected error, got nil")
				}
				if reason == "" {
					t.Errorf("DetectWithReason() expected reason, got empty string")
				}
				if reason != tt.wantReason {
					t.Errorf("DetectWithReason() reason = %q, want %q", reason, tt.wantReason)
				}
				return
			}

			if err != nil {
				t.Errorf("DetectWithReason() unexpected error: %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("DetectWithReason() backend = %v, want %v", got, tt.want)
			}

			if reason != tt.wantReason {
				t.Errorf("DetectWithReason() reason = %q, want %q", reason, tt.wantReason)
			}
		})
	}
}

func TestDefaultSystemOps_CommandExists(t *testing.T) {
	sysOps := backend.NewDefaultSystemOps()

	// Test with a command that should exist
	exists := sysOps.CommandExists("sh")
	if !exists {
		t.Error("CommandExists(sh) = false, expected true")
	}

	// Test with a command that should not exist
	exists = sysOps.CommandExists("nonexistent-command-xyz123")
	if exists {
		t.Error("CommandExists(nonexistent-command-xyz123) = true, expected false")
	}
}

func TestDefaultSystemOps_FileExists(t *testing.T) {
	sysOps := backend.NewDefaultSystemOps()

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Test file that exists
	if !sysOps.FileExists(tmpfile.Name()) {
		t.Errorf("FileExists(%s) = false, expected true", tmpfile.Name())
	}

	// Test file that doesn't exist
	if sysOps.FileExists("/nonexistent/file/path") {
		t.Error("FileExists(/nonexistent/file/path) = true, expected false")
	}
}

func TestDefaultSystemOps_IsRegularFile(t *testing.T) {
	sysOps := backend.NewDefaultSystemOps()

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Test regular file
	isRegular, err := sysOps.IsRegularFile(tmpfile.Name())
	if err != nil {
		t.Errorf("IsRegularFile(%s) unexpected error: %v", tmpfile.Name(), err)
	}
	if !isRegular {
		t.Errorf("IsRegularFile(%s) = false, expected true", tmpfile.Name())
	}

	// Test directory
	tmpdir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	isRegular, err = sysOps.IsRegularFile(tmpdir)
	if err != nil {
		t.Errorf("IsRegularFile(%s) unexpected error: %v", tmpdir, err)
	}
	if isRegular {
		t.Errorf("IsRegularFile(%s) = true, expected false", tmpdir)
	}
}

func TestDefaultSystemOps_IsSymlink(t *testing.T) {
	sysOps := backend.NewDefaultSystemOps()

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Test regular file (not a symlink)
	isSymlink, err := sysOps.IsSymlink(tmpfile.Name())
	if err != nil {
		t.Errorf("IsSymlink(%s) unexpected error: %v", tmpfile.Name(), err)
	}
	if isSymlink {
		t.Errorf("IsSymlink(%s) = true, expected false", tmpfile.Name())
	}

	// Create a symlink
	symlinkPath := tmpfile.Name() + ".link"
	err = os.Symlink(tmpfile.Name(), symlinkPath)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(symlinkPath)

	isSymlink, err = sysOps.IsSymlink(symlinkPath)
	if err != nil {
		t.Errorf("IsSymlink(%s) unexpected error: %v", symlinkPath, err)
	}
	if !isSymlink {
		t.Errorf("IsSymlink(%s) = false, expected true", symlinkPath)
	}
}
