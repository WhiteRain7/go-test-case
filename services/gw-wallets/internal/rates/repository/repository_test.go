package repository_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	pb "github.com/WhiteRain7/GoTestCase/proto"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/rates/repository"
	"google.golang.org/grpc"
)

type MockLogger struct{}

func (*MockLogger) Info(string) {}

type MockGRPC struct{}

var (
	ErrInvalidCurrency = errors.New("invalid currency")
	ErrInvalidExchange = errors.New("invalid exchange")
)

func (m *MockGRPC) GetRates(ctx context.Context, in *pb.Empty, opts ...grpc.CallOption) (*pb.GetRatesResponse, error) {
	return &pb.GetRatesResponse{Rates: map[string]float32{"RUB": 1, "USD": 100}}, nil
}

func (m *MockGRPC) Exchange(ctx context.Context, in *pb.ExchangeRequest, opts ...grpc.CallOption) (*pb.ExchangeResponse, error) {
	rates := map[string]float32{"RUB": 1, "USD": 100}
	if _, ok := rates[in.FromCurrency]; !ok {
		return nil, ErrInvalidCurrency
	}
	if _, ok := rates[in.ToCurrency]; !ok {
		return nil, ErrInvalidCurrency
	}
	if in.ToCurrency == in.FromCurrency {
		return nil, ErrInvalidExchange
	}
	return &pb.ExchangeResponse{
		FromCurrency: in.FromCurrency,
		ToCurrency:   in.ToCurrency,
		Rate:         rates[in.FromCurrency] / rates[in.ToCurrency],
	}, nil
}

func TestRepository_GetRates(t *testing.T) {
	tests := []struct {
		name      string
		wantRates map[string]float32
		wantErr   error
	}{
		{"get rates", map[string]float32{"RUB": 1, "USD": 100}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := repository.NewRepository(&MockGRPC{}, &MockLogger{})

			rates, err := repo.GetRates(context.Background())

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(rates, test.wantRates) {
				t.Errorf("got rates %v, want %v", rates, test.wantRates)
			}
		})
	}
}

func TestRepository_Exchange(t *testing.T) {
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
			repo := repository.NewRepository(&MockGRPC{}, &MockLogger{})

			rate, err := repo.Exchange(context.Background(), test.fromCurrency, test.toCurrency)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(rate, test.wantRate) {
				t.Errorf("got rate %v, want %v", rate, test.wantRate)
			}
		})
	}
}
