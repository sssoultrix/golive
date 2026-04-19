package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sssoultrix/golive/services/auth-service/internal/domain"
)

type UserCredentialsRepository interface {
	GetCredentialsByLogin(ctx context.Context, login string) (userID uuid.UUID, passwordHash string, err error)
}

type PasswordVerifier interface {
	Verify(password, passwordHash string) error
}

type LoginUser struct {
	users              UserCredentialsRepository
	verifier           PasswordVerifier
	refreshTokenRepo   RefreshTokenRepository
	refreshTokenCache  RefreshTokenCache
	tokens             TokenGenerator
	refreshTokenHasher RefreshTokenHasher
	refreshTTL         time.Duration
}

func NewLoginUser(
	users UserCredentialsRepository,
	verifier PasswordVerifier,
	refreshTokenRepo RefreshTokenRepository,
	refreshTokenCache RefreshTokenCache,
	tokens TokenGenerator,
	refreshTokenHasher RefreshTokenHasher,
	refreshTTL time.Duration,
) *LoginUser {
	return &LoginUser{
		users:              users,
		verifier:           verifier,
		refreshTokenRepo:   refreshTokenRepo,
		refreshTokenCache:  refreshTokenCache,
		tokens:             tokens,
		refreshTokenHasher: refreshTokenHasher,
		refreshTTL:         refreshTTL,
	}
}

func (l *LoginUser) Execute(ctx context.Context, login, password string) (uuid.UUID, TokenPair, error) {
	userID, passwordHash, err := l.users.GetCredentialsByLogin(ctx, login)
	if err != nil {
		return uuid.Nil, TokenPair{}, err
	}

	if err := l.verifier.Verify(password, passwordHash); err != nil {
		return uuid.Nil, TokenPair{}, domain.ErrInvalidCredentials
	}

	tokenPair, err := l.tokens.GeneratePair(userID)
	if err != nil {
		return uuid.Nil, TokenPair{}, err
	}

	refreshHash := l.refreshTokenHasher.Hash(tokenPair.RefreshToken)
	expiresAt := time.Now().UTC().Add(l.refreshTTL)

	if err := l.refreshTokenRepo.CreateToken(ctx, userID, refreshHash, expiresAt); err != nil {
		return uuid.Nil, TokenPair{}, err
	}

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return userID, tokenPair, nil
	}
	if err := l.refreshTokenCache.Set(ctx, refreshHash, userID, ttl); err != nil {
		return uuid.Nil, TokenPair{}, err
	}

	return userID, tokenPair, nil
}

