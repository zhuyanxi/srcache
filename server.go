package srcache

import (
	"net/http"
	"strings"
)

const defaultPathPrefix = "/_srcache/"

// Server :
type Server struct {
	localURL   string
	pathPrefix string

	cache *SRCache

	// callback is the function that will query data from a outside data source
	// when the cache is missed
	callback callbackFunc
}

type callbackFunc func(key string) ([]byte, error)

// NewServer :
func NewServer(localURL string, cacheCapacity uint, callback callbackFunc) *Server {
	return &Server{
		localURL:   localURL,
		pathPrefix: defaultPathPrefix,
		cache:      NewSRCache(uint(cacheCapacity)),
		callback:   callback,
	}
}

// Serve :
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, s.pathPrefix) {
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

	key := r.URL.Path[len(s.pathPrefix):]

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
