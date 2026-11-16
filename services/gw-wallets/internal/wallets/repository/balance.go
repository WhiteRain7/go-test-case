package repository

import (
	"context"
	"fmt"

	"github.com/WhiteRain7/GoTestCase/pkg/utils/coded-errors"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/domain"
)

func (r *Repository) GetBalance(ctx context.Context, walletId domain.UUID) (*domain.WalletBalance, error) {
	r.logger.Info(fmt.Sprintf("DB request - Wallets - Get wallet balance of %s", walletId))
	db := r.db.Get()

	query := "SELECT currency, amount FROM wallet_currencies WHERE wallet = $1"
	rows, err := db.QueryContext(ctx, query, walletId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	balance := domain.WalletBalance{walletId, make(map[string]float32)}

	for rows.Next() {
		var cur string
		var amount int
		if err := rows.Scan(&cur, &amount); err != nil {
			return nil, codederrors.Unknown(err)
		}
		balance.Balance[cur] = float32(amount) / 100
	}

	return &balance, nil
}

func (r *Repository) SetBalance(ctx context.Context, walletId domain.UUID, curr string, amount float32) error {
	r.logger.Info(fmt.Sprintf("DB request - Wallets - Update balance of %s", walletId))
	db := r.db.Get()

	query := "INSERT INTO wallet_currencies (wallet, currency, amount) VALUES ($1, $2, $3)"
	_, err := db.ExecContext(ctx, query, walletId, curr, int(amount*100))

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateBalance(ctx context.Context, walletId domain.UUID, curr string, delta float32) error {
	r.logger.Info(fmt.Sprintf("DB request - Wallets - Update balance of %s", walletId))
	db := r.db.Get()

	query := "UPDATE wallet_currencies SET amount = amount + $1 WHERE wallet = $2 AND currency = $3"
	_, err := db.ExecContext(ctx, query, int(delta*100), walletId, curr)

	if err != nil {
		return err
	}

	return nil
}
