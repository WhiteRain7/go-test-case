package usecase

import (
	"context"

	"github.com/WhiteRain7/GoTestCase/exchanger/internal/rates/domain"
)

type UseCase struct {
	repo   RatesRepository
	logger Logger
}

type Logger interface {
	Info(string)
}

type RatesRepository interface {
	Get(ctx context.Context, curr string) (*domain.Rate, error)
	All(ctx context.Context) ([]domain.Rate, error)
	Save(ctx context.Context, rate domain.Rate) error
}

func NewUseCase(repo RatesRepository, logger Logger) *UseCase {
	logger.Info("Use case init - Rates")
	return &UseCase{repo, logger}
}

func (uc *UseCase) Get(ctx context.Context, curr string) (*domain.Rate, error) {
	return uc.repo.Get(ctx, curr)
}

func (uc *UseCase) All(ctx context.Context) (map[string]float32, error) {
	rates, err := uc.repo.All(ctx)
	if err != nil {
		return nil, err
	}
	results := make(map[string]float32)

	for _, rate := range rates {
		results[rate.Currency] = rate.Rate
	}

	return results, nil
}
