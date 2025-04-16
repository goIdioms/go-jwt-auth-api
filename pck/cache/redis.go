package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	rdb *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisCache{rdb: rdb}
}

func (rc *RedisCache) Set(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration) error {
	return rc.rdb.Set(ctx, key, value, expiration).Err()
}

func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return rc.rdb.Get(ctx, key).Result()
}

func (rc *RedisCache) Delete(ctx context.Context, key ...string) error {
	return rc.rdb.Del(ctx, key...).Err()
}
