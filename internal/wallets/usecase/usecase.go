package usecase

import (
	"context"
	"errors"
	"net/http"

	"github.com/WhiteRain7/GoTestCase/internal/shared/infrastructure/apperr"
	"github.com/WhiteRain7/GoTestCase/internal/wallets/domain"
	"github.com/WhiteRain7/GoTestCase/internal/wallets/handler"
)

var (
	ErrInvalidOperation = errors.New("Invalid operation")
)

type UseCase struct {
	repo   WalletRepository
	logger Logger
}

type WalletRepository interface {
	Create(ctx context.Context) (*domain.Wallet, error)
	Get(ctx context.Context, id domain.UUID) (*domain.Wallet, error)
	Save(ctx context.Context, w *domain.Wallet) error
}

type Logger interface {
	Info(s string)
}

func NewUseCase(repo WalletRepository, logger Logger) *UseCase {
	logger.Info("Use case init - Wallets")
	return &UseCase{repo, logger}
}

func (uc *UseCase) Create(ctx context.Context) (*domain.Wallet, error) {
	return uc.repo.Create(ctx)
}

func (uc *UseCase) Get(ctx context.Context, id domain.UUID) (*domain.Wallet, error) {
	return uc.repo.Get(ctx, id)
}

func (uc *UseCase) Update(ctx context.Context, op handler.WalletOperation) (*domain.Wallet, error) {
	wallet, err := uc.repo.Get(ctx, op.WalletId)
	if err != nil {
		return nil, err
	}

	var operr error

	switch op.OperationType {
	case handler.OperationDeposit:
		operr = wallet.Deposit(op.Amount)
	case handler.OperationWithdraw:
		operr = wallet.Withdraw(op.Amount)
	default:
		return nil, apperr.New(http.StatusBadRequest, ErrInvalidOperation, "operation '%s' is not supported, choose DEPOSIT or WITHDRAW", op.OperationType)
	}

	if operr != nil {
		return nil, operr
	}

	if err = uc.repo.Save(ctx, wallet); err != nil {
		return nil, err
	}
	return wallet, nil
}
