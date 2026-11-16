package usecase

import (
	"context"

	"github.com/WhiteRain7/GoTestCase/notification/internal/transfers/domain"
)

type UseCase struct {
	repo Repository
}

type Repository interface {
	Add(context.Context, domain.Transfer) error
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo}
}

func (uc *UseCase) Add(ctx context.Context, transfer domain.Transfer) error {
	return uc.repo.Add(ctx, transfer)
}
