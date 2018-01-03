package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Routes attaches all the v1 routes to router
func Routes(router *mux.Router) {
	router.HandleFunc("/api/v1/register", RegisterUser).Methods("POST")
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

// errorResponse formats error and writes to ResponseWriter
func errorResponse(w http.ResponseWriter, err error, status int) {
	http.Error(w, err.Error(), status)
}
