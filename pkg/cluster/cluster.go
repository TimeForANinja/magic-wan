package cluster

import (
	log "github.com/sirupsen/logrus"
	"magic-wan/pkg/cluster/shared"
	"magic-wan/pkg/various"
	"time"
)

const announcementInterval = 60 * time.Second
const maxVoteAge = 5 * announcementInterval
const maxFailedSends = 2

type CoreConfig struct {
	version uint16
}

type Cluster[T CoreConfig] struct {
	self   *votingPeer
	peers  []*votingPeer
	votes  []*peerVote
	master *votingPeer
	config T
}

func NewCluster[T CoreConfig](selfIP string) *Cluster[T] {
	selfPeer := &votingPeer{ip: selfIP}
	c := &Cluster[T]{
		self:  selfPeer,
		peers: []*votingPeer{selfPeer},
		votes: make([]*peerVote, 0),
	}

	return c
}

func (c *Cluster[T]) StartAnnouncements() {
	// Start the loop that calls callVote every 60 seconds
	ticker := time.NewTicker(announcementInterval)
	defer ticker.Stop()
	for range ticker.C {
		// TODO: use actual config Version
		c.doVoting(255)
	}
}

func (c *Cluster[T]) AddPeer(ip string) {
	c.peers = append(c.peers, &votingPeer{
		ip: ip,
	})
}

func (c *Cluster[T]) RemovePeer(ip string) {
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

func (c *Cluster[T]) HasMaster() bool {
	return c.master != nil
}

func (c *Cluster[T]) getMaster() *votingPeer {
	return c.master
}

func (c *Cluster[T]) isMaster(ip string) bool {
	return c.master != nil && c.master.ip == ip
}

func (c *Cluster[T]) doVoting(cfgVersion uint16) {
	masterToVote := c.getMasterToVote()
	if masterToVote == nil {
		log.Info("skipping sending cluster vote, no master found")
		return
	}

	voteMessage := shared.VoteMessage{
		Voter:         c.self.ip,
		Vote:          masterToVote.ip,
		ConfigVersion: cfgVersion,
	}

	log.WithFields(log.Fields{
		"peers": c.peers,
		"vote":  voteMessage,
	}).Infof("Preparing to send cluster vote")

	for _, peer := range c.peers {
		go peer.sendVote(voteMessage)
	}
}

func (c *Cluster[T]) checkForNewMaster() {
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
