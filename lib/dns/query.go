package dns

import (
	"net"
)

// QueryDNS queries the given DNS servers for the provided hostname
func QueryDNS(hostname string) (ips []string, err error) {
	addresses, err := net.LookupIP(hostname)
	if err != nil {
		return nil, err
	}
	for _, ip := range addresses {
		if ipv4 := ip.To4(); ipv4 != nil {
			ips = append(ips, ipv4.String())
		}
	}

	return ips, nil
}
