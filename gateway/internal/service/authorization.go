package service

import (
	"context"
	"fmt"

	authclient "github.com/lekss361/curserv2/gateway/internal/clients/auth"
	"github.com/lekss361/curserv2/gateway/internal/repository"
)

type AuthService interface {
	Login(ctx context.Context, login, password string) (string, error)
	ValidateToken(ctx context.Context, token string) (string, bool, error)
}

type authService struct {
	client *authclient.Client
	repo   repository.UserRepo
}

func NewAuthService(client *authclient.Client, repo repository.UserRepo) AuthService {
	return &authService{client: client, repo: repo}
}

func (s *authService) Login(ctx context.Context, login, password string) (string, error) {
	stored, err := s.repo.GetPassword(login)
	if err != nil {
		return "", fmt.Errorf("user lookup failed: %w", err)
	}
	if stored != password {
		return "", fmt.Errorf("invalid credentials")
	}
	token, err := s.client.Login(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("auth service login failed: %w", err)
	}
	return token, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (string, bool, error) {
	valid, err := s.client.ValidateToken(ctx, token)
	if err != nil {
		return "", false, fmt.Errorf("auth service validate failed: %w", err)
	}
	if !valid {
		return "", false, nil
	}
	return "", true, nil
}
