package frr

import (
	"bytes"
	"fmt"
	"magic-wan/lib"
	"magic-wan/lib/osUtil"
	"os"
	"os/exec"
	"strings"
)

func GetFRRService() *osUtil.Service {
	return &osUtil.Service{
		Name: "frr",
	}
}

// EnableOSPF enables the OSPF daemon in the FRR daemons configuration file by setting the ospfd entry to "ospfd=yes".
func EnableOSPF() error {
	return lib.GenericFileProcessor(
		"/etc/frr/daemons",
		func(line string) string {
			if strings.HasPrefix(line, "ospfd=") {
				return "ospfd=yes"
			}
			return line
		})
}

// WriteFRRConfig writes a new FRR configuration file.
func WriteFRRConfig(configFilePath string, configContent string) error {
	file, err := os.Create(configFilePath)
	if err != nil {
		return fmt.Errorf("failed to create FRR config file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(configContent); err != nil {
		return fmt.Errorf("failed to write to FRR config file: %v", err)
	}

	// Restart FRR service to apply changes
	cmd := exec.Command("systemctl", "restart", "frr")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to restart FRR service: %v, output: %s", err, string(output))
	}

	return nil
}

// Debug runs specified OSPF debug commands and returns their output.
func Debug() (string, error) {
	commands := []string{
		"show ip ospf neighbor",
		"show ip ospf route",
		"show ip ospf interface",
	}
	var output bytes.Buffer

	for _, cmd := range commands {
		command := exec.Command("vtysh", "-c", cmd)
		cmdOutput, err := command.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to run command '%s': %v, output: %s", cmd, err, string(cmdOutput))
		}
		output.WriteString(fmt.Sprintf("Command: %s\nOutput:\n%s\n", cmd, string(cmdOutput)))
	}

	return output.String(), nil
}
