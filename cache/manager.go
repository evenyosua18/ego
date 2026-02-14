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
	manager      *CacheManager
	once         sync.Once
)

type CacheManager struct {
	store Store
	group singleflight.Group
}

func InitCacheManager(store Store) {
	once.Do(func() {
		manager = &CacheManager{
			store: store,
		}
	})
}

func GetManager() *CacheManager {
	if manager == nil {
		panic("cache manager is not initialized. call cache.InitCacheManager() first.")
	}
	return manager
}

func GetOrSet[T any](ctx context.Context, m *CacheManager, key string, ttl time.Duration, fetcher func() (T, error)) (T, error) {
	var zero T

	// get from cache
	valBytes, err := m.store.Get(ctx, key)
	if err == nil {
		var res T
		if json.Unmarshal(valBytes, &res) == nil {
			return res, nil
		}
	} else if err != ErrCacheMiss {
		// error with cache
		return zero, err
	}

	// get from db
	res, err, _ := m.group.Do(key, func() (interface{}, error) {
		data, err := fetcher()
		if err != nil {
			return zero, err
		}

		// set to cache
		bytes, _ := json.Marshal(data)
		err = m.store.Set(ctx, key, bytes, ttl)
		if err != nil {
			return zero, err
		}

		return data, nil
	})

	if err != nil {
		return zero, err
	}

	return res.(T), nil
}

func (m *CacheManager) Invalidate(ctx context.Context, keys ...string) error {
	for _, k := range keys {
		_ = m.store.Delete(ctx, k)
	}
	return nil
}
