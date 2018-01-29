package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/rubenv/sql-migrate"
)

func main() {
	migrationDirection := flag.String("direction", "up", "direction of migration (up or down)")
	flag.Parse()
	direction := migrate.Up
	if *migrationDirection == "down" {
		direction = migrate.Down
	}

	var err error
	dbStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbStr)
	if err != nil {
		log.Fatalf("Unable to connect to db: %s\n", err.Error())
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations/postgres",
	}

	n, err := migrate.Exec(db, "postgres", migrations, direction)
	if err != nil {
		log.Fatalf("Unable to apply migrations: %s\n", err.Error())
	}
	log.Printf("Applied %d migrations.\n", n)
}
