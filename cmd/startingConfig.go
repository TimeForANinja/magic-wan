package main

import (
	"golang.zx2c4.com/wireguard/wgctrl"
	"magic-wan/internal/cfg"
	"magic-wan/pkg/various"
	"magic-wan/pkg/wg"
)

func buildStateFromConfigs(private *cfg.PrivateConfig, shared *cfg.SharedConfig) {
	// build root config
	selfPeer := various.ArrayFind(shared.SharedWireGuard.Peers, func(peer *cfg.Peer) bool { return peer.UID == private.NodeID })
	if selfPeer == nil {
		panic("Self peer not found")
	}
	globalRunningState = &state{
		privateKey: &private.PrivateWireGuard.PrivateKey,
		name:       selfPeer.Name,
		startPort:  shared.SharedWireGuard.StartPort,
		selfIDX:    private.NodeID,
		subnet:     shared.Router.Subnet,
		peers:      make(map[uint8]*peerState),
	}

	// build peers
	for _, peer := range shared.SharedWireGuard.Peers {
		if peer.UID == private.NodeID {
			// can't peer with self
			continue
		}

		newPeer := &peerState{
			publicKey: &peer.PublicKey,
			hostname:  peer.Hostname,
			uid:       peer.UID,
			keepalive: peer.Keepalive,
			// populated in onPeerAdded
			_parent: nil,
			ip:      nil,
		}
		onPeerAdded(newPeer)
	}
}

func configureWGInterface(client *wgctrl.Client, peer *peerState) {
	ifcName := peer.getWGName()

	// check to only create new device if it does not already exist
	devices, err := wg.GetDevices(client)
	panicOn(err)
	includes := false
	for _, dev := range devices {
		if dev.Name == ifcName {
			includes = true
			break
		}
	}
	if !includes {
		// make sure interface exists and is up
		err := wg.CreateNewDevice(ifcName)
		panicOn(err)
	}

	// build & execute configuration
	err = wg.UpdateDevice(client, ifcName, peer.BuildWGConfig())
	panicOn(err)

	// set interface IPs
	selfIP, peerIP := peer.GetLinkIPs()
	err = wg.BaseConfigureInterface(ifcName, selfIP, peerIP)
	panicOn(err)
}

func unconfigureWGInterface(client *wgctrl.Client, peer *peerState) {
	err := wg.DisableDevice(client, peer.getWGName())
	panicOn(err)
}
