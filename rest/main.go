package rest

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"magic-wan/rest/cluster"
	"net/http"
	"time"
)

func StartRest() error {
	http.HandleFunc("/api/v1/debug", DebugV1Handler)
	http.HandleFunc("/api/v1/cluster/vote", cluster.VoteV1Handler)

	port := 80
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: 3 * time.Second,
	}
	log.Info(fmt.Sprintf("Starting server on :%d...\n", port))
	err := server.ListenAndServe()
	return err
}
