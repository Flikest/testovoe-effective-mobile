package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	var storagePath, migrationsPath string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")

	if storagePath == "" || migrationsPath == "" {
		panic("storage-path or migrations-path is required")
	}

	migratinos, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("postgres://%s:%s/%s?sslmode=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_SSL")),
	)
	if err != nil {
		panic(err)
	}

	if err := migratinos.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("migrations not apply")
			return
		}
		panic(err)
	}
	log.Println("migrations up")
}
