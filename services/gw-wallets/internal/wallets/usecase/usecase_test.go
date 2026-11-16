package usecase_test

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"testing"

	"github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/domain"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/repository"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/usecase"
)

type MockTXManager struct{}

func (*MockTXManager) BeginTx() error  { return nil }
func (*MockTXManager) Commit() error   { return nil }
func (*MockTXManager) Rollback() error { return nil }

type MockLogger struct{}

func (*MockLogger) Info(string) {}

type MockRepo struct{}

func (*MockRepo) Get(ctx context.Context, userId int) (*domain.Wallet, error) {
	wallets := []*domain.Wallet{
		&domain.Wallet{"abc", 1},
		&domain.Wallet{"bcd", 2},
	}
	for _, wallet := range wallets {
		if wallet.User == userId {
			return wallet, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (*MockRepo) AllCurrencies(ctx context.Context) ([]string, error) {
	return []string{"USD", "EUR", "RUB"}, nil
}

func (r *MockRepo) GetCurrency(ctx context.Context, curr string) (string, error) {
	supported, _ := r.AllCurrencies(ctx)
	if slices.Contains(supported, curr) {
		return curr, nil
	}
	return "", repository.ErrNotFound
}

func (*MockRepo) GetBalance(ctx context.Context, walletId domain.UUID) (*domain.WalletBalance, error) {
	balances := []*domain.WalletBalance{
		&domain.WalletBalance{"abc", map[string]float32{"RUB": 1}},
	}
	for _, balance := range balances {
		if balance.Wallet == walletId {
			return balance, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *MockRepo) SetBalance(ctx context.Context, walletId domain.UUID, curr string, amount float32) error {
	balance, err := r.GetBalance(ctx, walletId)
	if err != nil {
		return err
	}
	balance.Balance[curr] = amount
	return nil
}

func (r *MockRepo) UpdateBalance(ctx context.Context, walletId domain.UUID, curr string, delta float32) error {
	balance, err := r.GetBalance(ctx, walletId)
	if err != nil {
		return err
	}
	balance.Balance[curr] += delta
	return nil
}

type MockRates struct{}

func (*MockRates) Exchange(ctx context.Context, fromCurr string, toCurr string) (float32, error) {
	rates := map[string]float32{"RUB": 1, "USD": 100, "EUR": 50}
	if _, ok := rates[fromCurr]; !ok {
		return 0, repository.ErrNotFound
	}
	if _, ok := rates[toCurr]; !ok {
		return 0, repository.ErrNotFound
	}
	return rates[fromCurr] / rates[toCurr], nil
}

type MockProducer struct{}

func (p *MockProducer) Send(context.Context, []byte) error {
	return nil
}

func TestUseCase_Get(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		wantBalance *domain.WalletBalance
		wantErr     error
	}{
		{"get", 1, &domain.WalletBalance{"abc", map[string]float32{"RUB": 1}}, nil},
		{"not found", 3, nil, repository.ErrNotFound},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uc := usecase.NewUseCase(&MockTXManager{}, &MockRepo{}, &MockRates{}, &MockProducer{}, &MockLogger{})

			balance, err := uc.Get(context.Background(), test.id)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(balance, test.wantBalance) {
				t.Errorf("got balance %v, want %v", balance, test.wantBalance)
			}
		})
	}
}

func TestUseCase_Deposit(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		curr        string
		amount      float32
		wantBalance *domain.WalletBalance
		wantErr     error
	}{
		{
			"deposit used curr",
			1,
			"RUB",
			100.0,
			&domain.WalletBalance{"abc", map[string]float32{"RUB": 101}},
			nil,
		},
		{
			"deposit not used curr",
			1,
			"EUR",
			100.0,
			&domain.WalletBalance{"abc", map[string]float32{"RUB": 1, "EUR": 100}},
			nil,
		},
		{
			"deposit non-existing curr",
			1,
			"GRB",
			100.0,
			nil,
			usecase.ErrCurrencyNotSupported,
		},
		{
			"deposit non-existing wallet",
			2,
			"RUB",
			100.0,
			nil,
			repository.ErrNotFound,
		},
		{
			"deposit negative",
			1,
			"RUB",
			-100.0,
			nil,
			usecase.ErrNegativeAmount,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uc := usecase.NewUseCase(&MockTXManager{}, &MockRepo{}, &MockRates{}, &MockProducer{}, &MockLogger{})

			balance, err := uc.Deposit(context.Background(), test.id, test.curr, test.amount)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(balance, test.wantBalance) {
				t.Errorf("got balance %v, want %v", balance, test.wantBalance)
			}
		})
	}
}

func TestUseCase_Withdraw(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		curr        string
		amount      float32
		wantBalance *domain.WalletBalance
		wantErr     error
	}{
		{
			"withdraw used curr",
			1,
			"RUB",
			1.0,
			&domain.WalletBalance{"abc", map[string]float32{"RUB": 0}},
			nil,
		},
		{
			"withdraw non-existing curr",
			1,
			"GRB",
			100.0,
			nil,
			usecase.ErrCurrencyNotSupported,
		},
		{
			"withdraw non-existing wallet",
			2,
			"RUB",
			100.0,
			nil,
			repository.ErrNotFound,
		},
		{
			"withdraw negative",
			1,
			"RUB",
			-100.0,
			nil,
			usecase.ErrNegativeAmount,
		},
		{
			"withdraw too much",
			1,
			"RUB",
			100.0,
			nil,
			usecase.ErrNotEnoughFunds,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uc := usecase.NewUseCase(&MockTXManager{}, &MockRepo{}, &MockRates{}, &MockProducer{}, &MockLogger{})

			balance, err := uc.Withdraw(context.Background(), test.id, test.curr, test.amount)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(balance, test.wantBalance) {
				t.Errorf("got balance %v, want %v", balance, test.wantBalance)
			}
		})
	}
}

func TestUseCase_Exchange(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		fromCurr    string
		toCurr      string
		amount      float32
		wantBalance *domain.WalletBalance
		wantErr     error
	}{
		{
			"exchange used curr to unused",
			1,
			"RUB",
			"USD",
			1.0,
			&domain.WalletBalance{"abc", map[string]float32{"RUB": 0, "USD": 0.01}},
			nil,
		},
		{
			"exchange with unknown currency 1",
			1,
			"GRB",
			"RUB",
			1.0,
			nil,
			usecase.ErrCurrencyNotSupported,
		},
		{
			"exchange with unknown currency 2",
			1,
			"RUB",
			"GRB",
			1.0,
			nil,
			usecase.ErrCurrencyNotSupported,
		},
		{
			"exchange zero value",
			1,
			"RUB",
			"USD",
			0.0,
			nil,
			usecase.ErrNegativeAmount,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uc := usecase.NewUseCase(&MockTXManager{}, &MockRepo{}, &MockRates{}, &MockProducer{}, &MockLogger{})

			balance, _, err := uc.Exchange(context.Background(), test.id, test.fromCurr, test.toCurr, test.amount)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(balance, test.wantBalance) {
				t.Errorf("got balance %v, want %v", balance, test.wantBalance)
			}
		})
	}
}
