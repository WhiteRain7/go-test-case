package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/WhiteRain7/GoTestCase/analytics/internal/transfers/domain"
)

type Repository struct {
	ch     Client
	logger Logger
}

type Logger interface {
	Info(string)
}

type Client interface {
	PrepareBatch(context.Context, string, ...driver.PrepareBatchOption) (driver.Batch, error)
}

func NewRepository(client Client, logger Logger) *Repository {
	return &Repository{client, logger}
}

func (r *Repository) Add(ctx context.Context, transfer domain.Transfer) error {
	return r.AddBatch(ctx, []domain.Transfer{transfer})
}

func (r *Repository) AddBatch(ctx context.Context, transfers []domain.Transfer) error {
	r.logger.Info(fmt.Sprintf("transfers repository - add batch - %d", len(transfers)))
	batch, err := r.ch.PrepareBatch(ctx, "INSERT INTO analytics.transfers")
	if err != nil {
		return err
	}

	for _, transfer := range transfers {
		ts, err := time.ParseInLocation("20060102150405", transfer.Timestamp, time.Local)

		err = batch.Append(
			uint64(transfer.User),
			transfer.Operation,
			float64(transfer.Amount),
			transfer.Currency,
			ts,
		)
		if err != nil {
			return err
		}
	}

	if err := batch.Send(); err != nil {
		return err
	}

	return nil
}
