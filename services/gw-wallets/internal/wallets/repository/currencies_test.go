package repository_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/78bits/go-sqlmock-sqlx"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/repository"
)

func TestRepository_GetCurrency(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		rows    *sqlmock.Rows
		wantKey string
		wantErr error
	}{
		{
			"Get key",
			"RUB",
			sqlmock.NewRows([]string{"key"}).AddRow("RUB"),
			"RUB",
			nil,
		},
		{
			"Get key fails",
			"RUB",
			sqlmock.NewRows([]string{"key"}),
			"",
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

			mock.ExpectQuery("SELECT key FROM currency_types").WithArgs(test.key).WillReturnRows(test.rows)

			repo := repository.NewRepository(&MockDBGetter{db}, &MockLogger{})

			key, err := repo.GetCurrency(context.Background(), test.key)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(key, test.wantKey) {
				t.Errorf("got key %v, want %v", key, test.wantKey)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("got unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_AllCurrencies(t *testing.T) {
	tests := []struct {
		name     string
		rows     *sqlmock.Rows
		wantKeys []string
		wantErr  error
	}{
		{
			"Get keys",
			sqlmock.NewRows([]string{"key"}).AddRow("RUB").AddRow("EUR").AddRow("USD"),
			[]string{"RUB", "EUR", "USD"},
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

			mock.ExpectQuery("SELECT key FROM currency_types").WillReturnRows(test.rows)

			repo := repository.NewRepository(&MockDBGetter{db}, &MockLogger{})

			key, err := repo.AllCurrencies(context.Background())

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(key, test.wantKeys) {
				t.Errorf("got key %v, want %v", key, test.wantKeys)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("got unfulfilled expectations: %s", err)
			}
		})
	}
}
