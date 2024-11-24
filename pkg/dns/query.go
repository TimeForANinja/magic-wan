package dns

import (
	"net"
)

// QueryDNS queries the given DNS servers for the provided hostname
func QueryDNS(hostname string, ipv4only bool) (ips []string, err error) {
	addresses, err := net.LookupIP(hostname)
	if err != nil {
		return nil, err
	}
	for _, ip := range addresses {
		if !ipv4only || ip.To4() != nil {
			ips = append(ips, ip.String())
		}
	}

	return ips, nil
}
