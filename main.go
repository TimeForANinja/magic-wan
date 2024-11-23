package main

import (
	"fmt"
	"magic-wan/lib/cfg"
	"magic-wan/lib/osUtil"
	"magic-wan/lib/wg"
)

func main() {
	// TODO: uncomment
	// checkDependencies()

	err, privateCfg, sharedCfg := loadSettings()
	panicOn(err)
	wgPeers, err := buildConfigs(privateCfg, sharedCfg)
	fmt.Println(wgPeers)
	panicOn(err)

	client := wg.MustCreateController()
	defer client.Close()
	for _, peer := range wgPeers {
		err, config := peer.ToWGConfig()
		panicOn(err)
		wg.CreateNewDevice(client, peer.Name, config)
	}

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

func buildConfigs(private *cfg.PrivateConfig, shared *cfg.SharedConfig) ([]cfg.WGLink, error) {
	// find config for self
	var self *cfg.Peer
	for _, peer := range shared.SharedWireGuard.Peers {
		if private.NodeID == peer.UID {
			self = &peer
			continue
		}
	}
	if self == nil {
		return nil, fmt.Errorf("failed to find self in shared config")
	}

	peers := make([]cfg.WGLink, 0)

	// build peer != self connections
	for _, peer := range shared.SharedWireGuard.Peers {
		if private.NodeID == peer.UID {
			continue
		}

		peers = append(peers, cfg.WGLink{
			General:     shared,
			Name:        wg.BuildName(self.UID, peer.UID),
			Self:        self,
			SelfPrivate: private,
			Peer:        &peer,
		})
	}

	return peers, nil
}
