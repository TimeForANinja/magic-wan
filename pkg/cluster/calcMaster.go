package cluster

import (
	"magic-wan/pkg/various"
	"math/rand"
	"time"
)

func (c *Cluster) getMasterToVote() *votingPeer {
	// If we already have a master it's best to stick with it
	master := c.getMaster()
	if master != nil && !master.isStale() {
		return master
	}

	// Option 2 is favoring one which already has many votes
	candidates := c.masterCandidate()
	if len(candidates) != 0 {
		randomIndex := rand.Intn(len(candidates))
		return candidates[randomIndex]
	}

	// Option 3 is a random one
	randomIndex := rand.Intn(len(c.peers))
	return c.peers[randomIndex]
}

func (c *Cluster) getValidVotes() []*peerVote {
	// Filter out votes older than maxAge
	validVotes := various.ArrayFilter(c.votes, func(v *peerVote) bool {
		return time.Since(v.time) <= maxVoteAge
	})

	// Filter out votes with stale peers
	return various.ArrayFilter(validVotes, func(v *peerVote) bool {
		return !(v.vote.isStale() || v.voter.isStale())
	})
}

func (c *Cluster) calcMaster(knownPeers []*votingPeer) *votingPeer {
	validVotes := c.getValidVotes()

	// Ensure more than 50% of knownPeers have voted
	if len(validVotes) <= len(knownPeers)/2 {
		return nil
	}

	// Get a list of all Candidates that have the most votes
	candidates := c.calcTopCandidates(validVotes)
	// We need exactly one peer to have the most votes, to be considered master
	if len(candidates) != 1 {
		return nil
	}

	return candidates[0]
}

func (c *Cluster) masterCandidate() []*votingPeer {
	validVotes := c.getValidVotes()

	candidates := c.calcTopCandidates(validVotes)

	return candidates
}

func (c *Cluster) calcTopCandidates(votes []*peerVote) []*votingPeer {
	voteCounts := make(map[*votingPeer]int)

	// Explicit initialise with "0" for all Peers
	for _, peer := range c.peers {
		voteCounts[peer] = 0
	}

	// Count votes per peer
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
