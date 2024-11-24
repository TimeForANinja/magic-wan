package main

import (
	"golang.zx2c4.com/wireguard/wgctrl"
	"magic-wan/pkg/wg"
	"net"
)

var globalRunningState *state

var globalClient *wgctrl.Client

func onPeerAdded(newPeer *peerState) {
	// update state
	newPeer._parent = globalRunningState
	globalRunningState.peers[newPeer.uid] = newPeer

	// call action to add
	configureWGInterface(globalClient, newPeer)

	// update frr
	updateFRR()
}

func onPeerRemoved(oldPeer *peerState) {
	// update state
	delete(globalRunningState.peers, oldPeer.uid)

	// call action to remove
	unconfigureWGInterface(globalClient, oldPeer)

	// update frr
	updateFRR()
}

func onPeerChangeIP(peer *peerState, newIP *net.UDPAddr) {
	// update state
	peer.ip = newIP

	// update relevant running config
	err := wg.UpdateDevice(globalClient, peer.getWGName(), peer.BuildWGConfig())
	panicOn(err)
}
