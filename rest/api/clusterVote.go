package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"magic-wan/pkg/cluster/shared"
	"net/http"
)

func ClusterVoteV1HandlerFactory(cls shared.Cluster) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the vote from the body
		msg := &shared.VoteMessage{}
		err := json.NewDecoder(r.Body).Decode(msg)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			log.Debugf("Invalid request body: %v", err)
			return
		}

		// Update State with new Vote
		err = cls.OnVoteReceived(msg)
		if err != nil {
			http.Error(w, "Invalid vote", http.StatusBadRequest)
			log.Debugf("Invalid vote: %v", err)
			return
		}

		// Sent Reply
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("{ status: 200, message: \"vote counted\" }")))
	}
}
