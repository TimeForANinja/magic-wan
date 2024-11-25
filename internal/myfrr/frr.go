package myfrr

import (
	"fmt"
	"os"
	"os/exec"
)

// WriteFRRConfig writes a new FRR configuration file.
func WriteFRRConfig(configFilePath string, configContent string) error {
	file, err := os.Create(configFilePath)
	if err != nil {
		return fmt.Errorf("failed to create FRR config file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(configContent)
	if err != nil {
		return fmt.Errorf("failed to write to FRR config file: %w", err)
	}

	// Restart FRR service to apply changes
	cmd := exec.Command("systemctl", "restart", "frr")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to restart FRR service: %w, output: %s", err, string(output))
	}

	return nil
}
