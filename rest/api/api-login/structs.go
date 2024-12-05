package api_login

import "magic-wan/rest/shared"

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool         `json:"success"`
	Error   *string      `json:"error,omitempty"`
	User    *shared.User `json:"user,omitempty"`
}

func newFailedResponse(error string) *LoginResponse {
	return &LoginResponse{
		Success: false,
		Error:   &error,
	}
}

func newSuccessResponse(user *shared.User) *LoginResponse {
	return &LoginResponse{
		Success: true,
		User:    user,
	}
}
