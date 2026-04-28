package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sssoultrix/golive/services/profile-service/internal/adapter/postgres"
	"github.com/sssoultrix/golive/services/profile-service/internal/config"
	httpt "github.com/sssoultrix/golive/services/profile-service/internal/transport/http"
	"github.com/sssoultrix/golive/services/profile-service/internal/usecase"
)

func main() {
	if err := run(); err != nil {
		slog.Error("profile-service exited with error", slog.Any("err", err))
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	profileRepo := postgres.NewProfileRepository(db)

	createUC := usecase.NewCreateProfileUseCase(profileRepo)
	getUC := usecase.NewGetProfileUseCase(profileRepo)
	updateUC := usecase.NewUpdateProfileUseCase(profileRepo)
	deleteUC := usecase.NewDeleteProfileUseCase(profileRepo)

	router := httpt.NewRouter()
	router.SetupRoutes(createUC, getUC, updateUC, deleteUC)
	r := router.GetEngine()

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
		return nil
	case err := <-errCh:
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

