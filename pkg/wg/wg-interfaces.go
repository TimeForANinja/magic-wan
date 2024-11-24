package wg

import (
	"fmt"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"log"
	"os/exec"
	"strings"
)

func UpdateDevice(client *wgctrl.Client, ifcName string, config wgtypes.Config) error {
	// Apply the configuration to the interface
	err := client.ConfigureDevice(ifcName, config)
	if err != nil {
		return err
	}
	log.Printf("Successfully configured WireGuard interface %s", ifcName)
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

func BaseConfigureInterface(ifcName string, selfIP, peerIP string) error {
	commands := [][]string{
		{"ip", "address", "add", "dev", ifcName, selfIP, "peer", peerIP},
		{"ip", "link", "set", "up", "dev", ifcName},
	}

	for _, cmd := range commands {
		command := exec.Command(cmd[0], cmd[1:]...)
		cmdOutput, err := command.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to run command '%s': %w, output: %s", strings.Join(cmd, " "), err, string(cmdOutput))
		}
	}

	return nil
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

func startService(service string) error {
	cmd := exec.Command("systemctl", "enable", "--now", service)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to enable and start \"%s\" service: %w, output: %s", service, err, string(output))
	}
	return nil
}