package service

import (
	"context"
	"github.com/Ypxd/diplom/auth/internal/models"
	"github.com/Ypxd/diplom/auth/internal/repository"
	"github.com/jmoiron/sqlx"
)

type TagsService struct {
	repo *repository.Repository
	conn *sqlx.DB
}

func (t *TagsService) UpdateUnfavoriteTags(ctx context.Context, req []models.AllTags, userID string) error {
	return t.repo.Tags.UpdateUnfavoriteTagsTags(ctx, req, userID)
}

func (t *TagsService) GetUnfavoriteTags(ctx context.Context, userID string) ([]models.AllTags, error) {
	sTags, err := t.repo.Tags.UserUnfavoriteTags(ctx, userID)
	if err != nil {
		return nil, err
	}

	return t.repo.Tags.AllUnfavoriteTagsTags(ctx, sTags)
}

func (t *TagsService) GetAllTags(ctx context.Context) ([]models.AllTags, error) {
	return t.repo.Tags.AllTags(ctx)
}

func NewTagsService(repo *repository.Repository, conn *sqlx.DB) *TagsService {
	return &TagsService{
		repo: repo,
		conn: conn,
	}
}
