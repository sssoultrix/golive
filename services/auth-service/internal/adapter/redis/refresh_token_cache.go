package redis

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RefreshTokenCache struct {
	r *redis.Client
}

func NewRefreshTokenCache(r *redis.Client) *RefreshTokenCache {
	return &RefreshTokenCache{r: r}
}

func (c *RefreshTokenCache) Set(ctx context.Context, tokenHash string, userID uuid.UUID, ttl time.Duration) error {
	key := "rt:" + tokenHash
	if ttl <= 0 {
		// Avoid creating persistent keys (go-redis uses ttl=0 as "persist").
		return c.r.Del(ctx, key).Err()
	}
	return c.r.Set(ctx, key, userID.String(), ttl).Err()
}

func (c *RefreshTokenCache) Get(ctx context.Context, tokenHash string) (uuid.UUID, bool, error) {
	key := "rt:" + tokenHash
	s, err := c.r.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return uuid.Nil, false, nil
		}
		return uuid.Nil, false, err
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, false, err
	}
	return id, true, nil
}

func (c *RefreshTokenCache) Delete(ctx context.Context, tokenHash string) error {
	key := "rt:" + tokenHash
	return c.r.Del(ctx, key).Err()
}
