package middleware

import (
	"github.com/khisakuni/kommunicake/api/helpers"
	"github.com/khisakuni/kommunicake/models"
	"github.com/khisakuni/kommunicake/database"
	"net/http"
	"strings"
	"context"
)

const authKey contextKey = "auth key"

func WithAuth(next http.Handler, db *database.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		token, err := models.FromTokenString(tokenString, db)
		if err != nil {
			helpers.ErrorResponse(w, err, http.StatusUnauthorized)
			return
		}

		user, err := token.GetUser(db)
		if err != nil {
			helpers.ErrorResponse(w, err, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), authKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) *models.User {
	return ctx.Value(authKey).(*models.User)
}