package myfrr

import "fmt"

func BuildBaseConfig(selfName string, uid uint8, network string, interfaces []string) string {
	interfaceConfig := ""
	for _, iface := range interfaces {
		interfaceConfig += fmt.Sprintf(`interface %s
 ip router ospf area 0.0.0.0
exit
!
`, iface)
	}

	// TODO: allow for broadcasting (+ summarizing) additional networks
	// TODO: allow for including additional interfaces
	// TODO: allow for SNAT / DNAT to be done for external interfaces

	return fmt.Sprintf(`log syslog informational
hostname %s
!
router ospf
 ospf router-id 0.0.0.%d
 network %s area 0.0.0.0
 !
 area 0.0.0.0 range %s
exit
!
%s`, selfName, uid, network, network, interfaceConfig)
}
