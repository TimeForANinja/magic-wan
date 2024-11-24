package wg

import (
	"fmt"
	"os/exec"
)

// Debug runs specified OSPF debug commands and returns their output.
func Debug() (string, error) {
	command := exec.Command("wg")
	cmdOutput, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run command 'wg': %w, output: %s", err, string(cmdOutput))
	}
	return fmt.Sprintf("Command: wg\nOutput:\n%s\n", string(cmdOutput)), nil

}
