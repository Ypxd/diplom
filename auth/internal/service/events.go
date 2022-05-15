package service

import (
	"context"
	"github.com/Ypxd/diplom/auth/internal/models"

	"github.com/Ypxd/diplom/auth/internal/repository"
	"github.com/jmoiron/sqlx"
)

type EventsService struct {
	repo *repository.Repository
	conn *sqlx.DB
}

func (e *EventsService) GetAllEvents(ctx context.Context) ([]models.Events, error) {
	return e.repo.Events.AllEvents(ctx)
}

func NewEventsService(repo *repository.Repository, conn *sqlx.DB) *EventsService {
	return &EventsService{
		repo: repo,
		conn: conn,
	}
}
