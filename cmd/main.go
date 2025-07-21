package main

import (
	"jwt_auth_project/cmd/app"
	"jwt_auth_project/config"
	"log/slog"
	"os"
)

func main() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)

	conf, err := config.LoadConfig()
	if err != nil {
		app.Fatal("Error loading config", err)
	}

	pool, err := app.InitPostgres(conf.PostgresConfig, logger)
	if err != nil {
		app.Fatal("Error initializing postgres", err)
	}
	defer pool.Close()

	server := app.NewAPIServer(":8080", pool)
	if err := server.Run(); err != nil {
		app.Fatal("Error starting server", err)
	}
}
