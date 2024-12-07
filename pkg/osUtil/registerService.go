package osUtil

import (
	"fmt"
	"os"
	"path/filepath"
)

const WanServiceName = "magic-wan"

func InstallAsService() (*Service, error) {
	absPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	err = installNewService(WanServiceName, absPath)
	if err != nil {
		return nil, err
	}
	return &Service{Name: WanServiceName}, nil
}

func installNewService(serviceName, executablePath string) error {
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
	err := WriteFile(servicePath, serviceContent)
	if err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	return nil
}
