package postgres

import (
	"fmt"
	"github.com/Ypxd/diplom/auth/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const dbConnStr = "postgres://%s:%s@%s/%s?sslmode=disable"

func Connect() (*sqlx.DB, error) {
	cfg := utils.GetConfig().DB
	connStr := fmt.Sprintf(dbConnStr, cfg.User, cfg.Password,
		cfg.Address, cfg.DBName)

	db, err := sqlx.Open(cfg.DriverName, connStr)
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}
