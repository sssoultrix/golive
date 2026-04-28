package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sssoultrix/golive/services/auth-service/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	CreateUserWithOutbox(ctx context.Context, user *domain.User, event *domain.OutboxEvent) error
}

type OutboxRepository interface {
	CreateEvent(ctx context.Context, event *domain.OutboxEvent) error
}

type PasswordHasher interface {
	GenerateHash(password string) (string, error)
}

type TokenGenerator interface {
	GeneratePair(userID uuid.UUID) (TokenPair, error)
}

type RefreshTokenRepository interface {
	CreateToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
}

type RefreshTokenCache interface {
	Set(ctx context.Context, tokenHash string, userID uuid.UUID, ttl time.Duration) error
}

type RefreshTokenHasher interface {
	Hash(token string) string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type RegisterUser struct {
	userRepo           UserRepository
	outboxRepo         OutboxRepository
	refreshTokenRepo   RefreshTokenRepository
	refreshTokenCache  RefreshTokenCache
	passwordHasher     PasswordHasher
	tokens             TokenGenerator
	refreshTokenHasher RefreshTokenHasher
	refreshTTL         time.Duration
}

func NewRegisterUser(
	userRepo UserRepository,
	outboxRepo OutboxRepository,
	refreshTokenRepo RefreshTokenRepository,
	refreshTokenCache RefreshTokenCache,
	passwordHasher PasswordHasher,
	tokens TokenGenerator,
	refreshTokenHasher RefreshTokenHasher,
	refreshTTL time.Duration,
) *RegisterUser {
	return &RegisterUser{
		userRepo:           userRepo,
		outboxRepo:         outboxRepo,
		refreshTokenRepo:   refreshTokenRepo,
		refreshTokenCache:  refreshTokenCache,
		passwordHasher:     passwordHasher,
		tokens:             tokens,
		refreshTokenHasher: refreshTokenHasher,
		refreshTTL:         refreshTTL,
	}
}

func (r *RegisterUser) Execute(ctx context.Context, login, password string) (uuid.UUID, TokenPair, error) {
	passwordHash, err := r.passwordHasher.GenerateHash(password)
	if err != nil {
		return uuid.Nil, TokenPair{}, err
	}

	user, err := domain.NewUser(login, passwordHash)
	if err != nil {
		return uuid.Nil, TokenPair{}, err
	}

	event := domain.NewOutboxEvent("UserRegistered", map[string]interface{}{
		"user_id": user.ID.String(),
		"login":   user.Login,
	})

	if err := r.userRepo.CreateUserWithOutbox(ctx, user, event); err != nil {
		return uuid.Nil, TokenPair{}, err
	}

	tokenPair, err := r.tokens.GeneratePair(user.ID)
	if err != nil {
		return uuid.Nil, TokenPair{}, err
	}

	refreshHash := r.refreshTokenHasher.Hash(tokenPair.RefreshToken)
	expiresAt := time.Now().UTC().Add(r.refreshTTL)
	if err := r.refreshTokenRepo.CreateToken(ctx, user.ID, refreshHash, expiresAt); err != nil {
		return uuid.Nil, TokenPair{}, err
	}

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return user.ID, tokenPair, nil
	}
	if err := r.refreshTokenCache.Set(ctx, refreshHash, user.ID, ttl); err != nil {
		return uuid.Nil, TokenPair{}, err
	}

	return user.ID, tokenPair, nil
}
