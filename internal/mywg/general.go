package mywg

import "fmt"

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
