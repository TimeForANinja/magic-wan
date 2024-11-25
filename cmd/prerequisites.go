package main

import (
	"golang.zx2c4.com/wireguard/wgctrl"
	"magic-wan/internal"
	"magic-wan/internal/cfg"
	"magic-wan/pkg/frr"
	"magic-wan/pkg/osUtil"
	"magic-wan/pkg/wg"
)

func ensurePrerequisites() (*cfg.PrivateConfig, *cfg.SharedConfig, *wgctrl.Client) {
	checkDependencies()

	baseConfigureDependencies()

	privateCfg, sharedCfg, err := loadSettings()
	panicOn(err)

	// get wireguard controller
	wgClient := wg.MustCreateController()

	// as a preparation all existing wg interfaces will be removed
	removeAllWGDevices(wgClient)

	// register self as a service to auto-start
	service, err := osUtil.InstallAsService()
	panicOn(err)
	err = service.Enable()
	panicOn(err)

	return privateCfg, sharedCfg, wgClient
}

func removeAllWGDevices(client *wgctrl.Client) {
	devices, err := wg.GetDevices(client)
	panicOn(err)

	for _, device := range devices {
		err := wg.RemoveDevice(device.Name)
		panicOn(err)
	}
}

func baseConfigureDependencies() {
	err := osUtil.EnableIPV4Routing()
	panicOn(err)

	// frr is default-enabled after installation so let's disable it
	err = internal.FrrService.Disable()
	panicOn(err)
	err = internal.FrrService.Stop()
	panicOn(err)

	err = frr.EnableOSPF()
	panicOn(err)
}

func checkDependencies() {
	if !osUtil.IsLinuxArchitecture() {
		panic("Unsupported architecture")
	}

	err := osUtil.InstallPackages([]string{
		"wireguard",
		"frr",
	})
	panicOn(err)
}

func loadSettings() (*cfg.PrivateConfig, *cfg.SharedConfig, error) {
	shared, err := cfg.LoadSharedConfig("shared.yml")
	if err != nil {
		return nil, nil, err
	}
	private, err := cfg.LoadPrivateConfig("private.yml")
	if err != nil {
		return nil, nil, err
	}

	return private, shared, nil
}
