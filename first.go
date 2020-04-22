package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/zhuyanxi/sr-cache/core"
)

var wg sync.WaitGroup

func main() {
	t1 := time.Now()

	count := 10
	capacity := 8
	src := core.NewSRCache(uint(capacity))
	var wgT1 sync.WaitGroup
	for i := 0; i < count; i++ {
		wgT1.Add(1)
		go func(j int) {
			defer wgT1.Done()
			src.Add(strconv.Itoa(j), []byte(strconv.Itoa(j*j)))
		}(i)
	}
	wgT1.Wait()

	src.Get("3")
	src.Get("4")
	src.Get("5")

	for i := 0; i < count; i += 2 {
		wgT1.Add(1)
		go func(j int) {
			defer wgT1.Done()
			src.Add(strconv.Itoa(j), []byte(strconv.Itoa(j*j*j)))
		}(i)
	}
	wgT1.Wait()

	fmt.Println(src.Len())
	elapsed := time.Since(t1)
	fmt.Println("Time: ", elapsed)
	fmt.Println("main finished")
}
