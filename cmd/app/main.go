package main

import (
	"github.com/gorilla/mux"
	"github.com/khisakuni/kommunicake/api/v1"
	"github.com/khisakuni/kommunicake/app"
	"github.com/khisakuni/kommunicake/db"
)

func main() {
	d, err := db.NewDB()
	r := mux.NewRouter()
	v1.CreateV1Routes(r)

	a := app.NewApp(d, r)
	if err != nil {
		panic(err)
	}

	a.Run()
}
