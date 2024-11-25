package osUtil

import (
	"encoding/json"
	"fmt"
	"magic-wan/pkg/various"
	"os/exec"
)

// TODO: check which fields are actually required
type AddrInfo struct {
	Family            string `json:"family"`
	Local             string `json:"local"`
	Prefixlen         int    `json:"prefixlen"`
	Scope             string `json:"scope"`
	Label             string `json:"label,omitempty"`
	ValidLifeTime     int64  `json:"valid_life_time"`
	PreferredLifeTime int64  `json:"preferred_life_time"`
	Noprefixroute     bool   `json:"noprefixroute,omitempty"`
	Broadcast         string `json:"broadcast,omitempty"`
	Dynamic           bool   `json:"dynamic,omitempty"`
}
type NetworkInterface struct {
	Ifindex     int         `json:"ifindex"`
	Ifname      string      `json:"ifname"`
	Flags       []string    `json:"flags"`
	Mtu         int         `json:"mtu"`
	Qdisc       string      `json:"qdisc"`
	Operstate   string      `json:"operstate"`
	Group       string      `json:"group"`
	Txqlen      int         `json:"txqlen"`
	LinkType    string      `json:"link_type"`
	Address     string      `json:"address"`
	Broadcast   string      `json:"broadcast"`
	AddrInfo    []*AddrInfo `json:"addr_info"`
	LinkIndex   int         `json:"link_index,omitempty"`
	LinkNetnsid int         `json:"link_netnsid,omitempty"`
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
