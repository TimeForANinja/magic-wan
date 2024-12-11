package shared

import (
	"magic-wan/internal/appState"
	"net"
)

type ChangeNotifier interface {
	OnManualInterfaceAdd(iface *appState.ManualInterface, skipFRR bool) error
	OnWGInterfaceAdd(newPeer *appState.PeerState, skipFRR bool) error
	OnWGInterfaceRemove(oldPeer *appState.PeerState) error
	OnIPChange(peer *appState.PeerState, newIP *net.UDPAddr) error
}
