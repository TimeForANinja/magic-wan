package wg

import (
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"log"
)

func mustCreateController() *wgctrl.Client {
	// Create a new WireGuard controller
	client, err := wgctrl.New()
	if err != nil {
		log.Fatalf("Failed to create WireGuard client: %v", err)
	}
	return client
}

func mustConfigureDevice(client *wgctrl.Client, ifcName string, config wgtypes.Config) {
	err := client.ConfigureDevice(ifcName, config)
	if err != nil {
		log.Fatalf("Failed to configure WireGuard interface %s: %v", ifcName, err)
	}
}
