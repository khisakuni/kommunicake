package main

import (
	"fmt"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"github.com/khisakuni/kommunicake/api/v1"
	"github.com/khisakuni/kommunicake/app"
	"github.com/khisakuni/kommunicake/database"
)

func main() {

	fmt.Println("Starting web server")

	// Initialize DB connection
	fmt.Println("Initializing database...")
	db, err := database.NewDB()
	defer db.Conn.Close()

	// Create router
	fmt.Println("Initializing router...")
	r := mux.NewRouter()

	// Setup routes
	fmt.Println("Setting up routes...")
	v1.Routes(r, db)

	var port string
	if p := os.Getenv("PORT"); len(p) > 0 {
		port = p
	} else {
		port = "3000"
	}

	// Create App
	a := app.NewApp(db, r, app.WithPort(":"+port))
	if err != nil {
		panic(err)
	}

	// Run app on default port 3000
	a.Run()
}
