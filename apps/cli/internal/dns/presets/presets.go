package presets

import "github.com/junevm/cdns/apps/cli/internal/dns/models"

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

// AdGuardFamily returns the AdGuard Family preset
func AdGuardFamily() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"94.140.14.15", "94.140.15.16"},
		Description: "Ad blocking with adult content filtering",
	}
}

// OpenDNSFamily returns the OpenDNS FamilyShield preset
func OpenDNSFamily() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"208.67.222.123", "208.67.220.123"},
		Description: "Child-safe filtering by default",
	}
}

// CleanBrowsingFamily returns the CleanBrowsing Family preset
func CleanBrowsingFamily() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"185.228.168.168", "185.228.169.168"},
		Description: "Blocks adult content & malicious sites",
	}
}

// CleanBrowsingSecurity returns the CleanBrowsing Security preset
func CleanBrowsingSecurity() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"185.228.168.9", "185.228.169.9"},
		Description: "Malware & phishing protection",
	}
}

// YandexBasic returns the Yandex Basic preset
func YandexBasic() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"77.88.8.8", "77.88.8.1"},
		Description: "Reliable DNS with basic filtering",
	}
}

// YandexSafe returns the Yandex Safe preset
func YandexSafe() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"77.88.8.88", "77.88.8.2"},
		Description: "Protection from malware & phishing",
	}
}

// YandexFamily returns the Yandex Family preset
func YandexFamily() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"77.88.8.7", "77.88.8.3"},
		Description: "Safe search & adult content blocking",
	}
}

// Comodo returns the Comodo Secure DNS preset
func Comodo() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"8.26.56.26", "8.20.247.20"},
		Description: "Security-focused with threat detection",
	}
}

// Verisign returns the Verisign Public DNS preset
func Verisign() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"64.6.64.6", "64.6.65.6"},
		Description: "Stable, secure, and private",
	}
}

// DNSWatch returns the DNS.WATCH preset
func DNSWatch() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"84.200.69.80", "84.200.70.40"},
		Description: "Fast, non-profit, and no logging",
	}
}

// Level3 returns the Level3 (CenturyLink) preset
func Level3() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"4.2.2.1", "4.2.2.2"},
		Description: "Legacy ISP DNS infrastructure",
	}
}

// Tencent returns the Tencent DNSPod preset
func Tencent() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"119.29.29.29", "182.254.116.116"},
		Description: "Optimized for mainland China",
	}
}

// Alibaba returns the Alibaba DNS preset
func Alibaba() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"223.5.5.5", "223.6.6.6"},
		Description: "Fast & stable China-based resolver",
	}
}

// Neustar returns the Neustar UltraDNS preset
func Neustar() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"156.154.70.1", "156.154.71.1"},
		Description: "Reliable with security features",
	}
}

// OpenNIC returns the OpenNIC preset
func OpenNIC() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"104.131.0.16", "192.71.245.208"},
		Description: "Community-driven & volunteer-run",
	}
}

// HurricaneElectric returns the Hurricane Electric DNS preset
func HurricaneElectric() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"74.82.42.42"},
		Description: "Reliable public DNS provided by HE",
	}
}

// SafeServe returns the SafeServe DNS preset
func SafeServe() models.DNSServer {
	return models.DNSServer{
		IPv4:        []string{"104.155.237.225", "104.155.237.226"},
		Description: "Global privacy with threat filtering",
	}
}

// All returns a map of all available DNS presets
func All() map[string]models.DNSServer {
	return map[string]models.DNSServer{
		"cloudflare":             Cloudflare(),
		"google":                 Google(),
		"quad9":                  Quad9(),
		"opendns":                OpenDNS(),
		"opendns-family":         OpenDNSFamily(),
		"adguard":                AdGuard(),
		"adguard-family":         AdGuardFamily(),
		"cleanbrowsing-family":   CleanBrowsingFamily(),
		"cleanbrowsing-security": CleanBrowsingSecurity(),
		"yandex-basic":           YandexBasic(),
		"yandex-safe":            YandexSafe(),
		"yandex-family":          YandexFamily(),
		"comodo":                 Comodo(),
		"verisign":               Verisign(),
		"dnswatch":               DNSWatch(),
		"level3":                 Level3(),
		"tencent":                Tencent(),
		"alibaba":                Alibaba(),
		"neustar":                Neustar(),
		"opennic":                OpenNIC(),
		"he":                     HurricaneElectric(),
		"safeserve":              SafeServe(),
	}
}

// Get retrieves a preset by name, returns the preset and a boolean indicating if it was found
func Get(name string) (models.DNSServer, bool) {
	preset, ok := All()[name]
	return preset, ok
}
