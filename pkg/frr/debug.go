package frr

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Debug runs specified OSPF debug commands and returns their output.
func Debug() (string, error) {
	commands := []string{
		"show ip ospf neighbor",
		"show ip ospf route",
		"show ip ospf interface",
		"show running-config",
	}
	var output bytes.Buffer

	for _, cmd := range commands {
		command := exec.Command("vtysh", "-c", cmd)
		cmdOutput, err := command.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to run command '%s': %w, output: %s", cmd, err, string(cmdOutput))
		}
		output.WriteString(fmt.Sprintf("Command: %s\nOutput:\n%s\n", cmd, string(cmdOutput)))
	}

	return output.String(), nil
}
