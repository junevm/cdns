package models_test

import (
	"testing"

	"github.com/junevm/cdns/apps/cli/internal/dns/models"
)

func TestDNSServer(t *testing.T) {
	tests := []struct {
		name     string
		ipv4     []string
		ipv6     []string
		wantIPv4 []string
		wantIPv6 []string
	}{
		{
			name:     "valid IPv4 and IPv6 addresses",
			ipv4:     []string{"1.1.1.1", "1.0.0.1"},
			ipv6:     []string{"2606:4700:4700::1111", "2606:4700:4700::1001"},
			wantIPv4: []string{"1.1.1.1", "1.0.0.1"},
			wantIPv6: []string{"2606:4700:4700::1111", "2606:4700:4700::1001"},
		},
		{
			name:     "only IPv4 addresses",
			ipv4:     []string{"8.8.8.8", "8.8.4.4"},
			ipv6:     []string{},
			wantIPv4: []string{"8.8.8.8", "8.8.4.4"},
			wantIPv6: []string{},
		},
		{
			name:     "only IPv6 addresses",
			ipv4:     []string{},
			ipv6:     []string{"2001:4860:4860::8888"},
			wantIPv4: []string{},
			wantIPv6: []string{"2001:4860:4860::8888"},
		},
		{
			name:     "empty addresses",
			ipv4:     []string{},
			ipv6:     []string{},
			wantIPv4: []string{},
			wantIPv6: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dns := models.DNSServer{
				IPv4: tt.ipv4,
				IPv6: tt.ipv6,
			}

			if len(dns.IPv4) != len(tt.wantIPv4) {
				t.Errorf("IPv4 count = %d, want %d", len(dns.IPv4), len(tt.wantIPv4))
			}

			if len(dns.IPv6) != len(tt.wantIPv6) {
				t.Errorf("IPv6 count = %d, want %d", len(dns.IPv6), len(tt.wantIPv6))
			}

			for i, addr := range dns.IPv4 {
				if addr != tt.wantIPv4[i] {
					t.Errorf("IPv4[%d] = %s, want %s", i, addr, tt.wantIPv4[i])
				}
			}

			for i, addr := range dns.IPv6 {
				if addr != tt.wantIPv6[i] {
					t.Errorf("IPv6[%d] = %s, want %s", i, addr, tt.wantIPv6[i])
				}
			}
		})
	}
}

func TestBackend(t *testing.T) {
	tests := []struct {
		name    string
		backend models.Backend
		want    string
	}{
		{
			name:    "NetworkManager backend",
			backend: models.BackendNetworkManager,
			want:    "NetworkManager",
		},
		{
			name:    "systemd-resolved backend",
			backend: models.BackendSystemdResolved,
			want:    "systemd-resolved",
		},
		{
			name:    "resolv.conf backend",
			backend: models.BackendResolvConf,
			want:    "resolv.conf",
		},
		{
			name:    "netplan backend",
			backend: models.BackendNetplan,
			want:    "netplan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.backend) != tt.want {
				t.Errorf("Backend = %s, want %s", tt.backend, tt.want)
			}
		})
	}
}

func TestNetworkInterface(t *testing.T) {
	tests := []struct {
		name          string
		interfaceName string
		backend       models.Backend
		wantName      string
		wantBackend   models.Backend
	}{
		{
			name:          "ethernet interface with NetworkManager",
			interfaceName: "eth0",
			backend:       models.BackendNetworkManager,
			wantName:      "eth0",
			wantBackend:   models.BackendNetworkManager,
		},
		{
			name:          "wifi interface with systemd-resolved",
			interfaceName: "wlan0",
			backend:       models.BackendSystemdResolved,
			wantName:      "wlan0",
			wantBackend:   models.BackendSystemdResolved,
		},
		{
			name:          "loopback interface",
			interfaceName: "lo",
			backend:       models.BackendResolvConf,
			wantName:      "lo",
			wantBackend:   models.BackendResolvConf,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iface := models.NetworkInterface{
				Name:    tt.interfaceName,
				Backend: tt.backend,
			}

			if iface.Name != tt.wantName {
				t.Errorf("Name = %s, want %s", iface.Name, tt.wantName)
			}

			if iface.Backend != tt.wantBackend {
				t.Errorf("Backend = %s, want %s", iface.Backend, tt.wantBackend)
			}
		})
	}
}

func TestDNSConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    models.DNSConfig
		wantIface string
		wantDNS   models.DNSServer
	}{
		{
			name: "complete DNS configuration",
			config: models.DNSConfig{
				Interface: models.NetworkInterface{
					Name:    "eth0",
					Backend: models.BackendNetworkManager,
				},
				DNS: models.DNSServer{
					IPv4: []string{"1.1.1.1", "1.0.0.1"},
					IPv6: []string{"2606:4700:4700::1111"},
				},
			},
			wantIface: "eth0",
			wantDNS: models.DNSServer{
				IPv4: []string{"1.1.1.1", "1.0.0.1"},
				IPv6: []string{"2606:4700:4700::1111"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.Interface.Name != tt.wantIface {
				t.Errorf("Interface.Name = %s, want %s", tt.config.Interface.Name, tt.wantIface)
			}

			if len(tt.config.DNS.IPv4) != len(tt.wantDNS.IPv4) {
				t.Errorf("DNS.IPv4 count = %d, want %d", len(tt.config.DNS.IPv4), len(tt.wantDNS.IPv4))
			}

			if len(tt.config.DNS.IPv6) != len(tt.wantDNS.IPv6) {
				t.Errorf("DNS.IPv6 count = %d, want %d", len(tt.config.DNS.IPv6), len(tt.wantDNS.IPv6))
			}
		})
	}
}
