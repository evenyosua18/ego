package cache

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/evenyosua18/ego/code"
	"golang.org/x/sync/singleflight"
)

var (
	ErrCacheMiss = code.Get(code.CacheNotFound).SetErrorMessage("cache not found")
	manager      ICacheManager
	once         sync.Once
)

type ICacheManager interface {
	GetOrSetRaw(ctx context.Context, key string, ttl time.Duration, target any, fetcher func() (any, error)) error
	Invalidate(ctx context.Context, keys ...string) error
}

type cacheManagerImpl struct {
	store Store
	group singleflight.Group
}

func InitCacheManager(store Store) {
	once.Do(func() {
		manager = &cacheManagerImpl{
			store: store,
		}
	})
}

func GetManager() ICacheManager {
	if manager == nil {
		panic("cache manager is not initialized. call cache.InitCacheManager() first.")
	}
	return manager
}

func (m *cacheManagerImpl) GetOrSetRaw(ctx context.Context, key string, ttl time.Duration, target any, fetcher func() (any, error)) error {
	// get from cache
	valBytes, err := m.store.Get(ctx, key)
	if err == nil {
		if json.Unmarshal(valBytes, target) == nil {
			return nil
		}
	} else if err != ErrCacheMiss {
		// error with cache
		return err
	}

	// get from db
	res, err, _ := m.group.Do(key, func() (interface{}, error) {
		data, err := fetcher()
		if err != nil {
			return nil, err
		}

		// set to cache
		bytes, _ := json.Marshal(data)
		err = m.store.Set(ctx, key, bytes, ttl)
		if err != nil {
			return nil, err
		}

		return data, nil
	})

	if err != nil {
		return err
	}

	// We use json marshaling to assign the result safely into the pointer target interface
	if bytes, err := json.Marshal(res); err == nil {
		json.Unmarshal(bytes, target)
	}

	return nil
}

func GetOrSet[T any](ctx context.Context, m ICacheManager, key string, ttl time.Duration, fetcher func() (T, error)) (T, error) {
	var zero T
	var res T

	err := m.GetOrSetRaw(ctx, key, ttl, &res, func() (any, error) {
		return fetcher()
	})
	if err != nil {
		return zero, err
	}

	return res, nil
}

func (m *cacheManagerImpl) Invalidate(ctx context.Context, keys ...string) error {
	for _, k := range keys {
		_ = m.store.Delete(ctx, k)
	}
	return nil
}
