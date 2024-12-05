package rest

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"magic-wan/pkg/cluster"
	"magic-wan/rest/api"
	"magic-wan/rest/api/api-login"
	"magic-wan/rest/gui"
	"magic-wan/rest/shared"
	"net/http"
	"time"
)

const (
	AuthRequired = iota
	NoAuthRequired
	AuthIrrelevant
)

func padWithLogin(checkAuth int, redirectURL string, handler func(http.ResponseWriter, *http.Request, *shared.User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := shared.JWTManagerInstance.ParseFromRequest(r)

		if checkAuth == AuthRequired && user == nil {
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		} else if checkAuth == NoAuthRequired && user != nil {
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}

		handler(w, r, user)
	}
}

func StartRest(cluster *cluster.Cluster, errorChannel chan error) {
	// API Routes
	http.HandleFunc("/api/v1/debug", padWithLogin(AuthRequired, "/login", api.DebugV1Handler))
	clusterVoteV1Handler := api.ClusterVoteV1HandlerFactory(cluster)
	http.HandleFunc("/api/v1/cluster/vote", clusterVoteV1Handler)
	http.HandleFunc("/api/v1/wgkey", api.WireguardKeyGenV1Handler)
	http.HandleFunc("/api/v1/doLogin", padWithLogin(NoAuthRequired, "/", api_login.DoLoginV1Handler))
	http.HandleFunc("/api/v1/checkLogin", padWithLogin(AuthIrrelevant, "", api_login.CheckLoginV1Handler))

	// GUI
	http.HandleFunc("/login", padWithLogin(NoAuthRequired, "/", gui.LoginHandler))
	http.HandleFunc("/", padWithLogin(AuthRequired, "/login", gui.HomeHandler))
	http.HandleFunc("/logout", padWithLogin(AuthRequired, "/login", gui.LogoutHandler))

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
