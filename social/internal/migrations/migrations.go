package migrations

import (
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/moeryomenko/highload-architect-otus/social/internal/config"
)

func Up(cfg *config.Config, conn *sql.DB) error {
	instance, err := mysql.WithInstance(conn, &mysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(cfg.Database.MigrationDir, "mysql", instance)
	if err != nil {
		return err
	}
	err = m.Steps(2)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
