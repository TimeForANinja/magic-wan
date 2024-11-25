package main

import (
	log "github.com/sirupsen/logrus"
	"magic-wan/internal/cfg"
	"magic-wan/rest"
	"time"
)

func main() {
	var privateConfig *cfg.PrivateConfig
	var sharedConfig *cfg.SharedConfig
	privateConfig, sharedConfig, globalClient = ensurePrerequisites()
	defer globalClient.Close()
	log.Info("Checked Prerequisites")

	buildStateFromConfigs(privateConfig, sharedConfig)
	log.Info("Build Initial State")

	updateFRR()
	// TODO: check if "updateFRR" which includes a restart is enough
	startFRR()
	log.Info("Started FRR")

	// run in background, so that we can do other repeating tasks
	go rest.StartRest()

	// Until the end of time all we now do is check the DNS
	ticker := time.NewTicker(time.Minute * 1)
	defer ticker.Stop()
	for range ticker.C {
		checkDNS()
	}
}
