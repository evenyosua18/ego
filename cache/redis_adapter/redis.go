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

type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	MaxRetries   int
	MinIdleConns int
	PoolSize     int
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DialTimeout  time.Duration
	PoolTimeout  time.Duration
}

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(config RedisConfig) (*RedisAdapter, error) {
	// create new redis client
	client := redis.NewClient(&redis.Options{
		Addr:            config.Addr,
		Password:        config.Password,
		DB:              config.DB,
		MaxRetries:      config.MaxRetries,
		MinIdleConns:    config.MinIdleConns,
		PoolSize:        config.PoolSize,
		ConnMaxIdleTime: config.IdleTimeout,
		ReadTimeout:     config.ReadTimeout,
		WriteTimeout:    config.WriteTimeout,
		DialTimeout:     config.DialTimeout,
		PoolTimeout:     config.PoolTimeout,
	})

	// check client
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, code.Wrap(err, code.CacheError)
	}

	return &RedisAdapter{client: client}, nil
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
