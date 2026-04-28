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
	"github.com/redis/go-redis/v9"
	"github.com/sssoultrix/golive/services/auth-service/internal/adapter/bcrypt"
	"github.com/sssoultrix/golive/services/auth-service/internal/adapter/jwt_manager"
	"github.com/sssoultrix/golive/services/auth-service/internal/adapter/postgres"
	rediscache "github.com/sssoultrix/golive/services/auth-service/internal/adapter/redis"
	"github.com/sssoultrix/golive/services/auth-service/internal/adapter/refresh_token"
	"github.com/sssoultrix/golive/services/auth-service/internal/config"
	httpt "github.com/sssoultrix/golive/services/auth-service/internal/transport/http"
	"github.com/sssoultrix/golive/services/auth-service/internal/usecase"
)

func main() {
	if err := run(); err != nil {
		slog.Error("auth-service exited with error", slog.Any("err", err))
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

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	defer func() { _ = rdb.Close() }()

	userRepo := postgres.NewUserRepository(db)
	refreshRepo := postgres.NewRefreshTokenRepository(db)
	refreshCache := rediscache.NewRefreshTokenCache(rdb)

	passwordHasher := bcrypt.NewPasswordHasher(cfg.BcryptCost)
	refreshHasher := refresh_token.NewRefreshTokenHasher(cfg.RefreshTokenPepper)
	tokenGen := jwt_manager.NewPairGenerator([]byte(cfg.JWTSecret), cfg.AccessTTL, 32)
	accessValidator := jwt_manager.NewAccessValidator([]byte(cfg.JWTSecret))

	registerUC := usecase.NewRegisterUser(
		userRepo,
		refreshRepo,
		refreshCache,
		passwordHasher,
		tokenGen,
		refreshHasher,
		cfg.RefreshTTL,
	)
	loginUC := usecase.NewLoginUser(
		userRepo,
		passwordHasher,
		refreshRepo,
		refreshCache,
		tokenGen,
		refreshHasher,
		cfg.RefreshTTL,
	)
	refreshUC := usecase.NewRefreshTokens(
		refreshRepo,
		refreshRepo,
		refreshCache,
		refreshCache,
		refreshRepo,
		refreshCache,
		tokenGen,
		refreshHasher,
		cfg.RefreshTTL,
	)
	logoutUC := usecase.NewLogout(
		refreshRepo,
		refreshRepo,
		refreshCache,
		refreshCache,
		refreshHasher,
	)

	router := httpt.NewRouter()
	router.SetupRoutes(registerUC, loginUC, refreshUC, logoutUC, accessValidator)
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
