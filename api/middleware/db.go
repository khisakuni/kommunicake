package middleware

import (
	"context"
	"net/http"

	"github.com/khisakuni/kommunicake/database"
)

const dbKey contextKey = "database key"

// WithDB adds ref to db to context
func WithDB(next http.Handler, db *database.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := newContextWithDB(r.Context(), db, r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetDBFromContext retrieves ref to DB from context
func GetDBFromContext(ctx context.Context) *database.DB {
	return ctx.Value(dbKey).(*database.DB)
}

func newContextWithDB(ctx context.Context, db *database.DB, r *http.Request) context.Context {
	return context.WithValue(ctx, dbKey, db)
}
