package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhuyanxi/srcache"
	"github.com/zhuyanxi/srcache/consistenthash"
	"gopkg.in/yaml.v2"
)

// Peer :
type Peer struct {
	URL  string `yaml:"url"`
	Port string `yaml:"port"`
}

// ConfigYaml :
type ConfigYaml struct {
	CacheCapacity uint   `yaml:"cachecapacity"`
	Replicate     int    `yaml:"replicate"`
	PathPrefix    string `yaml:"pathprefix"`

	DataReadEntry string `yaml:"datareadentry"`

	Local string `yaml:"local"`
	Peers []Peer `yaml:"peers"`

	UseGRPC bool `yaml:"usegrpc"`
}

var (
	addr   = flag.String("listen-address", "", "The address to listen on for HTTP requests.")
	config = flag.String("config", "config.yaml", "Config file in yaml format.")
)

var peerConfig ConfigYaml

func initConfig() {
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
	initConfig()
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

	listenEntry := peerConfig.Local
	if len(*addr) != 0 {
		listenEntry = *addr
	}

	opts := &srcache.ServerOptions{
		LocalURL:      listenEntry,
		CacheCapacity: peerConfig.CacheCapacity,
		Replicate:     peerConfig.Replicate,
		PathPrefix:    peerConfig.PathPrefix,
	}
	srServer := srcache.NewServer(func(key string) ([]byte, error) {
		time.Sleep(200 * time.Millisecond)
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
			Addr: v.URL + ":" + v.Port,
		})
	}

	srServer.SetPeers(peers...)

	if peerConfig.UseGRPC {
		logrus.Infoln("grpc server is running at: ", listenEntry)
		srServer.ServeGRPC(strings.Split(peerConfig.Local, ":")[1])
	} else {
		logrus.Infoln("http server is running at: ", listenEntry)
		logrus.Fatalln(http.ListenAndServe(listenEntry, srServer))
	}
}
