package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/khisakuni/kommunicake/api/helpers"
	"github.com/khisakuni/kommunicake/database"
	"github.com/khisakuni/kommunicake/models"
)

const authKey contextKey = "auth key"

func WithAuth(next http.Handler, db *database.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := models.FromTokenString(parseTokenString(r), db)
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

func parseTokenString(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}
