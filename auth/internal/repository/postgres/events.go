package postgres

import (
	"context"
	"github.com/Ypxd/diplom/auth/internal/models"
	"github.com/jmoiron/sqlx"
)

type Events struct {
	db *sqlx.DB
}

func (e *Events) AllEvents(ctx context.Context) ([]models.Events, error) {
	var res []models.Events
	const query = `
		SELECT title, address, png
		FROM tags.events
`

	err := e.db.SelectContext(ctx, &res, query)
	if err != nil {
		return res, err
	}

	return res, nil
}

func NewEventsRepo(db *sqlx.DB) *Events {
	return &Events{
		db: db,
	}
}
