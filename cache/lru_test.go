package cache

import "testing"

var getTests = []struct {
	name       string
	keyToAdd   string
	keyToGet   string
	expectedOk bool
}{
	{"string_hit", "myKey", "myKey", true},
	{"string_miss", "myKey", "nonsense", false},
	{"string_hit_two", "two", "two", true},
	{"string_miss_two", "two", "noway", false},
	{"string_hit_three", "three", "three", true},
}

func TestGetAndRemoveOldest(t *testing.T) {
	// init the lrucache with 1 capacity
	lru := New(1)

	for _, test := range getTests {
		// if len > capacity, the cache will remove the oldest automatically
		lru.Add(test.keyToAdd, []byte("1234"))

		val, ok := lru.Get(test.keyToGet)
		if ok != test.expectedOk {
			t.Fatalf("%s: cache hit = %v; want %v\n", test.name, ok, !ok)
		} else if ok && string(val) != "1234" {
			t.Fatalf("%s expected get to return 1234 but got %v\n", test.name, val)
		}
	}
}
