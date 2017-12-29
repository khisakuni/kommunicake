package app

import (
	"log"

	"github.com/khisakuni/kommunicake/db"
	"github.com/khisakuni/kommunicake/router"
)

// App has a DB, Router, and Port
type App struct {
	DB     *db.DB
	Router *router.Router
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
func NewApp(d *db.DB, r *router.Router, options ...Option) *App {
	a := &App{d, r, ":3000"}
	for _, option := range options {
		option(a)
	}
	return a
}

// Run starts the app on given port. Default port is ":3000"
func (a *App) Run() {
	log.Fatal(a.Router.Serve(a.Port))
}
