package repository

import (
	"github.com/go-mysql-org/go-mysql/client"

	"github.com/moeryomenko/highload-architect-otus/social/internal/config"
)

// InitConnPool initialize and configure db connection pool.
func InitConnPool(cfg *config.Config) (*client.Pool, error) {
	pool := client.NewPool(
		noopLog,
		cfg.Database.Pool.MaxOpenConns/2,
		cfg.Database.Pool.MaxOpenConns,
		cfg.Database.Pool.MaxIdleConns,
		cfg.Database.Addr(),
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)
	return pool, nil
}

var noopLog = func(format string, args ...interface{}) {}
