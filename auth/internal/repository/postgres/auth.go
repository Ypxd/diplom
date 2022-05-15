package postgres

import (
	"context"
	"errors"
	"github.com/Ypxd/diplom/auth/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strconv"
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
	i, err := strconv.ParseInt(*req.Age, 10, 64)
	if err != nil {
		return 0, errors.New("неправильное поле возраста")
	}
	err = a.db.SelectContext(ctx, &age, query, i)
	if err != nil {
		return 0, err
	}
	if len(age) == 1 {
		return age[0], nil
	}

	return 0, errors.New("неправильное поле возраста")
}

func (a *Auth) ChangePass(ctx context.Context, req models.ChangePassReq, userID string) error {
	const query = `
		UPDATE auth.users SET password = $1
		WHERE user_id = $3 AND password = $2
		RETURNING user_id
`
	usrID := make([]uuid.UUID, 0)
	err := a.db.SelectContext(ctx, &usrID, query, req.NewPassword, req.OldPassword, userID)
	if err != nil {
		return err
	}
	if len(usrID) == 1 {
		return nil
	}

	return errors.New("неправильный пароль")
}

func (a *Auth) Register(ctx context.Context, req models.AuthReq, age int64) (string, error) {
	const query = `
		INSERT INTO auth."users" (user_id, login, password, email, name, age_id)
		VALUES ($1, $2, $3, $4, $5, $6)
`

	_, err := a.db.ExecContext(ctx, query, uuid.New(), req.Login, req.Password, req.Email, req.Name, age)
	if errPq, ok := err.(*pq.Error); ok {
		myError := errors.New("такой " + errPq.Constraint + " уже существует")
		return "", myError
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

	return nil, errors.New("неправильный логин или пароль")
}

func (a *Auth) UserInfo(ctx context.Context, userID string) (*models.UserInfo, error) {
	var usr []models.UserInfo
	const query = `
		SELECT u.login, u.name, u.email, u.f_tags, u.uf_tags, a.title AS age FROM auth.users u
		JOIN auth.age a ON a.id = u.age_id
		WHERE user_id = $1
`

	err := a.db.SelectContext(ctx, &usr, query, userID)
	if err != nil {
		return nil, err
	}
	if len(usr) == 1 {
		return &usr[0], nil
	}

	return nil, errors.New("unexpected error")
}

func NewAuthRepo(db *sqlx.DB) *Auth {
	return &Auth{
		db: db,
	}
}
