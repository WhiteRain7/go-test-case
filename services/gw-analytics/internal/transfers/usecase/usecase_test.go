package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/WhiteRain7/GoTestCase/analytics/internal/transfers/domain"
	"github.com/WhiteRain7/GoTestCase/analytics/internal/transfers/usecase"
)

type MockLogger struct{}

func (*MockLogger) Info(string) {}

type MockRepo struct{}

func (m MockRepo) Add(ctx context.Context, transfer domain.Transfer) error {
	return nil
}

func TestUseCase_Add(t *testing.T) {
	tests := []struct {
		name     string
		transfer domain.Transfer
		wantErr  error
	}{
		{
			"transfer",
			domain.Transfer{1, "DEPOSIT", 30, "RUB", "11111111111"},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uc := usecase.NewUseCase(&MockRepo{})

			err := uc.Add(context.Background(), test.transfer)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}
		})
	}
}
