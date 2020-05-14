package srcache

import (
	"fmt"
	"net"
	"net/rpc"

	"github.com/sirupsen/logrus"
	"github.com/zhuyanxi/srcache/grpc"
)

const grpcCacheServiceName = "srcache/grpc.CacheService"

type grpcCacheServiceInterface = interface {
	Get(in *grpc.Request, out *grpc.Response) error
}

func registerCacheService(ics grpcCacheServiceInterface) error {
	return rpc.RegisterName(grpcCacheServiceName, ics)
}

type grpcService struct {
	server *Server
}

// Get :
func (grpcs *grpcService) Get(in *grpc.Request, out *grpc.Response) error {
	s := grpcs.server
	key := in.Key
	peer := s.GetPeer(key)
	logrus.Infof("LocalAddr: %s. Query from server: %s .\n", s.opts.LocalURL, peer.Addr)
	if peer.Addr != s.opts.LocalURL {
		dataFromPeer, errPeer := s.getFromGrpcPeer(key)
		if errPeer == nil {
			// err := proto.Unmarshal(dataFromPeer, out)
			// if err != nil {
			// 	fmt.Println("unmarshaling error: ", err)
			// 	return err
			// }
			logrus.Infoln("Data from peer(", peer.Addr, "):", string(dataFromPeer))
			out.Value = dataFromPeer
		} else {
			// http.Error(w, "bad request: "+errPeer.Error(), http.StatusBadRequest)
			out.Value = []byte(errPeer.Error())
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
			fmt.Println(string(dataSource))
			out.Value = dataSource
			// err := proto.Unmarshal(dataSource, out)
			// if err != nil {
			// 	fmt.Println("unmarshaling error: ", err)
			// 	return err
			// }
		}
	}

	return nil
}

// Get :
// func (s *Server) Get(in *grpc.Request, out *grpc.Response) error {
// 	key := in.Key
// 	peer := s.GetPeer(key)
// 	logrus.Infof("LocalAddr: %s. Query from server: %s .\n", s.opts.LocalURL, peer.Addr)
// 	if peer.Addr != s.opts.LocalURL {
// 		dataFromPeer, errPeer := s.getFromGrpcPeer(key)
// 		if errPeer == nil {
// 			// err := proto.Unmarshal(dataFromPeer, out)
// 			// if err != nil {
// 			// 	fmt.Println("unmarshaling error: ", err)
// 			// 	return err
// 			// }
// 			logrus.Infoln("Data from peer(", peer.Addr, "):", string(dataFromPeer))
// 			out.Value = dataFromPeer
// 		} else {
// 			// http.Error(w, "bad request: "+errPeer.Error(), http.StatusBadRequest)
// 			out.Value = []byte(errPeer.Error())
// 		}
// 		return nil
// 	}

// 	data, ok := s.cache.Get(key)
// 	if ok {
// 		// err := proto.Unmarshal(data, out)
// 		// if err != nil {
// 		// 	fmt.Println("unmarshaling error: ", err)
// 		// 	return err
// 		// }
// 		out.Value = data
// 	} else {
// 		dataSource, err := s.callback(key)
// 		if err == nil {
// 			s.cache.Add(key, dataSource)
// 			fmt.Println(string(dataSource))
// 			out.Value = dataSource
// 			// err := proto.Unmarshal(dataSource, out)
// 			// if err != nil {
// 			// 	fmt.Println("unmarshaling error: ", err)
// 			// 	return err
// 			// }
// 		}
// 	}

// 	return nil
// }

// ServeGRPC :
func (s *Server) ServeGRPC(port string) {
	grpcs := &grpcService{
		server: s,
	}
	registerCacheService(grpcs)

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
			continue
		}
		fmt.Println("Accept from ", conn.RemoteAddr().String(), " success")
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
		return nil, err
	}
	defer client.Close()

	req := &grpc.Request{
		Key: key,
	}

	var reply grpc.Response
	err = client.Call(grpcCacheServiceName+".Get", req, &reply)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("Peer (", peer.Addr, ") return value:", string(reply.Value))
	return reply.Value, nil

	// data, err := proto.Marshal(&reply)
	// if err != nil {
	// 	fmt.Println("marshaling error: ", err)
	// 	return nil, err
	// }
	// fmt.Println("Peer (", peer.Addr, ") return marshaled proto value:", string(data))
	// return data, nil
}
