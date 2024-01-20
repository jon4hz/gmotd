package config

import (
	"github.com/spf13/viper"
)

func Load(cfg *Config) error {
	viper.SetConfigName("gmotd")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/gmotd")
	viper.AddConfigPath("$HOME/.gmotd")
	viper.AddConfigPath(".")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}
		return err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}

	return nil
}

func setDefaults() {
	viper.SetDefault("zpool.disabled", true)
}
