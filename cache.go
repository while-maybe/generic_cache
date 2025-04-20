package cache

// Cache is key-value pair storage.
type Cache[K comparable, V any] struct {
	data map[K]V
}

// New creates a usable cache.
func New[K comparable, V any]() Cache[K, V] {
	return Cache[K, V]{
		data: make(map[K]V),
	}
}

// Read returns the associated value for a key, and a boolean true if the key is absent.
func (c *Cache[K, V]) Read(key K) (V, bool) {
	v, found := c.data[key]
	return v, found
}

// Upsert overrides the value for a given key.
func (c *Cache[K, V]) Upsert(key K, value V) error {
	c.data[key] = value

	// Don't return an error for the moment, possibly in the future
	return nil
}

// Delete removes the entry for the given key.
func (c *Cache[K, V]) Delete(key K) {
	delete(c.data, key)
}
