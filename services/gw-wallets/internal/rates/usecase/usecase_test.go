package usecase_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/WhiteRain7/GoTestCase/wallets/internal/rates/usecase"
)

var (
	ErrInvalidCurrency = errors.New("invalid currency")
	ErrInvalidExchange = errors.New("invalid exchange")
)

type MockLogger struct{}

func (MockLogger) Info(string) {}

type MockRepo struct{}

func (m *MockRepo) GetRates(ctx context.Context) (map[string]float32, error) {
	return map[string]float32{"RUB": 1, "USD": 100}, nil
}

func (m *MockRepo) Exchange(ctx context.Context, fromCurr string, toCurr string) (float32, error) {
	rates := map[string]float32{"RUB": 1, "USD": 100}
	if _, ok := rates[fromCurr]; !ok {
		return 0, ErrInvalidCurrency
	}
	if _, ok := rates[toCurr]; !ok {
		return 0, ErrInvalidCurrency
	}
	if fromCurr == toCurr {
		return 0, ErrInvalidExchange
	}
	return rates[fromCurr] / rates[toCurr], nil
}

func TestUseCase_All(t *testing.T) {
	tests := []struct {
		name      string
		wantRates map[string]float32
		wantErr   error
	}{
		{"all", map[string]float32{"RUB": 1, "USD": 100}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uc := usecase.NewUseCase(&MockRepo{}, &MockLogger{})

			rates, err := uc.All(context.Background())

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(rates, test.wantRates) {
				t.Errorf("got rates %v, want %v", rates, test.wantRates)
			}
		})
	}
}

func TestUseCase_Exchange(t *testing.T) {
	tests := []struct {
		name         string
		fromCurrency string
		toCurrency   string
		wantRate     float32
		wantErr      error
	}{
		{
			"RUB -> USD",
			"RUB",
			"USD",
			0.01,
			nil,
		},
		{
			"USD -> RUB",
			"USD",
			"RUB",
			100.0,
			nil,
		},
		{
			"invalid currency 1",
			"GRB",
			"RUB",
			0,
			ErrInvalidCurrency,
		},
		{
			"invalid currency 2",
			"USD",
			"GRB",
			0,
			ErrInvalidCurrency,
		},
		{
			"RUB -> RUB",
			"RUB",
			"RUB",
			0,
			ErrInvalidExchange,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uc := usecase.NewUseCase(&MockRepo{}, &MockLogger{})

			rates, err := uc.Exchange(context.Background(), test.fromCurrency, test.toCurrency)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(rates, test.wantRate) {
				t.Errorf("got rates %v, want %v", rates, test.wantRate)
			}
		})
	}
}
