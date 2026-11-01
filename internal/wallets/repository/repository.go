package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/WhiteRain7/GoTestCase/internal/shared/infrastructure/apperr"
	"github.com/WhiteRain7/GoTestCase/internal/wallets/domain"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFound = errors.New("Wallet not found")
)

type Repository struct {
	db     *sqlx.DB
	logger Logger
}

type Logger interface {
	Info(string)
}

func NewRepository(db *sqlx.DB, logger Logger) *Repository {
	logger.Info("New repository - Wallets")
	return &Repository{db, logger}
}

func (r *Repository) Get(ctx context.Context, id domain.UUID) (*domain.Wallet, error) {
	r.logger.Info(fmt.Sprintf("DB request - Wallets - Get %s", id))
	row := r.db.QueryRowContext(ctx, "SELECT amount FROM wallets WHERE id = $1", id)
	var amount int
	if err := row.Scan(&amount); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, apperr.New(http.StatusNotFound, ErrNotFound, "Wallet with id '%s' doesn't exists.", id)
		default:
			return nil, apperr.Unknown(err)
		}
	}
	return &domain.Wallet{ID: id, Amount: amount}, nil
}

func (r *Repository) Create(ctx context.Context) (*domain.Wallet, error) {
	r.logger.Info("DB request - Wallets - Create")
	row := r.db.QueryRowContext(ctx, "INSERT INTO wallets (amount) VALUES (0) RETURNING id")
	var id string
	if err := row.Scan(&id); err != nil {
		switch {
		default:
			return nil, apperr.Unknown(err)
		}
	}
	return &domain.Wallet{ID: id, Amount: 0}, nil
}

func (r *Repository) Save(ctx context.Context, w *domain.Wallet) error {
	r.logger.Info(fmt.Sprintf("DB request - Wallets - Save %s", w.ID))
	_, err := r.db.ExecContext(ctx, "UPDATE wallets SET amount = $2 WHERE id = $1", w.ID, w.Amount)
	if err != nil {
		switch {
		default:
			return apperr.Unknown(err)
		}
	} else {
		return nil
	}
}
