package transferNetwork

import (
	"testing"
)

func Test_CalculatePorts(t *testing.T) {
	usage := make(map[uint8]map[int]int)
	startPort := uint16(1000)
	totalNodes := uint8(128)

	for i := uint8(0); i < totalNodes; i++ {
		usage[i] = make(map[int]int) // Initialize map for each node
	}

	for i := uint8(0); i < totalNodes; i++ {

		for j := uint8(0); j < totalNodes; j++ {
			myPort, peerPort := CalculatePorts(startPort, i, j)
			usage[i][myPort]++
			usage[j][peerPort]++
		}
	}

	// Check that each port is used exactly twice for each node
	// once as "source" and once as "destination"
	for node, ports := range usage {
		for port, count := range ports {
			if count != 2 {
				t.Errorf("Port %d on node %d is used %d times; expected 2 times", port, node, count)
			}
		}
	}
}
