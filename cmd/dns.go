package main

import (
	log "github.com/sirupsen/logrus"
	"magic-wan/internal/appState"
	"time"
)

const dnsCheckInterval = 1 * time.Minute

func startDNSChecks(state *configState.ApplicationState) {
	// Until the end of time all we now do is check the DNS
	ticker := time.NewTicker(dnsCheckInterval)
	defer ticker.Stop()
	for range ticker.C {
		doDNSCheck(state)
	}
}

func doDNSCheck(state *configState.ApplicationState) {
	log.Info("Checking for DNS Changes")

	for _, pState := range state.GetPeers() {
		// fetch new ip
		newIP, err := pState.ResolveAddr()
		if err != nil {
			// simply ignore if we can't resolve
			// code will keep running with the last resolved ip
			continue
		}

		// compare against state ip and trigger event if required
		if newIP.String() != pState.CachedAddr().String() {
			err = pState.NotifyIPChange(newIP)
			log.Errorf("Unable to update Interface IP: %v", err)
		}
	}
}
