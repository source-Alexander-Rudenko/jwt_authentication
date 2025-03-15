package delivery

import (
	"github.com/gorilla/mux"
	"jwt/internal/domain"
	"jwt/internal/usecase"
	"jwt/internal/utils"
	"net/http"
)

type Handler struct {
	userUseCase usecase.UserUseCase
}

func NewHandler(userUseCase usecase.UserUseCase) *Handler {
	return &Handler{userUseCase: userUseCase}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/register", h.handleRegister).Methods(http.MethodPost)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload domain.RegisterUserPayload
	if err := utils.ParceJSON(r, payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

}
