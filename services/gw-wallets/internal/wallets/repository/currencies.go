package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/WhiteRain7/GoTestCase/pkg/utils/coded-errors"
)

func (r *Repository) GetCurrency(ctx context.Context, key string) (string, error) {
	r.logger.Info(fmt.Sprintf("DB request - Wallets - Get currency %s", key))
	db := r.db.Get()

	query := "SELECT key FROM currency_types WHERE key = $1"
	row := db.QueryRowxContext(ctx, query, key)

	var _key string

	if err := row.Scan(&_key); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", codederrors.New(http.StatusNotFound, ErrNotFound, "Currency with key '%s' doesn't exist.", key)
		default:
			return "", codederrors.Unknown(err)
		}
	}

	return _key, nil
}

func (r *Repository) AllCurrencies(ctx context.Context) ([]string, error) {
	r.logger.Info("DB request - Wallets - Get currencies")
	db := r.db.Get()

	query := "SELECT key FROM currency_types"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		// TODO: errNoRows
		return nil, codederrors.Unknown(err)
	}
	defer rows.Close()

	result := make([]string, 0)

	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, codederrors.Unknown(err)
		}
		result = append(result, key)
	}

	return result, nil
}
