package cache_test

import (
	"fmt"
	cache "genericcache"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
TestCache is an integration test.

- Create a cache and check that it is empty?
- Upsert a new key <5, fünf> in the cache
- Read the value for this entry
- Upsert for the same entry with the new value
- Read the new value
- Upsert another key <3, drei> and check that is doesn't override
- Delete 5 and check that 3 still exists
*/

func TestCache(t *testing.T) {
	c := cache.New[int, string](5, time.Minute)

	c.Upsert(5, "fünf")
	v, found := c.Read(5)
	assert.True(t, found)
	assert.Equal(t, "fünf", v)

	c.Upsert(5, "pum")
	v, found = c.Read(5)
	assert.True(t, found)
	assert.Equal(t, "pum", v)

	c.Upsert(3, "drei")
	v, found = c.Read(3)
	assert.True(t, found)
	assert.Equal(t, "drei", v)

	c.Upsert(5, "pum")
	v, found = c.Read(5)
	assert.True(t, found)
	assert.Equal(t, "pum", v)

	c.Delete(5)
	v, found = c.Read(5)
	assert.False(t, found)
	assert.Equal(t, "", v)

	v, found = c.Read(3)
	assert.True(t, found)
	assert.Equal(t, "drei", v)
}

// TestCache_Parallel_goroutines simulates a number of parallel tasks each operating on the cache.
// It passes if we only use "go test .", but we see the error as soon as we use "go test -race ."
func TestCache_Parallel_goroutines(t *testing.T) {
	c := cache.New[int, string](3, time.Second)

	const parallelTasks = 1000
	wg := sync.WaitGroup{}
	wg.Add(parallelTasks)

	for i := range parallelTasks {
		go func(j int) {
			defer wg.Done()
			c.Upsert(4, fmt.Sprint(j))
		}(i)
	}
}

// TestCache_Parallel runs two routines that have concurrent access to write to the cache.
func TestCache_Parallel(t *testing.T) {
	c := cache.New[int, string](5, time.Second)

	t.Run("write six", func(t *testing.T) {
		t.Parallel()
		c.Upsert(6, "six")
	})
	t.Run("write kuus", func(t *testing.T) {
		t.Parallel()
		c.Upsert(6, "kuus")
	})
}

// TestCache_TTL Upserts a new value to the cache, checks if value exists immediately after and checks a second time after the value expires
func TestCache_TTL(t *testing.T) {
	t.Parallel()

	c := cache.New[string, string](5, time.Millisecond*100)
	c.Upsert("Norwegian", "Blue")

	// Check if item exists
	got, found := c.Read("Norwegian")
	assert.True(t, found)
	assert.Equal(t, "Blue", got)

	time.Sleep(time.Millisecond * 200)
	got, found = c.Read("Norwegian")

	assert.False(t, found)
	assert.Equal(t, "", got)
}

// TestCache_MaxSize tests the maximum capacity feature of a cache. It checks that updated items are properly requeued as "new" items and that we make room by removing the oldest item for the new one.
func TestCache_MaxSize(t *testing.T) {
	t.Parallel()

	c := cache.New[int, int](3, time.Minute)

	c.Upsert(1, 1)
	c.Upsert(2, 2)
	c.Upsert(3, 3)

	got, found := c.Read(1)
	assert.True(t, found)
	assert.Equal(t, got, 1)

	// Update the oldest element making it the "newest"
	c.Upsert(1, 10)

	// Adding another element will push the current "oldest"(this is 2 at the moment) out
	c.Upsert(4, 4)

	// Attmempt to retrieve an element that should be in the cache at this moment
	got, found = c.Read(2)
	assert.False(t, found)
	assert.Equal(t, got, 0)
}
