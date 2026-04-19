package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) CreateToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, uuid.New(), userID, tokenHash, expiresAt, time.Now().UTC())
	return err
}

func (r *RefreshTokenRepository) GetToken(ctx context.Context, tokenHash string) (uuid.UUID, time.Time, bool, error) {
	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE token_hash = $1`
	var (
		userID    uuid.UUID
		expiresAt time.Time
	)
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(&userID, &expiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, time.Time{}, false, nil
		}
		return uuid.Nil, time.Time{}, false, err
	}
	return userID, expiresAt, true, nil
}

func (r *RefreshTokenRepository) DeleteToken(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`
	_, err := r.db.Exec(ctx, query, tokenHash)
	return err
}
