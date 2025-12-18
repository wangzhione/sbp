// Package localcache provides a wrapper around the Ristretto cache for simple cache operations.
package localcache

import (
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

// Cache is a wrapper around Ristretto cache
// It provides simple methods for cache operations

type Cache[K ristretto.Key, V any] struct {
	R *ristretto.Cache[K, V]
}

// NewCache creates a new Ristretto cache instance , num 预估保存多少活跃对象
func NewCache[K ristretto.Key, V any](num int) (*Cache[K, V], error) {
	c, err := ristretto.NewCache(&ristretto.Config[K, V]{
		NumCounters: int64(num) * 10, // 10x num items to optimize eviction
		MaxCost:     int64(num),      // Total cost of cache
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}

	return &Cache[K, V]{R: c}, nil
}

// SetTTL adds a key-value pair to the cache with a given cost and TTL | 0*time.Second is c.R.Set(key, value, 1)
func (c *Cache[K, V]) SetTTL(key K, value V, ttl time.Duration) bool {
	return c.R.SetWithTTL(key, value, 1, ttl)
}

// Get retrieves a value from the cache
func (c *Cache[K, V]) Get(key K) (V, bool) {
	return c.R.Get(key)
}

// Del removes a key from the cache
func (c *Cache[K, V]) Del(key K) {
	c.R.Del(key)
}
