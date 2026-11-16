package config

import (
	"github.com/WhiteRain7/GoTestCase/analytics/internal/shared/infrastructure/clickhouseserver"
	"github.com/WhiteRain7/GoTestCase/analytics/internal/shared/infrastructure/promanalytics"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/broker"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/config"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/logger"
)

type Config struct {
	Logger     logger.Config           `mapstructure:"logger"`
	Broker     broker.ConConfig        `mapstructure:"broker"`
	ClickHouse clickhouseserver.Config `mapstructure:"clickhouse"`
	Analytics  promanalytics.Config    `mapstructure:"analytics"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	err := config.LoadTo(cfg)
	return cfg, err
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}
