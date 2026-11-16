package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/WhiteRain7/GoTestCase/pkg/utils/coded-errors"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/domain"
)

var (
	ErrNegativeAmount       = errors.New("negative amount")
	ErrNotEnoughFunds       = errors.New("not enough funds")
	ErrCurrencyNotSupported = errors.New("currency is not supported")
	ErrInvalidExchange      = errors.New("invalid exchange")
)

type HighValueTransfer struct {
	User      int     `json:"user"`
	Operation string  `json:"operation"`
	Amount    float32 `json:"amount"`
	Currency  string  `json:"currency"`
	Timestamp string  `json:"timestamp"`
}

type UseCase struct {
	txm      TXManager
	repo     WalletRepository
	rates    RatesRepository
	producer Producer
	logger   Logger
}

type TXManager interface {
	BeginTx() error
	Commit() error
	Rollback() error
}

type Producer interface {
	Send(context.Context, []byte) error
}

type Logger interface {
	Info(string)
}

type WalletRepository interface {
	Get(ctx context.Context, userId int) (*domain.Wallet, error)
	GetCurrency(ctx context.Context, curr string) (string, error)
	AllCurrencies(ctx context.Context) ([]string, error)
	GetBalance(ctx context.Context, walletId domain.UUID) (*domain.WalletBalance, error)
	SetBalance(ctx context.Context, wallet domain.UUID, curr string, amount float32) error
	UpdateBalance(ctx context.Context, wallet domain.UUID, curr string, delta float32) error
}

type RatesRepository interface {
	Exchange(ctx context.Context, fromCurr string, toCurr string) (float32, error)
}

func NewUseCase(txm TXManager, repo WalletRepository, rates RatesRepository, producer Producer, logger Logger) *UseCase {
	logger.Info("Use case init - Wallets")
	return &UseCase{txm, repo, rates, producer, logger}
}

