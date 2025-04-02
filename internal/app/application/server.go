package application

import (
	"github.com/aviseu/go-sample/internal/app/application/api"
	"github.com/aviseu/go-sample/internal/app/domain"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"time"
)

type Config struct {
	Host            string        `default:"0.0.0.0:8080"`
	ShutdownTimeout time.Duration `default:"30s"`
}

func SetupServer(cfg Config, h http.Handler) http.Server {
	return http.Server{
		Addr:    cfg.Host,
		Handler: h,
	}
}

func APIHandler(log *slog.Logger, s *domain.Service, r api.Repository) http.Handler {
	router := chi.NewRouter()

	h := api.NewHandler(log, s, r)
	router.Mount("/api/tasks", h.Routes())

	return router
}
