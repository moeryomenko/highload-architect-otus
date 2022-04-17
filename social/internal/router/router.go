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
	filesDir := http.FileServer(http.Dir(cfg.AssetsDir))
	router.Handle("/*", filesDir)

	return &http.Server{
		Handler: HandlerWithOptions(&Service{
			auth:   auth,
			login:  login,
			users:  repo,
			logger: logger,
		}, ChiServerOptions{
			BaseURL:     cfg.APIBaseURL,
			BaseRouter:  router,
			Middlewares: []MiddlewareFunc{auth.Auth(), Json()},
		}),
		Addr: cfg.Addr(),
	}
}

func Json() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			next(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
