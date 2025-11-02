package service

import (
	"context"
	"errors"
	"time"

	"github.com/IlyaChern12/rtce/internal/models"
	"github.com/IlyaChern12/rtce/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.UserRepository
	jwtSecret string
}

// конструктор для аутентификатора
func NewAuthService(repo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		repo: repo,
		jwtSecret: jwtSecret,
	}
}

// регистрация юзера
func (a *AuthService) Register(ctx context.Context, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Email: email,
		PasswordHash: string(hash),
	}

	return a.repo.Create(ctx, user)
}

// авторизация юзера (вернет jwt токен)
func (a *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := a.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", errors.New("invalid email or password")
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.jwtSecret))
}