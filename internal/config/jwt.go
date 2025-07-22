package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type JWTConfig struct {
	Secret   string
	Lifetime time.Duration
}

func LoadJWT() (JWTConfig, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return JWTConfig{}, fmt.Errorf("JWT_SECRET is required")
	}

	// Время жизни в секундах, по умолчанию 3600s
	raw := os.Getenv("JWT_LIFETIME")
	if raw == "" {
		raw = "3600"
	}
	secs, err := strconv.Atoi(raw)
	if err != nil {
		return JWTConfig{}, fmt.Errorf("invalid JWT_LIFETIME: %w", err)
	}

	return JWTConfig{
		Secret:   secret,
		Lifetime: time.Duration(secs) * time.Second,
	}, nil
}
