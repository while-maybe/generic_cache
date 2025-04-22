package cache

import (
	"maps"
	"slices"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := New[int, int](3, time.Second)

	if c.data == nil {
		t.Error("Data content not initialised")
	}
	if c.chronologicalKeys == nil {
		t.Error("List of keys not initialised")
	}

	c.Upsert(1, 10)
	c.Upsert(2, 20)
	c.Upsert(3, 30)

	expectedKeys := []int{1, 2, 3}
	if !slices.Equal(c.chronologicalKeys, expectedKeys) {
		t.Errorf("List of keys should be %v, got %v", expectedKeys, c.chronologicalKeys)
	}

	// keys := make([]string, 0, len(myMap))
	// for k := range myMap {
	//     keys = append(keys, k)
	// }
	// A different way of doing the for loop above with the iter pattern for a change
	dataKeys := make([]int, 0, len(c.data))
	keyIter := maps.Keys(c.data)
	keyIter(func(k int) bool {
		dataKeys = append(dataKeys, k)
		return true
	})

	slices.Sort(dataKeys)
	if !slices.Equal(dataKeys, expectedKeys) {
		t.Errorf("Keys of the map should be %v, got %v", expectedKeys, dataKeys)
	}

	c.Upsert(2, 31)
	expectedKeys = []int{1, 3, 2}
	if !slices.Equal(c.chronologicalKeys, expectedKeys) {
		t.Errorf("After upserting: Keys of the map should be %v, got %v", expectedKeys, c.chronologicalKeys)
	}

	c.Delete(3)

	// see comment above
	dataKeys = make([]int, 0, len(c.data))
	keyIter = maps.Keys(c.data)
	keyIter(func(k int) bool {
		dataKeys = append(dataKeys, k)
		return true
	})

	slices.Sort(dataKeys)
	expectedKeys = []int{1, 2}
	if !slices.Equal(dataKeys, expectedKeys) {
		t.Errorf("After deleting: Keys of the map should be %v, got %v", expectedKeys, dataKeys)
	}

	value, found := c.Read(1)
	if !found {
		t.Error("Key 1 should be present")
	}

	if value != 10 {
		t.Errorf("Value of Key 1 should be 10, got %v", value)
	}

	// test TTL
	time.Sleep(time.Second)
	value, found = c.Read(1)
	if found {
		t.Error("Key 1 should have expired but still exists")
	}

	if value != 0 {
		t.Errorf("Value of Key 1 should be 0, got %v", value)
	}
}
