package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/WhiteRain7/GoTestCase/pkg/utils/coded-errors"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/users/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	ErrNotFound       = errors.New("User not found")
	ErrEmailExists    = errors.New("Email already exists")
	ErrUsernameExists = errors.New("Username already exists")
)

type Repository struct {
	db     DBGetter
	logger Logger
}

type DBGetter interface {
	Get() sqlx.ExtContext
}

type Logger interface {
	Info(string)
}

func NewRepository(db DBGetter, logger Logger) *Repository {
	logger.Info("New repository - Users")
	return &Repository{db, logger}
}

func (r *Repository) getBy[T string | int](ctx context.Context, field string, val T) (*domain.User, error) {
	r.logger.Info(fmt.Sprintf("DB request - Users - Get by %s with %v", field, val))
	db := r.db.Get()
	row := db.QueryRowxContext(ctx, fmt.Sprintf("SELECT id, username, email FROM users WHERE %s = $1", field), val)

	user := domain.User{}

	if err := row.Scan(&user.ID, &user.Username, &user.Email); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, codederrors.New(http.StatusNotFound, ErrNotFound, "User with %s '%v' doesn't exist.", field, val)
		default:
			return nil, codederrors.Unknown(err)
		}
	}
	return &user, nil
}

func (r *Repository) GetById(ctx context.Context, id int) (*domain.User, error) {
	return r.getBy(ctx, "id", id)
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return r.getBy(ctx, "username", username)
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.getBy(ctx, "email", email)
}

func (r *Repository) GetCredentials(ctx context.Context, username string) (*domain.Credentials, error) {
	r.logger.Info(fmt.Sprintf("DB request - Users - Get credentials of %s", username))
	db := r.db.Get()
	row := db.QueryRowxContext(ctx, "SELECT password FROM users WHERE username = $1", username)

	var password string

	if err := row.Scan(&password); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, codederrors.New(http.StatusNotFound, ErrNotFound, "User with username '%v' doesn't exist.", username)
		default:
			return nil, codederrors.Unknown(err)
		}
	}
	return &domain.Credentials{username, password}, nil
}

func (r *Repository) Create(ctx context.Context, username string, password domain.Hashed, email string) (*domain.User, error) {
	r.logger.Info("DB request - Users - Create")
	db := r.db.Get()
	row := db.QueryRowxContext(ctx,
		"INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING id",
		username, password, email,
	)

	var id int
	var pqErr *pq.Error

	if err := row.Scan(&id); err != nil {
		switch {
		case errors.As(err, &pqErr) && pqErr.Code == "23505":
			switch pqErr.Constraint {
			case "users_username_key":
				return nil, codederrors.New(http.StatusBadRequest, ErrUsernameExists, "Username or email already exists")
			case "users_email_key":
				return nil, codederrors.New(http.StatusBadRequest, ErrEmailExists, "Username or email already exists")
			}
		default:
			return nil, codederrors.Unknown(err)
		}
	}

	return &domain.User{id, username, email}, nil
}
