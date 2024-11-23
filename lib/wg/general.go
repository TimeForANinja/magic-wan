package wg

import "fmt"

func BuildName(selfIDX uint8, peerIDX uint8) string {
	return "wg1" + fmt.Sprintf("%02d", selfIDX) + fmt.Sprintf("%02d", peerIDX)
}

func BuildBaseConfig(selfName, privateKey, address string, listenPort, mtu uint16, peerName, publicKey, endpoint, allowedIPs string) string {
	return fmt.Sprintf(`[Interface]
# Name = %s
PrivateKey = %s
Address = %s
ListenPort = %d
Table = off
MTU = %d

[Peer]
# Name = %s
PublicKey = %s
Endpoint = %s
AllowedIPs = %s
`, selfName, privateKey, address, listenPort, mtu, peerName, publicKey, endpoint, allowedIPs)
}

func CalculatePort(startPort uint16, peerIDX uint8) *int {
	// TODO: check if this makes sense
	port := int(startPort + uint16(peerIDX))
	return &port
}
