package config

import (
	"github.com/WhiteRain7/GoTestCase/exchanger/internal/shared/infrastructure/promanalytics"
	"github.com/WhiteRain7/GoTestCase/exchanger/internal/shared/infrastructure/protoserver"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/config"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/logger"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/postgres"
)

type Config struct {
	Proto     protoserver.Config   `mapstructure:"protoserver"`
	DB        postgres.Config      `mapstructure:"db"`
	Logger    logger.Config        `mapstructure:"logger"`
	Analytics promanalytics.Config `mapstructure:"analytics"`
}

func MustLoad() (*Config, error) {
	cfg := &Config{}
	err := config.LoadTo(cfg)
	return cfg, err
}
