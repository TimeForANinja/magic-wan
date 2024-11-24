package main

import (
	"golang.zx2c4.com/wireguard/wgctrl"
	"magic-wan/internal/cfg"
	"magic-wan/pkg/transferNetwork"
	"magic-wan/pkg/wg"
)

func doStartingConfig(client *wgctrl.Client, private *cfg.PrivateConfig, shared *cfg.SharedConfig) []string {
	configs := collectDataForInterface(private, shared)

	createdInterfaces := configureWG(client, configs)

	return createdInterfaces
}

func collectDataForInterface(private *cfg.PrivateConfig, shared *cfg.SharedConfig) []cfg.P2P {
	self, err := findSelf(private, shared)
	panicOn(err)

	peers := make([]cfg.P2P, 0)
	for _, peer := range shared.SharedWireGuard.Peers {
		if private.NodeID == peer.UID {
			// Skip creating a connection to self
			continue
		}

		peers = append(peers, cfg.P2P{
			General:     shared,
			Self:        self,
			SelfPrivate: private,
			Peer:        &peer,
		})
	}
	return peers
}

func configureWG(client *wgctrl.Client, configs []cfg.P2P) []string {
	createdInterfaces := make([]string, 0)

	for _, config := range configs {
		name := transferNetwork.BuildWireguardInterfaceName(config.Self.UID, config.Peer.UID)

		// TODO: add a chack to unly create new device if it does not already exist
		// make sure interface exists and is up
		err := wg.CreateNewDevice(name)
		panicOn(err)

		// build & execute configuration
		myPort, peerPort := transferNetwork.CalculatePorts(config.General.SharedWireGuard.StartPort, config.Self.UID, config.Peer.UID)
		wgConfig, err := config.ToConfig(myPort, peerPort)
		panicOn(err)
		err = wg.UpdateDevice(client, name, wgConfig)
		panicOn(err)

		// store in list of new interfaces
		createdInterfaces = append(createdInterfaces, name)
	}

	return createdInterfaces
}
