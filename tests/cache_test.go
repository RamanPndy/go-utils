package goutils_test

import (
	"context"
	"errors"
	"testing"
	"time"

	goutils "github.com/RamanPndy/go-utils/utils"
)

type fakeExternalCacheClient struct {
	items map[string][]byte
}

func newFakeExternalCacheClient() *fakeExternalCacheClient {
	return &fakeExternalCacheClient{items: map[string][]byte{}}
}

func (f *fakeExternalCacheClient) Set(_ context.Context, key string, value []byte, _ time.Duration) error {
	f.items[key] = append([]byte(nil), value...)
	return nil
}

func (f *fakeExternalCacheClient) Get(_ context.Context, key string) ([]byte, error) {
	v, ok := f.items[key]
	if !ok {
		return nil, goutils.ErrCacheMiss
	}
	return append([]byte(nil), v...), nil
}

func (f *fakeExternalCacheClient) Delete(_ context.Context, key string) error {
	delete(f.items, key)
	return nil
}

func (f *fakeExternalCacheClient) Close() error {
	return nil
}

func TestMemoryCacheSetGetDelete(t *testing.T) {
	cache := goutils.NewMemoryCache()
	ctx := context.Background()

	if err := cache.Set(ctx, "k1", []byte("v1"), 0); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	got, err := cache.Get(ctx, "k1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if string(got) != "v1" {
		t.Fatalf("Get() value = %q, want %q", got, "v1")
	}

	if err := cache.Delete(ctx, "k1"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = cache.Get(ctx, "k1")
	if !errors.Is(err, goutils.ErrCacheMiss) {
		t.Fatalf("Get() after Delete() error = %v, want ErrCacheMiss", err)
	}
}

func TestMemoryCacheTTL(t *testing.T) {
	cache := goutils.NewMemoryCache()
	ctx := context.Background()

	if err := cache.Set(ctx, "ttl", []byte("value"), 30*time.Millisecond); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	time.Sleep(60 * time.Millisecond)

	_, err := cache.Get(ctx, "ttl")
	if !errors.Is(err, goutils.ErrCacheMiss) {
		t.Fatalf("Get() after ttl error = %v, want ErrCacheMiss", err)
	}
}

func TestMultiCacheFallbackAndWarmup(t *testing.T) {
	primary := goutils.NewMemoryCache()
	fallback := goutils.NewMemoryCache()
	ctx := context.Background()

	if err := fallback.Set(ctx, "name", []byte("raman"), 0); err != nil {
		t.Fatalf("fallback Set() error = %v", err)
	}

	multi, err := goutils.NewMultiCache(primary, fallback)
	if err != nil {
		t.Fatalf("NewMultiCache() error = %v", err)
	}

	value, err := multi.Get(ctx, "name")
	if err != nil {
		t.Fatalf("multi Get() error = %v", err)
	}
	if string(value) != "raman" {
		t.Fatalf("multi Get() value = %q, want %q", value, "raman")
	}

	// Value should be promoted into primary cache.
	primVal, err := primary.Get(ctx, "name")
	if err != nil {
		t.Fatalf("primary Get() after warmup error = %v", err)
	}
	if string(primVal) != "raman" {
		t.Fatalf("primary Get() value = %q, want %q", primVal, "raman")
	}
}

func TestNewCacheProviderSelection(t *testing.T) {
	if _, err := goutils.NewCache(goutils.CacheConfig{Provider: "invalid"}); err == nil {
		t.Fatalf("NewCache() expected unsupported provider error")
	}

	redisClient := newFakeExternalCacheClient()
	redisCache, err := goutils.NewCache(goutils.CacheConfig{
		Provider: goutils.CacheProviderRedis,
		Redis: goutils.RedisCacheConfig{
			Client: redisClient,
		},
	})
	if err != nil {
		t.Fatalf("NewCache(redis) error = %v", err)
	}

	ctx := context.Background()
	if err := redisCache.Set(ctx, "rk", []byte("rv"), 0); err != nil {
		t.Fatalf("redis Set() error = %v", err)
	}
	rv, err := redisCache.Get(ctx, "rk")
	if err != nil {
		t.Fatalf("redis Get() error = %v", err)
	}
	if string(rv) != "rv" {
		t.Fatalf("redis Get() value = %q, want %q", rv, "rv")
	}

	memcacheClient := newFakeExternalCacheClient()
	memcacheCache, err := goutils.NewCache(goutils.CacheConfig{
		Provider: goutils.CacheProviderMemcache,
		Memcache: goutils.MemcacheCacheConfig{
			Client: memcacheClient,
		},
	})
	if err != nil {
		t.Fatalf("NewCache(memcache) error = %v", err)
	}

	if err := memcacheCache.Set(ctx, "mk", []byte("mv"), 0); err != nil {
		t.Fatalf("memcache Set() error = %v", err)
	}
	mv, err := memcacheCache.Get(ctx, "mk")
	if err != nil {
		t.Fatalf("memcache Get() error = %v", err)
	}
	if string(mv) != "mv" {
		t.Fatalf("memcache Get() value = %q, want %q", mv, "mv")
	}
}
