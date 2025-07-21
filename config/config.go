package config

import (
	"errors"
	"os"
)

const (
	DefaultPort = "8080"
)

type Config struct {
	Port           string
	PostgresConfig PostgresConfig
	//RedisConfig    sessionRepository.RedisConfig
}

func LoadConfig() (Config, error) {
	var conf Config

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}
	conf.Port = port

	//redisConfig, err := sessionRepository.GetRedisConfig()
	//if err != nil {
	//	return Config{}, errors.New("Failed to connect to the redis database")
	//}
	//conf.RedisConfig = *redisConfig

	postgresConfig, err := GetPostgresConfig()
	if err != nil {
		return Config{}, errors.New("Failed to connect to the postgres database")
	}
	conf.PostgresConfig = postgresConfig

	return conf, nil
}
