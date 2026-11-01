package postgres

import (
	"fmt"

	"github.com/WhiteRain7/GoTestCase/internal/shared/infrastructure/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(cfg config.DBConfig) (*sqlx.DB, error) {
	var sslmode string
	if cfg.SSLMode {
		sslmode = "enable"
	} else {
		sslmode = "disable"
	}

	db, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			"user=%s password=%s host=%s dbname=%s port=%d sslmode=%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Database, cfg.Port, sslmode))

	if err != nil {
		return nil, err
	}

	return db, nil
}

func MustConnect(cfg config.DBConfig) *sqlx.DB {
	db, err := Connect(cfg)
	if err != nil {
		panic(err)
	}
	return db
}
