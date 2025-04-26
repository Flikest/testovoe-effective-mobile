package service

import (
	"github.com/Flikest/testovoe-effective-mobile/internal/storage"
)

type Service struct {
	Storage *storage.Storage
}

func NewServices(s *storage.Storage) *Service {
	return &Service{Storage: s}
}
