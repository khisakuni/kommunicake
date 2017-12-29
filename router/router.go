package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router is wrapper for mux router
type Router struct {
	m *mux.Router
}

// NewRouter returns new router
func NewRouter() *Router {
	return &Router{m: mux.NewRouter()}
}

// Serve starts the server
func (r *Router) Serve(port string) error {
	return http.ListenAndServe(port, r.m)
}

// AddRoute hooks up a route to the router
func (r *Router) AddRoute(path string, handler http.HandlerFunc) {
	r.m.HandleFunc(path, handler)
}
