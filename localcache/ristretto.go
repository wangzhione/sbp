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
		BufferItems: 256,
	})
	if err != nil {
		return nil, err
	}

	return &Cache[K, V]{R: c}, nil
}

// Set adds a key-value pair to the cache
func (c *Cache[K, V]) Set(key K, value V) bool {
	return c.R.Set(key, value, 1)
}

// Set adds a key-value pair to the cache with a given cost and TTL
func (c *Cache[K, V]) SetTTL(key K, value V, ttl time.Duration) bool {
	return c.R.SetWithTTL(key, value, 1, ttl)
}

// Get retrieves a value from the cache
func (c *Cache[K, V]) Get(key K) (V, bool) {
	return c.R.Get(key)
}

// Delete removes a key from the cache
func (c *Cache[K, V]) Delete(key K) {
	c.R.Del(key)
}

// Clear clears the entire cache
func (c *Cache[K, V]) Clear() {
	c.R.Clear()
}

// Wait waits for all cache operations to complete
func (c *Cache[K, V]) Wait() {
	c.R.Wait()
}

// Metrics returns cache statistics
func (c *Cache[K, V]) Metrics() *ristretto.Metrics {
	return c.R.Metrics
}
