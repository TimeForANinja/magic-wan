package osUtil

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

type Service struct {
	Name string
}

func (s *Service) Start() error {
	cmd := exec.Command("systemctl", "start", s.Name)
	log.WithFields(log.Fields{
		"cmd":     "systemctl start <service>",
		"service": s.Name,
	}).Debug("Calling systemctl")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start %s service: %w, output: %s", s.Name, err, string(output))
	}
	return nil
}

func (s *Service) StartEnable() error {
	cmd := exec.Command("systemctl", "enable", "--now", s.Name)
	log.WithFields(log.Fields{
		"cmd":     "systemctl enable --now <service>",
		"service": s.Name,
	}).Debug("Calling systemctl")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to enable and start %s service: %w, output: %s", s.Name, err, string(output))
	}
	return nil
}

func (s *Service) Enable() error {
	cmd := exec.Command("systemctl", "enable", s.Name)
	log.WithFields(log.Fields{
		"cmd":     "systemctl enable <service>",
		"service": s.Name,
	}).Debug("Calling systemctl")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to enable %s service: %w, output: %s", s.Name, err, string(output))
	}
	return nil
}

func (s *Service) Stop() error {
	cmd := exec.Command("systemctl", "stop", s.Name)
	log.WithFields(log.Fields{
		"cmd":     "systemctl stop <service>",
		"service": s.Name,
	}).Debug("Calling systemctl")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop %s service: %w, output: %s", s.Name, err, string(output))
	}
	return nil
}

func (s *Service) Disable() error {
	cmd := exec.Command("systemctl", "disable", s.Name)
	log.WithFields(log.Fields{
		"cmd":     "systemctl disable <service>",
		"service": s.Name,
	}).Debug("Calling systemctl")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to disable %s service: %w, output: %s", s.Name, err, string(output))
	}
	return nil
}

func (s *Service) Reload() error {
	// IMPROVMENT: check if this works well for frr / ospfd
	cmd := exec.Command("systemctl", "reload", s.Name)
	log.WithFields(log.Fields{
		"cmd":     "systemctl reload <service>",
		"service": s.Name,
	}).Debug("Calling systemctl")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to disable %s service: %w, output: %s", s.Name, err, string(output))
	}
	return nil
}

func (s *Service) GetStatus() (string, error) {
	// IMPROVMENT: check if this works and what possible values are
	cmd := exec.Command("systemctl", "show", "--no-pager", s.Name)
	log.WithFields(log.Fields{
		"cmd":     "systemctl show --no-pager <service>",
		"service": s.Name,
	}).Debug("Calling systemctl")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get status of %s service: %w, output: %s", s.Name, err, string(output))
	}
	// Extract ActiveState from output
	var activeState string
	for _, line := range strings.Split(string(output), "\n") {
		if strings.HasPrefix(line, "ActiveState=") {
			activeState = strings.TrimPrefix(line, "ActiveState=")
			break
		}
	}
	return activeState, nil
}
