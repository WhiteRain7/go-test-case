package repository

import (
	"context"
	"time"

	"github.com/WhiteRain7/GoTestCase/pkg/utils/coded-errors"
	pb "github.com/WhiteRain7/GoTestCase/proto"
)

type Repository struct {
	grpc   pb.ExchangeServiceClient
	logger Logger
}

type Logger interface {
	Info(string)
}

func NewRepository(grpc pb.ExchangeServiceClient, logger Logger) *Repository {
	logger.Info("New repository - Wallets")
	return &Repository{grpc, logger}
}

func (uc *Repository) GetRates(ctx context.Context) (map[string]float32, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := uc.grpc.GetRates(ctx, &pb.Empty{})
	if err != nil {
		return nil, codederrors.Unknown(err)
	}
	return resp.Rates, nil
}

func (uc *Repository) Exchange(ctx context.Context, fromCurr string, toCurr string) (float32, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := uc.grpc.Exchange(ctx, &pb.ExchangeRequest{FromCurrency: fromCurr, ToCurrency: toCurr})
	if err != nil {
		return 0, codederrors.Unknown(err)
	}
	return resp.Rate, nil
}
