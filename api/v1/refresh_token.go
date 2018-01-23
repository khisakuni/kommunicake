package v1

import (
	"net/http"
	"github.com/khisakuni/kommunicake/api/helpers"
	"github.com/khisakuni/kommunicake/api/middleware"
	"github.com/khisakuni/kommunicake/models"
)

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var params struct { Token string `json:"token"`}
	if ok := decodeJSON(w, r.Body, &params); !ok {
		return
	}

	db := middleware.GetDBFromContext(r.Context())
	token, err := models.FromTokenString(params.Token, db)
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	user, err := token.GetUser(db)
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	newToken, err := models.GenerateToken(db, user)
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	res := struct {
		User *models.User
		Token string
	}{User: user, Token: newToken.Value}
	jsonResponse(w, res, http.StatusOK)
}