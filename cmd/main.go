package main

import (
	"magic-wan/pkg/wg"
	"magic-wan/rest"
	"time"
)

func main() {
	privateCfg, sharedCfg := ensurePrerequisites()

	// get wireguard controller
	client := wg.MustCreateController()
	defer client.Close()

	createdInterfaces := doStartingConfig(client, privateCfg, sharedCfg)

	startFRR(privateCfg, sharedCfg, createdInterfaces)

	// run in background, so that we can do other repeating tasks
	go rest.StartRest()

	// Until the end of time all we now do is check the DNS
	for range time.Tick(time.Minute * 1) {
		checkDNS()
	}
}
