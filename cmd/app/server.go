package app

import (
	"database/sql"
	"github.com/gorilla/mux"
	"jwt/internal/delivery"
	"jwt/internal/repo"
	"jwt/internal/usecase"
	"log"
	"net/http"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/app/v1").Subrouter()

	userRepo := repo.NewRepo(s.db)
	userUseCase := usecase.NewUserUsecase(userRepo)
	userHandler := delivery.NewHandler(userUseCase)
	userHandler.RegisterRoutes(subrouter)

	log.Println("listening on address", s.addr)

	return http.ListenAndServe(s.addr, router)
}
