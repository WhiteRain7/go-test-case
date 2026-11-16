package handler

import (
	"context"
	"encoding/json"

	"github.com/WhiteRain7/GoTestCase/analytics/internal/transfers/domain"
)

type Handler struct {
	uc UseCase
}

type UseCase interface {
	Add(ctx context.Context, transfer domain.Transfer) error
}

func NewHandler(uc UseCase) *Handler {
	return &Handler{uc}
}

func (h *Handler) Handle(ctx context.Context, record []byte) error {
	transfer := domain.Transfer{}
	err := json.Unmarshal(record, &transfer)
	if err != nil {
		return err
	}
	return h.uc.Add(ctx, transfer)
}
