package presets_test

import (
	"testing"

	"github.com/junevm/cdns/internal/dns/presets"
)

func TestCloudflare(t *testing.T) {
	dns := presets.Cloudflare()

	wantIPv4 := []string{"1.1.1.1", "1.0.0.1"}
	wantIPv6 := []string{"2606:4700:4700::1111", "2606:4700:4700::1001"}

	if len(dns.IPv4) != len(wantIPv4) {
		t.Errorf("Cloudflare IPv4 count = %d, want %d", len(dns.IPv4), len(wantIPv4))
	}

	if len(dns.IPv6) != len(wantIPv6) {
		t.Errorf("Cloudflare IPv6 count = %d, want %d", len(dns.IPv6), len(wantIPv6))
	}

	for i, addr := range dns.IPv4 {
		if addr != wantIPv4[i] {
			t.Errorf("Cloudflare IPv4[%d] = %s, want %s", i, addr, wantIPv4[i])
		}
	}

	for i, addr := range dns.IPv6 {
		if addr != wantIPv6[i] {
			t.Errorf("Cloudflare IPv6[%d] = %s, want %s", i, addr, wantIPv6[i])
		}
	}
}

func TestGoogle(t *testing.T) {
	dns := presets.Google()

	wantIPv4 := []string{"8.8.8.8", "8.8.4.4"}
	wantIPv6 := []string{"2001:4860:4860::8888", "2001:4860:4860::8844"}

	if len(dns.IPv4) != len(wantIPv4) {
		t.Errorf("Google IPv4 count = %d, want %d", len(dns.IPv4), len(wantIPv4))
	}

	if len(dns.IPv6) != len(wantIPv6) {
		t.Errorf("Google IPv6 count = %d, want %d", len(dns.IPv6), len(wantIPv6))
	}

	for i, addr := range dns.IPv4 {
		if addr != wantIPv4[i] {
			t.Errorf("Google IPv4[%d] = %s, want %s", i, addr, wantIPv4[i])
		}
	}

	for i, addr := range dns.IPv6 {
		if addr != wantIPv6[i] {
			t.Errorf("Google IPv6[%d] = %s, want %s", i, addr, wantIPv6[i])
		}
	}
}

func TestQuad9(t *testing.T) {
	dns := presets.Quad9()

	wantIPv4 := []string{"9.9.9.9", "149.112.112.112"}
	wantIPv6 := []string{"2620:fe::fe", "2620:fe::9"}

	if len(dns.IPv4) != len(wantIPv4) {
		t.Errorf("Quad9 IPv4 count = %d, want %d", len(dns.IPv4), len(wantIPv4))
	}

	if len(dns.IPv6) != len(wantIPv6) {
		t.Errorf("Quad9 IPv6 count = %d, want %d", len(dns.IPv6), len(wantIPv6))
	}

	for i, addr := range dns.IPv4 {
		if addr != wantIPv4[i] {
			t.Errorf("Quad9 IPv4[%d] = %s, want %s", i, addr, wantIPv4[i])
		}
	}

	for i, addr := range dns.IPv6 {
		if addr != wantIPv6[i] {
			t.Errorf("Quad9 IPv6[%d] = %s, want %s", i, addr, wantIPv6[i])
		}
	}
}

func TestOpenDNS(t *testing.T) {
	dns := presets.OpenDNS()

	wantIPv4 := []string{"208.67.222.222", "208.67.220.220"}
	wantIPv6 := []string{"2620:119:35::35", "2620:119:53::53"}

	if len(dns.IPv4) != len(wantIPv4) {
		t.Errorf("OpenDNS IPv4 count = %d, want %d", len(dns.IPv4), len(wantIPv4))
	}

	if len(dns.IPv6) != len(wantIPv6) {
		t.Errorf("OpenDNS IPv6 count = %d, want %d", len(dns.IPv6), len(wantIPv6))
	}

	for i, addr := range dns.IPv4 {
		if addr != wantIPv4[i] {
			t.Errorf("OpenDNS IPv4[%d] = %s, want %s", i, addr, wantIPv4[i])
		}
	}

	for i, addr := range dns.IPv6 {
		if addr != wantIPv6[i] {
			t.Errorf("OpenDNS IPv6[%d] = %s, want %s", i, addr, wantIPv6[i])
		}
	}
}

