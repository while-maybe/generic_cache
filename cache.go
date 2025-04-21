package cache

import (
	"sync"
	"time"
)

// Cache is key-value pair storage.
type Cache[K comparable, V any] struct {
	ttl time.Duration

	mu   sync.Mutex
	data map[K]entryWithTimeout[V]
}

type entryWithTimeout[V any] struct {
	value   V
	expires time.Time // After this time, value is useless
}

// New creates a usable cache.
func New[K comparable, V any](ttl time.Duration) Cache[K, V] {
	return Cache[K, V]{
		ttl:  ttl,
		data: make(map[K]entryWithTimeout[V]),
	}
}

// Read returns the associated value for a key, and a boolean true if the key is absent.
func (c *Cache[K, V]) Read(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zeroV V

	v, ok := c.data[key]

	switch {
	case !ok:
		return zeroV, false
	case v.expires.Before(time.Now()):
		// this value is expired
		delete(c.data, key)
		return zeroV, false
	default:
		return v.value, ok
	}
}

// Upsert overrides the value for a given key.
func (c *Cache[K, V]) Upsert(key K, value V) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = entryWithTimeout[V]{
		value:   value,
		expires: time.Now().Add(c.ttl),
	}

	// Don't return an error for the moment, possibly in the future
	return nil
}

// Delete removes the entry for the given key.
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}
