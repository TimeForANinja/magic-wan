package shared

import (
	"encoding/json"
	"net/http"
)

// SendResponse helps to send JSON responses
func SendResponse(w http.ResponseWriter, statusCode int, response any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
