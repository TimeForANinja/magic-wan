package main

import (
	"errors"
	"golang.zx2c4.com/wireguard/wgctrl"
	"magic-wan/internal/appState"
	"magic-wan/internal/cfg"
	"magic-wan/pkg/osUtil"
	"magic-wan/pkg/transferNetwork"
	"magic-wan/pkg/various"
	"os"
)

func initState(wgClient *wgctrl.Client) (*configState.ApplicationState, error) {
	private, shared, err := loadConfigs()
	if err != nil {
		return nil, err
	}

	// build root config
	selfPeer := various.ArrayFind(shared.SharedWireGuard.Peers, func(peer *cfg.Peer) bool { return peer.UID == private.NodeID })
	if selfPeer == nil {
		return nil, errors.New("Self peer not found")
	}

	newState := configState.NewState(
		wgClient,
		selfPeer.Name,
		&private.PrivateWireGuard.PrivateKey,
		shared.SharedWireGuard.StartPort,
		selfPeer.UID,
		shared.Router.Subnet,
	)

	for _, peer := range shared.SharedWireGuard.Peers {
		if peer.UID == private.NodeID {
			// can't peer with self
			continue
		}

		err := newState.AddWireguardInterface(
			peer.UID,
			peer.Hostname,
			&peer.PublicKey,
			peer.Keepalive,
		)
		if err != nil {
			return nil, err
		}
	}

	// build loopback manual peer
	// since 0 as a unique id is not allowed, we can use those ranges as global unique identifiers for each node.
	loopbackIP, _, _, err := transferNetwork.GetPeerToPeerNet(private.NodeID, 0, shared.Router.Subnet)
	if err != nil {
		return nil, err
	}
	err = osUtil.EnsureInterfaceHasAddress("lo", loopbackIP.String())
	if err != nil {
		return nil, err
	}
	newState.AddManualInterface("lo", &loopbackIP, true)

	return newState, nil
}

func loadConfigs() (*cfg.PrivateConfig, *cfg.SharedConfig, error) {
	shared, err := cfg.LoadSharedConfig("shared.yml")
	if err != nil {
		return nil, nil, err
	}

	// IMPROVEMENT: probably remove this
	privFilePath := os.Getenv("PRIV_FILE")
	if privFilePath == "" {
		return nil, nil, errors.New("PRIV_FILE environment variable is not set")
	}
	private, err := cfg.LoadPrivateConfig(privFilePath)
	if err != nil {
		return nil, nil, err
	}

	return private, shared, nil
}
