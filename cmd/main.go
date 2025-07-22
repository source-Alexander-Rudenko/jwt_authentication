package main

import (
	"github.com/gorilla/mux"
	"jwt_auth_project/internal/config"
	"jwt_auth_project/internal/db"
	"jwt_auth_project/internal/delivery"
	"jwt_auth_project/internal/logger"
	"jwt_auth_project/internal/repo"
	"jwt_auth_project/internal/usecase"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	newLog := slog.New(handler)

	conf, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Error loading config", err)

	}

	pool, err := db.InitPostgres(conf.PostgresConfig, newLog)
	if err != nil {
		logger.Fatal("Error initializing postgres", err)
	}
	defer pool.Close()

	router := mux.NewRouter()

	userRepo := repo.NewUserRepo(pool)
	userUseCase := usecase.NewUserUsecase(userRepo, conf.JWT.Secret, conf.JWT.Lifetime)
	userHandler := delivery.NewHandler(userUseCase)
	userHandler.RegisterRoutes(router)

	slog.Info("listening on address", conf.Port)
	err = http.ListenAndServe(conf.Port, router)
	if err != nil {
		logger.Fatal("Error starting server", err)
	}
}
