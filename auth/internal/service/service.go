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
	UserInfo(ctx context.Context, userID string) (models.UserInfo, error)
	ChangePass(ctx context.Context, req models.ChangePassReq, userID string) error
}

type Events interface {
	GetAllEvents(ctx context.Context) ([]models.Events, error)
}

type Tags interface {
	GetAllTags(ctx context.Context) ([]models.AllTags, error)
	GetUnfavoriteTags(ctx context.Context, userID string) ([]models.AllTags, error)
	UpdateUnfavoriteTags(ctx context.Context, req []models.AllTags, userID string) error
}

type Service struct {
	Auth   Auth
	Events Events
	Tags   Tags
}

func NewService(repo *repository.Repository, conn *sqlx.DB) *Service {
	authService := NewAuthService(repo, conn)
	eventsService := NewEventsService(repo, conn)
	tagsService := NewTagsService(repo, conn)

	return &Service{
		Auth:   authService,
		Events: eventsService,
		Tags:   tagsService,
	}
}
