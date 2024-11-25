package cluster

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
)

type VoteMessage struct {
	Voter string `json:"voter"`
	Vote  string `json:"vote"`
}

func (c *Cluster) getMasterToVote() *votingPeer {
	// If we already have a master it's best to stick with it
	if c.HasMaster() {
		return c.GetMaster()
	}

	// Option 2 is favoring one which already has many votes
	candidates := c.masterCandidate(MaxVoteAge)
	if len(candidates) != 0 {
		randomIndex := rand.Intn(len(candidates))
		return candidates[randomIndex]
	}

	// Option 3 is a random one
	randomIndex := rand.Intn(len(c.peers))
	return c.peers[randomIndex]
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

	// TODO: if request fails, we should repick
	for _, peer := range c.peers {
		go sendVote(peer, voteMessage)
	}
}

func sendVote(peer *votingPeer, voteMessage VoteMessage) {
	url := "http://" + peer.ip + "/api/v1/cluster/vote"

	jsonData, err := json.Marshal(voteMessage)
	if err != nil {
		log.Debugf("Failed to marshal vote message: %v", err)
		// IMPROVEMENT: Handle error
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Debugf("Failed to create request: %v", err)
		// IMPROVEMENT: Handle error
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Debugf("Failed to send vote: %v", err)
		// IMPROVEMENT: Handle error
		return
	}
	if resp.StatusCode != 200 {
		log.Debugf("Failed to send, got response code %s", resp.Status)
		// IMPROVEMENT: Handle error
		return
	}
	defer resp.Body.Close()
}
