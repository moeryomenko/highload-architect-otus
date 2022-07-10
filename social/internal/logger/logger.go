package logger

import (
	"go.uber.org/zap"

	"github.com/moeryomenko/highload-architect-otus/social/internal/config"
)

// InitLogger returns logger.
// TODO: rewrite for custom and more detail configurations.
func InitLogger(cfg *config.Config) (logger *zap.Logger, err error) {
	if cfg.Logger.IsDevelopment {
		return zap.NewDevelopment(zap.AddCaller())
	}

	return zap.NewProduction()
}
