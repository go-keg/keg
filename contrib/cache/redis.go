package cache

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-keg/keg/contrib/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	rdb     *redis.Client
	options []func(options *redis.Options)
	prefix  string
}

type RedisOptionFunc func(*Redis)

func SetPassword(password string) RedisOptionFunc {
	return func(r *Redis) {
		r.SetOption(func(options *redis.Options) {
			options.Password = password
		})
	}
}

func SetDB(db int) RedisOptionFunc {
	return func(r *Redis) {
		r.SetOption(func(options *redis.Options) {
			options.DB = db
		})
	}
}

func SetPrefix(prefix string) RedisOptionFunc {
	return func(r *Redis) {
		r.prefix = prefix
	}
}

func (r *Redis) SetOption(opt func(*redis.Options)) {
	r.options = append(r.options, opt)
}

func NewRedisFromConfig(config config.Redis) Cache {
	var opts []RedisOptionFunc
	if config.Password != "" {
		opts = append(opts, SetPassword(config.Password))
	}
	if config.DB != "" {
		db, err := strconv.Atoi(config.DB)
		if err != nil {
			panic(fmt.Errorf("redis db:[%s] invalid", config.DB))
		}
		opts = append(opts, SetDB(db))
	}
	if config.Prefix != "" {
		opts = append(opts, SetPrefix(config.Prefix))
	}
	return NewRedis(config.Addr, opts...)
}

func NewRedis(addr string, opts ...RedisOptionFunc) Cache {
	cache := &Redis{
		prefix: "",
	}

	for _, fn := range opts {
		fn(cache)
	}

	options := &redis.Options{
		Addr: addr,
	}
	for _, fn := range cache.options {
		fn(options)
	}

	cache.rdb = redis.NewClient(options)

	return cache
}

func (r *Redis) Has(ctx context.Context, key string) (bool, error) {
	result, err := r.rdb.Exists(ctx, r.getKey(key)).Result()
	if result == 0 {
		return false, err
	}
	return true, err
}

func (r *Redis) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return r.rdb.Set(ctx, r.getKey(key), value, expiration).Err()
}

func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	result, err := r.rdb.Get(ctx, r.getKey(key)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, ErrNotExist
	}
	return result, err
}

func (r *Redis) Forget(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, r.getKey(key)).Err()
}

func (r *Redis) Remember(ctx context.Context, key string, ttl time.Duration, fn func() ([]byte, error)) ([]byte, error) {
	result, err := r.Get(ctx, key)
	if err == nil {
		return result, err
	}

	result, err = fn()
	if err != nil {
		return nil, err
	}
	if err = r.Set(ctx, key, result, ttl); err != nil {
		return nil, err
	}
	return result, err
}

func (r *Redis) getKey(key string) string {
	if r.prefix != "" {
		return r.prefix + ":" + key
	}
	return key
}
