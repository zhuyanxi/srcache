package srcache

import (
	"fmt"
	"net"
	"net/rpc"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/zhuyanxi/srcache/grpc"
)

const grpcCacheServiceName = "srcache/grpc.CacheService"

type grpcCacheServiceInterface = interface {
	Get(request *grpc.Request, response *grpc.Response) error
}

func registerCacheService(ics grpcCacheServiceInterface) error {
	return rpc.RegisterName(grpcCacheServiceName, ics)
}

// Get :
func (s *Server) Get(in *grpc.Request, out *grpc.Response) error {
	key := in.Key
	peer := s.GetPeer(key)
	logrus.Infof("LocalAddr: %s. Query from server: %s .\n", s.opts.LocalURL, peer.Addr)
	if peer.Addr != s.opts.LocalURL {
		dataFromPeer, errPeer := s.getFromGrpcPeer(key)
		if errPeer == nil {
			// w.Write(dataFromPeer)
			// err := proto.Unmarshal(dataFromPeer, out)
			// if err != nil {
			// 	fmt.Println("unmarshaling error: ", err)
			// 	return err
			// }
			out.Value = dataFromPeer
		} else {
			// http.Error(w, "bad request: "+errPeer.Error(), http.StatusBadRequest)
		}
		return nil
	}

	data, ok := s.cache.Get(key)
	if ok {
		// err := proto.Unmarshal(data, out)
		// if err != nil {
		// 	fmt.Println("unmarshaling error: ", err)
		// 	return err
		// }
		out.Value = data
	} else {
		dataSource, err := s.callback(key)
		if err == nil {
			s.cache.Add(key, dataSource)
			out.Value = dataSource
			// err := proto.Unmarshal(data, out)
			// if err != nil {
			// 	fmt.Println("unmarshaling error: ", err)
			// 	return err
			// }
		}
	}

	return nil
}

// ServeGRPC :
func (s *Server) ServeGRPC(port string) {
	// grpc.InitServer()
	registerCacheService(s)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Listen TCP error:", err)
	}

	counter := 0
	// go func() {
	for {
		counter++
		fmt.Println("RPC Server Start at port ", port, "--Counter:", counter)
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
		}

		rpc.ServeConn(conn)
	}
	// }()
}

func (s *Server) getFromGrpcPeer(key string) ([]byte, error) {
	peer := s.GetPeer(key)
	addr := peer.Addr
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		fmt.Println("dialing:", err)
	}

	req := &grpc.Request{
		Key: key,
	}

	var reply grpc.Response
	err = client.Call(grpcCacheServiceName+".Get", req, &reply)
	if err != nil {
		fmt.Println(err)
	}

	data, err := proto.Marshal(&reply)
	if err != nil {
		fmt.Println("marshaling error: ", err)
		return nil, err
	}

	return data, nil
}
