package postgres

import (
	"context"
	"errors"
	"github.com/Ypxd/diplom/auth/internal/models"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type Events struct {
	db *sqlx.DB
}

type selRes struct {
	Selected string `db:"selected"`
}

func (e *Events) AllEvents(ctx context.Context) ([]models.Events, error) {
	var res []models.Events
	const query = `
		SELECT title, address, tags, png, age_id, selected
		FROM tags.events
`

	err := e.db.SelectContext(ctx, &res, query)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (e *Events) GetEventsTag(ctx context.Context, t string, userID string) (string, error) {
	q := `
		SELECT selected FROM tags.events
		WHERE title = $1
`
	sel := make([]selRes, 0)
	err := e.db.SelectContext(ctx, &sel, q, t)
	if err != nil {
		return "", err
	}
	if len(sel) != 1 {
		return "", errors.New("null selected")
	}
	val, err := strconv.ParseInt(sel[0].Selected, 10, 64)
	if err != nil {
		return "", err
	}

	val++

	q = `UPDATE tags.events SET selected = $1
		WHERE title = $2
`
	_, err = e.db.Exec(q, val, t)
	if err != nil {
		return "", err
	}

	var res []models.Events
	const query = `
		SELECT title, address, tags, png
		FROM tags.events WHERE title = $1
`

	err = e.db.SelectContext(ctx, &res, query, t)
	if err != nil {
		return "", err
	}

	return res[0].Tags, nil
}

func NewEventsRepo(db *sqlx.DB) *Events {
	return &Events{
		db: db,
	}
}
