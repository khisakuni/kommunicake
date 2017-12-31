package v1

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetMessages(r *mux.Router) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hi")
	}
	r.HandleFunc("/api/v1/messages", handler).Methods("GET")
}
