package main

import (
	"magic-wan/internal"
	"magic-wan/internal/myfrr"
	"magic-wan/pkg/frr"
	"magic-wan/pkg/various"
)

func startFRR() {
	// FRR was stopped as part of the baseConfigureDependencies, so we restart it now after configuration
	err := internal.FrrService.Start()
	panicOn(err)
}

func updateFRR() {
	frrConfigString := buildFRRBaseConfig(globalRunningState)

	err := myfrr.WriteFRRConfig(frr.DEFAULT_CONFIG_PATH, frrConfigString)
	panicOn(err)
}

func buildFRRBaseConfig(state *state) string {
	interfaces := various.ReflectMap(state.peers, func(peer *peerState) string {
		return peer.getWGName()
	})

	startCfg := myfrr.BuildBaseConfig(
		state.name,
		state.selfIDX,
		state.subnet,
		interfaces,
	)

	return startCfg
}
