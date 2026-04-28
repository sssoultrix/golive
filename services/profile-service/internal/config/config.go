package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
}

func Load() (Config, error) {
	v := viper.New()

	v.SetDefault("http_addr", ":8081")

	v.SetConfigFile("services/profile-service/configs/config.yaml")
	_ = v.ReadInConfig()

	v.AutomaticEnv()
	_ = v.BindEnv("database_url", "DATABASE_URL")
	_ = v.BindEnv("http_addr", "HTTP_ADDR")

	cfg := Config{
		HTTPAddr:    v.GetString("http_addr"),
		DatabaseURL: v.GetString("database_url"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

