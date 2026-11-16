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
	"github.com/WhiteRain7/GoTestCase/wallets/internal/users/domain"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/users/handler"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/users/repository"
	"golang.org/x/crypto/bcrypt"
)

type MockUseCase struct{}

func (m *MockUseCase) Create(ctx context.Context, username string, password string, email string) (*domain.User, error) {
	if username == "user" {
		return nil, codederrors.New(400, repository.ErrUsernameExists, "Username already exists")
	}
	if email == "mail@example.com" {
		return nil, codederrors.New(400, repository.ErrEmailExists, "Email already exists")
	}
	return &domain.User{2, username, email}, nil
}

func (m *MockUseCase) Validate(ctx context.Context, username string, password string) error {
	if username != "user" {
		return codederrors.New(404, repository.ErrNotFound, "User not found")
	}
	if password != "123" {
		return codederrors.New(400, bcrypt.ErrMismatchedHashAndPassword, "Password mismatch")
	}
	return nil
}

func (m *MockUseCase) GetID(ctx context.Context, username string) (int, error) {
	if username == "user" {
		return 1, nil
	}
	return 0, codederrors.New(404, repository.ErrNotFound, "User not found")
}

type MockAuth struct{}

func (m *MockAuth) GenerateJWT(id int, username string) (string, error) {
	return "token", nil
}

type MockLogger struct{}

func (*MockLogger) Info(string) {}

func toRes(v any) string {
	res, _ := json.Marshal(v)
	return string(res)
}

func TestHandler_SignUp(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			"sign up",
			toRes(handler.SignUpRequest{"newuser", "123", "newmail@example.com"}),
			201,
			toRes(handler.SignUpResponse{"User registered successfully"}),
		},
		{
			"username already exists",
			toRes(handler.SignUpRequest{"user", "123", "newmail@example.com"}),
			400,
			"Username already exists",
		},
		{
			"email already exists",
			toRes(handler.SignUpRequest{"newuser", "123", "mail@example.com"}),
			400,
			"Email already exists",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := handler.NewHandler(&MockUseCase{}, &MockAuth{}, &MockLogger{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/register", bytes.NewBufferString(test.body))
			rec := httptest.NewRecorder()

			h.SignUp(rec, req)

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

func TestHandler_SignIn(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			"sign in",
			toRes(handler.SignInRequest{"user", "123"}),
			200,
			toRes(handler.SignInResponse{"token"}),
		},
		{
			"username not exists",
			toRes(handler.SignInRequest{"username", "123"}),
			404,
			"User not found",
		},
		{
			"pwd mismatch",
			toRes(handler.SignInRequest{"user", "1234"}),
			400,
			"Password mismatch",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := handler.NewHandler(&MockUseCase{}, &MockAuth{}, &MockLogger{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/login", bytes.NewBufferString(test.body))
			rec := httptest.NewRecorder()

			h.SignIn(rec, req)

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
