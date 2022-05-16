package repository

import (
	"context"
	"github.com/Ypxd/diplom/auth/internal/models"
	"github.com/Ypxd/diplom/auth/internal/repository/postgres"
	"github.com/Ypxd/diplom/auth/internal/repository/redis"
	redcon "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Redis interface {
	UpdateUserTags(t string, userID string) (map[int64]int64, error)
}

type Auth interface {
	Age(ctx context.Context, req models.AuthReq) (int64, error)
	Register(ctx context.Context, req models.AuthReq, age int64) (string, error)
	Auth(ctx context.Context, req models.AuthReq) (*uuid.UUID, error)
	UserInfo(ctx context.Context, userID string) (*models.UserInfo, error)
	ChangePass(ctx context.Context, req models.ChangePassReq, userID string) error
}

type Events interface {
	AllEvents(ctx context.Context) ([]models.Events, error)
	GetEventsTag(ctx context.Context, t string) (string, error)
}

type Tags interface {
	AllTags(ctx context.Context) ([]models.AllTags, error)
	UserUnfavoriteTags(ctx context.Context, userID string) ([]string, error)
	AllUnfavoriteTagsTags(ctx context.Context, s []string) ([]models.AllTags, error)
	UpdateUnfavoriteTags(ctx context.Context, req []models.AllTags, userID string) error
	UpdateFavoriteTags(ctx context.Context, s string, userID string) error
}

type Repository struct {
	Redis  Redis
	Auth   Auth
	Events Events
	Tags   Tags
}

func NewRepo(redisCon *redcon.Client) (*Repository, *sqlx.DB, error) {
	db, err := postgres.Connect()
	if err != nil {
		return nil, nil, err
	}

	events := postgres.NewEventsRepo(db)
	auth := postgres.NewAuthRepo(db)
	tags := postgres.NewTagsRepo(db)
	redisCli := redis.NewRedis(redisCon)

	return &Repository{
		Redis:  redisCli,
		Auth:   auth,
		Events: events,
		Tags:   tags,
	}, db, nil
}
