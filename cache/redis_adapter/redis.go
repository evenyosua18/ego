package redis_adapter

import (
	"context"
	"time"

	"github.com/evenyosua18/ego/code"
	redis "github.com/redis/go-redis/v9"
)

var (
	ErrCacheMiss      = code.Get(code.CacheNotFound).SetErrorMessage("cache not found")
	ErrRedisClientNil = code.Get(code.CacheError).SetErrorMessage("redis client is nil")
	ErrCache          = code.Get(code.CacheError)
)

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(client *redis.Client) *RedisAdapter {
	return &RedisAdapter{client: client}
}

func (r *RedisAdapter) Get(ctx context.Context, key string) ([]byte, error) {
	if r.client == nil {
		return nil, ErrRedisClientNil
	}

	val, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, ErrCacheMiss
	}
	return val, err
}

func (r *RedisAdapter) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if r.client == nil {
		return ErrRedisClientNil
	}

	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return code.Wrap(err, code.CacheError)
	}

	return nil
}

func (r *RedisAdapter) Delete(ctx context.Context, key string) error {
	if r.client == nil {
		return ErrRedisClientNil
	}
	return r.client.Del(ctx, key).Err()
}
