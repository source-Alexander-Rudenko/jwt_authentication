package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/multitracer"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/pressly/goose/v3"
	"jwt_auth_project/config"
	"log/slog"
	"os"
)

func InitPostgres(config config.PostgresConfig, slogger *slog.Logger) (*pgxpool.Pool, error) {
	dbConf, err := pgxpool.ParseConfig(config.URL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse db URL: %v", err)
	}

	queryTracer := NewTracer(slogger, tracelog.LogLevelDebug)

	dbConf.ConnConfig.Tracer = multitracer.New(queryTracer)

	pool, err := pgxpool.NewWithConfig(context.Background(), dbConf)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %v", err)
	}

	postgresPing := pool.Ping(context.Background())
	if postgresPing != nil {
		return nil, fmt.Errorf("unable to ping db: %v", err)
	}

	if err := RunMigrations(config.URL); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}
	slog.Info("Successfully connected to db")
	return pool, nil
}

func RunMigrations(dbURL string) error {
	migrationsDir := os.Getenv("MIGRATION_FOLDER")
	if migrationsDir == "" {
		return fmt.Errorf("MIGRATION_FOLDER environment variable is not set")
	}

	sqlDB, err := sql.Open("pgx", dbURL)
	if err != nil {
		return fmt.Errorf("unable to open db: %v", err)
	}
	defer sqlDB.Close()

	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}
	slog.Info("Successfully migrated db")
	return nil
}
