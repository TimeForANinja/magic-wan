package cluster

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type votingPeer struct {
	ip          string
	failedSends int
}

func (peer *votingPeer) sendVote(voteMessage VoteMessage) {
	url := fmt.Sprintf("http://%s/api/v1/cluster/vote", peer.ip)

	jsonData, err := json.Marshal(voteMessage)
	if err != nil {
		log.Errorf("Failed to marshal vote message: %v", err)
		// IMPROVEMENT: Handle error
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Debugf("Failed to create request: %v", err)
		peer.updateStale(false)
		// IMPROVEMENT: Handle error
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Debugf("Failed to send vote: %v", err)
		peer.updateStale(false)
		// IMPROVEMENT: Handle error
		return
	}
	if resp.StatusCode != 200 {
		log.Debugf("Failed to send, got response code %s", resp.Status)
		peer.updateStale(false)
		// IMPROVEMENT: Handle error
		return
	}
	defer resp.Body.Close()
	log.Debugf("Vote sent to %s", peer.ip)
	peer.updateStale(true)
}

func (peer *votingPeer) isStale() bool {
	return peer.failedSends > maxFailedSends
}

func (peer *votingPeer) updateStale(successful bool) {
	if successful {
		peer.failedSends = 0
	} else {
		peer.failedSends++
	}
}
