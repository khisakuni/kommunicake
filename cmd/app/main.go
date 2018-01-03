package main

import (
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"github.com/khisakuni/kommunicake/api/v1"
	"github.com/khisakuni/kommunicake/app"
	"github.com/khisakuni/kommunicake/database"
)

func main() {
	// Initialize DB connection
	db, err := database.NewDB()
	defer db.Conn.Close()

	// Create router
	r := mux.NewRouter()

	// Setup routes
	v1.Routes(r)

	// Create App
	a := app.NewApp(db, r)
	if err != nil {
		panic(err)
	}

	// Run app on default port 3000
	a.Run()
}
