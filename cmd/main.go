package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"jwt_auth_project/internal/config"
	"jwt_auth_project/internal/db"
	"jwt_auth_project/internal/delivery"
	middleware "jwt_auth_project/internal/delivery/middleware"
	"jwt_auth_project/internal/logger"
	"jwt_auth_project/internal/repo"
	"jwt_auth_project/internal/usecase"
)

func main() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	newLog := slog.New(handler)

	conf, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("load config failed", err)
	}

	s3Cfg := config.NewConfigFromEnv()
	s3Client := config.NewS3Client(s3Cfg)

	pool, err := db.InitPostgres(conf.PostgresConfig, newLog)
	if err != nil {
		logger.Fatal("init postgres failed", err)
	}
	defer pool.Close()

	userRepo := repo.NewUserRepo(pool)
	userUC := usecase.NewUserUsecase(userRepo, conf.JWT.Secret, conf.JWT.Lifetime)

	adsRepo := repo.NewAdsRepo(pool)
	adsUC := usecase.NewAdsUsecase(adsRepo, s3Client, s3Cfg.Bucket, 5<<20) // макс 5MiB, например

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := adsUC.InitBucket(ctx); err != nil {
		logger.Fatal("init S3 bucket failed", err)
	}

	router := mux.NewRouter()

	userHandler := delivery.NewHandler(userUC)
	userHandler.RegisterRoutes(router)

	adsHandler := delivery.NewAdsHandler(adsUC)
	router.Use(middleware.AuthMiddleware(userUC))
	adsHandler.RegisterRoutes(router)

	slog.Info("listening on", "port", conf.Port)
	if err := http.ListenAndServe(conf.Port, router); err != nil {
		logger.Fatal("server error", err)
	}
}
