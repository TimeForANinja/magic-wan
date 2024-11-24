package transferNetwork

func CalculatePorts(startPort uint16, selfIDX uint8, peerIDX uint8) (int, int) {
	// since the start port is shared, we can ensure uniqueness using the peer indexes

	// since we need one listener per remote peer connecting we offset our port by that
	myPort := int(startPort + uint16(peerIDX))
	// for the peer port we do it the other way around
	peerPort := int(startPort + uint16(selfIDX))

	return myPort, peerPort
}
