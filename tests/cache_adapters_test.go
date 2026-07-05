package goutils_test

import (
	"errors"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/redis/go-redis/v9"

	goutils "github.com/RamanPndy/go-utils/utils"
)

func TestNewRedisCacheFromGoRedis(t *testing.T) {
	if _, err := goutils.NewRedisCacheFromGoRedis(nil); !errors.Is(err, goutils.ErrCacheClientNil) {
		t.Fatalf("NewRedisCacheFromGoRedis(nil) error = %v, want ErrCacheClientNil", err)
	}

	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"})
	cache, err := goutils.NewRedisCacheFromGoRedis(client)
	if err != nil {
		t.Fatalf("NewRedisCacheFromGoRedis(client) error = %v", err)
	}
	if cache == nil {
		t.Fatalf("NewRedisCacheFromGoRedis(client) returned nil cache")
	}
}

func TestNewMemcacheCacheFromClient(t *testing.T) {
	if _, err := goutils.NewMemcacheCacheFromClient(nil); !errors.Is(err, goutils.ErrCacheClientNil) {
		t.Fatalf("NewMemcacheCacheFromClient(nil) error = %v, want ErrCacheClientNil", err)
	}

	client := memcache.New("127.0.0.1:11211")
	cache, err := goutils.NewMemcacheCacheFromClient(client)
	if err != nil {
		t.Fatalf("NewMemcacheCacheFromClient(client) error = %v", err)
	}
	if cache == nil {
		t.Fatalf("NewMemcacheCacheFromClient(client) returned nil cache")
	}
}
