package repository

import (
	"context"
	"github.com/Ypxd/diplom/auth/internal/models"
	"github.com/Ypxd/diplom/auth/internal/repository/postgres"
	"github.com/Ypxd/diplom/auth/internal/repository/redis"
	redcon "github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type Redis interface {
}

type Auth interface {
	Age(ctx context.Context, req models.AuthReq) (int64, error)
	Register(ctx context.Context, req models.AuthReq, age int64) (string, error)
	Auth(ctx context.Context, req models.AuthReq) error
}

type Repository struct {
	Redis Redis
	Auth  Auth
}

func NewRepo(redisCon *redcon.Client) (*Repository, *sqlx.DB, error) {
	db, err := postgres.Connect()
	if err != nil {
		return nil, nil, err
	}

	auth := postgres.NewAuthRepo(db)
	redisCli := redis.NewRedis(redisCon)

	return &Repository{
		Redis: redisCli,
		Auth:  auth,
	}, db, nil
}
