package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/WhiteRain7/GoTestCase/internal/shared/infrastructure/apperr"
	"github.com/WhiteRain7/GoTestCase/internal/wallets/domain"
)

type Handler struct {
	uc     WalletUseCase
	logger Logger
}

type OperationType string

const (
	OperationDeposit  OperationType = "DEPOSIT"
	OperationWithdraw OperationType = "WITHDRAW"
)

func (o OperationType) Valid() bool {
	return o == OperationDeposit || o == OperationWithdraw
}

type WalletOperation struct {
	WalletId      domain.UUID   `json:"walletId"`
	OperationType OperationType `json:"operationType"`
	Amount        int           `json:"amount"`
}

type WalletUseCase interface {
	Create(ctx context.Context) (*domain.Wallet, error)
	Get(ctx context.Context, id domain.UUID) (*domain.Wallet, error)
	Update(ctx context.Context, op WalletOperation) (*domain.Wallet, error)
}

type Logger interface {
	Info(s string)
	Error(s string)
}

func NewHandler(uc WalletUseCase, logger Logger) *Handler {
	logger.Info("Handler init - Wallets")
	return &Handler{uc, logger}
}

func (h *Handler) Get(res http.ResponseWriter, req *http.Request) {
	wallet, err := h.uc.Get(req.Context(), req.PathValue("id"))
	if err != nil {
		code := apperr.ErrorResponse(err, res)
		h.log(req, code)
	} else {
		h.log(req, 200)
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		json.NewEncoder(res).Encode(wallet)
	}
}

func (h *Handler) Create(res http.ResponseWriter, req *http.Request) {
	wallet, err := h.uc.Create(req.Context())
	if err != nil {
		code := apperr.ErrorResponse(err, res)
		h.log(req, code)
	} else {
		h.log(req, 201)
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(201)
		json.NewEncoder(res).Encode(wallet)
	}
}

func (h *Handler) Update(res http.ResponseWriter, req *http.Request) {
	var op WalletOperation
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&op); err != nil {
		code := apperr.ErrorResponse(apperr.New(400, err, err.Error()), res)
		h.log(req, code)
	}

	wallet, err := h.uc.Update(req.Context(), op)
	if err != nil {
		code := apperr.ErrorResponse(err, res)
		h.log(req, code)
	} else {
		h.log(req, 200)
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		json.NewEncoder(res).Encode(wallet)
	}
}

func (h *Handler) log(req *http.Request, code int) {
	if code < 400 {
		h.logger.Info(fmt.Sprintf("[%d] HTTP %s %s", code, req.Method, req.RequestURI))
	} else if code < 500 {
		h.logger.Info(fmt.Sprintf("[%d] HTTP %s %s", code, req.Method, req.RequestURI))
	} else {
		h.logger.Error(fmt.Sprintf("[%d] HTTP %s %s", code, req.Method, req.RequestURI))
	}
}
