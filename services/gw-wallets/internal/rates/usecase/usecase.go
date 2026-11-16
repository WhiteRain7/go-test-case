package usecase

import (
	"context"
	"errors"
)

var (
	ErrInvalidOperation = errors.New("Invalid operation")
)

type UseCase struct {
	repo   RatesRepo
	logger Logger
}

type Logger interface {
	Info(string)
}

type RatesRepo interface {
	GetRates(ctx context.Context) (map[string]float32, error)
	Exchange(ctx context.Context, fromCurr string, toCurr string) (float32, error)
}

func NewUseCase(repo RatesRepo, logger Logger) *UseCase {
	logger.Info("Use case init - Wallets")
	return &UseCase{repo, logger}
}

func (uc *UseCase) All(ctx context.Context) (map[string]float32, error) {
	return uc.repo.GetRates(ctx)
}

func (uc *UseCase) Exchange(ctx context.Context, fromCurr string, toCurr string) (float32, error) {
	return uc.repo.Exchange(ctx, fromCurr, toCurr)
}
