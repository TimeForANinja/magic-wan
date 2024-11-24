package osUtil

import (
	"fmt"
	"os/exec"
)

type Service struct {
	Name string
}

func (s *Service) Start() error {
	cmd := exec.Command("systemctl", "start", s.Name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start %s service: %w, output: %s", s.Name, err, string(output))
	}
	return nil
}

func (s *Service) StartEnable() error {
	cmd := exec.Command("systemctl", "enable", "--now", s.Name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to enable and start %s service: %w, output: %s", s.Name, err, string(output))
	}
	return nil
}

func (s *Service) Enable() error {
	cmd := exec.Command("systemctl", "enable", s.Name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to enable %s service: %w, output: %s", s.Name, err, string(output))
	}
	return nil
}

func (s *Service) Stop() error {
	cmd := exec.Command("systemctl", "stop", s.Name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop %s service: %w, output: %s", s.Name, err, string(output))
	}
	return nil
}

func (s *Service) Disable() error {
	cmd := exec.Command("systemctl", "disable", s.Name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to disable %s service: %w, output: %s", s.Name, err, string(output))
	}
	return nil
}

func (s *Service) Reload() error {
	// TODO: check if this works well for frr / ospfd
	cmd := exec.Command("systemctl", "reload", s.Name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to disable %s service: %w, output: %s", s.Name, err, string(output))
	}
	return nil
}