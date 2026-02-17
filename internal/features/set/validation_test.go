package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDNSAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		// Valid IPv4 addresses
		{name: "valid IPv4 - Google", address: "8.8.8.8", wantErr: false},
		{name: "valid IPv4 - Cloudflare", address: "1.1.1.1", wantErr: false},
		{name: "valid IPv4 - Quad9", address: "9.9.9.9", wantErr: false},
		{name: "valid IPv4 - private", address: "192.168.1.1", wantErr: false},
		{name: "valid IPv4 - localhost", address: "127.0.0.1", wantErr: false},

		// Valid IPv6 addresses
		{name: "valid IPv6 - Cloudflare full", address: "2606:4700:4700::1111", wantErr: false},
		{name: "valid IPv6 - Google full", address: "2001:4860:4860::8888", wantErr: false},
		{name: "valid IPv6 - Quad9", address: "2620:fe::fe", wantErr: false},
		{name: "valid IPv6 - localhost", address: "::1", wantErr: false},
		{name: "valid IPv6 - full format", address: "2001:0db8:85a3:0000:0000:8a2e:0370:7334", wantErr: false},

		// Invalid addresses
		{name: "invalid - empty", address: "", wantErr: true},
		{name: "invalid - not an IP", address: "google.com", wantErr: true},
		{name: "invalid - malformed IPv4", address: "256.256.256.256", wantErr: true},
		{name: "invalid - incomplete IPv4", address: "192.168.1", wantErr: true},
		{name: "invalid - too many octets", address: "192.168.1.1.1", wantErr: true},
		{name: "invalid - malformed IPv6", address: "::gggg", wantErr: true},
		{name: "invalid - text", address: "not-an-ip-address", wantErr: true},
		{name: "invalid - partial", address: "8.8", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDNSAddress(tt.address)
			if tt.wantErr {
				assert.Error(t, err, "expected error for address: %s", tt.address)
			} else {
				assert.NoError(t, err, "expected no error for address: %s", tt.address)
			}
		})
	}
}

func TestValidateDNSAddresses(t *testing.T) {
	t.Run("validates multiple valid addresses", func(t *testing.T) {
		addresses := []string{"8.8.8.8", "8.8.4.4", "2606:4700:4700::1111"}
		err := ValidateDNSAddresses(addresses)
		assert.NoError(t, err)
	})

	t.Run("fails on empty list", func(t *testing.T) {
		addresses := []string{}
		err := ValidateDNSAddresses(addresses)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least one DNS address required")
	})

	t.Run("fails on single invalid address", func(t *testing.T) {
		addresses := []string{"8.8.8.8", "invalid", "8.8.4.4"}
		err := ValidateDNSAddresses(addresses)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid")
	})

	t.Run("fails on all invalid addresses", func(t *testing.T) {
		addresses := []string{"invalid1", "invalid2"}
		err := ValidateDNSAddresses(addresses)
		assert.Error(t, err)
	})
}

func TestValidatePresetName(t *testing.T) {
	tests := []struct {
		name       string
		presetName string
		wantErr    bool
	}{
		{name: "valid - cloudflare", presetName: "cloudflare", wantErr: false},
		{name: "valid - google", presetName: "google", wantErr: false},
		{name: "valid - quad9", presetName: "quad9", wantErr: false},
		{name: "valid - opendns", presetName: "opendns", wantErr: false},
		{name: "valid - adguard", presetName: "adguard", wantErr: false},
		{name: "valid - yandex-family", presetName: "yandex-family", wantErr: false},

		{name: "invalid - empty", presetName: "", wantErr: true},
		{name: "invalid - unknown", presetName: "unknown", wantErr: true},
		{name: "valid - case insensitive", presetName: "Cloudflare", wantErr: false},
		{name: "invalid - typo", presetName: "gogle", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePresetName(tt.presetName)
			if tt.wantErr {
				assert.Error(t, err, "expected error for preset: %s", tt.presetName)
			} else {
				assert.NoError(t, err, "expected no error for preset: %s", tt.presetName)
			}
		})
	}
}

func TestValidateInterfaceName(t *testing.T) {
	tests := []struct {
		name          string
		interfaceName string
		wantErr       bool
	}{
		// Valid interface names
		{name: "valid - eth0", interfaceName: "eth0", wantErr: false},
		{name: "valid - wlan0", interfaceName: "wlan0", wantErr: false},
		{name: "valid - enp3s0", interfaceName: "enp3s0", wantErr: false},
		{name: "valid - wlp2s0", interfaceName: "wlp2s0", wantErr: false},
		{name: "valid - lo", interfaceName: "lo", wantErr: false},
		{name: "valid - docker0", interfaceName: "docker0", wantErr: false},
		{name: "valid - br-123abc", interfaceName: "br-123abc", wantErr: false},

		// Invalid interface names
		{name: "invalid - empty", interfaceName: "", wantErr: true},
		{name: "invalid - spaces", interfaceName: "eth 0", wantErr: true},
		{name: "invalid - special chars", interfaceName: "eth@0", wantErr: true},
		{name: "invalid - too long", interfaceName: "this_is_a_very_long_interface_name_that_exceeds_limits", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateInterfaceName(tt.interfaceName)
			if tt.wantErr {
				assert.Error(t, err, "expected error for interface: %s", tt.interfaceName)
			} else {
				assert.NoError(t, err, "expected no error for interface: %s", tt.interfaceName)
			}
		})
	}
}

func TestSeparateIPv4AndIPv6(t *testing.T) {
	t.Run("separates mixed addresses", func(t *testing.T) {
		addresses := []string{
			"8.8.8.8",
			"2606:4700:4700::1111",
			"8.8.4.4",
			"2001:4860:4860::8888",
		}

		ipv4, ipv6 := SeparateIPv4AndIPv6(addresses)

		assert.Len(t, ipv4, 2)
		assert.Contains(t, ipv4, "8.8.8.8")
		assert.Contains(t, ipv4, "8.8.4.4")

		assert.Len(t, ipv6, 2)
		assert.Contains(t, ipv6, "2606:4700:4700::1111")
		assert.Contains(t, ipv6, "2001:4860:4860::8888")
	})

	t.Run("handles only IPv4", func(t *testing.T) {
		addresses := []string{"8.8.8.8", "8.8.4.4"}
		ipv4, ipv6 := SeparateIPv4AndIPv6(addresses)

		assert.Len(t, ipv4, 2)
		assert.Len(t, ipv6, 0)
	})

	t.Run("handles only IPv6", func(t *testing.T) {
		addresses := []string{"2606:4700:4700::1111", "::1"}
		ipv4, ipv6 := SeparateIPv4AndIPv6(addresses)

		assert.Len(t, ipv4, 0)
		assert.Len(t, ipv6, 2)
	})

	t.Run("handles empty list", func(t *testing.T) {
		addresses := []string{}
		ipv4, ipv6 := SeparateIPv4AndIPv6(addresses)

		assert.Len(t, ipv4, 0)
		assert.Len(t, ipv6, 0)
	})
}
