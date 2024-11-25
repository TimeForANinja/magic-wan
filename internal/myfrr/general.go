package myfrr

import (
	"fmt"
	"net"
)

func buildInterfaceConfig(iface string, passive bool) string {
	if !passive {
		return fmt.Sprintf(`interface %s
 ip router ospf area 0.0.0.0
exit
!
`, iface)
	}
	return fmt.Sprintf(`interface %s
 ip ospf passive
 ip router ospf area 0.0.0.0
exit
!
`, iface)
}

func BuildBaseConfig(selfName string, uid uint8, network *net.IPNet, activeInterfaces []string, passiveInterfaces []string) string {
	interfaceConfig := ""
	for _, iface := range activeInterfaces {
		interfaceConfig += buildInterfaceConfig(iface, false)
	}
	for _, iface := range passiveInterfaces {
		interfaceConfig += buildInterfaceConfig(iface, true)
	}

	// TODO: allow for broadcasting (+ summarizing) additional networks
	// TODO: allow for including additional interfaces
	// TODO: allow for SNAT / DNAT to be done for external interfaces
	// TODO: remove debug log

	return fmt.Sprintf(`log file /var/log/frr/debug.log debugging
log syslog informational
hostname %s
!
router ospf
 ospf router-id 0.0.0.%d
 network %s area 0.0.0.0
 !
 area 0.0.0.0 range %s
exit
!
%s`, selfName, uid, network.String(), network.String(), interfaceConfig)
}
