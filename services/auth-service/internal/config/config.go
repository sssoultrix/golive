package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	RedisAddr   string

	JWTSecret  string
	AccessTTL  time.Duration
	RefreshTTL time.Duration

	BcryptCost         int
	RefreshTokenPepper string

	ProfileServiceURL string
}

func Load() (Config, error) {
	v := viper.New()

	v.SetDefault("http_addr", ":8080")
	v.SetDefault("access_ttl", "15m")
	v.SetDefault("refresh_ttl", "720h")
	v.SetDefault("bcrypt_cost", 10)
	v.SetDefault("profile_service_url", "http://localhost:8081")

	v.SetConfigFile("services/auth-service/configs/config.yaml")
	_ = v.ReadInConfig()

	v.AutomaticEnv()
	_ = v.BindEnv("database_url", "DATABASE_URL")
	_ = v.BindEnv("redis_addr", "REDIS_ADDR")
	_ = v.BindEnv("jwt_secret", "JWT_SECRET")
	_ = v.BindEnv("http_addr", "HTTP_ADDR")
	_ = v.BindEnv("access_ttl", "ACCESS_TTL")
	_ = v.BindEnv("refresh_ttl", "REFRESH_TTL")
	_ = v.BindEnv("bcrypt_cost", "BCRYPT_COST")
	_ = v.BindEnv("refresh_token_pepper", "REFRESH_TOKEN_PEPPER")
	_ = v.BindEnv("profile_service_url", "PROFILE_SERVICE_URL")

	accessTTL, err := time.ParseDuration(v.GetString("access_ttl"))
	if err != nil {
		return Config{}, fmt.Errorf("ACCESS_TTL: %w", err)
	}
	refreshTTL, err := time.ParseDuration(v.GetString("refresh_ttl"))
	if err != nil {
		return Config{}, fmt.Errorf("REFRESH_TTL: %w", err)
	}

	cfg := Config{
		HTTPAddr:    v.GetString("http_addr"),
		DatabaseURL: v.GetString("database_url"),
		RedisAddr:   v.GetString("redis_addr"),

		JWTSecret:  v.GetString("jwt_secret"),
		AccessTTL:  accessTTL,
		RefreshTTL: refreshTTL,

		BcryptCost:         v.GetInt("bcrypt_cost"),
		RefreshTokenPepper: v.GetString("refresh_token_pepper"),

		ProfileServiceURL: v.GetString("profile_service_url"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.RedisAddr == "" {
		return Config{}, fmt.Errorf("REDIS_ADDR is required")
	}
	if len(cfg.JWTSecret) < 32 {
		return Config{}, fmt.Errorf("JWT_SECRET must be at least 32 bytes")
	}
	if cfg.BcryptCost < 10 || cfg.BcryptCost > 14 {
		return Config{}, fmt.Errorf("BCRYPT_COST must be between 10 and 14")
	}
	if cfg.RefreshTokenPepper == "" {
		return Config{}, fmt.Errorf("REFRESH_TOKEN_PEPPER is required")
	}

	return cfg, nil
}
