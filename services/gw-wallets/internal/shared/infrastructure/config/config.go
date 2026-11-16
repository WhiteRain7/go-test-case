package config

import (
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/broker"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/config"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/logger"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/postgres"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/protoclient"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/redis"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/shared/infrastructure/auth"
	httpserver "github.com/WhiteRain7/GoTestCase/wallets/internal/shared/infrastructure/httpserver"
	"github.com/WhiteRain7/GoTestCase/wallets/internal/shared/infrastructure/swagger"
)

type Config struct {
	HTTP    httpserver.Config  `mapstructure:"http"`
	Proto   protoclient.Config `mapstructure:"protoclient"`
	DB      postgres.Config    `mapstructure:"db"`
	Logger  logger.Config      `mapstructure:"logger"`
	Auth    auth.Config        `mapstructure:"auth"`
	Swagger swagger.Config     `mapstructure:"swagger"`
	Redis   redis.Config       `mapstructure:"redis"`
	Broker  broker.ProdConfig  `mapstructure:"broker"`
	Debug   bool               `mapstructure:"debug" default:"0"`
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
