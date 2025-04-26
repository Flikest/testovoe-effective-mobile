package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	db      pgx.Conn
	context context.Context
}

func InitStorage(s Storage) *Storage {
	return &Storage{
		db:      s.db,
		context: s.context,
	}
}
