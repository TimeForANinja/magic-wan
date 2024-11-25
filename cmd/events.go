package main

import (
	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl"
	"magic-wan/pkg/wg"
	"net"
)

var globalRunningState *state

var globalClient *wgctrl.Client

func onPeerAdded(newPeer *peerState) {
	log.WithFields(log.Fields{
		"peer": newPeer,
	}).Info("onPeerAdded")
	// update state
	newPeer._parent = globalRunningState
	globalRunningState.peers[newPeer.uid] = newPeer

	// call action to add
	configureWGInterface(globalClient, newPeer)

	// update frr
	updateFRR()
}

func onPeerRemoved(oldPeer *peerState) {
	log.WithFields(log.Fields{
		"peer": oldPeer,
	}).Info("onPeerRemoved")
	// update state
	delete(globalRunningState.peers, oldPeer.uid)

	// call action to remove
	unconfigureWGInterface(globalClient, oldPeer)

	// update frr
	updateFRR()
}

func onPeerChangeIP(peer *peerState, newIP *net.UDPAddr) {
	log.WithFields(log.Fields{
		"peer":  peer,
		"newIP": newIP,
	}).Info("onPeerRemoved")
	// update state
	peer.ip = newIP

	// update relevant running config
	err := wg.UpdateDevice(globalClient, peer.getWGName(), peer.BuildWGConfig())
	panicOn(err)
}