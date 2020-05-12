package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/zhuyanxi/srcache/grpc"
)

const HelloServiceName = "srcache/grpc.HelloService" // service name

// HelloServiceInterface : method list
type HelloServiceInterface = interface {
	Hello(request string, reply *string) error
}

// RegisterHelloService : regist hello service
func RegisterHelloService(svc HelloServiceInterface) error {
	return rpc.RegisterName(HelloServiceName, svc)
}

type HelloService struct{}

func (p *HelloService) Hello(request string, reply *string) error {
	*reply = "hello:" + request
	return nil
}

func main() {
	test := &grpc.Request{
		Key: "192.168.1.11",
	}
	data, err := proto.Marshal(test)
	if err != nil {
		fmt.Println("marshaling error: ", err)
	}
	newTest := &grpc.Response{}
	err = proto.Unmarshal(data, newTest)
	if err != nil {
		fmt.Println("unmarshaling error: ", err)
	}
	fmt.Println(string(newTest.Value))

	// rpc.RegisterName("HelloService", new(HelloService))
	RegisterHelloService(new(HelloService))

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Listen TCP error:", err)
	}

	counter := 0
	go func() {
		for {
			counter++
			fmt.Println("RPC Server Start at port 1234--Counter:", counter)
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Accept error:", err)
			}

			rpc.ServeConn(conn)
		}
	}()

	// All URLs will be handled by this function
	// http.HandleFunc uses the DefaultServeMux
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	logrus.Infoln("HTTP Server is running at: 8080")
	// Continue to process new requests until an error occurs
	logrus.Fatalln(http.ListenAndServe(":8080", nil))

}
