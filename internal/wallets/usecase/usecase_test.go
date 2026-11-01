package usecase_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/WhiteRain7/GoTestCase/internal/wallets/domain"
	"github.com/WhiteRain7/GoTestCase/internal/wallets/handler"
	"github.com/WhiteRain7/GoTestCase/internal/wallets/repository"
	"github.com/WhiteRain7/GoTestCase/internal/wallets/usecase"
)

type mockRepo struct{}

func (r *mockRepo) Create(ctx context.Context) (*domain.Wallet, error) {
	return &domain.Wallet{ID: "abc", Amount: 0}, nil
}

func (r *mockRepo) Get(ctx context.Context, id domain.UUID) (*domain.Wallet, error) {
	wallets_ := []domain.Wallet{
		{ID: "abc", Amount: 0},
		{ID: "bcd", Amount: 200},
	}

	for _, wallet := range wallets_ {
		if wallet.ID == id {
			return &wallet, nil
		}
	}

	return nil, repository.ErrNotFound
}

func (r *mockRepo) Save(ctx context.Context, w *domain.Wallet) error {
	return nil
}

type MockLogger struct{}

func (*MockLogger) Info(string) {}

func TestCreate(t *testing.T) {
	tests := []struct {
		name       string
		wantWallet *domain.Wallet
		wantErr    error
	}{
		{"create new wallet", &domain.Wallet{"abc", 0}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uc := usecase.NewUseCase(&mockRepo{}, &MockLogger{})

			wallet, err := uc.Create(context.Background())

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(wallet, test.wantWallet) {
				t.Errorf("got wallet %v, want %v", wallet, test.wantWallet)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		ID         string
		wantWallet *domain.Wallet
		wantErr    error
	}{
		{"get empty wallet", "abc", &domain.Wallet{ID: "abc", Amount: 0}, nil},
		{"get non-empty wallet", "bcd", &domain.Wallet{ID: "bcd", Amount: 200}, nil},
		{"get non-existing wallet", "aaa", nil, repository.ErrNotFound},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uc := usecase.NewUseCase(&mockRepo{}, &MockLogger{})

			wallet, err := uc.Get(context.Background(), test.ID)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(wallet, test.wantWallet) {
				t.Errorf("got wallet %v, want %v", wallet, test.wantWallet)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name       string
		op         handler.WalletOperation
		wantWallet *domain.Wallet
		wantErr    error
	}{
		{"deposit sum", handler.WalletOperation{WalletId: "abc", OperationType: handler.OperationDeposit, Amount: 100}, &domain.Wallet{ID: "abc", Amount: 100}, nil},
		{"deposit negative", handler.WalletOperation{WalletId: "abc", OperationType: handler.OperationDeposit, Amount: -100}, nil, domain.ErrNegativeAmount},
		{"deposit zero", handler.WalletOperation{WalletId: "abc", OperationType: handler.OperationDeposit, Amount: 0}, &domain.Wallet{ID: "abc", Amount: 0}, nil},

		{"withdraw sum", handler.WalletOperation{WalletId: "bcd", OperationType: handler.OperationWithdraw, Amount: 100}, &domain.Wallet{ID: "bcd", Amount: 100}, nil},
		{"withdraw full sum", handler.WalletOperation{WalletId: "bcd", OperationType: handler.OperationWithdraw, Amount: 200}, &domain.Wallet{ID: "bcd", Amount: 0}, nil},
		{"withdraw more than exists", handler.WalletOperation{WalletId: "bcd", OperationType: handler.OperationWithdraw, Amount: 1000}, nil, domain.ErrNotEnough},
		{"withdraw negative", handler.WalletOperation{WalletId: "bcd", OperationType: handler.OperationWithdraw, Amount: -100}, nil, domain.ErrNegativeAmount},
		{"withdraw zero", handler.WalletOperation{WalletId: "abc", OperationType: handler.OperationWithdraw, Amount: 0}, &domain.Wallet{ID: "abc", Amount: 0}, nil},

		{"unsupported operation", handler.WalletOperation{WalletId: "abc", OperationType: "MULTIPLY", Amount: 0}, nil, usecase.ErrInvalidOperation},
		{"non-existing wallet", handler.WalletOperation{WalletId: "aaa", OperationType: handler.OperationDeposit, Amount: 100}, nil, repository.ErrNotFound},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uc := usecase.NewUseCase(&mockRepo{}, &MockLogger{})

			wallet, err := uc.Update(context.Background(), test.op)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(wallet, test.wantWallet) {
				t.Errorf("got wallet %v, want %v", wallet, test.wantWallet)
			}
		})
	}
}
