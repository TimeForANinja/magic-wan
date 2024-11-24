package transferNetwork

import (
	"fmt"
	"strconv"
)

func BuildWireguardInterfaceName(selfIDX uint8, peerIDX uint8) string {
	// 3 digits since uint8 has a max value of 255
	return fmt.Sprintf("wg1%03d%03d", selfIDX, peerIDX)
}

func SplitWireguardInterfaceName(name string) (uint8, uint8, error) {
	if len(name) != 11 || name[:3] != "wg1" {
		return 0, 0, fmt.Errorf("invalid interface name format")
	}

	selfIDXPart := name[3:6]
	peerIDXPart := name[6:9]

	selfIDX, err := strconv.Atoi(selfIDXPart)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing selfIDX: %w", err)
	}

	peerIDX, err := strconv.Atoi(peerIDXPart)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing peerIDX: %w", err)
	}

	return uint8(selfIDX), uint8(peerIDX), nil
}
