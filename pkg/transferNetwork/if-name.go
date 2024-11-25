package transferNetwork

import (
	"fmt"
	"strconv"
)

func BuildWireguardInterfaceName(selfIDX uint8, peerIDX uint8) string {
	// 3 digits since uint8 has a max value of 255
	return fmt.Sprintf("wg1%03d%03d", selfIDX, peerIDX)
}

func MatchesWireguardInterfaceName(name string) bool {
	if len(name) != 11 || name[:3] != "wg1" {
		return false
	}

	// split and parse idx
	selfIDXPart := name[3:6]
	peerIDXPart := name[6:9]

	selfIDX, err1 := strconv.Atoi(selfIDXPart)
	peerIDX, err2 := strconv.Atoi(peerIDXPart)
	if err1 != nil || err2 != nil {
		return false
	}

	// check that indexes are uint8
	if selfIDX < 0 || selfIDX > 255 || peerIDX < 0 || peerIDX > 255 {
		return false
	}

	return true
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
