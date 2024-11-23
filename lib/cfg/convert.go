package cfg

import (
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"magic-wan/lib/transfer"
	"magic-wan/lib/wg"
	"net"
)

type WGLink struct {
	General     *SharedConfig
	Name        string
	Self        *Peer
	SelfPrivate *PrivateConfig
	Peer        *Peer
}

func (w *WGLink) ToWGConfig() (error, wgtypes.Config) {
	selfPrivateKey, err := wgtypes.ParseKey(w.SelfPrivate.PrivateWireGuard.PrivateKey)
	if err != nil {
		return err, wgtypes.Config{}
	}

	peerPubKey, err := wgtypes.ParseKey(w.Peer.PublicKey)
	if err != nil {
		return err, wgtypes.Config{}
	}

	_, _, sharedNW, err := transfer.GetPeerToPeerNet(w.Self.UID, w.Peer.UID, w.General.Router.Subnet)
	if err != nil {
		return err, wgtypes.Config{}
	}

	return nil, wgtypes.Config{
		PrivateKey:   &selfPrivateKey,                                                   // Generate a private key for this interface
		ListenPort:   wg.CalculatePort(w.General.SharedWireGuard.StartPort, w.Self.UID), // Specify a listen port, nil to randomize
		ReplacePeers: true,                                                              // Replace existing peers with the provided ones
		Peers: []wgtypes.PeerConfig{
			{
				PublicKey: peerPubKey,
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
	}
}
