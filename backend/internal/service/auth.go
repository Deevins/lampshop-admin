// internal/service/auth_service.go

package service

import (
	"context"
	"errors"
	"github.com/Deevins/lampshop-admin-backend/internal/infra"
	"github.com/Deevins/lampshop-admin-backend/internal/repository"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Login(_ context.Context, username, password string) (string, error) {
	if username != "admin" || password != "password123" {
		return "", errors.New("invalid credentials")
	}

	token, err := infra.GenerateJWT(username)
	if err != nil {
		return "", err
	}
	return token, nil
}
