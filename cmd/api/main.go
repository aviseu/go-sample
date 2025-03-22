package main

import (
	"context"
	"fmt"
	"github.com/aviseu/go-sample/internal/app/application"
	"github.com/aviseu/go-sample/internal/app/domain"
	"github.com/aviseu/go-sample/internal/app/infrastructure"
	"github.com/aviseu/go-sample/internal/app/infrastructure/postgres"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type config struct {
	Log struct {
		Level slog.Level `default:"info"`
	}
	DB  infrastructure.Config
	API application.Config
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})))

	ctx := context.Background()
	if err := run(ctx); err != nil {
		slog.Error(err.Error())
	}
}

func run(ctx context.Context) error {
	// load environment variables
	slog.Info("loading environment variables...")
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return fmt.Errorf("failed to process env vars: %w", err)
	}

	// set logging
	slog.Info("setting logging...")
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.Log.Level}))
	slog.SetDefault(log)

	// setup database
	log.Info("setting up database...")
	db, err := infrastructure.SetupDatabase(cfg.DB)
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	// Setup services & repositories
	log.Info("setting up services & repositories...")
	tr := postgres.NewTaskRepository(db)
	ts := domain.NewService(tr)

	// setup server
	log.Info("setting up server...")
	server := application.SetupServer(cfg.API, application.APIHandler(log, ts))
	serverErrors := make(chan error, 1)

	go func() {
		log.Info("starting up server...")
		serverErrors <- server.ListenAndServe()
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-done:
		log.Info("shutting down server...")

		ctx, cancel := context.WithTimeout(ctx, cfg.API.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}
	}

	return nil
}
