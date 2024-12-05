package api_login

import (
	"encoding/json"
	"fmt"
	"magic-wan/rest/shared"
	"net/http"
)

const DEFAULT_PW = "password"

func DoLoginV1Handler(w http.ResponseWriter, r *http.Request, _ *shared.User) {
	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		fmt.Println(err)
		shared.SendResponse(w, http.StatusBadRequest, newFailedResponse("Invalid Request Body"))
		return
	}

	// check credentials
	if creds.Username != "root" || creds.Password != DEFAULT_PW {
		shared.SendResponse(w, http.StatusUnauthorized, newFailedResponse("Invalid Credentials"))
		return
	}

	// build user and token
	user := &shared.User{Name: creds.Username}
	token, err := shared.JWTManagerInstance.CreateToken(user)
	if err != nil {
		fmt.Println(err)
		shared.SendResponse(w, http.StatusInternalServerError, newFailedResponse("Failed to generate Token"))
		return
	}

	// set token
	http.SetCookie(w, &http.Cookie{
		Name:   "Authorization",
		Value:  token,
		Path:   "/",
		MaxAge: int(shared.TokenExpirationDuration.Seconds()),
	})

	shared.SendResponse(w, http.StatusOK, newSuccessResponse(user))
}
