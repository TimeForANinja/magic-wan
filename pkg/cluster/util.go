package cluster

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"magic-wan/pkg/various"
	"time"
)

func (c *Cluster) parseVote(from, to string) (*vote, error) {
	fromPeer := various.ArrayFind(c.peers, func(p *votingPeer) bool {
		return p.ip == from
	})
	toPeer := various.ArrayFind(c.peers, func(p *votingPeer) bool {
		return p.ip == to
	})
	log.WithFields(log.Fields{
		"from": from, "to": to,
		"fromLookup": fromPeer, "toLookup": toPeer,
	}).Debugf("Parsing Vote")
	if fromPeer == nil || toPeer == nil {
		return nil, fmt.Errorf("unknown peer(s)")
	}

	return &vote{
		voter: fromPeer,
		vote:  toPeer,
		time:  time.Now(),
	}, nil
}
