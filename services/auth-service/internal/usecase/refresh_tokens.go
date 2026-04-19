package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sssoultrix/golive/services/auth-service/internal/domain"
)

type RefreshTokenReader interface {
	GetToken(ctx context.Context, tokenHash string) (userID uuid.UUID, expiresAt time.Time, found bool, err error)
}

type RefreshTokenDeleter interface {
	DeleteToken(ctx context.Context, tokenHash string) error
}

type RefreshTokenCacheReader interface {
	Get(ctx context.Context, tokenHash string) (userID uuid.UUID, found bool, err error)
}

type RefreshTokenCacheDeleter interface {
	Delete(ctx context.Context, tokenHash string) error
}

type RefreshTokens struct {
	reader             RefreshTokenReader
	deleter            RefreshTokenDeleter
	cacheReader        RefreshTokenCacheReader
	cacheDeleter       RefreshTokenCacheDeleter
	creator            RefreshTokenRepository
	cacheWriter        RefreshTokenCache
	tokens             TokenGenerator
	refreshTokenHasher RefreshTokenHasher
	refreshTTL         time.Duration
}

func NewRefreshTokens(
	reader RefreshTokenReader,
	deleter RefreshTokenDeleter,
	cacheReader RefreshTokenCacheReader,
	cacheDeleter RefreshTokenCacheDeleter,
	creator RefreshTokenRepository,
	cacheWriter RefreshTokenCache,
	tokens TokenGenerator,
	refreshTokenHasher RefreshTokenHasher,
	refreshTTL time.Duration,
) *RefreshTokens {
	return &RefreshTokens{
		reader:             reader,
		deleter:            deleter,
		cacheReader:        cacheReader,
		cacheDeleter:       cacheDeleter,
		creator:            creator,
		cacheWriter:        cacheWriter,
		tokens:             tokens,
		refreshTokenHasher: refreshTokenHasher,
		refreshTTL:         refreshTTL,
	}
}

func (r *RefreshTokens) Execute(ctx context.Context, refreshToken string) (TokenPair, error) {
	oldHash := r.refreshTokenHasher.Hash(refreshToken)

	userID, found, err := r.cacheReader.Get(ctx, oldHash)
	if err != nil {
		return TokenPair{}, err
	}
	if !found {
		dbUserID, expiresAt, found, err := r.reader.GetToken(ctx, oldHash)
		if err != nil {
			return TokenPair{}, err
		}
		if !found || time.Now().UTC().After(expiresAt) {
			return TokenPair{}, domain.ErrInvalidRefreshToken
		}
		userID = dbUserID
	}

	// Rotate refresh token on every refresh.
	pair, err := r.tokens.GeneratePair(userID)
	if err != nil {
		return TokenPair{}, err
	}

	newHash := r.refreshTokenHasher.Hash(pair.RefreshToken)
	expiresAt := time.Now().UTC().Add(r.refreshTTL)
	if err := r.creator.CreateToken(ctx, userID, newHash, expiresAt); err != nil {
		return TokenPair{}, err
	}

	ttl := time.Until(expiresAt)
	if ttl > 0 {
		if err := r.cacheWriter.Set(ctx, newHash, userID, ttl); err != nil {
			return TokenPair{}, err
		}
	}

	// Best-effort revocation of the previous refresh token.
	_ = r.deleter.DeleteToken(ctx, oldHash)
	_ = r.cacheDeleter.Delete(ctx, oldHash)

	return pair, nil
}

