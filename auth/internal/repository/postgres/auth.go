package postgres

import (
	"context"
	"errors"
	"github.com/Ypxd/diplom/auth/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Auth struct {
	db *sqlx.DB
}

func (a *Auth) Age(ctx context.Context, req models.AuthReq) (int64, error) {
	const query = `
		SELECT id FROM auth.age
		WHERE $1 > min AND $1 < max
`
	age := make([]int64, 0)
	err := a.db.SelectContext(ctx, &age, query, req.Age)
	if err != nil {
		return 0, err
	}
	if len(age) == 1 {
		return age[0], nil
	}

	return 0, errors.New("wrong Age")
}

func (a *Auth) Register(ctx context.Context, req models.AuthReq, age int64) (string, error) {
	const query = `
		INSERT INTO auth."users" (user_id, login, password, email, name, age_id)
		VALUES ($1, $2, $3, $4, $5, $6)
`

	_, err := a.db.ExecContext(ctx, query, uuid.New(), req.Login, req.Password, req.Email, req.Name, age)
	if errPq, ok := err.(*pq.Error); ok {
		myError := errors.New("duplicate " + errPq.Constraint)
		return myError.Error(), err
	}
	return "", err
}

func (a *Auth) Auth(ctx context.Context, req models.AuthReq) (*uuid.UUID, error) {
	const query = `
		SELECT user_id FROM auth."users"
		WHERE login = $1 AND password = $2
`
	userID := make([]uuid.UUID, 0)
	err := a.db.SelectContext(ctx, &userID, query, req.Login, req.Password)
	if err != nil {
		return nil, err
	}
	if len(userID) == 1 {
		return &userID[0], nil
	}

	return nil, errors.New("duplicate user id")
}

func NewAuthRepo(db *sqlx.DB) *Auth {
	return &Auth{
		db: db,
	}
}
