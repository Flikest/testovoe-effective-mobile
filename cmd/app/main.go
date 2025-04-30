package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Flikest/testovoe-effective-mobile/internal/database/postgresql"
	"github.com/Flikest/testovoe-effective-mobile/internal/handler"
	"github.com/Flikest/testovoe-effective-mobile/internal/service"
	"github.com/Flikest/testovoe-effective-mobile/internal/storage"
	"github.com/Flikest/testovoe-effective-mobile/pkg/logger"
	"github.com/joho/godotenv"
)

// @title           тестовое задание
// @version         1.0

// @BasePath  /v1/user/
// @contact.email  grecmanviktor6@gmail.com

func main() {
	var port string
	flag.StringVar(&port, "port", "", "Port for the application")

	flag.Parse()

	godotenv.Load()

	log := logger.InitLogger(os.Getenv("LVL_DEPLOYMENT"))

	db, err := postgresql.NewDatabase(&postgresql.PostgresConfig{
		DBPath:  os.Getenv("DB_PATH"),
		Context: context.Background(),
	})
	if err != nil {
		log.Error("error connect to database: ", err)
	}

	storage := storage.InitStorage(&storage.Storage{
		DB:      db,
		Context: context.Background(),
		Log:     log,
	})
	service := service.NewServices(storage)
	handler := handler.NewHandler(service)
	router := handler.InitRouter()

	go func() {
		if err := router.Run(":" + port); err != nil {
			log.Error("servoer not started: ", err)
		}
	}()

	closer := make(chan os.Signal, 1)
	signal.Notify(closer, syscall.SIGTERM, syscall.SIGINT)
	<-closer

	log.Info("server shoting down...")

	if err := db.Close(storage.Context); err != nil {
		log.Error("error db close")
	}
}
