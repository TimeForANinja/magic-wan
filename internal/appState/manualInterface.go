package configState

import "net"

type ManualInterface struct {
	interfaceName string
	ip            *net.IP
	ospfPassive   bool
}