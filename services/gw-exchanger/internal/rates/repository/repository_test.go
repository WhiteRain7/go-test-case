package repository_test

import (
	"context"
	"database/sql/driver"
	"errors"
	"reflect"
	"testing"

	"github.com/78bits/go-sqlmock-sqlx"
	"github.com/WhiteRain7/GoTestCase/exchanger/internal/rates/domain"
	"github.com/WhiteRain7/GoTestCase/exchanger/internal/rates/repository"
)

type MockLogger struct{}

func (*MockLogger) Info(string) {}

func TestRepository_Get(t *testing.T) {
	tests := []struct {
		name     string
		curr     string
		rows     *sqlmock.Rows
		wantRate *domain.Rate
		wantErr  error
	}{
		{
			"get",
			"RUB",
			sqlmock.NewRows([]string{"rate"}).AddRow(1),
			&domain.Rate{"RUB", 1},
			nil,
		},
		{
			"not exists",
			"GRB",
			sqlmock.NewRows([]string{"rate"}),
			nil,
			repository.ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.Newx()
			if err != nil {
				t.Fatalf("got error '%s' when tried to mock sql", err)
			}
			defer db.Close()

			mock.ExpectQuery("SELECT rate FROM Rates").
				WithArgs(test.curr).
				WillReturnRows(test.rows)

			repo := repository.NewRepository(db, &MockLogger{})

			rate, err := repo.Get(context.Background(), test.curr)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(rate, test.wantRate) {
				t.Errorf("got rate %v, want %v", rate, test.wantRate)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("got unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_All(t *testing.T) {
	tests := []struct {
		name      string
		rows      *sqlmock.Rows
		wantRates []domain.Rate
		wantErr   error
	}{
		{
			"all",
			sqlmock.NewRows([]string{"currency", "rate"}).AddRow("RUB", 1).AddRow("USD", 100),
			[]domain.Rate{{"RUB", 1}, {"USD", 100}},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.Newx()
			if err != nil {
				t.Fatalf("got error '%s' when tried to mock sql", err)
			}
			defer db.Close()

			mock.ExpectQuery("SELECT currency, rate FROM Rates").WillReturnRows(test.rows)

			repo := repository.NewRepository(db, &MockLogger{})

			rates, err := repo.All(context.Background())

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(rates, test.wantRates) {
				t.Errorf("got rates %v, want %v", rates, test.wantRates)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("got unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_Save(t *testing.T) {
	tests := []struct {
		name    string
		rate    domain.Rate
		result  driver.Result
		wantErr error
	}{
		{
			"save",
			domain.Rate{"RUB", 1},
			sqlmock.NewResult(1, 1),
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.Newx()
			if err != nil {
				t.Fatalf("got error '%s' when tried to mock sql", err)
			}
			defer db.Close()

			mock.ExpectExec("INSERT INTO rates").
				WithArgs(test.rate.Currency, test.rate.Rate).
				WillReturnResult(test.result)

			repo := repository.NewRepository(db, &MockLogger{})

			err = repo.Save(context.Background(), test.rate)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("got unfulfilled expectations: %s", err)
			}
		})
	}
}
