package helpers

import (
	"net/http"
	"encoding/json"
)

// ErrorResponse formats error and writes to ResponseWriter
func ErrorResponse(w http.ResponseWriter, err error, status int) {
	type errorMessage struct {
		Message string `json:"message"`
	}
	js, err := json.Marshal(errorMessage{Message: err.Error()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Error(w, string(js), status)
}