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
	keyIter(func (k int) bool {
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
}
