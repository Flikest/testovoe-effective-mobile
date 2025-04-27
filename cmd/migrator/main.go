package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	var migrationsPath string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")

	flag.Parse()

	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	migratinos, err := migrate.New(
		"file://"+migrationsPath,
		os.Getenv("DB_PATH"),
	)
	if err != nil {
		panic(err)
	}

	if err := migratinos.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("migrations not apply")
			return
		}
		panic(fmt.Sprintf("error migrations up: %s", err))
	}
	log.Println("migrations up")
}
