package main

import (
	"context"
	"log"

	"github.com/moeryomenko/healing"
	"github.com/moeryomenko/squad"
	"go.uber.org/zap"

	"github.com/moeryomenko/highload-architect-otus/social/internal/config"
	"github.com/moeryomenko/highload-architect-otus/social/internal/logger"
	"github.com/moeryomenko/highload-architect-otus/social/internal/repository"
	"github.com/moeryomenko/highload-architect-otus/social/internal/router"
	"github.com/moeryomenko/highload-architect-otus/social/internal/services"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	zapLog, err := logger.InitLogger(cfg)
	if err != nil {
		log.Fatalf("could not init logger: %v", err)
	}
	defer zapLog.Sync()

	writePool, readPool, err := repository.InitConnPool(context.Background(), cfg)
	if err != nil {
		zapLog.With(zap.Error(err)).Fatal("could not init database connection pool")
	}

	login := services.NewLogin(cfg, repository.NewLogin(writePool, readPool))

	server := router.NewRouter(cfg, zapLog, login, repository.NewUsers(writePool, readPool))

	healthController := healing.New(cfg.Health.Port,
		healing.WithCheckPeriod(cfg.Health.Period),
		healing.WithHealthzEndpoint(cfg.Health.LiveEndpoint),
		healing.WithReadyEndpoint(cfg.Health.ReadyEndpoint),
	)

	healthController.AddReadyChecker("mysql_master", writePool.CheckReadinessProber)
	healthController.AddReadyChecker("mysql_slave", readPool.CheckReadinessProber)

	group, err := squad.New(squad.WithSignalHandler())
	if err != nil {
		zapLog.With(zap.Error(err)).Fatal("could not create execution group")
	}

	group.RunGracefully(squad.RunServer(server))
	group.RunGracefully(healthController.Heartbeat, healthController.Stop)

	errs := group.Wait()
	for _, err := range errs {
		zapLog.With(zap.Error(err)).Error("service down with error")
	}
}
