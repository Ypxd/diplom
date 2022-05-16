package postgres

import (
	"context"
	"github.com/Ypxd/diplom/auth/internal/models"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

type Tags struct {
	db *sqlx.DB
}

func (t *Tags) UpdateUnfavoriteTags(ctx context.Context, req []models.AllTags, userID string) error {
	const query = `
		UPDATE auth.users SET uf_tags = $1 WHERE user_id = $2
`
	s := ""
	for _, r := range req {
		if s == "" {
			s = s + strconv.FormatInt(r.ID, 10)
		} else {
			s = s + ";" + strconv.FormatInt(r.ID, 10)
		}
	}

	_, err := t.db.ExecContext(ctx, query, s, userID)
	if err != nil {
		return err
	}

	return nil
}

func (t *Tags) UpdateFavoriteTags(ctx context.Context, s string, userID string) error {
	const query = `
		UPDATE auth.users SET f_tags = $1 WHERE user_id = $2
`

	_, err := t.db.ExecContext(ctx, query, s, userID)
	if err != nil {
		return err
	}

	return nil
}

func (t *Tags) UserUnfavoriteTags(ctx context.Context, userID string) ([]string, error) {
	var (
		res []string
		s   []string
	)
	const query = `
		SELECT uf_tags FROM auth.users WHERE user_id = $1
`

	err := t.db.SelectContext(ctx, &s, query, userID)
	if err != nil {
		return res, err
	}

	if len(s) != 0 {
		if s[0] != "" {
			res = strings.Split(s[0], ";")
		}
	}

	return res, nil
}

func (t *Tags) AllUnfavoriteTagsTags(ctx context.Context, str []string) ([]models.AllTags, error) {
	const query = `
		SELECT id, title FROM tags.tags WHERE id in ($1)
`

	result := make([]models.AllTags, 0)
	for _, s := range str {
		var res []models.AllTags
		err := t.db.SelectContext(ctx, &res, query, s)
		if err != nil {
			return result, err
		}
		for _, r := range res {
			result = append(result, r)
		}
	}

	return result, nil
}

func (t *Tags) AllTags(ctx context.Context) ([]models.AllTags, error) {
	var res []models.AllTags
	const query = `
		SELECT id, title, val FROM tags.tags
`

	err := t.db.SelectContext(ctx, &res, query)
	if err != nil {
		return res, err
	}

	return res, nil
}

func NewTagsRepo(db *sqlx.DB) *Tags {
	return &Tags{
		db: db,
	}
}
