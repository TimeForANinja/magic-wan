package cluster

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
)

type voteMessage struct {
	Vote string `json:"vote"`
}

func getMasterToVote() *votingPeer {
	// If we already have a master it's best to stick with it
	if HasMaster() {
		return GetMaster()
	}

	// Option 2 is favoring one which already has many votes
	candidates := masterCandidate(MaxVoteAge)
	if len(candidates) != 0 {
		randomIndex := rand.Intn(len(candidates))
		return candidates[randomIndex]
	}

	// Option 3 is a random one
	randomIndex := rand.Intn(len(globalPeers))
	return globalPeers[randomIndex]
}

func callVote() {
	masterToVote := getMasterToVote()
	if masterToVote == nil {
		return
	}

	voteMessage := voteMessage{Vote: masterToVote.ip}

	for _, peer := range globalPeers {
		go sendVote(peer, voteMessage)
	}
}

func sendVote(peer *votingPeer, voteMessage voteMessage) {
	url := "http://" + peer.ip + "/api/v1/vote"

	jsonData, err := json.Marshal(voteMessage)
	if err != nil {
		// TODO: Handle error
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		// TODO: Handle error
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// TODO: Handle error
		return
	}
	defer resp.Body.Close()
}
