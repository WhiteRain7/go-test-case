package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/WhiteRain7/GoTestCase/pkg/utils/coded-errors"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/domain"
)

func (r *Repository) Create(ctx context.Context, userId int) (*domain.Wallet, error) {
	r.logger.Info(fmt.Sprintf("DB request - Wallets - Create wallet for user %d", userId))
	db := r.db.Get()

	query := "INSERT INTO wallets (\"user\") VALUES ($1) RETURNING id"
	row := db.QueryRowxContext(ctx, query, userId)

	var uuid domain.UUID
	if err := row.Scan(&uuid); err != nil {
		switch {
		default:
			return nil, codederrors.Unknown(err)
		}
	}

	return &domain.Wallet{uuid, userId}, nil
}

func (r *Repository) Get(ctx context.Context, userId int) (*domain.Wallet, error) {
	r.logger.Info(fmt.Sprintf("DB request - Wallets - Get wallet of user %d", userId))
	db := r.db.Get()

	query := "SELECT id FROM wallets WHERE \"user\" = $1"
	row := db.QueryRowxContext(ctx, query, userId)

	var uuid domain.UUID
	if err := row.Scan(&uuid); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, codederrors.New(http.StatusNotFound, ErrNotFound, "Wallet of user #%d doesn't exist.", userId)
		default:
			return nil, codederrors.Unknown(err)
		}
	}
	return &domain.Wallet{uuid, userId}, nil
}
