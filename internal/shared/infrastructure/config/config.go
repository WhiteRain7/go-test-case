package config

import (
	"github.com/spf13/viper"
)

type HTTPConfig struct {
	Port int `mapstructure:"port"`
}

type DBConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Database string `mapstructure:"name"`
	Port     int    `mapstructure:"port"`
	SSLMode  bool   `mapstructure:"ssl"`
}

type LoggerConfig struct {
	Path    string `mapstructure:"path"`
	Console bool   `mapstructure:"console"`
}

type Config struct {
	HTTP   HTTPConfig   `mapstructure:"http"`
	DB     DBConfig     `mapstructure:"db"`
	Logger LoggerConfig `mapstructure:"logger"`
}

func Load(dirs ...string) (*Config, error) {
	v := viper.New()

	for _, path := range dirs {
		v.AddConfigPath(path)
	}
	if len(dirs) == 0 {
		v.AddConfigPath(".")
	}

	v.SetConfigName(".env")
	v.SetConfigType("dotenv")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func MustLoad(dirs ...string) *Config {
	cfg, err := Load(dirs...)
	if err != nil {
		panic(err)
	}
	return cfg
}
