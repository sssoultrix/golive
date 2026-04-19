package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sssoultrix/golive/services/auth-service/internal/domain"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, login, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`

	if _, err := ur.db.Exec(ctx, query, user.ID, user.Login, user.PasswordHash, user.CreatedAt, user.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetCredentialsByLogin(ctx context.Context, login string) (uuid.UUID, string, error) {
	query := `SELECT id, password_hash FROM users WHERE login = $1`
	var (
		id           uuid.UUID
		passwordHash string
	)
	err := ur.db.QueryRow(ctx, query, login).Scan(&id, &passwordHash)
	if err != nil {
		return uuid.Nil, "", err
	}
	return id, passwordHash, nil
}
