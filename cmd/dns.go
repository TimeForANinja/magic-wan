package main

import log "github.com/sirupsen/logrus"

func checkDNS() {
	log.Info("Checking for DNS Changes")

	for _, pState := range globalRunningState.peers {
		// fetch new ip
		newIP := pState.resolveIP()

		// compare against state ip and trigger event if required
		if newIP.String() != pState.ip.String() {
			onPeerChangeIP(pState, newIP)
		}
	}
}
