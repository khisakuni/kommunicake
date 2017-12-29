package app

import (
	"log"

	"github.com/khisakuni/kommunicake/db"
	"github.com/khisakuni/kommunicake/router"
)

type App struct {
	DB     *db.DB
	Router *router.Router
}

type AppOption func(a *App)

func NewApp(...AppOption) *App {
	return &App{db.NewDB(), router.NewRouter()}
}

func (a *App) Run() {
	log.Fatal(a.Router.Serve(":3000"))
}
