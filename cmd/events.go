package main

import (
	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl"
	"magic-wan/pkg/osUtil"
	"magic-wan/pkg/wg"
	"net"
)

var globalRunningState *state

var globalClient *wgctrl.Client

func onPeerAdded(newPeer *peerState, skipFRR bool) {
	log.WithFields(log.Fields{
		"peer": newPeer,
	}).Info("onPeerAdded")
	// update state
	newPeer._parent = globalRunningState
	newPeer.ip = newPeer.resolveIP()
	globalRunningState.peers[newPeer.uid] = newPeer

	// call action to add
	configureWGInterface(globalClient, newPeer)

	// update other modules (if requested)
	if !skipFRR {
		// update frr
		updateFRR()
	}
	addPeerToCluster(newPeer)
}

func onPeerRemoved(oldPeer *peerState) {
	log.WithFields(log.Fields{
		"peer": oldPeer,
	}).Info("onPeerRemoved")
	// update state
	delete(globalRunningState.peers, oldPeer.uid)

	// call action to remove
	unconfigureWGInterface(globalClient, oldPeer)

	// update other modules
	updateFRR()
	removePeerFromCluster(oldPeer)
}

func onManualInterfaceAdded(iface *ManualInterface) {
	log.WithFields(log.Fields{
		"iface": iface,
	}).Info("onManualInterfaceAdded")

	// add IP to ip if not already there
	if !osUtil.InterfaceHasAddress(iface.interfaceName, iface.ip.String()) {
		err := osUtil.SetInterfaceAddress(iface.interfaceName, iface.ip.String())
		panicOn(err)
	}
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
