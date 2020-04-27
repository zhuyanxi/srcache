package srcache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/zhuyanxi/srcache/consistenthash"
)

const defaultLocalURL = "localhost:9099"
const defaultPathPrefix = "/_srcache/"
const defaultReplicate = 50
const defaultCapacity = 10000

// Server :
type Server struct {
	// localURL string

	// pathPrefix string

	// replicate int

	opts ServerOptions

	cache *SRCache

	peers *consistenthash.ConsistentHash

	hashFunc consistenthash.HashFunc

	// callback is the function that will query data from a outside data source
	// when the cache is missed
	callback callbackFunc
}

// ServerOptions :
type ServerOptions struct {
	LocalURL string

	CacheCapacity uint

	PathPrefix string

	Replicate int
}

type callbackFunc func(key string) ([]byte, error)

// NewServer :
func NewServer(callback callbackFunc, hashFunc consistenthash.HashFunc, opts *ServerOptions) *Server {
	server := new(Server)

	if opts == nil {
		defaultOpts := ServerOptions{
			LocalURL:      defaultLocalURL,
			CacheCapacity: defaultCapacity,
			PathPrefix:    defaultPathPrefix,
			Replicate:     defaultReplicate,
		}
		server.opts = defaultOpts
	} else {
		server.opts = *opts
	}

	server.callback = callback
	server.hashFunc = hashFunc
	server.cache = NewSRCache(uint(server.opts.CacheCapacity))

	return server
}

// SetPeers :
func (s *Server) SetPeers(peers ...consistenthash.Node) {
	s.peers = consistenthash.NewConsistentHash(s.opts.Replicate, s.hashFunc)
	s.peers.Add(peers...)
}

// GetPeer :
func (s *Server) GetPeer(key string) consistenthash.Node {
	peer := s.peers.Get(key)
	return peer
}

// Serve :
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, s.opts.PathPrefix) {
		http.Error(w, "bad request: "+"Serving unexpected path: "+r.URL.Path, http.StatusBadRequest)
		//panic("Serving unexpected path: " + r.URL.Path)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	// /_srcache/key
	// parts := strings.SplitN(r.URL.Path[len(s.pathPrefix):], "/", 2)
	// if len(parts) != 2 {
	// 	http.Error(w, "bad request", http.StatusBadRequest)
	// 	return
	// }

	key := r.URL.Path[len(s.opts.PathPrefix):]

	peer := s.GetPeer(key)
	logrus.Infof("LocalAddr: %s. Query from server %s: .\n", s.opts.LocalURL, peer.Addr)
	if peer.Addr != s.opts.LocalURL {
		dataFromPeer, errPeer := s.getFromPeer(key)
		if errPeer == nil {
			w.Write(dataFromPeer)
		} else {
			//w.Write([]byte(errPeer.Error()))
			http.Error(w, "bad request: "+errPeer.Error(), http.StatusBadRequest)
		}
		return
	}

	//capacity := 8
	//cache := NewSRCache(uint(capacity))
	data, ok := s.cache.Get(key)
	if ok {
		w.Write(data)
	} else {
		dataSource, err := s.callback(key)
		if err == nil {
			s.cache.Add(key, dataSource)
			w.Write(dataSource)
		} else {
			//w.Write([]byte(fmt.Sprintf("%s not exist", key)))
			w.Write([]byte(err.Error()))
		}
	}
}

func (s *Server) getFromPeer(key string) ([]byte, error) {
	peer := s.GetPeer(key)
	url := fmt.Sprintf("http://%v%v%v", peer.Addr, s.opts.PathPrefix, key)
	// url := fmt.Sprintf("%v/%v", url.QueryEscape(peer.Addr), url.QueryEscape(key))
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get from peer %s: %v", peer.Addr, res.Status)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from peer %s: %v", peer.Addr, res.Status)
	}

	return resBody, nil
}
