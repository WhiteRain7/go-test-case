package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFound = errors.New("Not found")
)

type Repository struct {
	db     DBGetter
	logger Logger
}

type DBGetter interface {
	Get() sqlx.ExtContext
}

type Logger interface {
	Info(string)
}

func NewRepository(db DBGetter, logger Logger) *Repository {
	logger.Info("New repository - Wallets")
	return &Repository{db, logger}
}
