package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"log/slog"

	"jwt_auth_project/internal/usecase"
	"jwt_auth_project/internal/utils"
)

// AuthMiddleware проверяет JWT в cookie и добавляет userID в контекст
func AuthMiddleware(userUC usecase.UserUseCase) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Извлекаем JWT из cookie
			cookie, err := r.Cookie("jwt")
			if err != nil {
				slog.Error("auth: missing or invalid cookie", "error", err)
				utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("authentication required"))
				return
			}

			// Валидируем токен и получаем userID
			userID, err := userUC.ValidateToken(cookie.Value)
			if err != nil {
				slog.Error("auth: token validation failed", "error", err)
				utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
				return
			}

			// Сохраняем userID в контексте запроса
			ctx := context.WithValue(r.Context(), utils.ContextKeyUserID, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
