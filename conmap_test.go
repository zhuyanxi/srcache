package main

import (
	"strconv"
	"sync"
	"testing"
)

func TestConMapAddAsync(t *testing.T) {
	cm := NewConMap()
	var wg1 sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg1.Add(1)
		go func(j int) {
			defer wg1.Done()
			cm.Add(strconv.Itoa(j), []byte(strconv.Itoa(j*j)))
			// cm.AddSync(strconv.Itoa(j), []byte(strconv.Itoa(j*j)))
		}(i)
	}
	wg1.Wait()

	wantCount := 10
	if len(cm.m) != wantCount {
		t.Errorf("test failed: want %d; current %d\n", wantCount, len(cm.m))
	} else {
		t.Logf("test passed")
	}
}

func BenchmarkConMapAddSync(b *testing.B) {
	cm := NewConMap()
	for i := 0; i < b.N; i++ {
		for i := 0; i < 100; i++ {
			cm.AddSync(strconv.Itoa(i), []byte(strconv.Itoa(i*i)))
		}
	}
}

func BenchmarkConMapAddAsync(b *testing.B) {
	cm := NewConMap()
	var wg1 sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for i := 0; i < 100; i++ {
			wg1.Add(1)
			go func(j int) {
				defer wg1.Done()
				cm.Add(strconv.Itoa(j), []byte(strconv.Itoa(j*j)))
			}(i)
		}
		wg1.Wait()
	}
}
