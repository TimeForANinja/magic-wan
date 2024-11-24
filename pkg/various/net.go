package various

import (
	"fmt"
	"net"
)

func ResolveHostname(hostname string, port int) (*net.UDPAddr, error) {
	var endpoint *net.UDPAddr
	var err error

	if hostname != "" {
		addr := fmt.Sprintf("%s:%d", hostname, port)
		endpoint, err = net.ResolveUDPAddr("udp", addr)
		if err != nil {
			return endpoint, err
		}
	}
	return endpoint, nil
}
