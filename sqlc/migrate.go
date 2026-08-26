package sqlc

import (
	"database/sql"
	"embed"
	"log"

	goose "github.com/pressly/goose/v3"
	_ "modernc.org/sqlite" // pure-Go sqlite driver (no CGO needed)
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func RunMigrations(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal(err)

	}

	return nil
}
