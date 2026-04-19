package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sssoultrix/golive/services/auth-service/internal/config"
)

func main() {
	if err := run(); err != nil {
		slog.Error("cleanup-refresh-tokens failed", slog.Any("err", err))
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	const q = `DELETE FROM refresh_tokens WHERE expires_at < now()`
	res, err := db.Exec(ctx, q)
	if err != nil {
		return err
	}

	slog.Info("cleanup-refresh-tokens done", slog.Int64("deleted", res.RowsAffected()))
	return nil
}

