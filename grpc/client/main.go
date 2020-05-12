package main

import (
	"fmt"
	"net/rpc"
)

const HelloServiceName = "srcache/grpc.HelloService" // service name

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("dialing:", err)
	}

	var reply string
	err = client.Call(HelloServiceName+".Hello", "zhuyanxi", &reply)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(reply)
}
