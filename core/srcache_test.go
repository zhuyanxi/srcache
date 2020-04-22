package core

import (
	"strconv"
	"sync"
	"testing"
)

func TestAddSRCache(t *testing.T) {
	count := 10
	sc := NewSRCache(4)
	var wgT1 sync.WaitGroup
	for i := 0; i < count; i++ {
		wgT1.Add(1)
		go func(j int) {
			defer wgT1.Done()
			sc.Add(strconv.Itoa(j), []byte(strconv.Itoa(j*j)))
		}(i)
	}
	wgT1.Wait()
	if sc.info.lru.Len() != count {
		t.Errorf("test failed: want %d; current %d\n", count, 1)
	} else {
		t.Logf("test passed")
	}
}
