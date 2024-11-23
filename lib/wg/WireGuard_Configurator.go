package wg

import (
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"log"
)

func CreateNewDevice(client *wgctrl.Client, ifcName string, cfg wgtypes.Config) {
	// Apply the configuration to the interface
	mustConfigureDevice(client, ifcName, cfg)

	log.Printf("Successfully configured WireGuard interface %s", ifcName)
}

func RemoveDevice(client *wgctrl.Client, ifcName string) {
	zero := 0
	cfg := wgtypes.Config{
		// Clear device config.
		PrivateKey:   &wgtypes.Key{},
		ListenPort:   &zero,
		FirewallMark: &zero,

		// Clear all peers.
		ReplacePeers: true,
	}

	mustConfigureDevice(client, ifcName, cfg)
}
