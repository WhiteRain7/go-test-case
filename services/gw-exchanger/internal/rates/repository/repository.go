package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/WhiteRain7/GoTestCase/exchanger/internal/rates/domain"
	"github.com/WhiteRain7/GoTestCase/pkg/utils/coded-errors"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFound = errors.New("Rate not found")
)

type Repository struct {
	db     *sqlx.DB
	logger Logger
}

type Logger interface {
	Info(string)
}

func NewRepository(db *sqlx.DB, logger Logger) *Repository {
	logger.Info("New repository - Rates")
	return &Repository{db, logger}
}

func (r *Repository) Get(ctx context.Context, curr string) (*domain.Rate, error) {
	r.logger.Info(fmt.Sprintf("DB request - Rates - Get %s", curr))
	row := r.db.QueryRowContext(ctx, "SELECT rate FROM Rates WHERE currency = $1", curr)
	var rate float32
	if err := row.Scan(&rate); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, codederrors.New(http.StatusNotFound, ErrNotFound, "Rate of currency '%s' doesn't exist.", curr)
		default:
			return nil, codederrors.Unknown(err)
		}
	}
	return &domain.Rate{Currency: curr, Rate: rate}, nil
}

func (r *Repository) All(ctx context.Context) ([]domain.Rate, error) {
	r.logger.Info("DB request - Rates - All")
	rows, err := r.db.QueryContext(ctx, "SELECT currency, rate FROM Rates")
	if err != nil {
		return nil, codederrors.Unknown(err)
	}
	defer rows.Close()

	var results []domain.Rate

	for rows.Next() {
		var cur string
		var rate float32
		if err := rows.Scan(&cur, &rate); err != nil {
			return nil, codederrors.Unknown(err)
		}
		results = append(results, domain.Rate{cur, rate})
	}

	return results, nil
}

func (r *Repository) Save(ctx context.Context, rate domain.Rate) error {
	r.logger.Info(fmt.Sprintf("DB request - Wallets - Save %s rate %f", rate.Currency, rate.Rate))
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO rates (currency, rate) VALUES ($1, $2) ON CONFLICT (currency) DO UPDATE SET rate = EXCLUDED.rate)",
		rate.Currency, rate.Rate,
	)
	if err != nil {
		switch {
		default:
			return codederrors.Unknown(err)
		}
	} else {
		return nil
	}
}
