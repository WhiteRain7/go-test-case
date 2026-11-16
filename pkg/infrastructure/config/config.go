package config

import (
	"github.com/spf13/viper"
)

func LoadTo[T any](config *T) error {
	v := viper.New()

	v.AddConfigPath(".")
	v.SetConfigFile(".env")
	v.SetConfigType("dotenv")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.AutomaticEnv()

	if err := v.Unmarshal(config); err != nil {
		return err
	}

	return nil
}
