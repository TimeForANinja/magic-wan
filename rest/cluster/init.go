package cluster

import (
	"time"
)

var globalPeers = make([]*VotingPeer, 0)
var globalVotes = make([]vote, 0)
var globalMaster *VotingPeer = nil

const AnnouncementInterval = 60 * time.Second
const MaxVoteAge = 5 * AnnouncementInterval

func InitCluster(peers []string) {
	// populate list of peers
	for _, p := range peers {
		AddPeer(VotingPeer{ip: p})
	}

	// Start the loop that calls callVote every 60 seconds
	go func() {
		ticker := time.NewTicker(AnnouncementInterval)
		defer ticker.Stop()
		for range ticker.C {
			callVote()
		}
	}()
}

func HasMaster() bool {
	return globalMaster != nil
}

func GetMaster() *VotingPeer {
	return globalMaster
}
