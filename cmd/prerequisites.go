package main

import (
	"golang.zx2c4.com/wireguard/wgctrl"
	"magic-wan/pkg"
	"magic-wan/pkg/frr"
	"magic-wan/pkg/osUtil"
	"magic-wan/pkg/wg"
)

func ensurePrerequisites() (*wgctrl.Client, error) {
	err := checkDependencies()
	if err != nil {
		return nil, err
	}

	err = baseConfigureRouting()
	if err != nil {
		return nil, err
	}

	// get wireguard controller
	wgClient := wg.MustCreateController()

	// as a preparation all existing wg interfaces will be removed
	err = removeAllWGDevices(wgClient)
	if err != nil {
		return nil, err
	}

	// register self as a service to auto-start
	/*
		service, err := osUtil.InstallAsService()
		panicOn(err)
		err = service.Enable()
		panicOn(err)
	*/

	return wgClient, nil
}

func removeAllWGDevices(client *wgctrl.Client) error {
	devices, err := wg.GetDevices(client)
	if err != nil {
		return err
	}

	for _, device := range devices {
		err := wg.RemoveDevice(device.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func baseConfigureRouting() error {
	err := osUtil.EnableIPV4Routing()
	if err != nil {
		return err
	}

	// frr is default-enabled after installation so let's disable it
	err = pkg.FrrService.Disable()
	if err != nil {
		return err
	}
	err = pkg.FrrService.Stop()
	if err != nil {
		return err
	}

	// enable the ospf daemon inside the frr library
	return frr.EnableOSPF()
}

func checkDependencies() error {
	if !osUtil.IsLinuxArchitecture() {
		panic("Unsupported architecture")
	}

	return osUtil.InstallPackages([]string{
		"wireguard",
		"frr",
		"frr-pythontools",
	})
}
