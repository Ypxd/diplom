package service

import (
	"context"
	"github.com/Ypxd/diplom/auth/internal/models"
	"github.com/Ypxd/diplom/auth/internal/repository"
	"github.com/jmoiron/sqlx"
)

type Auth interface {
	Auth(ctx context.Context, req models.AuthReq) (string, error)
	Register(ctx context.Context, req models.AuthReq) (string, error)
}

type Service struct {
	Auth Auth
}

func NewService(repo *repository.Repository, conn *sqlx.DB) *Service {
	authService := NewAuthService(repo, conn)

	return &Service{
		Auth: authService,
	}
}
