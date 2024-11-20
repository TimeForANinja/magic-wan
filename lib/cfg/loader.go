package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"net"
	"os"
)

// Config represents the configuration structure
type Config struct {
	ConfigVersion string    `yaml:"config-version"`
	WireGuard     WireGuard `yaml:"wireguard"`
	Router        Router    `yaml:"router"`
	NodeID        uint8     `yaml:"node-id"`
}

// Router represents the router configuration
type Router struct {
	Subnet string `yaml:"subnet"`
}

// WireGuard represents the wireguard configuration
type WireGuard struct {
	PrivateKey string `yaml:"privkey"`
	Peers      []Peer `yaml:"peers"`
}

type Peer struct {
	PublicKey string `yaml:"pubkey"`
	UID       uint8  `yaml:"uid"`
	Hostname  string `yaml:"host"`
}

// normalizeConfig takes a Config object and returns a normalized instance of it.
func normalizeConfig(config Config) Config {
	return config
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// normalize and validate config
	config = normalizeConfig(config)
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// validateConfig validates the given configuration
func validateConfig(config *Config) error {
	// Ensure Router subnet is a valid IPv4 network in CIDR notation
	if ip, _, err := net.ParseCIDR(config.Router.Subnet); err != nil || ip.To4() == nil {
		return fmt.Errorf("invalid CIDR / IPv4 notation for Router subnet: %s", config.Router.Subnet)
	}

	return nil
}
