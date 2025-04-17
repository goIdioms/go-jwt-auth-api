package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheValue struct {
	UserID       string `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

const refreshPrefix = "refresh_token:"

func (r *RedisCache) SaveRefreshToken(ctx context.Context, sessionID string, value CacheValue, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}
	key := refreshPrefix + sessionID
	return r.rdb.Set(ctx, key, data, ttl).Err()
}

func (r *RedisCache) GetRefreshToken(ctx context.Context, sessionID string) (*CacheValue, error) {
	key := refreshPrefix + sessionID
	data, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}
	var value CacheValue
	if err := json.Unmarshal([]byte(data), &value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache value: %w", err)
	}

	return &value, nil
}

func (r *RedisCache) DeleteRefreshToken(ctx context.Context, sessionID string) error {
	key := refreshPrefix + sessionID
	return r.rdb.Del(ctx, key).Err()
}
