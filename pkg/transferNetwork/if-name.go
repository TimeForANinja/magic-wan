package transferNetwork

import "fmt"

func BuildWireguardInterfaceName(selfIDX uint8, peerIDX uint8) string {
	// TODO: define a "max" idx and add some check for it
	// if we just keep uint8 we should change is to 3d
	return fmt.Sprintf("wg1%02d%02d", selfIDX, peerIDX)
}
