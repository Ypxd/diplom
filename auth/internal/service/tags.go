package service

import (
	"context"
	"github.com/Ypxd/diplom/auth/internal/models"
	"github.com/Ypxd/diplom/auth/internal/repository"
	"github.com/Ypxd/diplom/auth/utils"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type TagsService struct {
	repo *repository.Repository
	conn *sqlx.DB
}

func (t *TagsService) UpdateUnfavoriteTags(ctx context.Context, req []models.AllTags, userID string) error {
	return t.repo.Tags.UpdateUnfavoriteTags(ctx, req, userID)
}

func (t *TagsService) GetUnfavoriteTags(ctx context.Context, userID string) ([]models.AllTags, error) {
	sTags, err := t.repo.Tags.UserUnfavoriteTags(ctx, userID)
	if err != nil {
		return nil, err
	}

	return t.repo.Tags.AllUnfavoriteTagsTags(ctx, sTags)
}

func (t *TagsService) UpdateFavoriteTags(ctx context.Context, req models.MyEvents, userID string) error {
	tags, err := t.repo.Events.GetEventsTag(ctx, req.Title)
	if err != nil {
		return err
	}

	res, err := t.repo.Redis.UpdateUserTags(tags, userID)
	if err != nil {
		return err
	}
	s := ""
	for r, count := range res {
		if count >= utils.GetConfig().RepeatCount {
			if s == "" {
				s = s + strconv.FormatInt(r, 10)
			} else {
				s = s + ";" + strconv.FormatInt(r, 10)
			}
		}
	}
	err = t.repo.Tags.UpdateFavoriteTags(ctx, s, userID)
	if err != nil {
		return err
	}

	return nil
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
