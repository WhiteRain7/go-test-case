package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/WhiteRain7/GoTestCase/analytics/internal/transfers/domain"
	"github.com/WhiteRain7/GoTestCase/analytics/internal/transfers/handler"
)

type MockLogger struct{}

func (*MockLogger) Info(string) {}

type MockUseCase struct{}

func (m MockUseCase) Add(ctx context.Context, transfer domain.Transfer) error {
	return nil
}

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name    string
		record  []byte
		wantErr error
	}{
		{
			"handler",
			[]byte(`{"user":1,"operation":"WITHDRAW","amount":30,"currency":"RUB","timestamp":"11111111"}`),
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := handler.NewHandler(&MockUseCase{})

			err := h.Handle(context.Background(), test.record)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}
		})
	}
}
