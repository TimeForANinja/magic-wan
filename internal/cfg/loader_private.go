package cfg

import (
	"fmt"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gopkg.in/yaml.v2"
	"os"
)

// PrivateConfig represents the configuration structure
type PrivateConfig struct {
	ConfigVersion    uint16           `yaml:"config-version"`
	NodeID           uint8            `yaml:"node-id"`
	PrivateWireGuard PrivateWireGuard `yaml:"wireguard"`
}

// PrivateWireGuard represents the wireguard configuration
type PrivateWireGuard struct {
	PrivateKey    wgtypes.Key `yaml:"-"`
	RawPrivateKey string      `yaml:"privkey"`
}

// normalizePrivateConfig takes a Config object and returns a normalized instance of it.
func normalizePrivateConfig(config *PrivateConfig) (*PrivateConfig, error) {
	var err error

	// Parse Wireguard Keys
	config.PrivateWireGuard.PrivateKey, err = wgtypes.ParseKey(config.PrivateWireGuard.RawPrivateKey)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// LoadPrivateConfig loads configuration from a YAML file
func LoadPrivateConfig(filename string) (*PrivateConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &PrivateConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	// normalize and validate config
	config, err = normalizePrivateConfig(config)
	if err != nil {
		return nil, err
	}
	err = validatePrivateConfig(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// validatePrivateConfig validates the given configuration
func validatePrivateConfig(config *PrivateConfig) error {
	if config.NodeID == 0 {
		return fmt.Errorf("invalid node id")
	}

	return nil
}
