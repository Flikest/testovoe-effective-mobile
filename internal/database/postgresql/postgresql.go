package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PostgresConfig struct {
	DBPath  string
	Context context.Context
}

func NewDatabase(cfg *PostgresConfig) (*pgx.Conn, error) {
	db, err := pgx.Connect(cfg.Context, cfg.DBPath)
	if err != nil {
		return nil, err
	}
	return db, err
}
