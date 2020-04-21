package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	// ipRecords, _ := net.LookupIP("baidu.com")
	// for _, ip := range ipRecords {
	// 	fmt.Println(ip)
	// }
	cm := NewConMap()
	// cm.Add("one", []byte("data one"))

	wg.Add(1)
	go func() {
		defer wg.Done()
		cm.Add("one", []byte("data one"))
	}()
	wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if v, ok := cm.Get("one"); ok {
			fmt.Println(string(v))
		} else {
			fmt.Println("not found")
		}
	}()
	wg.Wait()
	//time.Sleep(time.Second)
	fmt.Println("Done")
}
