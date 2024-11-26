package osUtil

import (
	"encoding/json"
	"fmt"
	"magic-wan/pkg/various"
	"os/exec"
)

// IMPROVEMENT: check which fields are actually required
type AddrInfo struct {
	Family    string `json:"family"`
	Local     string `json:"local"`
	Prefixlen int    `json:"prefixlen"`
}
type NetworkInterface struct {
	Ifname    string      `json:"ifname"`
	Address   string      `json:"address"`
	Broadcast string      `json:"broadcast"`
	AddrInfo  []*AddrInfo `json:"addr_info"`
}

func GetInterfaces() ([]*NetworkInterface, error) {
	cmd := exec.Command("ip", "--json", "address", "show")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to parse ip addresses: %w, output: %s", err, string(output))
	}

	var interfaces []*NetworkInterface
	err = json.Unmarshal(output, &interfaces)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON output: %w", err)
	}
	return interfaces, nil
}

func InterfaceHasAddress(interfaceName string, ip string) bool {
	interfaces, err := GetInterfaces()
	if err != nil {
		return false
	}

	iface := various.ArrayFind(interfaces, func(iface *NetworkInterface) bool {
		return iface.Ifname == interfaceName
	})
	if iface == nil {
		return false
	}

	return various.ArrayIncludes(iface.AddrInfo, func(addrInfo *AddrInfo) bool {
		return addrInfo.Local == ip
	})
}

func SetInterfaceAddress(interfaceName string, ip string) error {
	cmd := exec.Command("ip", "addr", "add", ip, "dev", interfaceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set ip addresses: %w, output: %s", err, string(output))
	}
	return nil
}

func EnsureInterfaceHasAddress(interfaceName string, ip string) error {
	if !InterfaceHasAddress(interfaceName, ip) {
		return SetInterfaceAddress(interfaceName, ip)
	}
	return nil
}

func SetInterfaceUp(interfaceName string) error {
	cmd := exec.Command("ip", "link", "set", "dev", interfaceName, "up")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set ip addresses: %w, output: %s", err, string(output))
	}
	return nil
}
