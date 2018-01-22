package v1

import (
	"fmt"
	"net/http"
	"os"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/khisakuni/kommunicake/api/helpers"
)

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var params struct { RefreshToken string `json:"refreshToken"`}
	if ok := decodeJSON(w, r.Body, &params); !ok {
		return
	}

	token, err := jwt.Parse(params.RefreshToken, func(tokwn *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	if token.Valid {
		// issue new access & refresh token

	} else if ve, ok := err.(*jwt.ValidationError); ok {
		helpers.ErrorResponse(w, ve, http.StatusUnauthorized)
		return
	} else {
		helpers.ErrorResponse(w, fmt.Errorf("Invalid token"), http.StatusUnauthorized)
		return
	}
}