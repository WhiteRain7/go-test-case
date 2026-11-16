package repository_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/78bits/go-sqlmock-sqlx"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/domain"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/repository"
)

func TestRepository_Create(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		rows       *sqlmock.Rows
		wantWallet *domain.Wallet
		wantErr    error
	}{
		{
			"create new wallet",
			1,
			sqlmock.NewRows([]string{"id"}).AddRow("abc"),
			&domain.Wallet{"abc", 1},
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

			mock.ExpectQuery("INSERT INTO wallets").
				WithArgs(test.id).
				WillReturnRows(test.rows)

			repo := repository.NewRepository(&MockDBGetter{db}, &MockLogger{})

			wallet, err := repo.Create(context.Background(), test.id)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(wallet, test.wantWallet) {
				t.Errorf("got wallet %v, want %v", wallet, test.wantWallet)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("got unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_Get(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		rows       *sqlmock.Rows
		wantWallet *domain.Wallet
		wantErr    error
	}{
		{
			"get wallet",
			1,
			sqlmock.NewRows([]string{"id"}).AddRow("abc"),
			&domain.Wallet{"abc", 1},
			nil,
		},
		{
			"get no wallet",
			2,
			sqlmock.NewRows([]string{"id"}),
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

			mock.ExpectQuery("SELECT id FROM wallets").
				WithArgs(test.id).
				WillReturnRows(test.rows)

			repo := repository.NewRepository(&MockDBGetter{db}, &MockLogger{})

			wallet, err := repo.Get(context.Background(), test.id)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(wallet, test.wantWallet) {
				t.Errorf("got wallet %v, want %v", wallet, test.wantWallet)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("got unfulfilled expectations: %s", err)
			}
		})
	}
}
