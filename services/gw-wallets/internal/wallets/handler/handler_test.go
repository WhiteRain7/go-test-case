package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	codederrors "github.com/WhiteRain7/GoTestCase/pkg/utils/coded-errors"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/domain"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/handler"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/repository"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/usecase"
)

type mockUseCase struct{}

func (*mockUseCase) Get(ctx context.Context, userId int) (*domain.WalletBalance, error) {
	if userId == 1 {
		return &domain.WalletBalance{"abc", map[string]float32{"RUB": 100}}, nil
	}
	return nil, repository.ErrNotFound
}

func (uc *mockUseCase) Deposit(ctx context.Context, userId int, curr string, amount float32) (*domain.WalletBalance, error) {
	if curr != "RUB" && curr != "USD" {
		return nil, codederrors.New(400, usecase.ErrCurrencyNotSupported, "Currency %s is not supported", curr)
	}
	if amount <= 0 {
		return nil, codederrors.New(400, usecase.ErrNegativeAmount, "Amount must be positive, got %f", amount)
	}

	balance, err := uc.Get(ctx, userId)
	if err != nil {
		return nil, err
	}

	if was, ok := balance.Balance[curr]; ok {
		balance.Balance[curr] = was + amount
	} else {
		balance.Balance[curr] = amount
	}
	return balance, nil
}

func (uc *mockUseCase) Withdraw(ctx context.Context, userId int, curr string, amount float32) (*domain.WalletBalance, error) {
	if curr != "RUB" && curr != "USD" {
		return nil, codederrors.New(400, usecase.ErrCurrencyNotSupported, "Currency %s is not supported", curr)
	}
	if amount <= 0 {
		return nil, codederrors.New(400, usecase.ErrNegativeAmount, "Amount must be positive, got %f", amount)
	}

	balance, err := uc.Get(ctx, userId)
	if err != nil {
		return nil, err
	}

	if was, ok := balance.Balance[curr]; ok {
		if was < amount {
			return nil, codederrors.New(400, usecase.ErrNotEnoughFunds, "Not enough funds to withdraw %f %s", amount, curr)
		}
		balance.Balance[curr] = was - amount
	} else {
		return nil, codederrors.New(400, usecase.ErrNotEnoughFunds, "Not enough funds to withdraw %f %s", amount, curr)
	}
	return balance, nil
}

func (uc *mockUseCase) Exchange(ctx context.Context, userId int, fromCurr, toCurr string, amount float32) (*domain.WalletBalance, float32, error) {
	if fromCurr == toCurr {
		return nil, 0, codederrors.New(http.StatusBadRequest, usecase.ErrInvalidExchange, "Invalid exchange %s to %s", fromCurr, toCurr)
	}
	if fromCurr != "RUB" && fromCurr != "USD" {
		return nil, 0, codederrors.New(400, usecase.ErrCurrencyNotSupported, "Currency %s is not supported", fromCurr)
	}
	if toCurr != "RUB" && toCurr != "USD" {
		return nil, 0, codederrors.New(400, usecase.ErrCurrencyNotSupported, "Currency %s is not supported", toCurr)
	}
	if amount <= 0 {
		return nil, 0, codederrors.New(400, usecase.ErrNegativeAmount, "Amount must be positive, got %f", amount)
	}

	balance, err := uc.Get(ctx, userId)
	if err != nil {
		return nil, 0, err
	}

	rate := map[string]float32{"RUB": 1, "USD": 100}

	if was, ok := balance.Balance[fromCurr]; ok {
		if was < amount {
			return nil, 0, codederrors.New(400, usecase.ErrNotEnoughFunds, "Not enough funds to exchange %f %s", amount, fromCurr)
		}
		balance.Balance[fromCurr] = was - amount
	} else {
		return nil, 0, codederrors.New(400, usecase.ErrNotEnoughFunds, "Not enough funds to exchange %f %s", amount, fromCurr)
	}

	exchanged := amount * (rate[fromCurr] / rate[toCurr])

	if was, ok := balance.Balance[toCurr]; ok {
		balance.Balance[toCurr] = was + exchanged
	} else {
		balance.Balance[toCurr] = exchanged
	}

	return balance, exchanged, nil
}

type MockLogger struct{}

func (*MockLogger) Info(string) {}

func toRes(v any) string {
	res, _ := json.Marshal(v)
	return string(res)
}

func TestHandler_Get(t *testing.T) {
	tests := []struct {
		name       string
		userId     int
		wantStatus int
		wantBody   string
		wantPanic  bool
	}{
		{
			"get balance",
			1,
			200,
			toRes(handler.GetResponse{map[string]float32{"RUB": 100}}),
			false,
		},
		{
			"user not exists",
			2,
			500,
			"Not found",
			false,
		},
		{
			"not authenticated",
			-1,
			200,
			"",
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := handler.NewHandler(&mockUseCase{}, &MockLogger{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/balance", bytes.NewBufferString(""))
			rec := httptest.NewRecorder()

			if test.userId != -1 {
				req = req.WithContext(context.WithValue(req.Context(), "user_id", test.userId))
			}

			defer func() {
				if err := recover(); err != nil {
					if test.wantPanic {
						return
					}

					t.Errorf("want stable, got panic")
				}
			}()

			h.Get(rec, req)

			if test.wantPanic {
				t.Errorf("want panic, got stable")
			}

			if rec.Code != test.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", rec.Code, test.wantStatus, rec.Body.String())
			}

			result := strings.Trim(rec.Body.String(), " \n\t")
			if strings.Compare(result, test.wantBody) != 0 {
				t.Errorf("got body %s, want %s", result, test.wantBody)
			}
		})
	}
}

