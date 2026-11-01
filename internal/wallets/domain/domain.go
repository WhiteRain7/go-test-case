package domain

import (
	"errors"
	"net/http"

	"github.com/WhiteRain7/GoTestCase/internal/shared/infrastructure/apperr"
)

type UUID = string

var (
	ErrNotEnough      = errors.New("Not enough money")
	ErrNegativeAmount = errors.New("Negative amount")
)

type Wallet struct {
	ID     UUID `json:"id"`
	Amount int  `json:"amount"`
}

func (w *Wallet) Deposit(amount int) error {
	if amount < 0 {
		return apperr.New(http.StatusBadRequest, ErrNegativeAmount, "Amount to deposit must be 0 or greater, got %d", amount)
	}
	w.Amount += amount
	return nil
}

func (w *Wallet) Withdraw(amount int) error {
	if amount < 0 {
		return apperr.New(http.StatusBadRequest, ErrNegativeAmount, "Amount to withdraw must be 0 or greater, got %d", amount)
	}
	if w.Amount < amount {
		return apperr.New(http.StatusBadRequest, ErrNotEnough, "Not enough money to perform withdraw")
	}
	w.Amount -= amount
	return nil
}
