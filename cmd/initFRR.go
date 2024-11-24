package main

import (
	"magic-wan/internal"
	"magic-wan/internal/cfg"
	"magic-wan/internal/myfrr"
	"magic-wan/pkg/frr"
)

func startFRR(private *cfg.PrivateConfig, shared *cfg.SharedConfig, interfaces []string) {
	startCfg := buildFRRBaseConfig(private, shared, interfaces)

	err := myfrr.WriteFRRConfig(frr.DEFAULT_CONFIG_PATH, startCfg)
	panicOn(err)

	// FRR was stopped as part of the baseConfigureDependencies, so we restart it now after configuration
	err = internal.FrrService.Start()
	panicOn(err)
}

func buildFRRBaseConfig(private *cfg.PrivateConfig, shared *cfg.SharedConfig, interfaces []string) string {
	self, err := findSelf(private, shared)
	panicOn(err)

	// TODO: get rid of the "interfaces" and parse them from the "wg#getDevices()"
	startCfg := myfrr.BuildBaseConfig(
		self.Name,
		self.UID,
		shared.Router.Subnet,
		interfaces,
	)

	return startCfg
}
