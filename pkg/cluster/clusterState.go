package cluster

import (
	log "github.com/sirupsen/logrus"
	"magic-wan/pkg/various"
	"time"
)

type votingPeer struct {
	ip string
}

type vote struct {
	voter *votingPeer
	vote  *votingPeer
	time  time.Time
}

func (c *Cluster) updateVotes(newVote *vote) {
	log.WithFields(log.Fields{
		"vote":  newVote,
		"votes": c.votes,
	}).Debugf("Updating Votes")

	// remove any old votes
	votes := various.ArrayFilter(c.votes, func(oldVote *vote) bool {
		return newVote.voter.ip != oldVote.voter.ip
	})

	// merge the new vote with the filtered votes, and store in globalvotes
	c.votes = append(votes, newVote)

	c.updateMaster()
}

func (c *Cluster) updateMaster() {
	newMaster := c.calcMaster(c.peers, MaxVoteAge)
	if c.master == newMaster {
		// no change, so no action required
		return
	}

	log.WithFields(log.Fields{
		"new": newMaster,
		"old": c.master,
	}).Infof("New Cluster Master")

	c.master = newMaster
}

func (c *Cluster) calcMaster(knownPeers []*votingPeer, maxAge time.Duration) *votingPeer {
	// Filter out votes older than maxAge
	validVotes := various.ArrayFilter(c.votes, func(v *vote) bool {
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

func (c *Cluster) masterCandidate(maxAge time.Duration) []*votingPeer {
	// Filter out votes older than maxAge
	validVotes := various.ArrayFilter(c.votes, func(v *vote) bool {
		return time.Since(v.time) <= maxAge
	})

	candidates := calcTopCandidates(validVotes)

	return candidates
}

func calcTopCandidates(votes []*vote) []*votingPeer {
	// Count votes per peer
	voteCounts := make(map[*votingPeer]int)
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
