package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/WhiteRain7/GoTestCase/analytics/internal/transfers/domain"
	"github.com/WhiteRain7/GoTestCase/analytics/internal/transfers/repository"
	mockhouse "github.com/srikanthccv/ClickHouse-go-mock"
)

type MockLogger struct{}

func (*MockLogger) Info(string) {}

func TestRepository_Add(t *testing.T) {
	tests := []struct {
		name     string
		transfer domain.Transfer
		wantErr  error
	}{
		{
			"name",
			domain.Transfer{1, "DEPOSIT", 30, "RUB", "11111111111"},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := mockhouse.NewClickHouseNative(nil)
			if err != nil {
				t.Fatalf("got error '%s' when tried to mock sql", err)
			}

			mock.ExpectPrepareBatch("INSERT INTO analytics.transfers").ExpectAppend()

			repo := repository.NewRepository(mock, &MockLogger{})

			err = repo.Add(context.Background(), test.transfer)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("got unfulfilled expectations: %s", err)
			}
		})
	}
}
