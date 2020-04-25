package srcache

import (
	"net/http"
	"strings"

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
	localURL string

	cacheCapacity uint

	pathPrefix string

	replicate int
}

type callbackFunc func(key string) ([]byte, error)

// NewServer :
func NewServer(callback callbackFunc, hashFunc consistenthash.HashFunc, opts *ServerOptions) *Server {
	server := new(Server)

	if opts == nil {
		defaultOpts := ServerOptions{
			localURL:      defaultLocalURL,
			cacheCapacity: defaultCapacity,
			pathPrefix:    defaultPathPrefix,
			replicate:     defaultReplicate,
		}
		server.opts = defaultOpts
	} else {
		server.opts = *opts
	}

	server.callback = callback
	server.hashFunc = hashFunc
	server.cache = NewSRCache(uint(server.opts.cacheCapacity))

	return server
}

// SetPeers :
func (s *Server) SetPeers(peers ...consistenthash.Node) {
	s.peers = consistenthash.NewConsistentHash(s.opts.replicate, s.hashFunc)
	// s.peers.Add()
	s.peers.Add(peers...)
}

// Serve :
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, s.opts.pathPrefix) {
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

	key := r.URL.Path[len(s.opts.pathPrefix):]

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
