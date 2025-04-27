package storage

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	DB      *pgx.Conn
	Context context.Context
	Log     *slog.Logger
}

func InitStorage(s *Storage) *Storage {
	return &Storage{
		DB:      s.DB,
		Context: s.Context,
		Log:     s.Log,
	}
}
