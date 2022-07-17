package main

import (
	"context"
	"log"

	"github.com/moeryomenko/healing"
	"github.com/moeryomenko/healing/checkers"
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

	writePool, readPool := repository.InitConnPool(context.Background(), cfg)

	login := services.NewLogin(cfg, repository.NewLogin(writePool, readPool))

	server := router.NewRouter(cfg, zapLog, login, repository.NewUsers(writePool, readPool))

	healthController := healing.New(cfg.Health.Port,
		healing.WithCheckPeriod(cfg.Health.Period),
		healing.WithHealthzEndpoint(cfg.Health.LiveEndpoint),
		healing.WithReadyEndpoint(cfg.Health.ReadyEndpoint),
		healing.WithMetrics(`/metrics`),
		healing.WithPProf(),
	)

	healthController.AddReadyChecker("mysql_master", checkers.MySQLReadinessProber(writePool, cfg.Database.Pool.MaxOpenConns))
	healthController.AddReadyChecker("mysql_slaves", checkers.MySQLReadinessProber(readPool, cfg.Database.Pool.MaxOpenConns))

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
