package grpc_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/WhiteRain7/GoTestCase/exchanger/internal/adapters/grpc"
	"github.com/WhiteRain7/GoTestCase/exchanger/internal/rates/domain"
	pb "github.com/WhiteRain7/GoTestCase/proto"
)

var (
	ErrNotFound = errors.New("not found")
)

type MockUseCase struct{}

func (m *MockUseCase) Get(ctx context.Context, s string) (*domain.Rate, error) {
	if s == "RUB" {
		return &domain.Rate{"RUB", 1}, nil
	}
	if s == "USD" {
		return &domain.Rate{"USD", 100}, nil
	}
	return nil, ErrNotFound
}

func (m *MockUseCase) All(ctx context.Context) (map[string]float32, error) {
	return map[string]float32{"RUB": 1, "USD": 100}, nil
}

func TestExchangeService_GetRates(t *testing.T) {
	tests := []struct {
		name    string
		wantRes *pb.GetRatesResponse
		wantErr error
	}{
		{
			"all",
			&pb.GetRatesResponse{Rates: map[string]float32{"RUB": 1, "USD": 100}},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := grpc.NewExchangeService(&MockUseCase{})

			rates, err := service.GetRates(context.Background(), &pb.Empty{})

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(rates, test.wantRes) {
				t.Errorf("got rates %v, want %v", rates, test.wantRes)
			}
		})
	}
}

func TestExchangeService_Exchange(t *testing.T) {
	tests := []struct {
		name    string
		body    *pb.ExchangeRequest
		wantRes *pb.ExchangeResponse
		wantErr error
	}{
		{
			"RUB -> USD",
			&pb.ExchangeRequest{FromCurrency: "RUB", ToCurrency: "USD"},
			&pb.ExchangeResponse{FromCurrency: "RUB", ToCurrency: "USD", Rate: 0.01},
			nil,
		},
		{
			"USD -> RUB",
			&pb.ExchangeRequest{FromCurrency: "USD", ToCurrency: "RUB"},
			&pb.ExchangeResponse{FromCurrency: "USD", ToCurrency: "RUB", Rate: 100},
			nil,
		},
		{
			"RUB -> RUB",
			&pb.ExchangeRequest{FromCurrency: "RUB", ToCurrency: "RUB"},
			&pb.ExchangeResponse{FromCurrency: "RUB", ToCurrency: "RUB", Rate: 1},
			nil,
		},
		{
			"Invalid currency 1",
			&pb.ExchangeRequest{FromCurrency: "GRB", ToCurrency: "RUB"},
			&pb.ExchangeResponse{},
			ErrNotFound,
		},
		{
			"Invalid currency 2",
			&pb.ExchangeRequest{FromCurrency: "RUB", ToCurrency: "GRB"},
			&pb.ExchangeResponse{},
			ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := grpc.NewExchangeService(&MockUseCase{})

			rate, err := service.Exchange(context.Background(), test.body)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(rate, test.wantRes) {
				t.Errorf("got rate %v, want %v", rate, test.wantRes)
			}
		})
	}
}
