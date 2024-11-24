package cfg

import (
	"fmt"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gopkg.in/yaml.v2"
	"net"
	"os"
)

// SharedConfig represents the configuration structure
type SharedConfig struct {
	ConfigVersion   uint16          `yaml:"config-version"`
	MTU             uint16          `yaml:"mtu"`
	SharedWireGuard SharedWireGuard `yaml:"wireguard"`
	Router          Router          `yaml:"router"`
}

// Router represents the router configuration
type Router struct {
	Subnet    *net.IPNet `yaml:"-"`
	RawSubnet string     `yaml:"subnet"`
}

// SharedWireGuard represents the wireguard configuration
type SharedWireGuard struct {
	StartPort uint16  `yaml:"startPort"`
	Peers     []*Peer `yaml:"peers"`
}

type Peer struct {
	PublicKey    wgtypes.Key `yaml:"-"`
	RawPublicKey string      `yaml:"pubkey"`
	Name         string      `yaml:"name"`
	Keepalive    bool        `yaml:"keepalive"`
	UID          uint8       `yaml:"uid"`
	Hostname     string      `yaml:"host"`
}

// normalizeSharedConfig takes a Config object and returns a normalized instance of it.
func normalizeSharedConfig(config *SharedConfig) (*SharedConfig, error) {
	var err error

	// Parse Wireguard Keys
	for idx, peer := range config.SharedWireGuard.Peers {
		config.SharedWireGuard.Peers[idx], err = normalizeSharedPeerConfig(peer)
		if err != nil {
			return nil, err
		}
	}

	// default mtu to 1350
	if config.MTU == 0 {
		config.MTU = 1350
	}

	// Parse & Validate Router subnet -should be a valid network in CIDR notation
	_, config.Router.Subnet, err = net.ParseCIDR(config.Router.RawSubnet)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR notation for Router subnet \"%s\": %w", config.Router.Subnet, err)
	}

	return config, nil
}

func normalizeSharedPeerConfig(peer *Peer) (*Peer, error) {
	var err error

	// Parse Wireguard Keys
	peer.PublicKey, err = wgtypes.ParseKey(peer.RawPublicKey)
	if err != nil {
		return nil, err
	}

	return peer, nil
}

// LoadSharedConfig loads configuration from a YAML file
func LoadSharedConfig(filename string) (*SharedConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &SharedConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	// normalize and validate config
	config, err = normalizeSharedConfig(config)
	if err != nil {
		return nil, err
	}
	err = validateSharedConfig(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// validateSharedConfig validates the given configuration
func validateSharedConfig(config *SharedConfig) error {
	var err error

	// Ensure Wireguard Keys are valid
	for _, peer := range config.SharedWireGuard.Peers {
		err = validateSharedPeerConfig(peer)
		if err != nil {
			return err
		}
	}

	// check MTU and Port are in range
	if config.MTU == 0 || config.MTU > 1500 {
		return fmt.Errorf("invalid MTU value")
	}
	if config.SharedWireGuard.StartPort <= 1024 {
		return fmt.Errorf("invalid wg start port value")
	}

	return nil
}

func validateSharedPeerConfig(peer *Peer) error {
	if peer.UID == 0 {
		return fmt.Errorf("invalid peer id")
	}

	if peer.Name == "" {
		return fmt.Errorf("invalid peer name")
	}

	return nil
}
