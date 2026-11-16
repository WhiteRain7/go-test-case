package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Database string `mapstructure:"name"`
	Port     int    `mapstructure:"port"`
	SSLMode  bool   `mapstructure:"ssl"`
}

func Connect(cfg Config) (*sqlx.DB, error) {
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

func MustConnect(cfg Config) *sqlx.DB {
	db, err := Connect(cfg)
	if err != nil {
		panic(err)
	}
	return db
}
