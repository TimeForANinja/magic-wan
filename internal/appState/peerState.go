package configState

import (
	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"magic-wan/pkg/transferNetwork"
	"magic-wan/pkg/various"
	"magic-wan/pkg/wg"
	"net"
	"time"
)

type peerState struct {
	publicKey *wgtypes.Key
	_parent   *ApplicationState
	hostname  string
	uid       uint8
	keepalive bool
	ip        *net.UDPAddr
}

func (p *peerState) BuildWGConfig() (wgtypes.Config, error) {
	selfPort, _ := p.GetPorts()
	linkNetwork, err := p.GetLinkNetwork()
	if err != nil {
		return wgtypes.Config{}, err
	}
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
					*linkNetwork,
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
	}, nil
}

func (p *peerState) GetLinkNetwork() (*net.IPNet, error) {
	_, _, nw, err := p.GetConnectionTo(p._parent.selfIDX)
	return nw, err
}

func (p *peerState) GetConnectionTo(node uint8) (net.IP, net.IP, *net.IPNet, error) {
	me, peer, transferNet, err := transferNetwork.GetPeerToPeerNet(node, p.uid, p._parent.subnet)
	return me, peer, transferNet, err
}

func (p *peerState) GetLinkIPs() (string, string, error) {
	me, peer, _, err := p.GetConnectionTo(p._parent.selfIDX)
	return me.String(), peer.String(), err
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

func (p *peerState) ResolveAddr() (*net.UDPAddr, error) {
	// if no hostname was provided we obviously can't parse it
	if p.hostname == "" {
		return nil, nil
	}

	// else we build the full address from host and port and (try to) resolve it
	_, peerPort := p.GetPorts()
	ip, err := various.ResolveHostname(p.hostname, peerPort)
	return ip, err
}

func (p *peerState) CachedAddr() *net.UDPAddr {
	return p.ip
}

/*
func (p *peerState) Remove() {
	// IMPROVEMENT: implement

	// old implementation:
	err := wg.DisableDevice(client, peer.getWGName())
	panicOn(err)
}
*/

func (p *peerState) PushToWireguard() error {
	ifcName := p.getWGName()

	// check if we already have this device
	devices, err := wg.GetDevices(p._parent.wgClient)
	if err != nil {
		return err
	}
	exists := various.ArrayFind(devices, func(dev *wgtypes.Device) bool { return dev.Name == ifcName }) != nil
	if !exists {
		err = wg.CreateNewDevice(ifcName)
		if err != nil {
			return err
		}
	}

	var ifConf wgtypes.Config
	ifConf, err = p.BuildWGConfig()
	if err != nil {
		return err
	}
	err = wg.UpdateDevice(p._parent.wgClient, ifcName, ifConf)
	if err != nil {
		return err
	}

	// set interface IPs
	var selfIP string
	selfIP, _, err = p.GetLinkIPs()
	if err != nil {
		return err
	}
	err = wg.BaseConfigureInterface(ifcName, selfIP)

	return err
}

func (p *peerState) NotifyIPChange(newIP *net.UDPAddr) error {
	log.WithFields(log.Fields{
		"peer":  p,
		"oldIP": p.ip,
		"newIP": newIP,
	}).Info("onPeerRemoved")

	p.ip = newIP

	// update relevant running config
	conf, err := p.BuildWGConfig()
	if err != nil {
		return err
	}
	err = wg.UpdateDevice(p._parent.wgClient, p.getWGName(), conf)
	return err
}
