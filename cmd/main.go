package main

import (
	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl"
	"magic-wan/internal/appState"
	"magic-wan/pkg/cluster"
	"magic-wan/rest"
)

func main() {
	// configure logging and defer closing the logfile
	logfile, err := configureLogging()
	panicOn(err)
	defer logfile.Close()
	log.Info("Initialised Logging")

	var wgClient *wgctrl.Client
	wgClient, err = ensurePrerequisites()
	panicOn(err)
	defer wgClient.Close()
	log.Info("Checked Prerequisites")

	// build initial state from config
	var appState *configState.ApplicationState
	appState, err = initState(wgClient)
	panicOn(err)
	log.Info("Build Initial State")

	// setup frr
	frrConfig := appState.DeriveFRRState()
	err = frrConfig.StartFRR()
	panicOn(err)
	log.Info("Started FRR")

	// prepare cluster engine
	var configCluster *cluster.Cluster
	configCluster, err = appState.DeriveCluster()
	panicOn(err)
	log.Info("Initialised HA Cluster")

	// start background tasks
	failChannel := make(chan error, 1)
	go rest.StartRest(configCluster, failChannel)
	go startDNSChecks(appState)
	go configCluster.StartAnnouncements()

	// wait for failChannel to return an error
	select {
	case err = <-failChannel:
		panicOn(err)
	}
}
