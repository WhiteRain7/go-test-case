package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/WhiteRain7/GoTestCase/pkg/utils/coded-errors"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/domain"
	"github.com/WhiteRain7/GoTestCase/wallets/pkg/decoder"
)

type Handler struct {
	uc     WalletUseCase
	logger Logger
}

type Logger interface {
	Info(string)
}

type WalletUseCase interface {
	Get(ctx context.Context, userId int) (*domain.WalletBalance, error)
	Deposit(ctx context.Context, userId int, curr string, amount float32) (*domain.WalletBalance, error)
	Withdraw(ctx context.Context, userId int, curr string, amount float32) (*domain.WalletBalance, error)
	Exchange(ctx context.Context, userId int, fromCurr string, toCurr string, amount float32) (*domain.WalletBalance, float32, error)
}

func NewHandler(uc WalletUseCase, logger Logger) *Handler {
	logger.Info("Handler init - Wallets")
	return &Handler{uc, logger}
}

type GetResponse struct {
	Balance map[string]float32 `json:"balance"`
}

// Get
// @Summary      Get wallet balance
// @Description  Returns actual wallet balance of authorized user.
// @Tags         Wallets
// @Accept       json
// @Produce      json
// @Success      200 {object} GetResponse
// @Failure      401 {string} string
// @Router       /api/v1/balance [get]
// @Security     BearerAuth
func (h *Handler) Get(res http.ResponseWriter, req *http.Request) {
	balance, err := h.uc.Get(req.Context(), req.Context().Value("user_id").(int))
	if err != nil {
		codederrors.ErrorResponse(err, res)
	} else {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		json.NewEncoder(res).Encode(GetResponse{balance.Balance})
	}
}

// ChangeRequest
type ChangeRequest struct {
	Amount   float32 `json:"amount"`
	Currency string  `json:"currency"`
}

// ChangeResponse
type ChangeResponse struct {
	Message    string             `json:"message"`
	NewBalance map[string]float32 `json:"new_balance"`
}

// Deposit
// @Summary      Add sum to wallet
// @Description  Adds specified amount in currency to the wallet of authorized user.
// @Tags         Wallets
// @Accept       json
// @Produce      json
// @Param        request body ChangeRequest true "Body"
// @Success      200 {object} ChangeResponse
// @Failure      400 {string} string "Negative amount, use currencies not exists"
// @Failure      401 {string} string
// @Router       /api/v1/wallet/deposit [post]
// @Security     BearerAuth
func (h *Handler) Deposit(res http.ResponseWriter, req *http.Request) {
	body := &ChangeRequest{}
	if err := decoder.Decode(req, body); err != nil {
		codederrors.ErrorResponse(err, res)
	}

	balance, err := h.uc.Deposit(req.Context(), req.Context().Value("user_id").(int), body.Currency, body.Amount)
	if err != nil {
		codederrors.ErrorResponse(err, res)
	} else {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		json.NewEncoder(res).Encode(ChangeResponse{"Account topped up successfully", balance.Balance})
	}
}

// Withdraw
// @Summary      Subtract sum from wallet
// @Description  Substructs specified amount in currency from the wallet of authorized user, if there's enough funds.
// @Tags         Wallets
// @Accept       json
// @Produce      json
// @Param        request body ChangeRequest true "Body"
// @Success      200 {object} ChangeResponse
// @Failure      400 {string} string "Negative amount, use currencies not exists, not enough funds"
// @Failure      401 {string} string
// @Router       /api/v1/wallet/withdraw [post]
// @Security     BearerAuth
func (h *Handler) Withdraw(res http.ResponseWriter, req *http.Request) {
	body := &ChangeRequest{}
	if err := decoder.Decode(req, body); err != nil {
		codederrors.ErrorResponse(err, res)
	}

	balance, err := h.uc.Withdraw(req.Context(), req.Context().Value("user_id").(int), body.Currency, body.Amount)
	if err != nil {
		codederrors.ErrorResponse(err, res)
	} else {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		json.NewEncoder(res).Encode(ChangeResponse{"Withdraw successful", balance.Balance})
	}
}

// ExchangeRequest
type ExchangeRequest struct {
	FromCurr string  `json:"from_currency"`
	ToCurr   string  `json:"to_currency"`
	Amount   float32 `json:"amount"`
}

// ExchangeResponse
type ExchangeResponse struct {
	Message   string             `json:"message"`
	Exchanged float32            `json:"exchanged_amount"`
	Balance   map[string]float32 `json:"exchanged_balance"`
}

// Exchange
// @Summary      Exchange currencies
// @Description  Allows to convert some amount from one currency to another with current rate. If user asked for actual rates within some time before, used this cache for exchange.
// @Tags         Wallets
// @Accept       json
// @Produce      json
// @Param        request body ExchangeRequest true "Body"
// @Success      200 {object} ExchangeResponse
// @Failure      400 {string} string "Negative amount, use currencies not exists, not enough funds"
// @Failure      401 {string} string
// @Router       /api/v1/exchange [post]
// @Security     BearerAuth
func (h *Handler) Exchange(res http.ResponseWriter, req *http.Request) {
	body := &ExchangeRequest{}
	if err := decoder.Decode(req, body); err != nil {
		codederrors.ErrorResponse(err, res)
	}

	balance, exchanged, err := h.uc.Exchange(req.Context(), req.Context().Value("user_id").(int), body.FromCurr, body.ToCurr, body.Amount)
	if err != nil {
		codederrors.ErrorResponse(err, res)
	} else {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		json.NewEncoder(res).Encode(ExchangeResponse{"Exchange successful", exchanged, balance.Balance})
	}
}
