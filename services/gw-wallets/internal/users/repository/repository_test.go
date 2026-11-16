package repository_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/78bits/go-sqlmock-sqlx"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/users/domain"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/users/repository"
	"github.com/jmoiron/sqlx"
)

type MockLogger struct{}

func (*MockLogger) Info(string) {}

type MockDBGetter struct {
	db sqlx.ExtContext
}

func (g *MockDBGetter) Get() sqlx.ExtContext {
	return g.db
}

func TestRepository_GetById(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		rows     *sqlmock.Rows
		wantUser *domain.User
		wantErr  error
	}{
		{
			"get user by id",
			1,
			sqlmock.NewRows([]string{"id", "username", "email"}).AddRow(1, "user", "mail@example.com"),
			&domain.User{1, "user", "mail@example.com"},
			nil,
		},
		{
			"no user",
			2,
			sqlmock.NewRows([]string{"id", "username", "email"}),
			nil,
			repository.ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.Newx()
			if err != nil {
				t.Fatalf("got error '%s' when tried to mock sql", err)
			}
			defer db.Close()

			mock.ExpectQuery("SELECT id, username, email FROM users").
				WithArgs(test.id).
				WillReturnRows(test.rows)

			repo := repository.NewRepository(&MockDBGetter{db}, &MockLogger{})

			user, err := repo.GetById(context.Background(), test.id)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(user, test.wantUser) {
				t.Errorf("got user %v, want %v", user, test.wantUser)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("got unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_GetByEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		rows     *sqlmock.Rows
		wantUser *domain.User
		wantErr  error
	}{
		{
			"get user by email",
			"mail@example.com",
			sqlmock.NewRows([]string{"id", "username", "email"}).AddRow(1, "user", "mail@example.com"),
			&domain.User{1, "user", "mail@example.com"},
			nil,
		},
		{
			"no user",
			"unknown@example.com",
			sqlmock.NewRows([]string{"id", "username", "email"}),
			nil,
			repository.ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.Newx()
			if err != nil {
				t.Fatalf("got error '%s' when tried to mock sql", err)
			}
			defer db.Close()

			mock.ExpectQuery("SELECT id, username, email FROM users").
				WithArgs(test.email).
				WillReturnRows(test.rows)

			repo := repository.NewRepository(&MockDBGetter{db}, &MockLogger{})

			user, err := repo.GetByEmail(context.Background(), test.email)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(user, test.wantUser) {
				t.Errorf("got user %v, want %v", user, test.wantUser)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("got unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_GetByUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		rows     *sqlmock.Rows
		wantUser *domain.User
		wantErr  error
	}{
		{
			"get user by email",
			"user",
			sqlmock.NewRows([]string{"id", "username", "email"}).AddRow(1, "user", "mail@example.com"),
			&domain.User{1, "user", "mail@example.com"},
			nil,
		},
		{
			"no user",
			"unknown",
			sqlmock.NewRows([]string{"id", "username", "email"}),
			nil,
			repository.ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.Newx()
			if err != nil {
				t.Fatalf("got error '%s' when tried to mock sql", err)
			}
			defer db.Close()

			mock.ExpectQuery("SELECT id, username, email FROM users").
				WithArgs(test.username).
				WillReturnRows(test.rows)

			repo := repository.NewRepository(&MockDBGetter{db}, &MockLogger{})

			user, err := repo.GetByUsername(context.Background(), test.username)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(user, test.wantUser) {
				t.Errorf("got user %v, want %v", user, test.wantUser)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("got unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_GetCredentials(t *testing.T) {
	tests := []struct {
		name            string
		username        string
		rows            *sqlmock.Rows
		wantCredentials *domain.Credentials
		wantErr         error
	}{
		{
			"get credentials by username",
			"user",
			sqlmock.NewRows([]string{"password"}).AddRow("hashed"),
			&domain.Credentials{"user", "hashed"},
			nil,
		},
		{
			"no credentials",
			"unknown",
			sqlmock.NewRows([]string{"password"}),
			nil,
			repository.ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.Newx()
			if err != nil {
				t.Fatalf("got error '%s' when tried to mock sql", err)
			}
			defer db.Close()

			mock.ExpectQuery("SELECT password FROM users").
				WithArgs(test.username).
				WillReturnRows(test.rows)

			repo := repository.NewRepository(&MockDBGetter{db}, &MockLogger{})

			credentials, err := repo.GetCredentials(context.Background(), test.username)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(credentials, test.wantCredentials) {
				t.Errorf("got credentials %v, want %v", credentials, test.wantCredentials)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("got unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_Create(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		email    string
		rows     *sqlmock.Rows
		wantUser *domain.User
		wantErr  error
	}{
		{
			"create user",
			"user",
			"mail@example.com",
			"hashed",
			sqlmock.NewRows([]string{"id"}).AddRow(1),
			&domain.User{1, "user", "hashed"},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlmock.Newx()
			if err != nil {
				t.Fatalf("got error '%s' when tried to mock sql", err)
			}
			defer db.Close()

			mock.ExpectQuery("INSERT INTO users").
				WithArgs(test.username, test.password, test.email).
				WillReturnRows(test.rows)

			repo := repository.NewRepository(&MockDBGetter{db}, &MockLogger{})

			credentials, err := repo.Create(context.Background(), test.username, test.password, test.email)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %v, want %v", err, test.wantErr)
			}

			if !reflect.DeepEqual(credentials, test.wantUser) {
				t.Errorf("got users %v, want %v", credentials, test.wantUser)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("got unfulfilled expectations: %s", err)
			}
		})
	}
}
