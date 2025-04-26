package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewDatabase(cfg *Config) (*pgx.Conn, error) {
	db, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://host=%s/port=%s/user=%s/dbname=%s/password=%s?sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	return db, err
}
