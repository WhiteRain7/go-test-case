package repository_test

import (
	"context"
	"database/sql/driver"
	"errors"
	"reflect"
	"testing"

	"github.com/78bits/go-sqlmock-sqlx"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/domain"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/repository"
)

func TestRepository_GetBalance(t *testing.T) {
	tests := []struct {
		name              string
		id                string
		rows              *sqlmock.Rows
		wantWalletBalance *domain.WalletBalance
		wantErr           error
	}{
		{
			"get balance",
			"abc",
			sqlmock.NewRows([]string{"currency", "amount"}).AddRow("RUB", 100),
			&domain.WalletBalance{"abc", map[string]float32{"RUB": 1}},
			nil,
		},
		{
			"get balance with multiple currencies",
			"abc",
			sqlmock.NewRows([]string{"currency", "amount"}).AddRow("RUB", 100).AddRow("EUR", 100),
			&domain.WalletBalance{"abc", map[string]float32{"RUB": 1, "EUR": 1}},
			nil,
		},
		{
			"balance of unknown wallet",
			"abc",
			sqlmock.NewRows([]string{"currency", "amount"}),
			&domain.WalletBalance{"abc", map[string]float32{}},
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

			mock.ExpectQuery("SELECT currency, amount FROM wallet_currencies").
				WithArgs(test.id).
				WillReturnRows(test.rows)

			repo := repository.NewRepository(&MockDBGetter{db}, &MockLogger{})

			walletBalance, err := repo.GetBalance(context.Background(), test.id)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(walletBalance, test.wantWalletBalance) {
				t.Errorf("got walletBalance %v, want %v", walletBalance, test.wantWalletBalance)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("got unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_SetBalance(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		curr    string
		amount  float32
		result  driver.Result
		wantErr error
	}{
		{
			"set balance",
			"abc",
			"RUB",
			1,
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

			mock.ExpectExec("INSERT INTO wallet_currencies").
				WithArgs(test.id, test.curr, int(test.amount*100)).
				WillReturnResult(test.result)

			repo := repository.NewRepository(&MockDBGetter{db}, &MockLogger{})

			err = repo.SetBalance(context.Background(), test.id, test.curr, test.amount)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("got unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_UpdateBalance(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		curr    string
		amount  float32
		result  driver.Result
		wantErr error
	}{
		{
			"update balance",
			"abc",
			"RUB",
			1,
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

			mock.ExpectExec("UPDATE wallet_currencies").
				WithArgs(int(test.amount*100), test.id, test.curr).
				WillReturnResult(test.result)

			repo := repository.NewRepository(&MockDBGetter{db}, &MockLogger{})

			err = repo.UpdateBalance(context.Background(), test.id, test.curr, test.amount)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("got unfulfilled expectations: %s", err)
			}
		})
	}
}
