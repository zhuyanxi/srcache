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
	sc := core.NewSRCache(4)
	var wgT1 sync.WaitGroup
	for i := 0; i < count; i++ {
		wgT1.Add(1)
		go func(j int) {
			defer wgT1.Done()
			sc.Add(strconv.Itoa(j), []byte(strconv.Itoa(j*j)))
		}(i)
	}
	wgT1.Wait()

	elapsed := time.Since(t1)
	fmt.Println("Time: ", elapsed)
	fmt.Println("main finished")
}
