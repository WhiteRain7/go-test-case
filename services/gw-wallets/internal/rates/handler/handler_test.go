package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WhiteRain7/GoTestCase/wallets/internal/rates/handler"
)

type MockLogger struct{}

func (MockLogger) Info(string)  {}
func (MockLogger) Error(string) {}

type MockUseCase struct{}

func (m MockUseCase) All(ctx context.Context) (map[string]float32, error) {
	return map[string]float32{"RUB": 1, "USD": 100}, nil
}

func toRes(v any) string {
	res, _ := json.Marshal(v)
	return string(res)
}

func TestAll(t *testing.T) {
	tests := []struct {
		name       string
		wantStatus int
		wantBody   string
	}{
		{
			"all",
			200,
			toRes(handler.GetRatesResponse{map[string]float32{"RUB": 1, "USD": 100}}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := handler.NewHandler(&MockUseCase{}, &MockLogger{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/exchange/rates", bytes.NewBufferString(""))
			rec := httptest.NewRecorder()

			h.All(rec, req)

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
