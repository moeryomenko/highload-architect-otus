package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/moeryomenko/highload-architect-otus/social/internal/config"
	log "github.com/moeryomenko/highload-architect-otus/social/internal/logger"
	"github.com/moeryomenko/highload-architect-otus/social/internal/repository"
	"github.com/moeryomenko/highload-architect-otus/social/internal/services"
)

type Server struct {
	auth *Auth
}

func NewRouter(
	cfg *config.Config,
	logger *zap.Logger,
	login *services.Login,
	repo *repository.Users,
) *http.Server {
	router := chi.NewRouter()
	auth := NewAuth(cfg)

	router.Use(log.Logger(logger.Named("router")))
	// filesDir := http.FileServer(http.Dir("./assets"))
	// router.Handle("/swagger/*", http.StripPrefix("/swagger", filesDir))

	return &http.Server{
		Handler: HandlerWithOptions(&Service{
			auth:  auth,
			login: login,
			users: repo,
		}, ChiServerOptions{
			BaseURL:     cfg.APIBaseURL,
			BaseRouter:  router,
			Middlewares: []MiddlewareFunc{auth.Auth()},
		}),
		Addr: cfg.Addr(),
	}
}
