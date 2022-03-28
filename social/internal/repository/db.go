package repository

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/moeryomenko/highload-architect-otus/social/internal/config"
)

// InitConnPool initialize and configure db connection pool.
func InitConnPool(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.Database.DSN())
	if err != nil {
		return nil, err
	}
	setupPool(cfg.Database.Pool, db)

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// setupPool configures connection pool.
func setupPool(cfg *config.PoolConfig, pool *sql.DB) {
	pool.SetMaxOpenConns(cfg.MaxOpenConns)
	pool.SetMaxIdleConns(cfg.MaxIdleConns)
	pool.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
}
