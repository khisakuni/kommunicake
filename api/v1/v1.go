package v1

import (
	"github.com/gorilla/mux"
)

// Routes attaches all the v1 routes to router
func Routes(router *mux.Router) {
	router.HandleFunc("/api/v1/register", RegisterUser).Methods("POST")
}
