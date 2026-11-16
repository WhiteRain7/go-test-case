package usecase_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/WhiteRain7/GoTestCase/wallets/internal/users/domain"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/users/repository"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/users/usecase"
	wallets "github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/domain"
	"golang.org/x/crypto/bcrypt"
)

type MockTXManager struct{}

func (*MockTXManager) BeginTx() error  { return nil }
func (*MockTXManager) Commit() error   { return nil }
func (*MockTXManager) Rollback() error { return nil }

type MockLogger struct{}

func (*MockLogger) Info(string) {}

type MockRepo struct {
	hashed string
}

func (*MockRepo) GetById(ctx context.Context, id int) (*domain.User, error) {
	if id == 1 {
		return &domain.User{1, "user", "mail@example.com"}, nil
	}
	return nil, repository.ErrNotFound
}

func (*MockRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if email == "mail@example.com" {
		return &domain.User{1, "user", "mail@example.com"}, nil
	}
	return nil, repository.ErrNotFound
}

func (*MockRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	if username == "user" {
		return &domain.User{1, "user", "mail@example.com"}, nil
	}
	return nil, repository.ErrNotFound
}

func (*MockRepo) GetCredentials(ctx context.Context, username string) (*domain.Credentials, error) {
	pwd_123, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	if username == "user" {
		return &domain.Credentials{"user", string(pwd_123)}, nil
	}
	return nil, repository.ErrNotFound
}

func (r *MockRepo) Create(ctx context.Context, username, password, email string) (*domain.User, error) {
	if username == "user" {
		return nil, repository.ErrUsernameExists
	}
	if email == "mail@example.com" {
		return nil, repository.ErrEmailExists
	}
	r.hashed = password
	return &domain.User{2, username, email}, nil
}

type MockWallets struct {
	created *wallets.Wallet
}

func (w *MockWallets) Create(ctx context.Context, userId int) (*wallets.Wallet, error) {
	wallet := &wallets.Wallet{"abc", userId}
	w.created = wallet
	return wallet, nil
}

func TestUseCase_Get(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantUser *domain.User
		wantErr  error
	}{
		{"get", "user", &domain.User{1, "user", "mail@example.com"}, nil},
		{"not found", "unknown", nil, repository.ErrNotFound},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uc := usecase.NewUseCase(&MockTXManager{}, &MockRepo{}, &MockWallets{}, &MockLogger{})

			user, err := uc.Get(context.Background(), test.username)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(user, test.wantUser) {
				t.Errorf("got user %v, want %v", user, test.wantUser)
			}
		})
	}
}

func TestUseCase_GetID(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantId   int
		wantErr  error
	}{
		{"get", "user", 1, nil},
		{"not found", "unknown", 0, repository.ErrNotFound},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uc := usecase.NewUseCase(&MockTXManager{}, &MockRepo{}, &MockWallets{}, &MockLogger{})

			id, err := uc.GetID(context.Background(), test.username)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(id, test.wantId) {
				t.Errorf("got id %v, want %v", id, test.wantId)
			}
		})
	}
}

func TestUseCase_Create(t *testing.T) {
	tests := []struct {
		name       string
		username   string
		password   string
		email      string
		wantUser   *domain.User
		wantWallet *wallets.Wallet
		wantErr    error
	}{
		{
			"create",
			"newuser",
			"123",
			"newmail@example.com",
			&domain.User{2, "newuser", "newmail@example.com"},
			&wallets.Wallet{"abc", 2},
			nil,
		},
		{
			"user exists",
			"user",
			"123",
			"newmail@example.com",
			nil,
			nil,
			repository.ErrUsernameExists,
		},
		{
			"email exists",
			"newuser",
			"123",
			"mail@example.com",
			nil,
			nil,
			repository.ErrEmailExists,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := &MockRepo{}
			mockWallets := &MockWallets{}

			uc := usecase.NewUseCase(&MockTXManager{}, mockRepo, mockWallets, &MockLogger{})

			user, err := uc.Create(context.Background(), test.username, test.password, test.email)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(user, test.wantUser) {
				t.Errorf("got user %v, want %v", user, test.wantUser)
			}

			if !reflect.DeepEqual(mockWallets.created, test.wantWallet) {
				t.Errorf("got wallet %v, want %v", mockWallets.created, test.wantWallet)
			}

			if test.wantUser != nil {
				if err := bcrypt.CompareHashAndPassword([]byte(mockRepo.hashed), []byte(test.password)); err != nil {
					t.Errorf("got hash error %v, want nil", err)
				}
			}
		})
	}
}

func TestUseCase_Validate(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantErr  error
	}{
		{
			"validate user",
			"user",
			"123",
			nil,
		},
		{
			"invalid pwd",
			"user",
			"1234",
			bcrypt.ErrMismatchedHashAndPassword,
		},
		{
			"invalid user",
			"username",
			"123",
			repository.ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := &MockRepo{}
			mockWallets := &MockWallets{}

			uc := usecase.NewUseCase(&MockTXManager{}, mockRepo, mockWallets, &MockLogger{})

			err := uc.Validate(context.Background(), test.username, test.password)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}
		})
	}
}
