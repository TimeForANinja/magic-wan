package cluster

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"magic-wan/pkg/cluster/shared"
	"magic-wan/pkg/various"
	"time"
)

type peerVote struct {
	voter *votingPeer
	vote  *votingPeer
	time  time.Time
}

func (c *Cluster[T]) OnVoteReceived(vm *shared.VoteMessage) error {
	vote, err := c.parse(vm)
	if err != nil {
		return err
	}

	c.updateVotes(vote)
	return nil
}

func (c *Cluster[CoreConfig]) parse(vm *shared.VoteMessage) (*peerVote, error) {
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

func (c *Cluster[T]) updateVotes(newVote *peerVote) {
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
