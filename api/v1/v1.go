package v1

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	helpers "github.com/khisakuni/kommunicake/api/helpers"
	middleware "github.com/khisakuni/kommunicake/api/middleware"
	"github.com/khisakuni/kommunicake/database"
)

// Routes attaches all the v1 routes to router
func Routes(router *mux.Router, db *database.DB) {
	// Auth
	router.HandleFunc("/api/v1/register", RegisterUser).Methods("POST")
	router.HandleFunc("/api/v1/login", Login).Methods("POST")
	router.HandleFunc("/api/v1/token", middleware.WithAuth(http.HandlerFunc(RefreshToken), db).ServeHTTP).Methods("POST")

	// Gmail
	router.HandleFunc("/api/v1/gmail_login", middleware.WithAuth(http.HandlerFunc(GmailLoginURL), db).ServeHTTP).Methods("POST")
	router.HandleFunc("/api/v1/webhooks/gmail", GmailWebhook).Methods("GET")
	router.HandleFunc("/api/v1/webhooks/gmail/message", GmailWebhookNewMessage).Methods("POST")
	router.HandleFunc("/api/v1/gmail/subscribe", middleware.WithAuth(http.HandlerFunc(GmailSubscribeToNewMessage), db).ServeHTTP).Methods("POST")

	// Slack
	router.HandleFunc("/api/v1/slack_login", middleware.WithAuth(http.HandlerFunc(SlackLoginURL), db).ServeHTTP).Methods("POST")
	router.HandleFunc("/api/v1/webhooks/slack", SlackWebhook).Methods("GET")

	// Message Providers
	router.HandleFunc("/api/v1/message_providers", middleware.WithAuth(http.HandlerFunc(MessageProvidersIndex), db).ServeHTTP).Methods("GET")
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
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		ok = false
	}
	return ok
}

// validateParams runs validation function
func validateParams(w http.ResponseWriter, validate func() error) bool {
	ok := true
	err := validate()
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusBadRequest)
		ok = false
	}
	return ok
}