func TestAdGuard(t *testing.T) {
	dns := presets.AdGuard()

	wantIPv4 := []string{"94.140.14.14", "94.140.15.15"}
	wantIPv6 := []string{"2a10:50c0::ad1:ff", "2a10:50c0::ad2:ff"}

	if len(dns.IPv4) != len(wantIPv4) {
		t.Errorf("AdGuard IPv4 count = %d, want %d", len(dns.IPv4), len(wantIPv4))
	}

	if len(dns.IPv6) != len(wantIPv6) {
		t.Errorf("AdGuard IPv6 count = %d, want %d", len(dns.IPv6), len(wantIPv6))
	}

	for i, addr := range dns.IPv4 {
		if addr != wantIPv4[i] {
			t.Errorf("AdGuard IPv4[%d] = %s, want %s", i, addr, wantIPv4[i])
		}
	}

	for i, addr := range dns.IPv6 {
		if addr != wantIPv6[i] {
			t.Errorf("AdGuard IPv6[%d] = %s, want %s", i, addr, wantIPv6[i])
		}
	}
}

func TestAllPresets(t *testing.T) {
	all := presets.All()

	expectedCount := 22

	if len(all) != expectedCount {
		t.Errorf("All() returned %d presets, want %d", len(all), expectedCount)
	}

	expectedNames := []string{
		"cloudflare", "google", "quad9", "opendns", "opendns-family",
		"adguard", "adguard-family", "cleanbrowsing-family", "cleanbrowsing-security",
		"yandex-basic", "yandex-safe", "yandex-family", "comodo", "verisign",
		"dnswatch", "level3", "tencent", "alibaba", "neustar", "opennic",
		"he", "safeserve",
	}
	for _, name := range expectedNames {
		if _, ok := all[name]; !ok {
			t.Errorf("All() missing preset: %s", name)
		}
	}
}

func TestGetPreset(t *testing.T) {
	tests := []struct {
		name      string
		presetKey string
		wantIPv4  []string
		wantFound bool
	}{
		{
			name:      "cloudflare preset exists",
			presetKey: "cloudflare",
			wantIPv4:  []string{"1.1.1.1", "1.0.0.1"},
			wantFound: true,
		},
		{
			name:      "google preset exists",
			presetKey: "google",
			wantIPv4:  []string{"8.8.8.8", "8.8.4.4"},
			wantFound: true,
		},
		{
			name:      "nonexistent preset",
			presetKey: "nonexistent",
			wantIPv4:  nil,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dns, found := presets.Get(tt.presetKey)

			if found != tt.wantFound {
				t.Errorf("Get(%s) found = %v, want %v", tt.presetKey, found, tt.wantFound)
			}

			if found && tt.wantIPv4 != nil {
				if len(dns.IPv4) != len(tt.wantIPv4) {
					t.Errorf("Get(%s) IPv4 count = %d, want %d", tt.presetKey, len(dns.IPv4), len(tt.wantIPv4))
				}

				for i, addr := range dns.IPv4 {
					if addr != tt.wantIPv4[i] {
						t.Errorf("Get(%s) IPv4[%d] = %s, want %s", tt.presetKey, i, addr, tt.wantIPv4[i])
					}
				}
			}
		})
	}
}

func TestPresetsImmutability(t *testing.T) {
	// Get the same preset twice
	first := presets.Cloudflare()
	second := presets.Cloudflare()

	// Modify the first one
	first.IPv4[0] = "modified"

	// Check that the second one is not modified
	if second.IPv4[0] == "modified" {
		t.Error("Preset was modified, but should be immutable")
	}

	// Verify second still has correct value
	if second.IPv4[0] != "1.1.1.1" {
		t.Errorf("Second preset IPv4[0] = %s, want 1.1.1.1", second.IPv4[0])
	}
}
