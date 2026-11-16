package config

import (
	"github.com/WhiteRain7/GoTestCase/notification/internal/shared/infrastructure/mongoserver"
	"github.com/WhiteRain7/GoTestCase/notification/internal/shared/infrastructure/promanalytics"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/broker"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/config"
	"github.com/WhiteRain7/GoTestCase/pkg/infrastructure/logger"
)

type Config struct {
	Logger    logger.Config        `mapstructure:"logger"`
	Broker    broker.ConConfig     `mapstructure:"broker"`
	Mongo     mongoserver.Config   `mapstructure:"mongo"`
	Analytics promanalytics.Config `mapstructure:"analytics"`
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
