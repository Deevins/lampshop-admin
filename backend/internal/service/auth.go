package service

import (
	"errors"
	"github.com/Deevins/lampshop-admin-backend/internal/domain"
	"github.com/Deevins/lampshop-admin-backend/internal/infra"
)

// AuthService описывает метод логина.
type AuthService interface {
	Login(username, password string) (string, error)
}

type authService struct{}

// NewAuthService создаёт новый AuthService.
func NewAuthService() AuthService {
	return &authService{}
}

// Жёстко захардкоженный админ на время разработки.
var adminUser = domain.User{
	Username: "admin",
	Password: "password123",
}

func (s *authService) Login(username, password string) (string, error) {
	if username != adminUser.Username || password != adminUser.Password {
		return "", errors.New("invalid credentials")
	}
	// Генерируем JWT-токен
	token, err := infra.GenerateJWT(username)
	if err != nil {
		return "", err
	}
	return token, nil
}
