package usecase

import (
	"context"
	"time"

	"github.com/sssoultrix/golive/services/auth-service/internal/domain"
)

type Logout struct {
	reader       RefreshTokenReader
	deleter      RefreshTokenDeleter
	cacheReader  RefreshTokenCacheReader
	cacheDeleter RefreshTokenCacheDeleter
	hasher       RefreshTokenHasher
}

func NewLogout(
	reader RefreshTokenReader,
	deleter RefreshTokenDeleter,
	cacheReader RefreshTokenCacheReader,
	cacheDeleter RefreshTokenCacheDeleter,
	hasher RefreshTokenHasher,
) *Logout {
	return &Logout{
		reader:       reader,
		deleter:      deleter,
		cacheReader:  cacheReader,
		cacheDeleter: cacheDeleter,
		hasher:       hasher,
	}
}

func (l *Logout) Execute(ctx context.Context, refreshToken string) error {
	hash := l.hasher.Hash(refreshToken)

	_, found, err := l.cacheReader.Get(ctx, hash)
	if err != nil {
		return err
	}
	if !found {
		_, expiresAt, found, err := l.reader.GetToken(ctx, hash)
		if err != nil {
			return err
		}
		if !found || time.Now().UTC().After(expiresAt) {
			return domain.ErrInvalidRefreshToken
		}
	}

	if err := l.deleter.DeleteToken(ctx, hash); err != nil {
		return err
	}
	_ = l.cacheDeleter.Delete(ctx, hash)
	return nil
}

