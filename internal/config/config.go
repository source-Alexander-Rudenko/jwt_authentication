package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

const (
	DefaultPort = "8080"
)

type Config struct {
	Port           string
	PostgresConfig PostgresConfig
	//RedisConfig    sessionRepository.RedisConfig
	JWT JWTConfig
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found")
	}
	var err error
	cfg := &Config{}

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}
	cfg.Port = ":" + port

	//redisConfig, err := sessionRepository.GetRedisConfig()
	//if err != nil {
	//	return Config{}, errors.New("Failed to connect to the redis database")
	//}
	//conf.RedisConfig = *redisConfig

	cfg.PostgresConfig, err = GetPostgresConfig()
	if err != nil {
		return nil, errors.New("Failed to connect to the postgres database")
	}

	cfg.JWT, err = LoadJWT()
	if err != nil {
		return nil, fmt.Errorf("load jwt config: %w", err)
	}

	return cfg, nil
}
