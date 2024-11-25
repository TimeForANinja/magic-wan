package main

import (
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"magic-wan/pkg/transferNetwork"
	"magic-wan/pkg/various"
	"net"
	"time"
)

type peerState struct {
	publicKey *wgtypes.Key
	_parent   *state
	hostname  string
	uid       uint8
	keepalive bool
	ip        *net.UDPAddr
}

func (p *peerState) BuildWGConfig() wgtypes.Config {
	selfPort, _ := p.GetPorts()
	return wgtypes.Config{
		PrivateKey:   p._parent.privateKey, // Generate a private key for this interface
		ListenPort:   &selfPort,
		ReplacePeers: true, // Replace existing peers with the provided ones
		Peers: []wgtypes.PeerConfig{
			{
				PublicKey:                   *p.publicKey,
				Endpoint:                    p.ip,
				PersistentKeepaliveInterval: p.GetKeepalive(),
				AllowedIPs: []net.IPNet{
					*p.GetLinkNetwork(),
					{
						IP:   net.ParseIP("0.0.0.0"),
						Mask: net.CIDRMask(0, 32),
					},
					{
						IP:   net.ParseIP("::"),
						Mask: net.CIDRMask(0, 128),
					},
				},
			},
		},
	}
}

func (p *peerState) GetLinkNetwork() *net.IPNet {
	_, _, nw := p.GetConnectionTo(p._parent.selfIDX)
	return nw
}

func (p *peerState) GetConnectionTo(node uint8) (net.IP, net.IP, *net.IPNet) {
	me, peer, transferNet, err := transferNetwork.GetPeerToPeerNet(node, p.uid, p._parent.subnet)
	panicOn(err)
	return me, peer, transferNet
}

func (p *peerState) GetLinkIPs() (string, string) {
	me, peer, _ := p.GetConnectionTo(p._parent.selfIDX)
	return me.String(), peer.String()
}

func (p *peerState) getWGName() string {
	return transferNetwork.BuildWireguardInterfaceName(p._parent.selfIDX, p.uid)
}

func (p *peerState) GetKeepalive() *time.Duration {
	keepAlive := 0 * time.Second
	if p.keepalive {
		keepAlive = 10 * time.Second
	}
	return &keepAlive
}

func (p *peerState) GetPorts() (int, int) {
	return transferNetwork.CalculatePorts(p._parent.startPort, p._parent.selfIDX, p.uid)
}

func (p *peerState) resolveIP() *net.UDPAddr {
	// if no hostname was provided we obviously can't parse it
	if p.hostname == "" {
		return nil
	}

	// else we build the full address from host and port and (try to) resolve it
	_, peerPort := p.GetPorts()
	ip, err := various.ResolveHostname(p.hostname, peerPort)
	panicOn(err)

	return ip
}

type state struct {
	privateKey     *wgtypes.Key
	name           string
	startPort      uint16
	selfIDX        uint8
	peers          map[uint8]*peerState
	subnet         *net.IPNet
	otherInterface []*ManualInterface
}

type ManualInterface struct {
	interfaceName string
	ip            *net.IP
	ospfPassive   bool
}
