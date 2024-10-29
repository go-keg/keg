package cache

import (
	"context"
	"errors"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	ErrNotExist = errors.New("cache does not exist")
)

type Cache interface {
	Has(ctx context.Context, key string) (bool, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Forget(ctx context.Context, key string) error
	Remember(ctx context.Context, key string, ttl time.Duration, fn func() ([]byte, error)) ([]byte, error)
}

var goCache *cache.Cache
var DefaultExpiration = time.Hour

func init() {
	goCache = cache.New(DefaultExpiration, 3*DefaultExpiration)
}

func LocalRemember[T any](key string, ttl time.Duration, fn func() (T, error)) (T, error) {
	var (
		v   T
		err error
		ok  bool
	)
	if val, has := goCache.Get(key); has {
		if v, ok = val.(T); ok {
			return v, nil
		}
		return v, nil
	}
	v, err = fn()
	if err != nil {
		return v, err
	}
	goCache.Set(key, v, ttl)
	return v, nil
}

func LocalSet(key string, val any, ttl time.Duration) {
	goCache.Set(key, val, ttl)
}

func LocalGet(key string) (any, bool) {
	return goCache.Get(key)
}

func LocalClear(key string) {
	goCache.Delete(key)
}
