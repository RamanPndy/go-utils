package goutils

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrCacheMiss            = errors.New("cache: miss")
	ErrUnsupportedProvider  = errors.New("cache: unsupported provider")
	ErrCacheClientNil       = errors.New("cache: client is nil")
	ErrNoCacheBackend       = errors.New("cache: no backend configured")
	ErrNoCachePrimary       = errors.New("cache: primary backend is nil")
	ErrNoCacheFallbackChain = errors.New("cache: no fallback backends configured")
)

type CacheProvider string

const (
	CacheProviderMemory   CacheProvider = "memory"
	CacheProviderRedis    CacheProvider = "redis"
	CacheProviderMemcache CacheProvider = "memcache"
	CacheProviderMulti    CacheProvider = "multi"
)

type Cache interface {
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	Close() error
}

// RedisClient is a minimal client contract to plug in a Redis backend.
type RedisClient interface {
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	Close() error
}

// MemcacheClient is a minimal client contract to plug in a Memcache backend.
type MemcacheClient interface {
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	Close() error
}

type CacheConfig struct {
	Provider CacheProvider
	Redis    RedisCacheConfig
	Memcache MemcacheCacheConfig
	Multi    MultiCacheConfig
}

type RedisCacheConfig struct {
	Client RedisClient
}

type MemcacheCacheConfig struct {
	Client MemcacheClient
}

type MultiCacheConfig struct {
	Primary   Cache
	Fallbacks []Cache
}

func NewCache(cfg CacheConfig) (Cache, error) {
	switch cfg.Provider {
	case "", CacheProviderMemory:
		return NewMemoryCache(), nil
	case CacheProviderRedis:
		return NewRedisCache(cfg.Redis.Client)
	case CacheProviderMemcache:
		return NewMemcacheCache(cfg.Memcache.Client)
	case CacheProviderMulti:
		return NewMultiCache(cfg.Multi.Primary, cfg.Multi.Fallbacks...)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedProvider, cfg.Provider)
	}
}

func NewMemoryCache() Cache {
	return &memoryCache{items: map[string]memoryItem{}}
}

func NewRedisCache(client RedisClient) (Cache, error) {
	if client == nil {
		return nil, ErrCacheClientNil
	}
	return &redisCache{client: client}, nil
}

func NewMemcacheCache(client MemcacheClient) (Cache, error) {
	if client == nil {
		return nil, ErrCacheClientNil
	}
	return &memcacheCache{client: client}, nil
}

func NewMultiCache(primary Cache, fallbacks ...Cache) (Cache, error) {
	if primary == nil {
		return nil, ErrNoCachePrimary
	}
	if len(fallbacks) == 0 {
		return nil, ErrNoCacheFallbackChain
	}

	backends := make([]Cache, 0, len(fallbacks)+1)
	backends = append(backends, primary)
	for _, fallback := range fallbacks {
		if fallback != nil {
			backends = append(backends, fallback)
		}
	}

	if len(backends) < 2 {
		return nil, ErrNoCacheFallbackChain
	}

	return &multiCache{backends: backends}, nil
}

type memoryItem struct {
	value     []byte
	expiresAt time.Time
	hasTTL    bool
}

type memoryCache struct {
	mu    sync.RWMutex
	items map[string]memoryItem
	dead  bool
}

func (m *memoryCache) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.dead {
		return ErrNoCacheBackend
	}

	item := memoryItem{value: append([]byte(nil), value...)}
	if ttl > 0 {
		item.hasTTL = true
		item.expiresAt = time.Now().Add(ttl)
	}
	m.items[key] = item
	return nil
}

func (m *memoryCache) Get(_ context.Context, key string) ([]byte, error) {
	m.mu.RLock()
	item, ok := m.items[key]
	dead := m.dead
	m.mu.RUnlock()

	if dead {
		return nil, ErrNoCacheBackend
	}
	if !ok {
		return nil, ErrCacheMiss
	}
	if item.hasTTL && time.Now().After(item.expiresAt) {
		m.mu.Lock()
		delete(m.items, key)
		m.mu.Unlock()
		return nil, ErrCacheMiss
	}

	return append([]byte(nil), item.value...), nil
}

func (m *memoryCache) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.dead {
		return ErrNoCacheBackend
	}
	delete(m.items, key)
	return nil
}

func (m *memoryCache) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.dead = true
	m.items = map[string]memoryItem{}
	return nil
}

type redisCache struct {
	client RedisClient
}

func (r *redisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl)
}

func (r *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	return r.client.Get(ctx, key)
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	return r.client.Delete(ctx, key)
}

func (r *redisCache) Close() error {
	return r.client.Close()
}

type memcacheCache struct {
	client MemcacheClient
}

func (m *memcacheCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return m.client.Set(ctx, key, value, ttl)
}

func (m *memcacheCache) Get(ctx context.Context, key string) ([]byte, error) {
	return m.client.Get(ctx, key)
}

func (m *memcacheCache) Delete(ctx context.Context, key string) error {
	return m.client.Delete(ctx, key)
}

func (m *memcacheCache) Close() error {
	return m.client.Close()
}

type multiCache struct {
	backends []Cache
}

func (m *multiCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if len(m.backends) == 0 {
		return ErrNoCacheBackend
	}

	for _, backend := range m.backends {
		if err := backend.Set(ctx, key, value, ttl); err != nil {
			return err
		}
	}
	return nil
}

func (m *multiCache) Get(ctx context.Context, key string) ([]byte, error) {
	if len(m.backends) == 0 {
		return nil, ErrNoCacheBackend
	}

	for idx, backend := range m.backends {
		value, err := backend.Get(ctx, key)
		if err == nil {
			// Warm up higher-priority caches if value is found in fallback.
			if idx > 0 {
				for i := 0; i < idx; i++ {
					_ = m.backends[i].Set(ctx, key, value, 0)
				}
			}
			return value, nil
		}
		if !errors.Is(err, ErrCacheMiss) {
			return nil, err
		}
	}

	return nil, ErrCacheMiss
}

func (m *multiCache) Delete(ctx context.Context, key string) error {
	if len(m.backends) == 0 {
		return ErrNoCacheBackend
	}

	for _, backend := range m.backends {
		if err := backend.Delete(ctx, key); err != nil {
			return err
		}
	}
	return nil
}

func (m *multiCache) Close() error {
	if len(m.backends) == 0 {
		return nil
	}
	for _, backend := range m.backends {
		if err := backend.Close(); err != nil {
			return err
		}
	}
	return nil
}
