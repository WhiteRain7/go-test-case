package repository_test

import (
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
