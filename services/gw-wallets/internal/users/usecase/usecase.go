package usecase

import (
	"context"
	"net/http"

	"github.com/WhiteRain7/GoTestCase/pkg/utils/coded-errors"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/users/domain"
	"golang.org/x/crypto/bcrypt"

	wallets "github.com/WhiteRain7/GoTestCase/wallets/internal/wallets/domain"
)

type UseCase struct {
	txm     TXManager
	repo    UserRepository
	wallets WalletsRepository
	logger  Logger
}

type TXManager interface {
	BeginTx() error
	Commit() error
	Rollback() error
}

type Logger interface {
	Info(string)
}

type UserRepository interface {
	Create(ctx context.Context, username string, password domain.Hashed, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetCredentials(ctx context.Context, username string) (*domain.Credentials, error)
}

type WalletsRepository interface {
	Create(ctx context.Context, userId int) (*wallets.Wallet, error)
}

func NewUseCase(txm TXManager, repo UserRepository, wallets WalletsRepository, logger Logger) *UseCase {
	logger.Info("Use case init - Wallets")
	return &UseCase{txm, repo, wallets, logger}
}

func (uc *UseCase) Get(ctx context.Context, username string) (*domain.User, error) {
	return uc.repo.GetByUsername(ctx, username)
}

func (uc *UseCase) GetID(ctx context.Context, username string) (int, error) {
	user, err := uc.repo.GetByUsername(ctx, username)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (uc *UseCase) Create(ctx context.Context, username string, password string, email string) (*domain.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, codederrors.Unknown(err)
	}

	uc.txm.BeginTx()
	defer uc.txm.Rollback()

	user, err := uc.repo.Create(ctx, username, string(hashed), email)
	if err != nil {
		return nil, err
	}

	_, err = uc.wallets.Create(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	uc.txm.Commit()

	return user, nil
}

func (uc *UseCase) Validate(ctx context.Context, username string, password string) error {
	credentials, err := uc.repo.GetCredentials(ctx, username)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(credentials.Password), []byte(password))
	if err != nil {
		return codederrors.New(http.StatusUnauthorized, err, "Incorrect username or password")
	}

	return nil
}
