package rest

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func StartRest() error {
	http.HandleFunc("/api/v1/debug", DebugV1Handler)

	port := 80
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: 3 * time.Second,
	}
	log.Info(fmt.Sprintf("Starting server on :%d...\n", port))
	err := server.ListenAndServe()
	return err
}
