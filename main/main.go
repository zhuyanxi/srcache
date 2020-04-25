package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/zhuyanxi/srcache"
	"github.com/zhuyanxi/srcache/consistenthash"
)

func main() {
	// t1 := time.Now()

	// count := 10
	// capacity := 8
	// src := srcache.NewSRCache(uint(capacity))
	// var wgT1 sync.WaitGroup
	// for i := 0; i < count; i++ {
	// 	wgT1.Add(1)
	// 	go func(j int) {
	// 		defer wgT1.Done()
	// 		src.Add(strconv.Itoa(j), []byte(strconv.Itoa(j*j)))
	// 	}(i)
	// }
	// wgT1.Wait()

	// src.Get("3")
	// src.Get("4")
	// src.Get("5")

	// for i := 0; i < count; i += 2 {
	// 	wgT1.Add(1)
	// 	go func(j int) {
	// 		defer wgT1.Done()
	// 		src.Add(strconv.Itoa(j), []byte(strconv.Itoa(j*j*j)))
	// 	}(i)
	// }
	// wgT1.Wait()
	// fmt.Println(src.Len())

	// elapsed := time.Since(t1)
	// fmt.Println("Time: ", elapsed)
	// fmt.Println("main finished")

	addr := "localhost:9099"
	srServer := srcache.NewServer(func(key string) ([]byte, error) {
		if v, ok := data[key]; ok {
			logrus.Infof("query key '%s' from db.\n", key)
			b, _ := json.Marshal(v)
			return b, nil
		}
		return nil, fmt.Errorf("key '%s' not exist", key)
	}, nil, nil)

	peers := []consistenthash.Node{
		{Addr: "10.192.168.10"},
		{Addr: "10.192.168.11"},
		{Addr: "10.192.168.12"},
	}

	srServer.SetPeers(peers...)

	logrus.Infoln("server is running at: ", addr)
	logrus.Fatalln(http.ListenAndServe(addr, srServer))
}
