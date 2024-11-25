package main

import (
	"magic-wan/rest/cluster"
)

func addPeerToCluster(newPeer *peerState) {
	_, peerIP, _ := newPeer.GetConnectionTo(0)

	cluster.AddPeer(peerIP.String())
}

func removePeerFromCluster(oldPeer *peerState) {
	_, peerIP, _ := oldPeer.GetConnectionTo(0)

	cluster.RemovePeer(peerIP.String())
}
