package v1

import (
	"net/http"

	"github.com/khisakuni/kommunicake/api/helpers"
	"github.com/khisakuni/kommunicake/api/middleware"
	"github.com/khisakuni/kommunicake/database"
	"github.com/khisakuni/kommunicake/models"
)

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Token string `json:"token"`
	}
	if ok := decodeJSON(w, r.Body, &params); !ok {
		return
	}

	db := middleware.GetDBFromContext(r.Context())

	newToken, err := generateNewToken(params.Token, db)
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	user, err := newToken.GetUser(db)
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	res := struct {
		User  *models.User `json:"user"`
		Token string       `json:"token"`
	}{User: user, Token: newToken.Value}
	jsonResponse(w, res, http.StatusOK)
}

func generateNewToken(tokenString string, db *database.DB) (*models.Token, error) {
	token, err := models.FromTokenString(tokenString, db)
	if err != nil {
		return nil, err
	}

	user, err := token.GetUser(db)
	if err != nil {
		return nil, err
	}

	newToken, err := models.GenerateToken(db, user)
	if err != nil {
		return nil, err
	}

	return newToken, nil
}
