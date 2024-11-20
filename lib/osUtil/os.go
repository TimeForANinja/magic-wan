package osUtil

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// EnableIPV4Routing enables IPv4 routing on the Linux host OS.
func EnableIPV4Routing() error {
	const path = "/proc/sys/net/ipv4/ip_forward"
	const value = "1" // Value to set to enable IPv4 routing

	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", path, err)
	}
	defer file.Close()

	if _, err = file.WriteString(value); err != nil {
		return fmt.Errorf("failed to write to file %s: %v", path, err)
	}

	return nil
}

// InstallPackages installs the specified packages using 'apt' on a Linux system.
func InstallPackages(packages []string) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("not a Linux system")
	}
	for _, pkg := range packages {
		cmd := exec.Command("apt", "install", "-y", pkg)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to install %s: %v\nOutput: %s", pkg, err, output)
		}
	}
	return nil
}

// IsLinuxArchitecture checks if the OS is running on a specified architecture (arm64 or amd64).
func IsLinuxArchitecture() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	if runtime.GOARCH != "arm64" && runtime.GOARCH != "amd64" {
		return false
	}
	return true
}
