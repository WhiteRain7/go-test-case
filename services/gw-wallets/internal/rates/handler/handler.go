package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/WhiteRain7/GoTestCase/pkg/utils/coded-errors"
)

type Handler struct {
	uc     RatesUseCase
	logger Logger
}

type Logger interface {
	Info(string)
}

type RatesUseCase interface {
	All(ctx context.Context) (map[string]float32, error)
}

func NewHandler(uc RatesUseCase, logger Logger) *Handler {
	logger.Info("Handler init - Rates")
	return &Handler{uc, logger}
}

type GetRatesResponse struct {
	Rates map[string]float32 `json:"rates"`
}

// All
// @Summary      Actual rates
// @Description  Returns actual rates from other service via gRPC. Rates cached for some time in case user will exchange currencies.
// @Tags         Rates
// @Accept       json
// @Produce      json
// @Success      200 {object} GetRatesResponse
// @Failure      401 {string} string
// @Router       /api/v1/exchange/rates [get]
// @Security     BearerAuth
func (h *Handler) All(res http.ResponseWriter, req *http.Request) {
	rates, err := h.uc.All(req.Context())
	if err != nil {
		codederrors.ErrorResponse(err, res)
	} else {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		json.NewEncoder(res).Encode(GetRatesResponse{rates})
	}
}
