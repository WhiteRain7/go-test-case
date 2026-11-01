package domain_test

import (
	"errors"
	"testing"

	"github.com/WhiteRain7/GoTestCase/internal/wallets/domain"
)

func TestWallet_Deposit(t *testing.T) {
	tests := []struct {
		name       string
		start      int
		amount     int
		wantAmount int
		wantErr    error
	}{
		{"valid deposit", 1000, 100, 1100, nil},
		{"negative amount", 1000, -100, 1000, domain.ErrNegativeAmount},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet := &domain.Wallet{ID: "123", Amount: test.start}

			err := wallet.Deposit(test.amount)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if wallet.Amount != test.wantAmount {
				t.Errorf("got amount %d, want %d", wallet.Amount, test.wantAmount)
			}
		})
	}
}

func TestWallet_Withdraw(t *testing.T) {
	tests := []struct {
		name       string
		start      int
		amount     int
		wantAmount int
		wantErr    error
	}{
		{"valid withdraw", 1000, 100, 900, nil},
		{"full withdraw", 1000, 1000, 0, nil},
		{"not enough money", 1000, 2000, 1000, domain.ErrNotEnough},
		{"negative amount", 1000, -100, 1000, domain.ErrNegativeAmount},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet := &domain.Wallet{ID: "123", Amount: test.start}

			err := wallet.Withdraw(test.amount)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if wallet.Amount != test.wantAmount {
				t.Errorf("got amount %d, want %d", wallet.Amount, test.wantAmount)
			}
		})
	}
}
