package v1

import "github.com/gorilla/mux"

func CreateV1Routes(r *mux.Router) {
	GetMessages(r)
}
