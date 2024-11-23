package frr

import (
	"fmt"
	"math/big"
	"net"
)

func calcOffset(target, source uint8) (offset uint32) {
	if target > source {
		return calcOffset(source, target) + 1
	}

	return uint32(source-1)*uint32(source) + 2*uint32(target)
}

// getNthAddress calculates the nth IP address in a given IP network
func getNthAddress(network *net.IPNet, n int) (net.IP, int, error) {
	ip := network.IP
	var ipLength int
	var prefix int

	// Determine whether it's IPv4 or IPv6
	if ipv4 := ip.To4(); ipv4 != nil {
		ip = ipv4
		prefix = 31
		ipLength = 4
	} else if ipv6 := ip.To16(); ipv6 != nil && ip.To4() == nil {
		ip = ipv6
		prefix = 127
		ipLength = 16
	} else {
		return nil, -1, fmt.Errorf("only IPv4 or IPv6 is supported")
	}

	// Turn IP into a *big.Int
	ipInt := big.NewInt(0)
	ipInt.SetBytes(ip)

	// Create a big.Int for the increment
	nBigInt := big.NewInt(int64(n))

	// Add the nth increment to the IP address
	ipInt.Add(ipInt, nBigInt)

	// Create the new IP address
	newIP := ipInt.Bytes()
	if len(newIP) < ipLength {
		paddedIP := make([]byte, ipLength)
		copy(paddedIP[ipLength-len(newIP):], newIP)
		newIP = paddedIP
	}

	return newIP, prefix, nil
}

func getPeerToPeerNet(myIDX, peerIDX uint8, baseNet string) (myIP, peerIP net.IP, p2pNet *net.IPNet, err error) {
	_, ipNet, err := net.ParseCIDR(baseNet)
	if err != nil {
		return
	}

	myIP, prefix, err := getNthAddress(ipNet, int(calcOffset(myIDX, peerIDX)))
	if err != nil {
		return
	}

	peerIP, _, err = getNthAddress(ipNet, int(calcOffset(peerIDX, myIDX)))
	if err != nil {
		return
	}

	_, p2pNet, err = net.ParseCIDR(fmt.Sprintf("%s/%d", myIP.String(), prefix))
	if err != nil {
		return
	}

	return
}
