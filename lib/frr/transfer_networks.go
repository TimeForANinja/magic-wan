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
func getNthAddress(network *net.IPNet, n int) (net.IP, error) {
	ip := network.IP.To4()
	if ip == nil {
		return nil, fmt.Errorf("only IPv4 is supported")
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
	if len(newIP) < 4 {
		paddedIP := make([]byte, 4)
		copy(paddedIP[4-len(newIP):], newIP)
		newIP = paddedIP
	}

	return newIP, nil
}

func getPeerToPeerNet(myIDX, peerIDX uint8, baseNet string) (myIP, peerIP net.IP, p2pNet *net.IPNet, err error) {
	_, ipNet, err := net.ParseCIDR(baseNet)
	if err != nil {
		return
	}

	myIP, err = getNthAddress(ipNet, int(calcOffset(myIDX, peerIDX)))
	if err != nil {
		return
	}

	peerIP, err = getNthAddress(ipNet, int(calcOffset(peerIDX, myIDX)))
	if err != nil {
		return
	}

	_, p2pNet, err = net.ParseCIDR(myIP.String() + "/31")
	if err != nil {
		return
	}

	return
}
