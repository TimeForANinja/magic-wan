package rest

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"magic-wan/pkg/wg"
	"net/http"
)

func WGKeyGenV1Handler(w http.ResponseWriter, r *http.Request) {
	privkey, pubkey, err := wg.GenerateKeyPair()
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Debugf("Invalid request body: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("privkey: %s\npubkey: %s", privkey.String(), pubkey.String())))
}
