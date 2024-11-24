package wg

import (
	"fmt"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func MustCreateController() *wgctrl.Client {
	// Create a new WireGuard controller
	client, err := wgctrl.New()
	if err != nil {
		panic(fmt.Sprintf("Failed to create WireGuard client: %v", err))
	}
	return client
}

func MustConfigureDevice(client *wgctrl.Client, ifcName string, config wgtypes.Config) {
	err := client.ConfigureDevice(ifcName, config)
	if err != nil {
		panic(fmt.Sprintf("Failed to configure WireGuard interface %s: %v", ifcName, err))
	}
}
