package gui

import (
	"fmt"
	"magic-wan/rest/shared"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, _ *http.Request, user *shared.User) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(fmt.Sprintf("Home for %s", user.Name)))
}
