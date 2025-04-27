package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PostgresConfig struct {
	DBPath string
}

func NewDatabase(cfg *PostgresConfig) (*pgx.Conn, error) {
	db, err := pgx.Connect(context.Background(), cfg.DBPath)
	if err != nil {
		return nil, err
	}
	return db, err
}
