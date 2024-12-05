package api_login

import (
	"magic-wan/rest/shared"
	"net/http"
)

func CheckLoginV1Handler(w http.ResponseWriter, r *http.Request, user *shared.User) {
	if user == nil {
		shared.SendResponse(w, http.StatusUnauthorized, newFailedResponse("Not logged in"))
	} else {
		// Respond with the user data on successful parsing
		shared.SendResponse(w, http.StatusOK, newSuccessResponse(user))
	}
}
