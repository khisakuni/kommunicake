package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/khisakuni/kommunicake/db"
)

// App has a DB, Router, and Port
type App struct {
	DB     *db.DB
	Router *mux.Router
	Port   string
}

// Option is a func to mutate the app
type Option func(a *App)

// WithPort sets the port. Default is ":3000"
func WithPort(port string) Option {
	return func(a *App) {
		a.Port = port
	}
}

// NewApp initializes an App with a DB, Router, and Port
func NewApp(d *db.DB, r *mux.Router, options ...Option) *App {
	a := &App{d, r, ":3000"}
	for _, option := range options {
		option(a)
	}
	return a
}

// Run starts the app on given port. Default port is ":3000"
func (a *App) Run() {
	log.Printf("Listening on port: %s\n", a.Port)
	log.Fatal(http.ListenAndServe(a.Port, a.Router))
}
