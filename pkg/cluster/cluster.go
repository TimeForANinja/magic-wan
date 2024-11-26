package cluster

import (
	log "github.com/sirupsen/logrus"
	"magic-wan/pkg/various"
	"time"
)

const announcementInterval = 60 * time.Second
const maxVoteAge = 5 * announcementInterval
const maxFailedSends = 2

type Cluster struct {
	self   *votingPeer
	peers  []*votingPeer
	votes  []*peerVote
	master *votingPeer
}

func NewCluster(selfIP string) *Cluster {
	selfPeer := &votingPeer{ip: selfIP}
	c := &Cluster{
		self:  selfPeer,
		peers: []*votingPeer{selfPeer},
		votes: make([]*peerVote, 0),
	}

	return c
}

func (c *Cluster) StartAnnouncements() {
	// Start the loop that calls callVote every 60 seconds
	ticker := time.NewTicker(announcementInterval)
	defer ticker.Stop()
	for range ticker.C {
		c.doVoting()
	}
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
	c.votes = various.ArrayFilter(c.votes, func(v *peerVote) bool {
		return v.voter.ip != ip
	})
	if c.isMaster(ip) {
		c.master = nil
	}
}

func (c *Cluster) HasMaster() bool {
	return c.master != nil
}

func (c *Cluster) getMaster() *votingPeer {
	return c.master
}

func (c *Cluster) isMaster(ip string) bool {
	return c.master != nil && c.master.ip == ip
}

func (c *Cluster) doVoting() {
	masterToVote := c.getMasterToVote()
	if masterToVote == nil {
		log.Info("skipping sending cluster vote, no master found")
		return
	}

	voteMessage := VoteMessage{
		Voter: c.self.ip,
		Vote:  masterToVote.ip,
	}

	log.WithFields(log.Fields{
		"peers": c.peers,
		"vote":  voteMessage,
	}).Infof("Preparing to send cluster vote")

	for _, peer := range c.peers {
		go peer.sendVote(voteMessage)
	}
}

func (c *Cluster) checkForNewMaster() {
	newMaster := c.calcMaster(c.peers)
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
