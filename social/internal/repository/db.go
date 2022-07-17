package repository

import (
	"context"
	"fmt"

	"github.com/go-mysql-org/go-mysql/client"
	"github.com/moeryomenko/highload-architect-otus/social/internal/config"
)

// InitConnPool initialize and configure db connection pool.
func InitConnPool(ctx context.Context, cfg *config.Config) (write *client.Pool, read *client.Pool) {
	write = client.NewPool(
		func(format string, args ...interface{}) {},
		cfg.Database.Pool.MaxOpenConns,
		cfg.Database.Pool.MaxOpenConns,
		cfg.Database.Pool.MaxIdleConns,
		fmt.Sprintf("%s:%d", cfg.Database.Host, cfg.Database.MasterPort),
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.Name,
	)

	read = client.NewPool(
		func(format string, args ...interface{}) {},
		cfg.Database.Pool.MaxOpenConns,
		cfg.Database.Pool.MaxOpenConns,
		cfg.Database.Pool.MaxIdleConns,
		fmt.Sprintf("%s:%d", cfg.Database.Host, cfg.Database.SlavePort),
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.Name,
	)

	return write, read
}
