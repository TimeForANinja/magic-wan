package cluster

import (
	"magic-wan/pkg/various"
	"time"
)

var globalPeers = make([]*votingPeer, 0)
var globalVotes = make([]vote, 0)
var globalMaster *votingPeer = nil

const AnnouncementInterval = 60 * time.Second
const MaxVoteAge = 5 * AnnouncementInterval

func InitCluster() {
	// Start the loop that calls callVote every 60 seconds
	go func() {
		ticker := time.NewTicker(AnnouncementInterval)
		defer ticker.Stop()
		for range ticker.C {
			callVote()
		}
	}()
}

func AddPeer(ip string) {
	globalPeers = append(globalPeers, &votingPeer{
		ip: ip,
	})
}

func RemovePeer(ip string) {
	globalPeers = various.ArrayFilter(globalPeers, func(p *votingPeer) bool {
		return p.ip != ip
	})
}

func HasMaster() bool {
	return globalMaster != nil
}

func GetMaster() *votingPeer {
	return globalMaster
}
