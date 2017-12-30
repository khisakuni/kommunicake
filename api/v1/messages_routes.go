package v1

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetMessages(r *mux.Router) {
	r.HandleFunc("/api/v1/messages", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hi")
	})
}
