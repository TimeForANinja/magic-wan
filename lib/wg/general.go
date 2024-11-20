package wg

import "fmt"

func buildName(selfIDX uint8, peerIDX uint8) string {
	return "wg1" + fmt.Sprintf("%02d", selfIDX) + fmt.Sprintf("%02d", peerIDX)
}

func buildBaseConfig(selfName, privateKey, address string, listenPort, mtu int, peerName, publicKey, endpoint, allowedIPs string) string {
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
