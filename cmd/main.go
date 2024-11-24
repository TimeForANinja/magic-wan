package main

import (
	"magic-wan/rest"
	"time"
)

func main() {
	privConfig, globalConfig, globalClient := ensurePrerequisites()
	defer globalClient.Close()

	buildStateFromConfigs(privConfig, globalConfig)

	updateFRR()
	startFRR()

	// run in background, so that we can do other repeating tasks
	go rest.StartRest()

	// Until the end of time all we now do is check the DNS
	for range time.Tick(time.Minute * 1) {
		checkDNS()
	}
}
