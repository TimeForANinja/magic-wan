package cluster

import (
	"magic-wan/pkg/various"
	"time"
)

type VotingPeer struct {
	ip string
}

type vote struct {
	voter *VotingPeer
	vote  *VotingPeer
	time  time.Time
}

func UpdateVotes(v vote) {
	// remove any old votes
	votes := various.ArrayFilter(globalVotes, func(v vote) bool {
		return v.voter.ip != v.vote.ip
	})

	// Append the new vote.
	votes = append(votes, v)

	updateMaster()
}

func updateMaster() {
	globalMaster = calcMaster(globalPeers, MaxVoteAge)
}

func AddPeer(p VotingPeer) {
	globalPeers = append(globalPeers, &p)
}

func calcMaster(knownPeers []*VotingPeer, maxAge time.Duration) *VotingPeer {
	// Filter out votes older than maxAge
	validVotes := various.ArrayFilter(globalVotes, func(v vote) bool {
		return time.Since(v.time) <= maxAge
	})

	// Ensure more than 50% of knownPeers have voted
	if len(validVotes) <= len(knownPeers)/2 {
		return nil
	}

	candidates := calcTopCandidates(validVotes)
	if len(candidates) != 1 {
		return nil
	}

	return candidates[0]
}

func masterCandidate(maxAge time.Duration) []*VotingPeer {
	// Filter out votes older than maxAge
	validVotes := various.ArrayFilter(globalVotes, func(v vote) bool {
		return time.Since(v.time) <= maxAge
	})

	candidates := calcTopCandidates(validVotes)

	return candidates
}

func calcTopCandidates(votes []vote) []*VotingPeer {
	// Count votes per peer
	voteCounts := make(map[*VotingPeer]int)
	for _, v := range votes {
		voteCounts[v.vote]++
	}

	// Find the max votes for a peer
	maxVotes := 0
	for _, count := range voteCounts {
		if count > maxVotes {
			maxVotes = count
		}
	}

	peersWithMax := various.MapFilter(voteCounts, func(count int) bool {
		return count == maxVotes
	})

	return peersWithMax
}
