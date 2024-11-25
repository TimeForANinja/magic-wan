package cluster

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"magic-wan/pkg/various"
	"net/http"
	"time"
)

func VoteV1Handler(w http.ResponseWriter, r *http.Request) {
	// Parse the vote from the body
	var msg voteMessage
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Debugf("Invalid request body: %w", err)
		return
	}

	err = parseVote(r.RemoteAddr, msg.Vote)
	if err != nil {
		http.Error(w, "Invalid vote", http.StatusBadRequest)
		log.Debugf("Invalid vote: %w", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{ status: 200, message: \"vote counted\" }")))
}

func parseVote(from, to string) error {
	fromPeer := various.ArrayFind(globalPeers, func(p *VotingPeer) bool {
		return p.ip == from
	})
	toPeer := various.ArrayFind(globalPeers, func(p *VotingPeer) bool {
		return p.ip == to
	})
	if fromPeer == nil || toPeer == nil {
		return fmt.Errorf("unknown peer(s)")
	}

	UpdateVotes(vote{
		voter: fromPeer,
		vote:  toPeer,
		time:  time.Now(),
	})

	return nil
}
