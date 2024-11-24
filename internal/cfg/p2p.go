package cfg

import (
	"fmt"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"magic-wan/pkg/transferNetwork"
	"net"
	"time"
)

// P2P represents the WireGuard link configuration between peers.
type P2P struct {
	General *SharedConfig

	Self        *Peer
	SelfPrivate *PrivateConfig

	Peer *Peer
}

func (settings *P2P) GetKeepalive() *time.Duration {
	keepAlive := 0 * time.Second
	if settings.Self.Keepalive {
		keepAlive = 10 * time.Second
	}
	return &keepAlive
}

func (settings *P2P) GetEndpoint(peerPort int) (*net.UDPAddr, error) {
	// Resolve hostname to *net.UDPAddr
	var endpoint *net.UDPAddr
	var err error

	if settings.Peer.Hostname != "" {
		endpoint, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", settings.Peer.Hostname, peerPort))
		if err != nil {
			return endpoint, err
		}
	}
	return endpoint, nil
}

func (settings *P2P) ToConfig(myPort, peerPort int) (wgtypes.Config, error) {
	selfPrivateKey, err := wgtypes.ParseKey(settings.SelfPrivate.PrivateWireGuard.PrivateKey)
	if err != nil {
		return wgtypes.Config{}, err
	}

	peerPubKey, err := wgtypes.ParseKey(settings.Peer.PublicKey)
	if err != nil {
		return wgtypes.Config{}, err
	}

	_, _, sharedNW, err := transferNetwork.GetPeerToPeerNet(settings.Self.UID, settings.Peer.UID, settings.General.Router.Subnet)
	if err != nil {
		return wgtypes.Config{}, err
	}

	endpoint, err := settings.GetEndpoint(peerPort)
	if err != nil {
		return wgtypes.Config{}, err
	}

	return wgtypes.Config{
		PrivateKey:   &selfPrivateKey, // Generate a private key for this interface
		ListenPort:   &myPort,
		ReplacePeers: true, // Replace existing peers with the provided ones
		Peers: []wgtypes.PeerConfig{
			{
				PublicKey:                   peerPubKey,
				Endpoint:                    endpoint,
				PersistentKeepaliveInterval: settings.GetKeepalive(),
				AllowedIPs: []net.IPNet{
					*sharedNW,
					{
						IP:   net.ParseIP("0.0.0.0"),
						Mask: net.CIDRMask(0, 32),
					},
					{
						IP:   net.ParseIP("::"),
						Mask: net.CIDRMask(0, 128),
					},
				},
			},
		},
	}, err
}
