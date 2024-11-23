package cfg

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"net"
	"os"
)

// SharedConfig represents the configuration structure
type SharedConfig struct {
	ConfigVersion   string          `yaml:"config-version"`
	MTU             uint16          `yaml:"mtu"`
	SharedWireGuard SharedWireGuard `yaml:"wireguard"`
	Router          Router          `yaml:"router"`
}

// Router represents the router configuration
type Router struct {
	Subnet string `yaml:"subnet"`
}

// SharedWireGuard represents the wireguard configuration
type SharedWireGuard struct {
	StartPort uint16 `yaml:"startPort"`
	Peers     []Peer `yaml:"peers"`
}

type Peer struct {
	PublicKey string `yaml:"pubkey"`
	Name      string `yaml:"name"`
	Keepalive bool   `yaml:"keepalive"`
	UID       uint8  `yaml:"uid"`
	Hostname  string `yaml:"host"`
}

// normalizeSharedConfig takes a Config object and returns a normalized instance of it.
func normalizeSharedConfig(config SharedConfig) SharedConfig {
	return config
}

// LoadSharedConfig loads configuration from a YAML file
func LoadSharedConfig(filename string) (*SharedConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config SharedConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// normalize and validate config
	config = normalizeSharedConfig(config)
	if err := validateSharedConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// validateSharedConfig validates the given configuration
func validateSharedConfig(config *SharedConfig) error {
	// Ensure Router subnet is a valid network in CIDR notation
	if _, _, err := net.ParseCIDR(config.Router.Subnet); err != nil {
		return fmt.Errorf("invalid CIDR notation for Router subnet: %s", config.Router.Subnet)
	}

	return nil
}
