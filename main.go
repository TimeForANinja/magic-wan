package main

import (
	"fmt"
	"magic-wan/lib/cfg"
	"magic-wan/lib/osUtil"
	"magic-wan/lib/transfer"
	"magic-wan/lib/wg"
)

func main() {
	// TODO: uncomment
	// checkDependencies()

	err, privateCfg, sharedCfg := loadSettings()
	panicOn(err)
	err = buildConfigs(privateCfg, sharedCfg)
	panicOn(err)

	// TODO: remove
	return

	err, service := osUtil.InstallAsService()
	panicOn(err)
	panicOn(service.Enable())
	panicOn(service.Start())
}

func checkDependencies() {
	if !osUtil.IsLinuxArchitecture() {
		panic("Unsupported architecture")
	}

	err := osUtil.InstallPackages([]string{
		"wireguard",
		"frr",
	})
	panicOn(err)
}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}

func loadSettings() (error, *cfg.PrivateConfig, *cfg.SharedConfig) {
	shared, err := cfg.LoadSharedConfig("shared.yml")
	if err != nil {
		return err, nil, nil
	}
	private, err := cfg.LoadPrivateConfig("private.yml")
	if err != nil {
		return err, nil, nil
	}

	return nil, private, shared
}

func buildConfigs(private *cfg.PrivateConfig, shared *cfg.SharedConfig) error {
	// find config for self
	var self *cfg.Peer
	for _, peer := range shared.SharedWireGuard.Peers {
		if private.NodeID == peer.UID {
			self = &peer
			continue
		}
	}
	if self == nil {
		return fmt.Errorf("failed to find self in shared config")
	}

	// build peer != self connections
	for _, peer := range shared.SharedWireGuard.Peers {
		if private.NodeID == peer.UID {
			continue
		}

		wgName := wg.BuildName(private.NodeID, peer.UID)
		myIP, peerIP, sharedNW, err := transfer.GetPeerToPeerNet(private.NodeID, peer.UID, shared.Router.Subnet)
		if err != nil {
			return err
		}
		fmt.Printf("Tunnel %s for (%s -> %s)\n", wgName, myIP.String(), peerIP.String())

		cfg := wg.BuildBaseConfig(self.Name, private.PrivateWireGuard.PrivateKey, myIP.String(), shared.SharedWireGuard.StartPort+uint16(private.NodeID), shared.MTU, peer.Name, peer.PublicKey, peer.Hostname, sharedNW.String())
		fmt.Println(cfg)
	}

	// TODO: implement
	return nil
}
