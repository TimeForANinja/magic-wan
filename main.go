package main

import "magic-wan/lib/osUtil"

func main() {
	checkDependencies()

	err, settings := loadSettings()
	panicOn(err)
	err = buildConfigs(settings)
	panicOn(err)

	err, service := osUtil.InstallAsService()
	panicOn(err)
	panicOn(service.Enable())
	panicOn(service.Start())
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

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}

func loadSettings() (error, any) {
	// TODO: implement
	return nil, nil
}

func buildConfigs(settings any) error {
	// TODO: implement
	return nil
}
