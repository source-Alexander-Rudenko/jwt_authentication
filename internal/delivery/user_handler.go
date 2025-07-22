package delivery

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"jwt_auth_project/internal/domain"
	"jwt_auth_project/internal/usecase"
	"jwt_auth_project/internal/utils"
)

type Handler struct {
	userUseCase usecase.UserUseCase
}

func NewHandler(userUseCase usecase.UserUseCase) *Handler {
	return &Handler{userUseCase: userUseCase}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleRegister).Methods(http.MethodPost)
	router.HandleFunc("/login", h.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/logout", h.handleLogout).Methods(http.MethodPost)
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload domain.RegisterUserPayload
	if err := utils.ParceJSON(r, &payload); err != nil {
		slog.Error("register: invalid JSON", "error", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, token, err := h.userUseCase.Register(r.Context(), payload)
	if err != nil {
		slog.Error("register: usecase failed", "email", payload.Email, "error", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	slog.Info("user registered", "email", user.Email)

	resp := map[string]any{
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	}
	if err := utils.WriteJSON(w, http.StatusCreated, resp); err != nil {
		slog.Error("register: write response failed", "error", err)
	}
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload domain.LoginUserPayload
	if err := utils.ParceJSON(r, &payload); err != nil {
		slog.Error("login: invalid JSON", "error", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	token, err := h.userUseCase.Login(r.Context(), payload)
	if err != nil {
		slog.Error("login: usecase failed", "email", payload.Email, "error", err)
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid credentials"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	slog.Info("user logged in", "email", payload.Email)

	// Можно вернуть простое сообщение
	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "ok"}); err != nil {
		slog.Error("login: write response failed", "error", err)
	}
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	slog.Info("user logged out")

	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "logged out"}); err != nil {
		slog.Error("logout: write response failed", "error", err)
	}
}
