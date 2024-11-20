package osUtil

import (
	"fmt"
	"os"
	"path/filepath"
)

func InstallAsService() (error, *Service) {
	name := "magic-wan"

	absPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err), nil
	}

	err = installAsService(name, absPath)
	if err != nil {
		return err, nil
	}
	return nil, &Service{Name: name}
}

func installAsService(serviceName, executablePath string) error {
	servicePath := fmt.Sprintf("/etc/systemd/system/%s.service", serviceName)

	// Create the systemd service file content
	serviceContent := `[Unit]
Description=Magic-WAN Router
After=systemd-sysctl.service

[Service]
ExecStart=%s
Restart=always
User=root

[Install]
WantedBy=multi-user.target
`

	// Prepare the content with the provided executable path
	serviceContent = fmt.Sprintf(serviceContent, executablePath)

	// Write the service file to the systemd directory
	if err := os.WriteFile(servicePath, []byte(serviceContent), 0644); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	return nil
}

func main() {
	// Example usage
	exePath, err := filepath.Abs("path/to/your/executable")
	if err != nil {
		fmt.Printf("Failed to get absolute path: %v\n", err)
		return
	}

	if err := installAsService("magicwan", exePath); err != nil {
		fmt.Printf("Error installing service: %v\n", err)
	}
}
