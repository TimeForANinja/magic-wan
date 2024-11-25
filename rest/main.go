package rest

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"magic-wan/pkg/cluster"
	"net/http"
	"time"
)

func StartRest(cluster *cluster.Cluster, errorChannel chan error) {
	http.HandleFunc("/api/v1/debug", debugV1Handler)
	clusterVoteV1Handler := clusterVoteV1Handler_Factory(cluster)
	http.HandleFunc("/api/v1/cluster/vote", clusterVoteV1Handler)
	http.HandleFunc("/api/v1/wgkey", wireguardKeyGenV1Handler)

	port := 80
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: 3 * time.Second,
	}
	log.Info(fmt.Sprintf("Starting server on :%d...\n", port))

	err := server.ListenAndServe()
	// Since StartRest is called as a go routine,
	// there's no point in "returning" the error
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server failed to start: %v", err)
		errorChannel <- err
	}
}