func (uc *UseCase) Get(ctx context.Context, userId int) (*domain.WalletBalance, error) {
	wallet, err := uc.repo.Get(ctx, userId)
	if err != nil {
		return nil, err
	}

	balance, err := uc.repo.GetBalance(ctx, wallet.ID)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (uc *UseCase) currencySupported(ctx context.Context, key string) bool {
	_, err := uc.repo.GetCurrency(ctx, key)
	return err == nil
}

func (uc *UseCase) Deposit(
	ctx context.Context,
	userId int,
	curr string,
	amount float32,
) (*domain.WalletBalance, error) {
	if amount <= 0 {
		return nil, codederrors.New(http.StatusBadRequest, ErrNegativeAmount, "Amount must be positive, got %f", amount)
	}

	uc.txm.BeginTx()
	defer uc.txm.Rollback()

	wallet, err := uc.repo.Get(ctx, userId)
	if err != nil {
		return nil, err
	}

	balance, err := uc.repo.GetBalance(ctx, wallet.ID)
	if err != nil {
		return nil, err
	}

	if _, ok := balance.Balance[curr]; !ok {
		if supported := uc.currencySupported(ctx, curr); !supported {
			return nil, codederrors.New(http.StatusBadRequest, ErrCurrencyNotSupported, "Currency %s is not supported", curr)
		}

		err = uc.repo.SetBalance(ctx, wallet.ID, curr, +amount)
		if err != nil {
			return nil, err
		}
		balance.Balance[curr] = amount
	} else {
		err = uc.repo.UpdateBalance(ctx, wallet.ID, curr, +amount)
		if err != nil {
			return nil, err
		}
		balance.Balance[curr] += amount
	}

	if amount >= 30_000 {
		transfer := HighValueTransfer{userId, "DEPOSIT", amount, curr, time.Now().Format("20060102150405")}
		bytes, err := json.Marshal(transfer)
		if err != nil {
			return nil, err
		}
		err = uc.producer.Send(ctx, bytes)
		if err != nil {
			return nil, err
		}
	}

	uc.txm.Commit()

	return balance, nil
}

func (uc *UseCase) Withdraw(
	ctx context.Context,
	userId int,
	curr string,
	amount float32,
) (*domain.WalletBalance, error) {
	if amount <= 0 {
		return nil, codederrors.New(http.StatusBadRequest, ErrNegativeAmount, "Amount must be positive, got %f", amount)
	}

	uc.txm.BeginTx()
	defer uc.txm.Rollback()

	wallet, err := uc.repo.Get(ctx, userId)
	if err != nil {
		return nil, err
	}

	balance, err := uc.repo.GetBalance(ctx, wallet.ID)
	if err != nil {
		return nil, err
	}

	if _, ok := balance.Balance[curr]; !ok {
		if supported := uc.currencySupported(ctx, curr); !supported {
			return nil, codederrors.New(http.StatusBadRequest, ErrCurrencyNotSupported, "Currency %s is not supported", curr)
		}
		balance.Balance[curr] = 0
	}

	if balance.Balance[curr] < amount {
		return nil, codederrors.New(http.StatusBadRequest, ErrNotEnoughFunds, "Not enough funds to withdraw %f %s", amount, curr)
	}

	err = uc.repo.UpdateBalance(ctx, wallet.ID, curr, -amount)
	if err != nil {
		return nil, err
	}

	balance.Balance[curr] -= amount

	if amount >= 30_000 {
		transfer := HighValueTransfer{userId, "WITHDRAW", amount, curr, time.Now().Format("20060102150405")}
		bytes, err := json.Marshal(transfer)
		if err != nil {
			return nil, err
		}
		err = uc.producer.Send(ctx, bytes)
		if err != nil {
			return nil, err
		}
	}

	uc.txm.Commit()

	return balance, nil
}

func (uc *UseCase) Exchange(
	ctx context.Context,
	userId int,
	fromCurr string,
	toCurr string,
	amount float32,
) (*domain.WalletBalance, float32, error) {
	if amount <= 0 {
		return nil, 0, codederrors.New(http.StatusBadRequest, ErrNegativeAmount, "Amount must be positive, got %f", amount)
	}
	if fromCurr == toCurr {
		return nil, 0, codederrors.New(http.StatusBadRequest, ErrInvalidExchange, "Invalid exchange %s to %s", fromCurr, toCurr)
	}

	uc.txm.BeginTx()
	defer uc.txm.Rollback()

	wallet, err := uc.repo.Get(ctx, userId)
	if err != nil {
		return nil, 0, err
	}

	balance, err := uc.repo.GetBalance(ctx, wallet.ID)
	if err != nil {
		return nil, 0, err
	}

	setTo := false

	if _, ok := balance.Balance[fromCurr]; !ok {
		if supported := uc.currencySupported(ctx, fromCurr); !supported {
			return nil, 0, codederrors.New(http.StatusBadRequest, ErrCurrencyNotSupported, "Currency %s is not supported", fromCurr)
		}
		balance.Balance[fromCurr] = 0
	}

	if balance.Balance[fromCurr] < amount {
		return nil, 0, codederrors.New(http.StatusBadRequest, ErrNotEnoughFunds, "Not enough funds to withdraw %f %s", amount, fromCurr)
	}

	if _, ok := balance.Balance[toCurr]; !ok {
		if supported := uc.currencySupported(ctx, toCurr); !supported {
			return nil, 0, codederrors.New(http.StatusBadRequest, ErrCurrencyNotSupported, "Currency %s is not supported", toCurr)
		}
		setTo = true
		balance.Balance[toCurr] = 0
	}

	rate, err := uc.rates.Exchange(ctx, fromCurr, toCurr)
	if err != nil {
		return nil, 0, err
	}

	converted := amount * rate

	err = uc.repo.UpdateBalance(ctx, wallet.ID, fromCurr, -amount)
	if err != nil {
		return nil, 0, err
	}

	if setTo {
		err = uc.repo.SetBalance(ctx, wallet.ID, toCurr, converted)
	} else {
		err = uc.repo.UpdateBalance(ctx, wallet.ID, toCurr, +converted)
	}
	if err != nil {
		return nil, 0, err
	}

	balance.Balance[fromCurr] -= amount
	balance.Balance[toCurr] += converted

	b := make(map[string]float32)
	b[fromCurr] = balance.Balance[fromCurr]
	b[toCurr] = balance.Balance[toCurr]

	balance.Balance = b

	uc.txm.Commit()

	return balance, amount * rate, nil
}
