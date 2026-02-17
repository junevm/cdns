package presets

import "cli/internal/dns/models"

// Cloudflare returns the Cloudflare DNS preset
func Cloudflare() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"1.1.1.1", "1.0.0.1"},
		IPv6:        []string{"2606:4700:4700::1111", "2606:4700:4700::1001"},
		Description: "Fast & privacy-focused",
	}
}

// Google returns the Google DNS preset
func Google() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"8.8.8.8", "8.8.4.4"},
		IPv6:        []string{"2001:4860:4860::8888", "2001:4860:4860::8844"},
		Description: "Reliable & widely used",
	}
}

// Quad9 returns the Quad9 DNS preset
func Quad9() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"9.9.9.9", "149.112.112.112"},
		IPv6:        []string{"2620:fe::fe", "2620:fe::9"},
		Description: "Security-focused with threat intelligence",
	}
}

// OpenDNS returns the OpenDNS preset
func OpenDNS() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"208.67.222.222", "208.67.220.220"},
		IPv6:        []string{"2620:119:35::35", "2620:119:53::53"},
		Description: "Fast with content filtering options",
	}
}

// AdGuard returns the AdGuard DNS preset
func AdGuard() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"94.140.14.14", "94.140.15.15"},
		IPv6:        []string{"2a10:50c0::ad1:ff", "2a10:50c0::ad2:ff"},
		Description: "Focuses on ad and tracker blocking",
	}
}

// All returns a map of all available DNS presets
func All() map[string]models.DNSServer {
	return map[string]models.DNSServer{
		"cloudflare": Cloudflare(),
		"google":     Google(),
		"quad9":      Quad9(),
		"opendns":    OpenDNS(),
		"adguard":    AdGuard(),
	}
}

// Get retrieves a preset by name, returns the preset and a boolean indicating if it was found
func Get(name string) (models.DNSServer, bool) {
	preset, ok := All()[name]
	return preset, ok
}
