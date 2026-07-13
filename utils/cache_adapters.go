package goutils

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/redis/go-redis/v9"
)

// NewRedisCacheFromGoRedis wraps a go-redis client into this package Cache.
func NewRedisCacheFromGoRedis(client redis.UniversalClient) (Cache, error) {
	if client == nil {
		return nil, ErrCacheClientNil
	}
	return NewRedisCache(&goRedisClientAdapter{client: client})
}

// NewMemcacheCacheFromClient wraps a gomemcache client into this package Cache.
func NewMemcacheCacheFromClient(client *memcache.Client) (Cache, error) {
	if client == nil {
		return nil, ErrCacheClientNil
	}
	return NewMemcacheCache(&goMemcacheClientAdapter{client: client})
}

type goRedisClientAdapter struct {
	client redis.UniversalClient
}

func (r *goRedisClientAdapter) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *goRedisClientAdapter) Get(ctx context.Context, key string) ([]byte, error) {
	value, err := r.client.Get(ctx, key).Bytes()
	if err == nil {
		return value, nil
	}
	if errors.Is(err, redis.Nil) {
		return nil, ErrCacheMiss
	}
	return nil, err
}

func (r *goRedisClientAdapter) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *goRedisClientAdapter) Close() error {
	return r.client.Close()
}

type goMemcacheClientAdapter struct {
	client *memcache.Client
}

func (m *goMemcacheClientAdapter) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	return m.client.Set(&memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: ttlToMemcacheSeconds(ttl),
	})
}

func (m *goMemcacheClientAdapter) Get(_ context.Context, key string) ([]byte, error) {
	item, err := m.client.Get(key)
	if err == nil {
		return append([]byte(nil), item.Value...), nil
	}
	if errors.Is(err, memcache.ErrCacheMiss) {
		return nil, ErrCacheMiss
	}
	return nil, err
}

func (m *goMemcacheClientAdapter) Delete(_ context.Context, key string) error {
	err := m.client.Delete(key)
	if errors.Is(err, memcache.ErrCacheMiss) {
		return nil
	}
	return err
}

func (m *goMemcacheClientAdapter) Close() error {
	return nil
}

func ttlToMemcacheSeconds(ttl time.Duration) int32 {
	if ttl <= 0 {
		return 0
	}
	seconds := ttl.Seconds()
	if seconds > math.MaxInt32 {
		return math.MaxInt32
	}
	return int32(seconds)
}
