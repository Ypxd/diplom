package shared

import (
	"database/sql"
	"log"

	"github.com/pressly/goose/v3"
)

func MakeMigration(db *sql.DB, args ...string) {
	err := goose.Run(args[0], db, "migrations", args...)
	if err != nil {
		log.Fatalln(err)
	}
}
