package main

import (
	"fmt"
	"magic-wan/internal/cfg"
)

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}

func findSelf(private *cfg.PrivateConfig, shared *cfg.SharedConfig) (*cfg.Peer, error) {
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
	return self, nil
}
