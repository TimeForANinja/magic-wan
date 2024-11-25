package cluster

import (
	"magic-wan/pkg/various"
	"time"
)

const AnnouncementInterval = 60 * time.Second
const MaxVoteAge = 5 * AnnouncementInterval

func InitCluster(selfIP string) *Cluster {
	selfPeer := &votingPeer{ip: selfIP}
	c := &Cluster{
		self:  selfPeer,
		peers: []*votingPeer{selfPeer},
		votes: make([]*vote, 0),
	}

	// Start the loop that calls callVote every 60 seconds
	go func() {
		ticker := time.NewTicker(AnnouncementInterval)
		defer ticker.Stop()
		for range ticker.C {
			c.doVoting()
		}
	}()

	return c
}

func (c *Cluster) AddPeer(ip string) {
	c.peers = append(c.peers, &votingPeer{
		ip: ip,
	})
}

func (c *Cluster) RemovePeer(ip string) {
	c.peers = various.ArrayFilter(c.peers, func(p *votingPeer) bool {
		return p.ip != ip
	})
}

func (c *Cluster) HasMaster() bool {
	return c.master != nil
}

func (c *Cluster) GetMaster() *votingPeer {
	return c.master
}
