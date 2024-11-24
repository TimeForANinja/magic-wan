package cfg

import (
	"gopkg.in/yaml.v2"
	"os"
)

// PrivateConfig represents the configuration structure
type PrivateConfig struct {
	ConfigVersion    string           `yaml:"config-version"`
	NodeID           uint8            `yaml:"node-id"`
	PrivateWireGuard PrivateWireGuard `yaml:"wireguard"`
}

// PrivateWireGuard represents the wireguard configuration
type PrivateWireGuard struct {
	PrivateKey string `yaml:"privkey"`
}

// normalizePrivateConfig takes a Config object and returns a normalized instance of it.
func normalizePrivateConfig(config PrivateConfig) PrivateConfig {
	return config
}

// LoadPrivateConfig loads configuration from a YAML file
func LoadPrivateConfig(filename string) (*PrivateConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config PrivateConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// normalize and validate config
	config = normalizePrivateConfig(config)
	if err := validatePrivateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// validatePrivateConfig validates the given configuration
func validatePrivateConfig(config *PrivateConfig) error {
	// TODO: implement tests

	return nil
}
