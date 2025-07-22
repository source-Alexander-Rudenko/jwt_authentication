package config

import (
	"fmt"
	"os"
)

type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DB       string
	URL      string
}

func GetPostgresConfig() (PostgresConfig, error) {
	var config PostgresConfig
	config.User = os.Getenv("POSTGRES_USER")
	config.Password = os.Getenv("POSTGRES_PASSWORD")
	config.Host = os.Getenv("POSTGRES_HOST")
	config.Port = os.Getenv("POSTGRES_PORT")
	config.DB = os.Getenv("POSTGRES_DB")
	config.URL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.User, config.Password, config.Host, config.Port, config.DB)

	return config, nil
}
