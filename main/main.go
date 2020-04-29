package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/zhuyanxi/srcache"
	"github.com/zhuyanxi/srcache/consistenthash"
	"gopkg.in/yaml.v2"
)

// Peer :
type Peer struct {
	Target string `yaml:"target"`
}

// ConfigYaml :
type ConfigYaml struct {
	CacheCapacity uint   `yaml:"cachecapacity"`
	Replicate     int    `yaml:"replicate"`
	PathPrefix    string `yaml:"pathprefix"`

	DataReadEntry string `yaml:"datareadentry"`

	Local string `yaml:"local"`
	Peers []Peer `yaml:"peers"`
}

var (
	//addr   = flag.String("listen-address", "localhost:3001", "The address to listen on for HTTP requests.")
	config = flag.String("config", "config.yaml", "Config file in yaml format.")
)

var peerConfig ConfigYaml

func init() {
	f, err := os.Open(*config)
	if err != nil {
		logrus.Fatalf("os.Open failed with '%s'\n", err)
	}
	defer f.Close()

	doc := yaml.NewDecoder(f)

	//var yamlFile PeerYaml
	err = doc.Decode(&peerConfig)
	if err != nil {
		logrus.Fatalf("doc.Decode failed with '%s'\n", err)
	}

	logrus.Infof("Decode config YAML peers:%#v\n", peerConfig)
}

func main() {
	flag.Parse()
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

	//addr := "localhost:3001"
	opts := &srcache.ServerOptions{
		LocalURL:      peerConfig.Local,
		CacheCapacity: peerConfig.CacheCapacity,
		Replicate:     peerConfig.Replicate,
		PathPrefix:    peerConfig.PathPrefix,
	}
	srServer := srcache.NewServer(func(key string) ([]byte, error) {
		if v, ok := data[key]; ok {
			logrus.Infof("query key '%s' from db.\n", key)
			b, _ := json.Marshal(v)
			return b, nil
		}
		return nil, fmt.Errorf("key '%s' not exist", key)
	}, nil, opts)

	// peers := []consistenthash.Node{
	// 	{Addr: "localhost:3001"},
	// 	{Addr: "localhost:3002"},
	// 	{Addr: "localhost:3003"},
	// }

	var peers []consistenthash.Node
	for _, v := range peerConfig.Peers {
		peers = append(peers, consistenthash.Node{
			Addr: v.Target,
		})
	}

	srServer.SetPeers(peers...)

	logrus.Infoln("server is running at: ", peerConfig.Local)
	logrus.Fatalln(http.ListenAndServe(peerConfig.Local, srServer))
}
