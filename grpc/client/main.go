package main

import (
	"fmt"
	"net/rpc"

	"github.com/zhuyanxi/srcache/grpc"
)

const HelloServiceName = "srcache/grpc.HelloService" // service name
const grpcCacheServiceName = "srcache/grpc.CacheService"

func main() {
	client, err := rpc.Dial("tcp", "localhost:3001")
	if err != nil {
		fmt.Println("dialing:", err)
	}

	// var reply string
	// err = client.Call(HelloServiceName+".Hello", "zhuyanxi", &reply)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(reply)

	msg := &grpc.Request{
		// Key: "49.198.100.160",
		Key: "165.183.1.112",
		// Key: "114.224.139.201",
	}
	var reply grpc.Response
	err = client.Call(grpcCacheServiceName+".Get", msg, &reply)
	if err != nil {
		fmt.Println(err)
	}
	// data, err := proto.Marshal(&reply)
	// if err != nil {
	// 	fmt.Println("marshaling error: ", err)
	// 	return
	// }
	// fmt.Println(string(data))
	fmt.Println(string(reply.Value))
}
