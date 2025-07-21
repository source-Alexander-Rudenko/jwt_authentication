package app

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"jwt_auth_project/internal/delivery"
	"jwt_auth_project/internal/repo"
	"jwt_auth_project/internal/usecase"
	"log"
	"net/http"
)

type APIServer struct {
	addr string
	pool *pgxpool.Pool
}

func NewAPIServer(addr string, pool *pgxpool.Pool) *APIServer {
	return &APIServer{
		addr: addr,
		pool: pool,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	userRepo := repo.NewUserRepo(s.pool)
	userUseCase := usecase.NewUserUsecase(userRepo)
	userHandler := delivery.NewHandler(userUseCase)
	userHandler.RegisterRoutes(router)

	log.Println("listening on address", s.addr)

	return http.ListenAndServe(s.addr, router)
}