func TestHandler_Deposit(t *testing.T) {
	tests := []struct {
		name       string
		userId     int
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			"deposit",
			1,
			"{\"currency\":\"RUB\",\"amount\":100}",
			200,
			toRes(handler.ChangeResponse{"Account topped up successfully", map[string]float32{"RUB": 200}}),
		},
		{
			"negative amount",
			1,
			"{\"currency\":\"RUB\",\"amount\":-100}",
			400,
			"Amount must be positive, got -100.000000",
		},
		{
			"invalid currency",
			1,
			"{\"currency\":\"GRB\",\"amount\":100}",
			400,
			"Currency GRB is not supported",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := handler.NewHandler(&mockUseCase{}, &MockLogger{})

			req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet/deposit", bytes.NewBufferString(test.body))
			rec := httptest.NewRecorder()

			req = req.WithContext(context.WithValue(req.Context(), "user_id", test.userId))

			h.Deposit(rec, req)

			if rec.Code != test.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", rec.Code, test.wantStatus, rec.Body.String())
			}

			result := strings.Trim(rec.Body.String(), " \n\t")
			if strings.Compare(result, test.wantBody) != 0 {
				t.Errorf("got body %s, want %s", result, test.wantBody)
			}
		})
	}
}

func TestHandler_Withdraw(t *testing.T) {
	tests := []struct {
		name       string
		userId     int
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			"withdraw",
			1,
			"{\"currency\":\"RUB\",\"amount\":100}",
			200,
			toRes(handler.ChangeResponse{"Withdraw successful", map[string]float32{"RUB": 0}}),
		},
		{
			"negative amount",
			1,
			"{\"currency\":\"RUB\",\"amount\":-100}",
			400,
			"Amount must be positive, got -100.000000",
		},
		{
			"invalid currency",
			1,
			"{\"currency\":\"GRB\",\"amount\":100}",
			400,
			"Currency GRB is not supported",
		},
		{
			"not enough funds",
			1,
			"{\"currency\":\"RUB\",\"amount\":1000}",
			400,
			"Not enough funds to withdraw 1000.000000 RUB",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := handler.NewHandler(&mockUseCase{}, &MockLogger{})

			req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet/deposit", bytes.NewBufferString(test.body))
			rec := httptest.NewRecorder()

			req = req.WithContext(context.WithValue(req.Context(), "user_id", test.userId))

			h.Withdraw(rec, req)

			if rec.Code != test.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", rec.Code, test.wantStatus, rec.Body.String())
			}

			result := strings.Trim(rec.Body.String(), " \n\t")
			if strings.Compare(result, test.wantBody) != 0 {
				t.Errorf("got body %s, want %s", result, test.wantBody)
			}
		})
	}
}

func TestHandler_Exchange(t *testing.T) {
	tests := []struct {
		name       string
		userId     int
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			"exchange",
			1,
			"{\"from_currency\":\"RUB\",\"to_currency\":\"USD\",\"amount\":100}",
			200,
			toRes(handler.ExchangeResponse{"Exchange successful", 1, map[string]float32{"RUB": 0, "USD": 1}}),
		},
		{
			"negative amount",
			1,
			"{\"from_currency\":\"RUB\",\"to_currency\":\"USD\",\"amount\":-100}",
			400,
			"Amount must be positive, got -100.000000",
		},
		{
			"invalid currency 1",
			1,
			"{\"from_currency\":\"GRB\",\"to_currency\":\"USD\",\"amount\":100}",
			400,
			"Currency GRB is not supported",
		},
		{
			"invalid currency 2",
			1,
			"{\"from_currency\":\"RUB\",\"to_currency\":\"GRB\",\"amount\":100}",
			400,
			"Currency GRB is not supported",
		},
		{
			"not enough funds",
			1,
			"{\"from_currency\":\"RUB\",\"to_currency\":\"USD\",\"amount\":1000}",
			400,
			"Not enough funds to exchange 1000.000000 RUB",
		},
		{
			"invalid exchange",
			1,
			"{\"from_currency\":\"RUB\",\"to_currency\":\"RUB\",\"amount\":100}",
			400,
			"Invalid exchange RUB to RUB",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := handler.NewHandler(&mockUseCase{}, &MockLogger{})

			req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet/deposit", bytes.NewBufferString(test.body))
			rec := httptest.NewRecorder()

			req = req.WithContext(context.WithValue(req.Context(), "user_id", test.userId))

			h.Exchange(rec, req)

			if rec.Code != test.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", rec.Code, test.wantStatus, rec.Body.String())
			}

			result := strings.Trim(rec.Body.String(), " \n\t")
			if strings.Compare(result, test.wantBody) != 0 {
				t.Errorf("got body %s, want %s", result, test.wantBody)
			}
		})
	}
}
