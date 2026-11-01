package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WhiteRain7/GoTestCase/internal/shared/infrastructure/apperr"
	"github.com/WhiteRain7/GoTestCase/internal/wallets/domain"
	"github.com/WhiteRain7/GoTestCase/internal/wallets/handler"
	"github.com/WhiteRain7/GoTestCase/internal/wallets/repository"
	"github.com/WhiteRain7/GoTestCase/internal/wallets/usecase"
)

type mockUseCase struct{}

func (uc *mockUseCase) Create(ctx context.Context) (*domain.Wallet, error) {
	return &domain.Wallet{ID: "abc", Amount: 0}, nil
}

func (uc *mockUseCase) Get(ctx context.Context, id domain.UUID) (*domain.Wallet, error) {
	wallets_ := []domain.Wallet{
		{ID: "abc", Amount: 0},
		{ID: "bcd", Amount: 200},
	}

	for _, wallet := range wallets_ {
		if wallet.ID == id {
			return &wallet, nil
		}
	}

	return nil, apperr.New(404, repository.ErrNotFound, "wallet doesn't exist")
}

func (uc *mockUseCase) Update(ctx context.Context, op handler.WalletOperation) (*domain.Wallet, error) {
	wallet, err := uc.Get(ctx, op.WalletId)
	if err != nil {
		return nil, err
	}

	var operr error

	switch op.OperationType {
	case handler.OperationDeposit:
		operr = wallet.Deposit(op.Amount)
	case handler.OperationWithdraw:
		operr = wallet.Withdraw(op.Amount)
	default:
		return nil, apperr.New(400, usecase.ErrInvalidOperation, "operation not supported")
	}

	if operr != nil {
		return nil, operr
	}

	return wallet, nil
}

type MockLogger struct{}

func (*MockLogger) Info(string)  {}
func (*MockLogger) Error(string) {}

func toRes(v any) string {
	res, _ := json.Marshal(v)
	return string(res)
}

func TestHandlerGet(t *testing.T) {
	tests := []struct {
		name       string
		reqID      string
		wantStatus int
		wantBody   string
	}{
		{"get wallet", "abc", 200, toRes(domain.Wallet{ID: "abc", Amount: 0})},
		{"get non-existing wallet", "aaa", 404, "wallet doesn't exist"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := handler.NewHandler(&mockUseCase{}, &MockLogger{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/{id}", bytes.NewBufferString(""))
			rec := httptest.NewRecorder()

			req.SetPathValue("id", test.reqID)

			h.Get(rec, req)

			if rec.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, test.wantStatus, rec.Body.String())
			}

			result := strings.Trim(rec.Body.String(), " \n\t")
			if strings.Compare(result, test.wantBody) != 0 {
				t.Fatalf("got body %s, want %s", result, test.wantBody)
			}
		})
	}
}

func TestHandlerCreate(t *testing.T) {
	tests := []struct {
		name       string
		wantStatus int
		wantBody   string
	}{
		{"create new wallet", 201, toRes(domain.Wallet{ID: "abc", Amount: 0})},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := handler.NewHandler(&mockUseCase{}, &MockLogger{})

			req := httptest.NewRequest(http.MethodPost, "/api/v1/wallets", bytes.NewBufferString(""))
			rec := httptest.NewRecorder()

			h.Create(rec, req)

			if rec.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, test.wantStatus, rec.Body.String())
			}

			result := strings.Trim(rec.Body.String(), " \n\t")
			if strings.Compare(result, test.wantBody) != 0 {
				t.Fatalf("got body %s, want %s", result, test.wantBody)
			}
		})
	}
}

func TestHandlerUpdate(t *testing.T) {
	tests := []struct {
		name       string
		reqOp      handler.WalletOperation
		wantStatus int
		wantBody   string
	}{
		{"deposit wallet", handler.WalletOperation{"abc", handler.OperationDeposit, 100}, 200, toRes(domain.Wallet{ID: "abc", Amount: 100})},
		{"withdraw wallet", handler.WalletOperation{"bcd", handler.OperationWithdraw, 100}, 200, toRes(domain.Wallet{ID: "bcd", Amount: 100})},
		{"unacceptable operation", handler.WalletOperation{"abc", "MULTIPLY", 0}, 400, "operation not supported"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := handler.NewHandler(&mockUseCase{}, &MockLogger{})

			body, _ := json.Marshal(test.reqOp)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()

			h.Update(rec, req)

			if rec.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, test.wantStatus, rec.Body.String())
			}

			result := strings.Trim(rec.Body.String(), " \n\t")
			if strings.Compare(result, test.wantBody) != 0 {
				t.Fatalf("got body %s, want %s", result, test.wantBody)
			}
		})
	}
}
