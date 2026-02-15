package models

// Backend represents the DNS management system type
type Backend string

const (
	// BackendNetworkManager represents NetworkManager backend
	BackendNetworkManager Backend = "NetworkManager"
	// BackendSystemdResolved represents systemd-resolved backend
	BackendSystemdResolved Backend = "systemd-resolved"
	// BackendResolvConf represents resolv.conf backend
	BackendResolvConf Backend = "resolv.conf"
	// BackendNetplan represents netplan backend
	BackendNetplan Backend = "netplan"
)

// DNSServer holds DNS server addresses
type DNSServer struct {
	IPv4        []string
	IPv6        []string
	Description string
}

// NetworkInterface represents a network interface configuration
type NetworkInterface struct {
	Name    string
	Backend Backend
}

// DNSConfig represents a complete DNS configuration for an interface
type DNSConfig struct {
	Interface NetworkInterface
	DNS       DNSServer
}
