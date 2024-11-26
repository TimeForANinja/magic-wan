package cluster

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"magic-wan/pkg/various"
	"time"
)

type peerVote struct {
	voter *votingPeer
	vote  *votingPeer
	time  time.Time
}

type VoteMessage struct {
	Voter string `json:"voter"`
	Vote  string `json:"vote"`
}

func (c *Cluster) OnVoteReceived(vm *VoteMessage) error {
	vote, err := vm.parse(c)
	if err != nil {
		return err
	}

	c.updateVotes(vote)
	return nil
}

func (vm *VoteMessage) parse(c *Cluster) (*peerVote, error) {
	fromPeer := various.ArrayFind(c.peers, func(p *votingPeer) bool {
		return p.ip == vm.Voter
	})
	toPeer := various.ArrayFind(c.peers, func(p *votingPeer) bool {
		return p.ip == vm.Vote
	})
	log.WithFields(log.Fields{
		"from": vm.Voter, "to": vm.Vote,
		"fromLookup": fromPeer, "toLookup": toPeer,
	}).Debugf("Parsing Vote")
	if fromPeer == nil || toPeer == nil {
		return nil, fmt.Errorf("unknown peer(s)")
	}

	return &peerVote{
		voter: fromPeer,
		vote:  toPeer,
		time:  time.Now(),
	}, nil
}

func (c *Cluster) updateVotes(newVote *peerVote) {
	log.WithFields(log.Fields{
		"vote":  newVote,
		"votes": c.votes,
	}).Debugf("Updating Votes")

	// remove any old votes
	votes := various.ArrayFilter(c.votes, func(oldVote *peerVote) bool {
		return newVote.voter.ip != oldVote.voter.ip
	})

	// merge the new vote with the filtered votes
	// and store them
	c.votes = append(votes, newVote)

	c.checkForNewMaster()
}
