package cache

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStore is a mock implementation of the Store interface
type MockStore struct {
	mock.Mock
}

func (m *MockStore) Get(ctx context.Context, key string) ([]byte, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockStore) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockStore) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

type TestStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestGetOrSet(t *testing.T) {
	ctx := context.Background()
	key := "test-key"
	ttl := time.Minute
	expectedVal := TestStruct{Name: "Test", Age: 30}
	expectedBytes, _ := json.Marshal(expectedVal)

	t.Run("Cache Hit", func(t *testing.T) {
		mockStore := new(MockStore)
		mockStore.On("Get", ctx, key).Return(expectedBytes, nil)

		cm := &cacheManager{store: mockStore}

		val, err := GetOrSet(ctx, cm, key, ttl, func() (TestStruct, error) {
			return TestStruct{}, errors.New("should not be called")
		})

		assert.NoError(t, err)
		assert.Equal(t, expectedVal, val)
		mockStore.AssertExpectations(t)
	})

	t.Run("Cache Miss Success", func(t *testing.T) {
		mockStore := new(MockStore)
		mockStore.On("Get", ctx, key).Return(nil, ErrCacheMiss)
		mockStore.On("Set", ctx, key, expectedBytes, ttl).Return(nil)

		cm := &cacheManager{store: mockStore}

		val, err := GetOrSet(ctx, cm, key, ttl, func() (TestStruct, error) {
			return expectedVal, nil
		})

		assert.NoError(t, err)
		assert.Equal(t, expectedVal, val)
		mockStore.AssertExpectations(t)
	})

	t.Run("Cache Miss Fetcher Error", func(t *testing.T) {
		mockStore := new(MockStore)
		mockStore.On("Get", ctx, key).Return(nil, ErrCacheMiss)

		cm := &cacheManager{store: mockStore}
		fetchErr := errors.New("fetch error")

		val, err := GetOrSet(ctx, cm, key, ttl, func() (TestStruct, error) {
			return TestStruct{}, fetchErr
		})

		assert.Error(t, err)
		assert.Equal(t, fetchErr, err)
		assert.Equal(t, TestStruct{}, val)
		mockStore.AssertExpectations(t)
	})

	t.Run("Cache Error", func(t *testing.T) {
		mockStore := new(MockStore)
		cacheErr := errors.New("redis error")
		mockStore.On("Get", ctx, key).Return(nil, cacheErr)

		cm := &cacheManager{store: mockStore}

		val, err := GetOrSet(ctx, cm, key, ttl, func() (TestStruct, error) {
			return expectedVal, nil
		})

		assert.Error(t, err)
		assert.Equal(t, cacheErr, err)
		assert.Equal(t, TestStruct{}, val)
		mockStore.AssertExpectations(t)
	})

	t.Run("Cache Set Error", func(t *testing.T) {
		mockStore := new(MockStore)
		mockStore.On("Get", ctx, key).Return(nil, ErrCacheMiss)
		setErr := errors.New("set error")
		mockStore.On("Set", ctx, key, expectedBytes, ttl).Return(setErr)

		cm := &cacheManager{store: mockStore}

		val, err := GetOrSet(ctx, cm, key, ttl, func() (TestStruct, error) {
			return expectedVal, nil
		})

		assert.Error(t, err)
		assert.Equal(t, setErr, err)
		assert.Equal(t, TestStruct{}, val)
		mockStore.AssertExpectations(t)
	})

	t.Run("Cache Corrupt", func(t *testing.T) {
		mockStore := new(MockStore)
		mockStore.On("Get", ctx, key).Return([]byte("{invalid-json"), nil)
		mockStore.On("Set", ctx, key, expectedBytes, ttl).Return(nil)

		cm := &cacheManager{store: mockStore}

		val, err := GetOrSet(ctx, cm, key, ttl, func() (TestStruct, error) {
			return expectedVal, nil
		})

		assert.NoError(t, err)
		assert.Equal(t, expectedVal, val)
		mockStore.AssertExpectations(t)
	})
}

func TestInvalidate(t *testing.T) {
	ctx := context.Background()
	keys := []string{"k1", "k2"}

	mockStore := new(MockStore)
	mockStore.On("Delete", ctx, "k1").Return(nil)
	mockStore.On("Delete", ctx, "k2").Return(nil)

	cm := &cacheManager{store: mockStore}
	err := cm.Invalidate(ctx, keys...)

	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
}
