package wg

import (
	"fmt"
	"golang.zx2c4.com/wireguard/wgctrl"
)

func MustCreateController() *wgctrl.Client {
	// Create a new WireGuard controller
	client, err := wgctrl.New()
	if err != nil {
		panic(fmt.Sprintf("Failed to create WireGuard client: %v", err))
	}
	return client
}
