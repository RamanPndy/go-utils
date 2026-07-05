# go-utils

Reusable utility helpers for Go applications, including generic collection helpers, reflection helpers, JSON/YAML conversion, HTTP helpers, Kubernetes wrappers, and lightweight data structures.

[![Coverage Status](https://coveralls.io/repos/github/RamanPndy/go-utils/badge.svg)](https://coveralls.io/github/RamanPndy/go-utils)

## Requirements

- Go 1.24+

## Installation

```bash
go get github.com/RamanPndy/go-utils
```

Then import the package:

```go
import goutils "github.com/RamanPndy/go-utils/utils"
```

## What Is Included

### Collection and Functional Helpers

- `Map`, `Filter`, `Any`, `All`
- `Contains`, `Dedupe`, `Zip`, `CombineSlicesToMap`
- `Equals`, `EqualsSlice`, `EqualsMap`

### Set/Dict/Iterator Data Structures

- `NewSet`, `NewFrozenSet`, `NewOrderedSet`
- `NewImmutableDict`, `NewOrderedDict`
- `NewIterator`, `NewOrderedList`

### Struct/Reflection Utilities

- `HasAttr`, `SetAttr`, `Vars`
- `IsInstance`, `IsSubclass`, `IsNilInterface`
- `GetStructFieldNames`, `GetStructFieldValue`
- `MergeUniqueFields`, `SkipMergeUniqueFields`

### JSON/YAML and Encoding Helpers

- `DeepCopyJSON`, `MustJSON`
- `JsonEncode*` and `JsonDecode*` helpers
- `JsonToYaml*` and `YamlToJson*`
- `JSONEqual`, `JSONDeepEqual`, `JSONContains`
- `Base64Encode`, `Base64Decode`
- `EscapeJSONPointer`, `UnescapeJSONPointer`

### HTTP and API Helpers

- `NewAPIClient`, `NewClientWithToken`, `NewAPIRequest`
- URL/query helpers: `EncodeQueryParams`, `DecodeQueryParams`
- Validation helpers: `IsValidURL`, `IsValidStatusCode`, `IsValidContentType`, and related validators
- `GetContentTypeFromHeaders`, `IsValidAPIResult`

### Infrastructure Helpers

- Database and DSN: `BuildDSN`, `NewDBConn`, `OpenDB`, `NewDsnLoader`, `NewPostgresStore`
- Cache: `NewCache`, `NewMemoryCache`, `NewMultiCache`, `NewRedisCacheFromGoRedis`, `NewMemcacheCacheFromClient`
- File watching: `NewFileWatcher`, `NewWatcher`, `NotifyOnFileChange`
- Kubernetes command wrappers: `KubectlGet`, `KubectlApply`, `KubectlRaw`, `InClusterConfig`
- Misc: `ReadFile`, `Compute`, `UnixTimeToTimestamp`, `GetFirstNonEmpty`, `IndentBlock`

## Quick Examples

### Functional map/filter

```go
nums := []int{1, 2, 3, 4}
squares := goutils.Map(nums, func(n int) int { return n * n })
evens := goutils.Filter(nums, func(n int) bool { return n%2 == 0 })
```

### Struct field access via reflection

```go
type User struct {
		Name string
		Age  int
}

u := User{Name: "Alice", Age: 30}

hasName := goutils.HasAttr(u, "Name")
_ = hasName

name, _ := goutils.GetStructFieldValue(u, "Name")
_ = name
```

### JSON encode/decode

```go
type Payload struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
}

data := Payload{ID: 1, Name: "demo"}
jsonStr, _ := goutils.JsonEncodeToString(data)
var decoded Payload
_ = goutils.JsonDecodeFromString(jsonStr, &decoded)
_ = decoded
```

### Cache (memory/redis/memcache/multi)

```go
ctx := context.Background()

cache := goutils.NewMemoryCache()
_ = cache.Set(ctx, "user:1", []byte("Alice"), 5*time.Minute)
value, _ := cache.Get(ctx, "user:1")
_ = value
```

```go
redisClient := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
redisCache, _ := goutils.NewRedisCacheFromGoRedis(redisClient)

_ = redisCache.Set(context.Background(), "session:token", []byte("abc123"), 10*time.Minute)
```

```go
memc := memcache.New("127.0.0.1:11211")
memcacheCache, _ := goutils.NewMemcacheCacheFromClient(memc)

_ = memcacheCache.Set(context.Background(), "feature:flag", []byte("enabled"), time.Minute)
```

```go
primary := goutils.NewMemoryCache()
fallback := goutils.NewMemoryCache()

multi, _ := goutils.NewMultiCache(primary, fallback)
_ = multi.Set(context.Background(), "key", []byte("value"), 0)
```

## Scripts

The `scripts/` directory contains Kubernetes helper scripts:

- `scripts/k8s-cpu-memory.sh <namespace>`
	- Prints deployment/container CPU and memory requests/limits.
- `scripts/k8s-get-image.sh`
	- Interactive script to select namespace, deployment, and container, then display the current image.
- `scripts/k8s-set-image.sh`
	- Interactive script to update a selected container image in a deployment.

These scripts require `kubectl` configured for your cluster. `k8s-cpu-memory.sh` also requires `jq`.

## Running Tests

```bash
go test ./...
```

## Running Demo Entrypoint

The repository includes a sample `main` that exercises many helpers:

```bash
go run .
```

## License

MIT (see `LICENSE`).
