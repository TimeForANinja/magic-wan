package cfg

import (
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"magic-wan/pkg/transferNetwork"
	"magic-wan/pkg/various"
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

func (settings *P2P) GetName() string {
	return transferNetwork.BuildWireguardInterfaceName(settings.Self.UID, settings.Peer.UID)
}

func (settings *P2P) GetPorts() (int, int) {
	return transferNetwork.CalculatePorts(settings.General.SharedWireGuard.StartPort, settings.Peer.UID, settings.Self.UID)
}

func (settings *P2P) GetPeerAddr() (*net.UDPAddr, error) {
	_, peerPort := settings.GetPorts()
	return various.ResolveHostname(settings.Peer.Hostname, peerPort)
}

func (settings *P2P) GetEndpoint(peerPort int) (*net.UDPAddr, error) {
	return various.ResolveHostname(settings.Peer.Hostname, peerPort)
}

func (settings *P2P) ToConfig() (wgtypes.Config, error) {
	myPort, peerPort := settings.GetPorts()

	_, _, sharedNW, err := transferNetwork.GetPeerToPeerNet(settings.Self.UID, settings.Peer.UID, settings.General.Router.Subnet)
	if err != nil {
		return wgtypes.Config{}, err
	}

	endpoint, err := settings.GetEndpoint(peerPort)
	if err != nil {
		return wgtypes.Config{}, err
	}

	return wgtypes.Config{
		PrivateKey:   &settings.SelfPrivate.PrivateWireGuard.PrivateKey, // Generate a private key for this interface
		ListenPort:   &myPort,
		ReplacePeers: true, // Replace existing peers with the provided ones
		Peers: []wgtypes.PeerConfig{
			{
				PublicKey:                   settings.Peer.PublicKey,
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
