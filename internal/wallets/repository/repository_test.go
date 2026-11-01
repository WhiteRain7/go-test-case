package repository_test

import (
	"context"
	"database/sql/driver"
	"errors"
	"reflect"
	"testing"

	"github.com/78bits/go-sqlmock-sqlx"
	"github.com/WhiteRain7/GoTestCase/internal/wallets/domain"
	"github.com/WhiteRain7/GoTestCase/internal/wallets/repository"
)

type MockLogger struct{}

func (*MockLogger) Info(string) {}

func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		rows       *sqlmock.Rows
		wantWallet *domain.Wallet
		wantErr    error
	}{
		{"get wallet", "abc", sqlmock.NewRows([]string{"amount"}).AddRow(100), &domain.Wallet{"abc", 100}, nil},
		{"get no wallet", "bcd", sqlmock.NewRows([]string{"amount"}), nil, repository.ErrNotFound},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.Newx()
			if err != nil {
				t.Fatalf("got error '%s' when tried to mock sql", err)
			}
			defer db.Close()

			mock.ExpectQuery("SELECT amount FROM wallets").
				WithArgs(test.id).
				WillReturnRows(test.rows)

			repo := repository.NewRepository(db, &MockLogger{})

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

func TestCreate(t *testing.T) {
	tests := []struct {
		name       string
		Id         string
		rows       *sqlmock.Rows
		wantWallet *domain.Wallet
		wantErr    error
	}{
		{"create new wallet", "abc", sqlmock.NewRows([]string{"id"}).AddRow("abc"), &domain.Wallet{"abc", 0}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.Newx()
			if err != nil {
				t.Fatalf("got error '%s' when tried to mock sql", err)
			}
			defer db.Close()

			mock.ExpectQuery("INSERT INTO wallets").WillReturnRows(test.rows)

			repo := repository.NewRepository(db, &MockLogger{})

			wallet, err := repo.Create(context.Background())

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

func TestSave(t *testing.T) {
	tests := []struct {
		name    string
		wallet  *domain.Wallet
		result  driver.Result
		wantErr error
	}{
		{"save wallet", &domain.Wallet{"abc", 100}, sqlmock.NewResult(0, 1), nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.Newx()
			if err != nil {
				t.Fatalf("got error '%s' when tried to mock sql", err)
			}
			defer db.Close()

			mock.ExpectExec("UPDATE wallets").WillReturnResult(test.result)

			repo := repository.NewRepository(db, &MockLogger{})

			err = repo.Save(context.Background(), test.wallet)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("got unfulfilled expectations: %s", err)
			}
		})
	}
}
