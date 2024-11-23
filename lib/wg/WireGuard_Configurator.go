package wg

import (
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"log"
	"net"
)

func main() {
	client := mustCreateController()
	defer client.Close()

	masterKey := mustGeneratePrivateKey()

	// Define the configuration for the WireGuard interface
	config := wgtypes.Config{
		PrivateKey:   &masterKey, // Generate a private key for this interface
		ListenPort:   new(int),   // Specify a listen port, nil to randomize
		ReplacePeers: true,       // Replace existing peers with the provided ones
		Peers: []wgtypes.PeerConfig{
			{
				PublicKey: mustParsePublicKey("PUBLIC_KEY_OF_PEER_1"),
				AllowedIPs: []net.IPNet{
					{
						IP:   net.ParseIP("10.0.0.1"),
						Mask: net.CIDRMask(32, 32),
					},
				},
			},
			{
				PublicKey: mustParsePublicKey("PUBLIC_KEY_OF_PEER_2"),
				AllowedIPs: []net.IPNet{
					{
						IP:   net.ParseIP("10.0.0.2"),
						Mask: net.CIDRMask(32, 32),
					},
				},
			},
		},
	}

	// Apply the configuration to the interface
	ifcName := "wg0"
	mustConfigureDevice(client, ifcName, config)

	log.Printf("Successfully configured WireGuard interface %s", ifcName)
}