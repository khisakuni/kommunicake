package middleware

import (
	"context"
	"net/http"

	"github.com/khisakuni/kommunicake/database"
)

const dbKey contextKey = "database key"

// WithDB adds ref to DB to context
func WithDB(db *database.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), dbKey, db)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetDBFromContext retrieves ref to DB from context
func GetDBFromContext(ctx context.Context) *database.DB {
	return ctx.Value(dbKey).(*database.DB)
}
