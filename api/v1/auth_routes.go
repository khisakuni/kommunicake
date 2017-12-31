package v1

import (
	"encoding/json"
	"net/http"
)

type user struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var u user
	// db := middleware.GetDBFromContext(r.Context())
	// fmt.Printf("Got DB: %v\n", db)
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		panic(err)
	}
}
