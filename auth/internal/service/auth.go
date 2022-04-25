package service

import (
	"context"
	"github.com/Ypxd/diplom/auth/internal/models"
	"github.com/Ypxd/diplom/auth/internal/repository"
	"github.com/jmoiron/sqlx"
)

type AuthService struct {
	repo *repository.Repository
	conn *sqlx.DB
}

func (a *AuthService) Auth(ctx context.Context, req models.AuthReq) (string, error) {
	err := a.repo.Auth.Auth(ctx, req)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

func (a *AuthService) Register(ctx context.Context, req models.AuthReq) (string, error) {
	age, err := a.repo.Auth.Age(ctx, req)
	if err != nil {
		return "", err
	}

	m, err := a.repo.Auth.Register(ctx, req, age)
	if err != nil {
		return m, err
	}

	return "OK", nil
}

func NewAuthService(repo *repository.Repository, conn *sqlx.DB) *AuthService {
	return &AuthService{
		repo: repo,
		conn: conn,
	}
}
