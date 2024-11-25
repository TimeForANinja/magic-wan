package wg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"os/exec"
)

func MustConfigureDevice(client *wgctrl.Client, ifcName string, config wgtypes.Config) {
	err := client.ConfigureDevice(ifcName, config)
	if err != nil {
		panic(fmt.Sprintf("Failed to configure WireGuard interface %s: %v", ifcName, err))
	}
}

func UpdateDevice(client *wgctrl.Client, ifcName string, config wgtypes.Config) error {
	// Apply the configuration to the interface
	err := client.ConfigureDevice(ifcName, config)
	if err != nil {
		return err
	}
	log.Infof("Successfully configured WireGuard interface %s", ifcName)
	return nil
}

func CreateNewDevice(ifcName string) error {
	cmd := exec.Command("ip", "link", "add", "dev", ifcName, "type", "wireguard")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add new wg interface \"%s\": %w, output: %s", ifcName, err, string(output))
	}
	return nil
}

func GetDevices(client *wgctrl.Client) ([]*wgtypes.Device, error) {
	devices, err := client.Devices()
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func DisableDevice(client *wgctrl.Client, ifcName string) error {
	zero := 0
	config := wgtypes.Config{
		// Clear device config.
		PrivateKey:   &wgtypes.Key{},
		ListenPort:   &zero,
		FirewallMark: &zero,

		// Clear all peers.
		ReplacePeers: true,
	}

	err := client.ConfigureDevice(ifcName, config)
	return err
}

func RemoveDevice(ifcName string) error {
	cmd := exec.Command("ip", "link", "del", "dev", ifcName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove wg interface \"%s\": %w, output: %s", ifcName, err, string(output))
	}
	return nil
}
