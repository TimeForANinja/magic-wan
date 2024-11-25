package configState

import (
	"fmt"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"magic-wan/pkg/cluster"
	"magic-wan/pkg/frr"
	"magic-wan/pkg/transferNetwork"
	"magic-wan/pkg/various"
	"net"
)

type ApplicationState struct {
	// required data
	privateKey *wgtypes.Key
	name       string
	startPort  uint16
	selfIDX    uint8
	subnet     *net.IPNet

	// required connectors
	wgClient *wgctrl.Client

	// lists, that can be empty
	peers          map[uint8]*peerState
	otherInterface []*ManualInterface
}

func NewState(
	wgClient *wgctrl.Client,
	name string,
	privateKey *wgtypes.Key,
	startPort uint16,
	selfIDX uint8,
	subnet *net.IPNet,
) *ApplicationState {
	return &ApplicationState{
		// required
		name:       name,
		privateKey: privateKey,
		startPort:  startPort,
		selfIDX:    selfIDX,
		subnet:     subnet,
		wgClient:   wgClient,

		// optional
		peers:          make(map[uint8]*peerState),
		otherInterface: make([]*ManualInterface, 0),
	}
}

const rootArea = "0.0.0.0"

func (s *ApplicationState) AddWireguardInterface(
	uid uint8,
	hostname string,
	publicKey *wgtypes.Key,
	keepalive bool,
) error {
	newPeer := peerState{
		publicKey: publicKey,
		hostname:  hostname,
		uid:       uid,
		keepalive: keepalive,
		// calculated later on
		_parent: nil,
		ip:      nil,
	}
	newPeer._parent = s

	var err error
	newPeer.ip, err = newPeer.ResolveAddr()
	if err != nil {
		return err
	}

	s.peers[newPeer.uid] = &newPeer
	err = newPeer.PushToWireguard()
	// TODO: update frr
	return err
}

func (s *ApplicationState) AddManualInterface(ifName string, ip *net.IP, passive bool) {
	s.otherInterface = append(s.otherInterface, &ManualInterface{
		interfaceName: ifName,
		ip:            ip,
		ospfPassive:   passive,
	})
}

func (s *ApplicationState) GetLoopbackAddress() (*net.IP, *net.IPNet, error) {
	selfLoopbackIP, _, peerNet, err := transferNetwork.GetPeerToPeerNet(s.selfIDX, 0, s.subnet)
	return &selfLoopbackIP, peerNet, err
}

func (s *ApplicationState) GetPeers() []*peerState {
	return various.MapValues(s.peers)
}

func (s *ApplicationState) DeriveFRRState() *frr.Config {
	frrConf := frr.NewConfig(s.name, fmt.Sprintf("0.0.0.%d", s.selfIDX))

	frrConf.AddLogging("syslog", "informational")
	// IMPROVMENT: remove debug log to file
	frrConf.AddLogging("file /var/log/frr/debug.log", "debugging")

	for _, peer := range s.peers {
		frrConf.AddInterface(peer.getWGName(), rootArea, false)
	}
	for _, iface := range s.otherInterface {
		frrConf.AddInterface(iface.interfaceName, rootArea, iface.ospfPassive)
	}

	frrConf.Router.AddNetwork(s.subnet.String(), rootArea)
	frrConf.Router.AddArea(s.subnet.String(), rootArea)

	return frrConf
}

func (s *ApplicationState) DeriveCluster() (*cluster.Cluster, error) {
	myIP, _, err := s.GetLoopbackAddress()
	if err != nil {
		return nil, err
	}
	return cluster.InitCluster(myIP.String()), nil
}
