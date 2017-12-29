package main

import "github.com/khisakuni/kommunicake/app"
import "github.com/khisakuni/kommunicake/db"
import "github.com/khisakuni/kommunicake/router"

func main() {
	d, err := db.NewDB()
	r := router.NewRouter()

	a := app.NewApp(d, r)
	if err != nil {
		panic(err)
	}

	a.Run()
}
