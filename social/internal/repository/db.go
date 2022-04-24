package repository

import (
	"context"

	"github.com/moeryomenko/healing/decorators/mysql"

	"github.com/moeryomenko/highload-architect-otus/social/internal/config"
)

// InitConnPool initialize and configure db connection pool.
func InitConnPool(ctx context.Context, cfg *config.Config) (*mysql.Pool, error) {
	return mysql.New(ctx, mysql.Config{
		Host:     cfg.Database.Host,
		Port:     uint16(cfg.Database.Port),
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.Name,
	}, mysql.WithPoolConfig(mysql.PoolConfig{
		MinAlive: cfg.Database.Pool.MaxOpenConns / 2,
		MaxAlive: cfg.Database.Pool.MaxOpenConns,
		MaxIdle:  cfg.Database.Pool.MaxIdleConns,
	}))
}
