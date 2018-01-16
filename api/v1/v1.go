package v1

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

// Routes attaches all the v1 routes to router
func Routes(router *mux.Router) {
	router.HandleFunc("/api/v1/register", RegisterUser).Methods("POST")
	router.HandleFunc("/api/v1/login", Login).Methods("POST")
	router.HandleFunc("/api/v1/gmail_login", GmailLoginURL).Methods("POST")
	router.HandleFunc("/api/v1/webhooks/gmail", GmailWebhook).Methods("GET")
	router.HandleFunc("/api/v1/gmail/token", GmailExchangeCode).Methods("POST")
}

// jsonResponse marshals struct and writes to ResponseWriter
func jsonResponse(w http.ResponseWriter, jsonable interface{}, status int) {
	js, err := json.Marshal(jsonable)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}

func decodeJSON(w http.ResponseWriter, body io.ReadCloser, p interface{}) bool {
	ok := true
	err := json.NewDecoder(body).Decode(p)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		ok = false
	}
	return ok
}

// validateParams runs validation function
func validateParams(w http.ResponseWriter, validate func() error) bool {
	ok := true
	err := validate()
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		ok = false
	}
	return ok
}

// errorResponse formats error and writes to ResponseWriter
func errorResponse(w http.ResponseWriter, err error, status int) {
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
