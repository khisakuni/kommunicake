package middleware

import (
	"github.com/khisakuni/kommunicake/api/helpers"
	"fmt"
	"net/http"
	"os"
	"strings"
	jwt "github.com/dgrijalva/jwt-go"
)

func What() {
	fmt.Println("what")
}

func WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			helpers.ErrorResponse(w, err, http.StatusUnauthorized)
			return
		}
		if token.Valid {
			// TODO: get email or id and get user, then put user into context
			fmt.Printf("TOKEN: %v\n", token)
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			helpers.ErrorResponse(w, ve, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}