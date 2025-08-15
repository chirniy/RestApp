package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var JwtKey = []byte(viper.GetString("JWT_KEY"))

type Config struct {
	Port        string
	DataBaseURL string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("cfg.env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("ошибка чтения конфигурации: %w", err)
	}

	cfg := &Config{
		Port:        viper.GetString("PORT"),
		DataBaseURL: viper.GetString("DB_URL"),
	}

	return cfg, nil
}
