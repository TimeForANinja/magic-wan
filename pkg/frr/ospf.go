package frr

import (
	"magic-wan/pkg/various"
	"strings"
)

const (
	DAEMON_CONFIG_PATH  = "/etc/frr/daemons"
	DEFAULT_CONFIG_PATH = "/etc/frr/frr.conf"
)

// EnableOSPF enables the OSPF daemon in the FRR daemons configuration file by setting the ospfd entry to "ospfd=yes".
func EnableOSPF() error {
	return various.GenericFileProcessor(
		DAEMON_CONFIG_PATH,
		func(line string) string {
			if strings.HasPrefix(line, "ospfd=") {
				return "ospfd=yes"
			}
			return line
		})
}
