package usecase

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"jwt_auth_project/internal/domain"
	"jwt_auth_project/internal/repo"
	"jwt_auth_project/internal/utils"
	"strconv"
	"strings"
	"time"
)

// UserUseCase описывает операции регистрации, логина, валидации токена и получение юзера
type UserUseCase interface {
	Register(ctx context.Context, payload domain.RegisterUserPayload) (*domain.User, string, error)
	Login(ctx context.Context, payload domain.LoginUserPayload) (string, error)
	ValidateToken(tokenStr string) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
}

type userUseCase struct {
	repo      repo.UserRepository
	jwtSecret []byte
	ttl       time.Duration
}

// NewUserUsecase конструктор. jwtSecret — из конфига, ttl — время жизни токена.
func NewUserUsecase(r repo.UserRepository, jwtSecret string, ttl time.Duration) UserUseCase {
	return &userUseCase{
		repo:      r,
		jwtSecret: []byte(jwtSecret),
		ttl:       ttl,
	}
}

func (u *userUseCase) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	return u.repo.GetUserByID(ctx, id)
}

// Register валидация, хеширование, сохранение и генерация JWT
func (u *userUseCase) Register(ctx context.Context, payload domain.RegisterUserPayload) (*domain.User, string, error) {
	payload.Username = strings.TrimSpace(payload.Username)
	payload.Email = strings.TrimSpace(payload.Email)
	if err := utils.Validate.Struct(payload); err != nil {
		return nil, "", err
	}
	newID := uuid.New()
	hashed, err := hashPassword(payload.Password)
	if err != nil {
		return nil, "", err
	}

	user := domain.User{
		ID:        newID,
		Username:  payload.Username,
		Email:     payload.Email,
		Password:  hashed,
		CreatedAt: time.Now(),
	}

	if err := u.repo.CreateUser(ctx, user); err != nil {
		return nil, "", err
	}

	tokenStr, err := u.generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return &user, tokenStr, nil
}

// Login валидация, проверка пароля, генерация JWT
func (u *userUseCase) Login(ctx context.Context, payload domain.LoginUserPayload) (string, error) {
	if err := utils.Validate.Struct(payload); err != nil {
		return "", err
	}

	user, err := u.repo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return "", err
	}

	ok, err := verifyPassword(user.Password, payload.Password)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("invalid credentials")
	}

	return u.generateToken(user.ID)
}

// ValidateToken парсит и проверяет JWT, возвращает userID из claims.Subject
func (u *userUseCase) ValidateToken(tokenStr string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return u.jwtSecret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	uid, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, err
	}
	return uid, nil
}

// generateToken соберет JWT с полем Subject=userID и сроком ttl
func (u *userUseCase) generateToken(userID uuid.UUID) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(u.ttl)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(u.jwtSecret)
}

// hashPassword хеширует пароль Argon2id и возвращает строку в формате "salt$hash"
func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	bSalt := base64.RawStdEncoding.EncodeToString(salt)
	bHash := base64.RawStdEncoding.EncodeToString(hash)
	return fmt.Sprintf("%s$%s", bSalt, bHash), nil
}

// verifyPassword сравнивает raw-пароль с закешированным hash в формате "salt$hash"
func verifyPassword(encoded, password string) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 2 {
		return false, errors.New("invalid hash format")
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	newHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return subtle.ConstantTimeCompare(hash, newHash) == 1, nil
}
