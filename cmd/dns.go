package main

import (
	"fmt"
	"time"
)

func checkDNS() {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "Checking DNS")

	for _, pState := range globalRunningState.peers {
		// fetch new ip
		newIP := pState.resolveIP()

		// compare against state ip and trigger event if required
		if newIP.String() != pState.ip.String() {
			onPeerChangeIP(pState, newIP)
		}
	}
}
