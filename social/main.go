package main

import (
	"context"
	"errors"
	stdLog "log"
	"net/http"

	"github.com/moeryomenko/healing"
	"github.com/moeryomenko/squad"
	"go.uber.org/zap"

	"github.com/moeryomenko/highload-architect-otus/social/internal/config"
	log "github.com/moeryomenko/highload-architect-otus/social/internal/logger"
	"github.com/moeryomenko/highload-architect-otus/social/internal/repository"
	"github.com/moeryomenko/highload-architect-otus/social/internal/router"
	"github.com/moeryomenko/highload-architect-otus/social/internal/services"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		stdLog.Fatalf("could not load config: %v", err)
	}

	logger, err := log.InitLogger(cfg)
	if err != nil {
		stdLog.Fatalf("could not init logger: %v", err)
	}
	defer logger.Sync()

	writePool, readPool, err := repository.InitConnPool(context.Background(), cfg)
	if err != nil {
		logger.With(zap.Error(err)).Fatal("could not init database connection pool")
	}

	login := services.NewLogin(cfg, repository.NewLogin(writePool, readPool))

	server := router.NewRouter(cfg, logger, login, repository.NewUsers(writePool, readPool))

	healthController := healing.New(
		healing.WithCheckPeriod(cfg.Health.Period),
		healing.WithHealthzEndpoint(cfg.Health.LiveEndpoint),
		healing.WithReadyEndpoint(cfg.Health.ReadyEndpoint),
	)

	healthController.AddReadyChecker("mysql_master", writePool.CheckReadinessProber)
	healthController.AddReadyChecker("mysql_slave", readPool.CheckReadinessProber)

	group, err := squad.NewSquad(context.Background(), squad.WithSignalHandler())
	if err != nil {
		logger.With(zap.Error(err)).Fatal("could not create execution group")
	}

	group.RunGracefully(func(ctx context.Context) error {
		return server.ListenAndServe()
	}, func(ctx context.Context) error {
		err := server.Shutdown(ctx)
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	})

	// Run health controller for monitor liveness of service.
	group.Run(healthController.Heartbeat)
	group.RunGracefully(healthController.ListenAndServe(cfg.Health.Port), healthController.Shutdown)

	errs := group.Wait()
	for _, err := range errs {
		logger.With(zap.Error(err)).Error("service down with error")
	}
}
