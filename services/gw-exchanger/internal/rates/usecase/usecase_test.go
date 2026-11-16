package usecase_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/WhiteRain7/GoTestCase/exchanger/internal/rates/domain"
	"github.com/WhiteRain7/GoTestCase/exchanger/internal/rates/repository"
	"github.com/WhiteRain7/GoTestCase/exchanger/internal/rates/usecase"
)

type MockLogger struct{}

func (*MockLogger) Info(string) {}

type MockRepo struct{}

func (m MockRepo) Get(ctx context.Context, curr string) (*domain.Rate, error) {
	if curr == "RUB" {
		return &domain.Rate{"RUB", 1}, nil
	}
	if curr == "USD" {
		return &domain.Rate{"USD", 100}, nil
	}
	return nil, repository.ErrNotFound
}

func (m MockRepo) All(ctx context.Context) ([]domain.Rate, error) {
	return []domain.Rate{{"RUB", 1}, {"USD", 100}}, nil
}

func (m MockRepo) Save(ctx context.Context, rate domain.Rate) error {
	return nil
}

func TestUseCase_Get(t *testing.T) {
	tests := []struct {
		name     string
		curr     string
		wantRate *domain.Rate
		wantErr  error
	}{
		{
			"get existing rate",
			"RUB",
			&domain.Rate{"RUB", 1},
			nil,
		},
		{
			"get not existing rate",
			"GRB",
			nil,
			repository.ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uc := usecase.NewUseCase(&MockRepo{}, &MockLogger{})

			rate, err := uc.Get(context.Background(), test.curr)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(rate, test.wantRate) {
				t.Errorf("got rate %v, want %v", rate, test.wantRate)
			}
		})
	}
}

func TestUseCase_All(t *testing.T) {
	tests := []struct {
		name      string
		wantRates map[string]float32
		wantErr   error
	}{
		{
			"get all",
			map[string]float32{"RUB": 1, "USD": 100},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uc := usecase.NewUseCase(&MockRepo{}, &MockLogger{})

			rate, err := uc.All(context.Background())

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(rate, test.wantRates) {
				t.Errorf("got rate %v, want %v", rate, test.wantRates)
			}
		})
	}
}
